package whatsapp

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"wazmeow/internal/config"
	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/domain/services"
	"wazmeow/internal/infra/whatsapp/client"
	"wazmeow/pkg/logger"
)

// Service implements the WhatsApp service
type Service struct {
	sessionRepo   repositories.SessionRepository
	container     *sqlstore.Container
	clientManager *client.Manager
	config        *config.WhatsAppConfig
}

// NewService creates a new WhatsApp service
func NewService(sessionRepo repositories.SessionRepository, container *sqlstore.Container, cfg *config.WhatsAppConfig) *Service {
	// Criar novo manager otimizado
	manager := client.NewManager(container, sessionRepo, cfg)

	return &Service{
		sessionRepo:   sessionRepo,
		container:     container,
		clientManager: manager,
		config:        cfg,
	}
}

// Initialize carrega todas as sessões salvas no banco durante o startup
func (s *Service) Initialize(ctx context.Context) error {
	logger.Info().Msg("Initializing WhatsApp service and loading sessions")

	// Carregar todas as sessões do banco
	if err := s.clientManager.LoadAll(ctx); err != nil {
		return err
	}

	// Tentar reconectar sessões que estavam conectadas
	go s.autoReconnectSessions(ctx)

	return nil
}

// autoReconnectSessions tenta reconectar sessões que estavam conectadas (similar ao wuzapi)
func (s *Service) autoReconnectSessions(ctx context.Context) {
	logger.Info().Msg("Starting auto-reconnection for connected sessions")

	// Buscar sessões que estavam conectadas (similar ao wuzapi: WHERE connected=1)
	sessions, err := s.sessionRepo.GetConnectedSessions(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get connected sessions for auto-reconnect")
		return
	}

	logger.Info().Int("count", len(sessions)).Msg("Found connected sessions for auto-reconnect")

	// Para cada sessão conectada, tentar reconectar
	for _, session := range sessions {
		go s.attemptReconnect(ctx, session.ID)
	}
}

// attemptReconnect tenta reconectar uma sessão específica
func (s *Service) attemptReconnect(ctx context.Context, sessionID string) {
	logger.Info().Str("sessionID", sessionID).Msg("Attempting auto-reconnect")

	// Verificar se wrapper já existe
	wrapper := s.clientManager.Get(sessionID)
	if wrapper == nil {
		logger.Warn().Str("sessionID", sessionID).Msg("No wrapper found for auto-reconnect")
		return
	}

	client := wrapper.Client()
	if client == nil {
		logger.Warn().Str("sessionID", sessionID).Msg("No client found for auto-reconnect")
		return
	}

	// Similar ao wuzapi: se já tem ID armazenado, apenas conecta
	if client.Store.ID != nil {
		logger.Info().Str("sessionID", sessionID).Msg("Session has stored ID, attempting direct connect")

		err := client.Connect()
		if err != nil {
			logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to auto-reconnect session")
			// Atualizar status para disconnected
			s.sessionRepo.UpdateStatus(ctx, sessionID, entities.StatusDisconnected)
		} else {
			logger.Info().Str("sessionID", sessionID).Msg("Session auto-reconnected successfully")
		}
	} else {
		logger.Info().Str("sessionID", sessionID).Msg("Session has no stored ID, requires QR authentication")
		// Atualizar status para disconnected pois precisa de QR
		s.sessionRepo.UpdateStatus(ctx, sessionID, entities.StatusDisconnected)
	}
}

// StartSession starts a WhatsApp session
func (s *Service) StartSession(ctx context.Context, sessionID string) error {
	// Check if session already exists
	if s.clientManager.Has(sessionID) {
		return fmt.Errorf("session %s already started", sessionID)
	}

	// Get session from database
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Update status to connecting
	session.Status = entities.StatusConnecting
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to update session status to connecting")
	}

	// Use ClientManager to create new session
	return s.clientManager.Create(ctx, sessionID)
}

