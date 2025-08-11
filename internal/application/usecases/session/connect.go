package session

import (
	"context"
	"errors"

	"wazmeow/internal/application/dto"
	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/domain/services"
	"wazmeow/pkg/logger"
)

// ConnectSessionUseCase handles session connection
type ConnectSessionUseCase struct {
	sessionRepo   repositories.SessionRepository
	whatsappSvc   services.WhatsAppService
}

// NewConnectSessionUseCase creates a new ConnectSessionUseCase
func NewConnectSessionUseCase(
	sessionRepo repositories.SessionRepository,
	whatsappSvc services.WhatsAppService,
) *ConnectSessionUseCase {
	return &ConnectSessionUseCase{
		sessionRepo: sessionRepo,
		whatsappSvc: whatsappSvc,
	}
}

// Execute connects a session to WhatsApp
func (uc *ConnectSessionUseCase) Execute(ctx context.Context, sessionID string, req dto.ConnectSessionRequest) (*dto.APIResponse, error) {
	logger.Info().Str("sessionId", sessionID).Msg("Connecting session")

	// Get session from repository
	session, err := uc.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Msg("Session not found")
		return nil, errors.New("session not found")
	}

	// Update events if provided
	if req.Events != "" {
		session.SetWebhook(session.WebhookURL, req.Events)
		if err := uc.sessionRepo.Update(ctx, session); err != nil {
			logger.Error().Err(err).Str("sessionId", sessionID).Msg("Failed to update session events")
			return nil, err
		}
	}

	// Update session status to connecting
	session.UpdateStatus(entities.StatusConnecting)
	if err := uc.sessionRepo.UpdateStatus(ctx, sessionID, entities.StatusConnecting); err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Msg("Failed to update session status")
		return nil, err
	}

	// Start WhatsApp session
	if err := uc.whatsappSvc.StartSession(ctx, sessionID); err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Msg("Failed to start WhatsApp session")
		
		// Revert status to disconnected
		session.UpdateStatus(entities.StatusDisconnected)
		uc.sessionRepo.UpdateStatus(ctx, sessionID, entities.StatusDisconnected)
		
		return nil, err
	}

	logger.Info().Str("sessionId", sessionID).Msg("Session connection initiated")

	return &dto.APIResponse{
		Success: true,
		Message: "Session connection initiated",
	}, nil
}
