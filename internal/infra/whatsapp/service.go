package whatsapp

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"

	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/domain/services"
	"wazmeow/pkg/logger"
)

// Service implements the WhatsApp service
type Service struct {
	sessionRepo   repositories.SessionRepository
	container     *sqlstore.Container
	clientManager *ClientManager
	mu            sync.RWMutex
}

// NewService creates a new WhatsApp service
func NewService(sessionRepo repositories.SessionRepository, container *sqlstore.Container) *Service {
	return &Service{
		sessionRepo:   sessionRepo,
		container:     container,
		clientManager: NewClientManager(),
	}
}

// StartSession starts a WhatsApp session
func (s *Service) StartSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if session already exists
	if s.clientManager.HasSession(sessionID) {
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

	// Create device store like wuzapi does
	var deviceStore *store.Device

	// First try to get existing device if we have a JID
	if session.DeviceJID != "" {
		jid, parseErr := types.ParseJID(session.DeviceJID)
		if parseErr == nil {
			deviceStore, err = s.container.GetDevice(ctx, jid)
			if err != nil {
				logger.Warn().Err(err).Str("sessionID", sessionID).Str("jid", session.DeviceJID).Msg("Failed to get existing device, creating new one")
				deviceStore = s.container.NewDevice()
			} else {
				logger.Info().Str("sessionID", sessionID).Str("jid", session.DeviceJID).Msg("Reusing existing device")
			}
		} else {
			logger.Warn().Err(parseErr).Str("sessionID", sessionID).Str("jid", session.DeviceJID).Msg("Invalid JID, creating new device")
			deviceStore = s.container.NewDevice()
		}
	} else {
		logger.Info().Str("sessionID", sessionID).Msg("No JID found, creating new device")
		deviceStore = s.container.NewDevice()
	}

	if deviceStore == nil {
		logger.Warn().Str("sessionID", sessionID).Msg("No store found, creating new one")
		deviceStore = s.container.NewDevice()
	}

	client := whatsmeow.NewClient(deviceStore, nil)

	// Create MyClient wrapper
	myClient := NewMyClient(client, sessionID, "", nil)

	// Store in client manager
	s.clientManager.SetMyClient(sessionID, myClient)

	// Start connection process (MyClient already has its own event handler)
	go s.startConnection(ctx, sessionID, myClient)

	return nil
}

