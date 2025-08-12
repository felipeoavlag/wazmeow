package session

import (
	"context"
	"fmt"

	"wazmeow/internal/config"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/domain/services"
)

// ManagerInterface define interface para manager
type ManagerInterface interface {
	Create(ctx context.Context, sessionID string) error
	Remove(sessionID string) error
	Get(sessionID string) SessionWrapperInterface
}

// SessionWrapperInterface define interface para wrapper de sessão
type SessionWrapperInterface interface {
	IsConnected() bool
	IsLoggedIn() bool
	GetJIDString() string
}

// Service implementa interface pública para gerenciamento de sessões
type Service struct {
	manager     ManagerInterface
	sessionRepo repositories.SessionRepository
	config      *config.WhatsAppConfig
}

// NewService cria um novo serviço de sessões
func NewService(
	manager ManagerInterface,
	sessionRepo repositories.SessionRepository,
	config *config.WhatsAppConfig,
) *Service {
	return &Service{
		manager:     manager,
		sessionRepo: sessionRepo,
		config:      config,
	}
}

// Start inicia uma sessão
func (s *Service) Start(ctx context.Context, sessionID string) error {
	return s.manager.Create(ctx, sessionID)
}

// Stop para uma sessão
func (s *Service) Stop(ctx context.Context, sessionID string) error {
	return s.manager.Remove(sessionID)
}

// GetQR retorna QR code de uma sessão
func (s *Service) GetQR(ctx context.Context, sessionID string) (string, error) {
	wrapper := s.manager.Get(sessionID)
	if wrapper == nil {
		return "", fmt.Errorf("session %s not found", sessionID)
	}

	// TODO: Implementar GetQR usando QR processor
	return "", fmt.Errorf("QR code not available")
}

// GetInfo retorna informações da sessão
func (s *Service) GetInfo(sessionID string) (*services.SessionInfo, error) {
	wrapper := s.manager.Get(sessionID)
	if wrapper == nil {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	return &services.SessionInfo{
		SessionID: sessionID,
		Connected: wrapper.IsConnected(),
		LoggedIn:  wrapper.IsLoggedIn(),
		DeviceJID: wrapper.GetJIDString(),
	}, nil
}