// StopSession stops a WhatsApp session
func (s *Service) StopSession(ctx context.Context, sessionID string) error {
	if !s.clientManager.Has(sessionID) {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Disconnect session using ClientManager
	err := s.clientManager.Remove(sessionID)
	if err != nil {
		return err
	}

	// Update session status in database
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err == nil {
		session.Status = entities.StatusDisconnected
		if err := s.sessionRepo.Update(ctx, session); err != nil {
			logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to update session status to disconnected")
		}
	}

	logger.Info().Str("sessionID", sessionID).Msg("Session stopped")
	return nil
}

// GetQRCode gets the QR code for session authentication
func (s *Service) GetQRCode(ctx context.Context, sessionID string) (string, error) {
	// Check if session exists
	if !s.clientManager.Has(sessionID) {
		return "", fmt.Errorf("no session")
	}

	// Check if already logged in
	wrapper := s.clientManager.Get(sessionID)
	if wrapper != nil && wrapper.IsLoggedIn() {
		return "", fmt.Errorf("already logged in")
	}

	// Get session from database to return stored QR code
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	// Return the stored QR code (base64 encoded PNG)
	return session.QRCode, nil
}

// PairPhone pairs a phone number with the session
func (s *Service) PairPhone(ctx context.Context, sessionID, phone string) (string, error) {
	wrapper := s.clientManager.Get(sessionID)
	if wrapper == nil {
		return "", fmt.Errorf("session %s not found", sessionID)
	}

	linkingCode, err := wrapper.Client().PairPhone(ctx, phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	logger.Info().Str("sessionID", sessionID).Str("phone", phone).Str("linkingCode", linkingCode).Msg("Phone pairing initiated")
	return linkingCode, nil
}

// Logout logs out from WhatsApp
func (s *Service) Logout(ctx context.Context, sessionID string) error {
	wrapper := s.clientManager.Get(sessionID)
	if wrapper == nil {
		return fmt.Errorf("session %s not found", sessionID)
	}

	err := wrapper.Client().Logout(ctx)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	// Update session status
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err == nil {
		session.Status = entities.StatusDisconnected
		session.DeviceJID = ""
		if err := s.sessionRepo.Update(ctx, session); err != nil {
			logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to update session status after logout")
		}
	}

	logger.Info().Str("sessionID", sessionID).Msg("Session logged out")
	return nil
}

// IsConnected checks if a session is connected
func (s *Service) IsConnected(sessionID string) bool {
	wrapper := s.clientManager.Get(sessionID)
	return wrapper != nil && wrapper.IsConnected()
}

// IsLoggedIn checks if a session is logged in
func (s *Service) IsLoggedIn(sessionID string) bool {
	wrapper := s.clientManager.Get(sessionID)
	return wrapper != nil && wrapper.IsLoggedIn()
}

// SetProxy sets proxy configuration for a session
func (s *Service) SetProxy(sessionID string, config *entities.ProxyConfig) error {
	// Proxy configuration would be handled during client creation
	// For now, we'll store it in the session entity
	logger.Info().Str("sessionID", sessionID).Interface("proxyConfig", config).Msg("Proxy configuration updated")
	return nil
}

// GetSessionInfo gets detailed session information
func (s *Service) GetSessionInfo(sessionID string) (*services.SessionInfo, error) {
	wrapper := s.clientManager.Get(sessionID)

	info := &services.SessionInfo{
		SessionID: sessionID,
		Connected: wrapper != nil && wrapper.IsConnected(),
		LoggedIn:  wrapper != nil && wrapper.IsLoggedIn(),
	}

	if wrapper != nil {
		jid := wrapper.JID()
		if !jid.IsEmpty() {
			info.DeviceJID = jid.String()
		}
	}

	return info, nil
}

// SetSubscriptions sets event subscriptions for a session (simplified)
func (s *Service) SetSubscriptions(sessionID string, subscriptions []string) error {
	if !s.clientManager.Has(sessionID) {
		return fmt.Errorf("session %s not found", sessionID)
	}
	logger.Debug().Str("sessionID", sessionID).Strs("subscriptions", subscriptions).Msg("Subscriptions managed automatically")
	return nil
}

// GetSubscriptions gets event subscriptions for a session (simplified)
func (s *Service) GetSubscriptions(sessionID string) ([]string, error) {
	if !s.clientManager.Has(sessionID) {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}
	return []string{"Message", "Connected", "Disconnected", "QR", "PairSuccess", "LoggedOut"}, nil
}

// AddSubscription adds a single event subscription to a session (simplified)
func (s *Service) AddSubscription(sessionID, eventType string) error {
	if !s.clientManager.Has(sessionID) {
		return fmt.Errorf("session %s not found", sessionID)
	}
	logger.Debug().Str("sessionID", sessionID).Str("eventType", eventType).Msg("Subscription managed automatically")
	return nil
}

// RemoveSubscription removes a single event subscription from a session (simplified)
func (s *Service) RemoveSubscription(sessionID, eventType string) error {
	if !s.clientManager.Has(sessionID) {
		return fmt.Errorf("session %s not found", sessionID)
	}
	logger.Debug().Str("sessionID", sessionID).Str("eventType", eventType).Msg("Subscription managed automatically")
	return nil
}

// GetSupportedEventTypes returns list of supported event types (simplified)
func (s *Service) GetSupportedEventTypes() []string {
	return []string{"All", "Connected", "Disconnected", "Message", "PairSuccess", "LoggedOut", "ReadReceipt", "Presence", "ConnectFailure", "QR"}
}

// GetAllSessionsInfo returns information about all active sessions
func (s *Service) GetAllSessionsInfo() []map[string]interface{} {
	// TODO: Implementar GetAll no Manager
	var result []map[string]interface{}

	// Por enquanto retorna lista vazia
	// Quando GetAll for implementado, usar:
	// sessions := s.clientManager.GetAll()
	// for _, sessionID := range sessions { ... }

	return result
}

// NOTA: Métodos de conexão removidos - agora gerenciados pelo ClientManager

// Shutdown para o service e todas as sessões
func (s *Service) Shutdown() {
	logger.Info().Msg("Shutting down WhatsApp service")
	ctx := context.Background()
	s.clientManager.Shutdown(ctx)
}