// StopSession stops a WhatsApp session
func (s *Service) StopSession(ctx context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	myClient := s.clientManager.GetMyClient(sessionID)
	if myClient == nil {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Send kill signal
	myClient.Kill()

	// Disconnect client
	myClient.Disconnect()

	// Remove from client manager
	s.clientManager.DeleteMyClient(sessionID)

	// Update session status
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err == nil {
		session.Status = entities.StatusDisconnected
		session.DeviceJID = ""
		if err := s.sessionRepo.Update(ctx, session); err != nil {
			logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to update session status to disconnected")
		}
	}

	logger.Info().Str("sessionID", sessionID).Msg("Session stopped")
	return nil
}

// GetQRCode gets the QR code for session authentication
func (s *Service) GetQRCode(ctx context.Context, sessionID string) (string, error) {
	s.mu.RLock()
	myClient := s.clientManager.GetMyClient(sessionID)
	s.mu.RUnlock()

	// Check if session exists (like wuzapi)
	if myClient == nil {
		return "", fmt.Errorf("no session")
	}

	// Check if already logged in (like wuzapi)
	if myClient.WAClient.Store.ID != nil {
		return "", fmt.Errorf("already logged in")
	}

	// Get session from database to return stored QR code (like wuzapi)
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	// Return the stored QR code (base64 encoded PNG)
	return session.QRCode, nil
}

// PairPhone pairs a phone number with the session
func (s *Service) PairPhone(ctx context.Context, sessionID, phone string) (string, error) {
	s.mu.RLock()
	myClient := s.clientManager.GetMyClient(sessionID)
	s.mu.RUnlock()

	if myClient == nil {
		return "", fmt.Errorf("session %s not found", sessionID)
	}

	linkingCode, err := myClient.PairPhone(ctx, phone)
	if err != nil {
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	logger.Info().Str("sessionID", sessionID).Str("phone", phone).Str("linkingCode", linkingCode).Msg("Phone pairing initiated")
	return linkingCode, nil
}

// Logout logs out from WhatsApp
func (s *Service) Logout(ctx context.Context, sessionID string) error {
	s.mu.RLock()
	myClient := s.clientManager.GetMyClient(sessionID)
	s.mu.RUnlock()

	if myClient == nil {
		return fmt.Errorf("session %s not found", sessionID)
	}

	err := myClient.Logout(ctx)
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
	s.mu.RLock()
	defer s.mu.RUnlock()

	myClient := s.clientManager.GetMyClient(sessionID)
	if myClient == nil {
		return false
	}

	return myClient.IsConnected()
}

// IsLoggedIn checks if a session is logged in
func (s *Service) IsLoggedIn(sessionID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	myClient := s.clientManager.GetMyClient(sessionID)
	if myClient == nil {
		return false
	}

	return myClient.IsLoggedIn()
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
	s.mu.RLock()
	myClient := s.clientManager.GetMyClient(sessionID)
	s.mu.RUnlock()

	info := &services.SessionInfo{
		SessionID: sessionID,
		Connected: false,
		LoggedIn:  false,
	}

	if myClient != nil {
		info.Connected = myClient.IsConnected()
		info.LoggedIn = myClient.IsLoggedIn()

		if myClient.WAClient.Store.ID != nil {
			info.DeviceJID = myClient.WAClient.Store.ID.String()
		}
	}

	return info, nil
}

// SetSubscriptions sets event subscriptions for a session
func (s *Service) SetSubscriptions(sessionID string, subscriptions []string) error {
	s.mu.RLock()
	myClient := s.clientManager.GetMyClient(sessionID)
	s.mu.RUnlock()

	if myClient == nil {
		return fmt.Errorf("session %s not found", sessionID)
	}

	myClient.SetSubscriptions(subscriptions)
	return nil
}

// GetSubscriptions gets event subscriptions for a session
func (s *Service) GetSubscriptions(sessionID string) ([]string, error) {
	s.mu.RLock()
	myClient := s.clientManager.GetMyClient(sessionID)
	s.mu.RUnlock()

	if myClient == nil {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	return myClient.GetSubscriptions(), nil
}

// AddSubscription adds a single event subscription to a session
func (s *Service) AddSubscription(sessionID, eventType string) error {
	s.mu.RLock()
	myClient := s.clientManager.GetMyClient(sessionID)
	s.mu.RUnlock()

	if myClient == nil {
		return fmt.Errorf("session %s not found", sessionID)
	}

	myClient.AddSubscription(eventType)
	return nil
}

// RemoveSubscription removes a single event subscription from a session
func (s *Service) RemoveSubscription(sessionID, eventType string) error {
	s.mu.RLock()
	myClient := s.clientManager.GetMyClient(sessionID)
	s.mu.RUnlock()

	if myClient == nil {
		return fmt.Errorf("session %s not found", sessionID)
	}

	myClient.RemoveSubscription(eventType)
	return nil
}

// GetSupportedEventTypes returns list of supported event types
func (s *Service) GetSupportedEventTypes() []string {
	// Return static list - all sessions support the same event types
	return []string{
		"All",
		"Connected",
		"Disconnected",
		"Message",
		"PairSuccess",
		"LoggedOut",
		"ReadReceipt",
		"Presence",
		"ConnectFailure",
		"QR",
	}
}

// GetAllSessionsInfo returns information about all active sessions
func (s *Service) GetAllSessionsInfo() []map[string]interface{} {
	s.mu.RLock()
	sessions := s.clientManager.GetAllSessions()
	s.mu.RUnlock()

	var result []map[string]interface{}
	for _, sessionID := range sessions {
		myClient := s.clientManager.GetMyClient(sessionID)
		if myClient != nil {
			sessionInfo := map[string]interface{}{
				"sessionID":     sessionID,
				"connected":     myClient.IsConnected(),
				"loggedIn":      myClient.IsLoggedIn(),
				"subscriptions": myClient.GetSubscriptions(),
				"webhook":       myClient.GetWebhook(),
			}

			if myClient.WAClient.Store.ID != nil {
				sessionInfo["deviceJID"] = myClient.WAClient.Store.ID.String()
			}

			result = append(result, sessionInfo)
		}
	}

	return result
}

// startConnection handles the connection process for a session
func (s *Service) startConnection(ctx context.Context, sessionID string, myClient *MyClient) {
	logger.Info().Str("sessionID", sessionID).Msg("Starting WhatsApp connection")

	client := myClient.WAClient
	if client.Store.ID == nil {
		// No ID stored, new login - need QR code
		qrChan, err := client.GetQRChannel(ctx)
		if err != nil {
			logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to get QR channel")
			s.updateSessionStatus(ctx, sessionID, entities.StatusDisconnected)
			return
		}

		// Connect FIRST like wuzapi does (Si no conectamos no se puede generar QR)
		err = client.Connect()
		if err != nil {
			logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to connect client")
			s.updateSessionStatus(ctx, sessionID, entities.StatusDisconnected)
			return
		}

		// Store QR channel in MyClient
		myClient.SetQRChannel(qrChan)

		// Handle QR events
		s.handleQREvents(ctx, sessionID, qrChan)
	} else {
		// Already logged in, just connect
		logger.Info().Str("sessionID", sessionID).Msg("Already logged in, connecting...")
		err := client.Connect()
		if err != nil {
			logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to connect")
			s.updateSessionStatus(ctx, sessionID, entities.StatusDisconnected)
			return
		}
	}

	// Keep connection alive
	s.keepAlive(sessionID, myClient)
}

// handleQREvents processes QR code events
func (s *Service) handleQREvents(_ context.Context, sessionID string, qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		switch evt.Event {
		case "code":
			logger.Info().Str("sessionID", sessionID).Msg("QR code generated")

			// Display QR in terminal for development
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Printf("QR code for session %s:\n%s\n", sessionID, evt.Code)

			// Generate base64 PNG
			qrPNG, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
			if err != nil {
				logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to encode QR code")
				continue
			}

			base64QR := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrPNG)

			// Update QR code in database using background context
			if err := s.updateQRCodeInDB(sessionID, base64QR); err != nil {
				logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to update QR code in database")
			} else {
				logger.Info().Str("sessionID", sessionID).Msg("QR code saved to database successfully")
			}

		case "timeout":
			logger.Warn().Str("sessionID", sessionID).Msg("QR code timeout")
			bgCtx := context.Background()
			s.updateSessionStatus(bgCtx, sessionID, entities.StatusDisconnected)
			s.StopSession(bgCtx, sessionID)
			return

		case "success":
			logger.Info().Str("sessionID", sessionID).Msg("QR pairing successful")
			// Clear QR code and update status
			bgCtx := context.Background()
			session, err := s.sessionRepo.GetByID(bgCtx, sessionID)
			if err == nil {
				session.QRCode = ""
				session.Status = entities.StatusConnected
				s.sessionRepo.Update(bgCtx, session)
				logger.Info().Str("sessionID", sessionID).Msg("Session status updated to connected")
			}

		default:
			logger.Info().Str("sessionID", sessionID).Str("event", evt.Event).Msg("QR event")
		}
	}
}

// keepAlive keeps the connection alive until killed
func (s *Service) keepAlive(sessionID string, myClient *MyClient) {
	killChan := myClient.KillChannel()

	for {
		select {
		case <-killChan:
			logger.Info().Str("sessionID", sessionID).Msg("Received kill signal")
			return
		case <-time.After(30 * time.Second):
			// Keep alive ping
			continue
		}
	}
}

// updateSessionStatus updates the session status in the database
func (s *Service) updateSessionStatus(ctx context.Context, sessionID string, status entities.SessionStatus) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to get session for status update")
		return
	}

	session.Status = status
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		logger.Error().Err(err).Str("sessionID", sessionID).Str("status", string(status)).Msg("Failed to update session status")
	}
}

// updateQRCodeInDB updates QR code directly in database like wuzapi
func (s *Service) updateQRCodeInDB(sessionID, qrCode string) error {
	// Use background context to avoid cancellation
	ctx := context.Background()

	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	session.QRCode = qrCode
	return s.sessionRepo.Update(ctx, session)
}
