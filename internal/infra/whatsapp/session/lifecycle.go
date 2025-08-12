package session

import (
	"context"
	"fmt"
	"time"

	"wazmeow/internal/config"
	"wazmeow/internal/domain/entities"
	"wazmeow/pkg/logger"
)

// Lifecycle gerencia ciclo de vida de sessões com timeouts e reconexão
type Lifecycle struct {
	config *config.WhatsAppConfig
}

// NewLifecycle cria um novo gerenciador de lifecycle
func NewLifecycle(config *config.WhatsAppConfig) *Lifecycle {
	return &Lifecycle{
		config: config,
	}
}

// Start inicia o ciclo de vida de uma sessão
func (l *Lifecycle) Start(ctx context.Context, wrapper WrapperInterface) {
	sessionID := wrapper.SessionID()
	logger.Info().Str("sessionID", sessionID).Msg("Starting session lifecycle")

	// Context com timeout para conexão inicial
	connectCtx, cancel := context.WithTimeout(ctx, l.config.ConnectionTimeout)
	defer cancel()

	// Tentar conectar
	if err := l.connect(connectCtx, wrapper); err != nil {
		logger.Error().Str("sessionID", sessionID).Err(err).Msg("Failed to connect session")
		wrapper.SetStatus(entities.StatusDisconnected)
		return
	}

	// Monitorar conexão
	l.monitor(ctx, wrapper)
}

// connect estabelece conexão com WhatsApp
func (l *Lifecycle) connect(ctx context.Context, wrapper WrapperInterface) error {
	sessionID := wrapper.SessionID()
	client := wrapper.Client()

	wrapper.SetStatus(entities.StatusConnecting)
	logger.Info().Str("sessionID", sessionID).Msg("Connecting to WhatsApp")

	// Verificar se já está logado
	if client.GetStore().GetID() == nil {
		// Não está logado - precisa de QR code ou pairing
		logger.Info().Str("sessionID", sessionID).Msg("No stored ID, needs authentication")
		return l.handleAuthentication(ctx, wrapper)
	}

	// Já está logado - apenas conectar
	logger.Info().Str("sessionID", sessionID).Msg("Stored ID found, connecting directly")
	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	wrapper.SetConnected(true)
	wrapper.SetLoggedIn(true)
	wrapper.SetStatus(entities.StatusConnected)

	logger.Info().Str("sessionID", sessionID).Msg("Session connected successfully")
	return nil
}

// handleAuthentication gerencia processo de autenticação
func (l *Lifecycle) handleAuthentication(ctx context.Context, wrapper WrapperInterface) error {
	sessionID := wrapper.SessionID()

	// Context com timeout para QR
	qrCtx, cancel := context.WithTimeout(ctx, l.config.QRTimeout)
	defer cancel()

	logger.Info().Str("sessionID", sessionID).Msg("Starting QR authentication process")

	// Conectar para iniciar processo de QR
	client := wrapper.Client()
	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect for QR: %w", err)
	}

	// Aguardar autenticação ou timeout
	select {
	case <-qrCtx.Done():
		logger.Error().Str("sessionID", sessionID).Msg("QR authentication timeout")
		return fmt.Errorf("QR authentication timeout")
	case <-ctx.Done():
		return ctx.Err()
	}
}

// monitor monitora conexão e reconecta se necessário
func (l *Lifecycle) monitor(ctx context.Context, wrapper WrapperInterface) {
	sessionID := wrapper.SessionID()
	ticker := time.NewTicker(l.config.ReconnectInterval)
	defer ticker.Stop()

	reconnectAttempts := 0

	for {
		select {
		case <-ctx.Done():
			logger.Info().Str("sessionID", sessionID).Msg("Lifecycle monitoring stopped")
			return

		case <-ticker.C:
			if !wrapper.IsConnected() && reconnectAttempts < l.config.MaxReconnectAttempts {
				reconnectAttempts++
				logger.Warn().
					Str("sessionID", sessionID).
					Int("attempt", reconnectAttempts).
					Int("maxAttempts", l.config.MaxReconnectAttempts).
					Msg("Attempting to reconnect")

				if err := l.reconnect(ctx, wrapper); err != nil {
					logger.Error().
						Str("sessionID", sessionID).
						Int("attempt", reconnectAttempts).
						Err(err).
						Msg("Reconnection failed")
				} else {
					logger.Info().Str("sessionID", sessionID).Msg("Reconnection successful")
					reconnectAttempts = 0 // Reset counter on success
				}
			} else if reconnectAttempts >= l.config.MaxReconnectAttempts {
				logger.Error().
					Str("sessionID", sessionID).
					Int("maxAttempts", l.config.MaxReconnectAttempts).
					Msg("Max reconnection attempts reached, giving up")
				wrapper.SetStatus(entities.StatusDisconnected)
				return
			}
		}
	}
}

// reconnect tenta reconectar uma sessão
func (l *Lifecycle) reconnect(ctx context.Context, wrapper WrapperInterface) error {
	sessionID := wrapper.SessionID()
	client := wrapper.Client()

	// Context com timeout para reconexão
	_, cancel := context.WithTimeout(ctx, l.config.ConnectionTimeout)
	defer cancel()

	wrapper.SetStatus(entities.StatusConnecting)

	// Tentar reconectar
	if err := client.Connect(); err != nil {
		wrapper.SetStatus(entities.StatusDisconnected)
		return fmt.Errorf("failed to reconnect: %w", err)
	}

	wrapper.SetConnected(true)
	wrapper.SetStatus(entities.StatusConnected)

	logger.Info().Str("sessionID", sessionID).Msg("Session reconnected successfully")
	return nil
}

// Stop para o ciclo de vida de uma sessão
func (l *Lifecycle) Stop(wrapper WrapperInterface) {
	sessionID := wrapper.SessionID()
	logger.Info().Str("sessionID", sessionID).Msg("Stopping session lifecycle")

	wrapper.Disconnect()
	logger.Info().Str("sessionID", sessionID).Msg("Session lifecycle stopped")
}

// WrapperInterface define interface mínima para wrapper
type WrapperInterface interface {
	SessionID() string
	Client() ClientInterface
	IsConnected() bool
	IsLoggedIn() bool
	SetConnected(bool)
	SetLoggedIn(bool)
	SetStatus(entities.SessionStatus)
	Disconnect()
}

// ClientInterface define interface mínima para cliente
type ClientInterface interface {
	Connect() error
	Disconnect()
	IsConnected() bool
	GetStore() StoreInterface
}

// StoreInterface define interface mínima para store
type StoreInterface interface {
	GetID() interface{}
}
