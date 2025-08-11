package session

import (
	"context"

	"wazmeow/internal/application/dto"
	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/pkg/logger"
)

// CreateSessionUseCase handles session creation
type CreateSessionUseCase struct {
	sessionRepo repositories.SessionRepository
}

// NewCreateSessionUseCase creates a new CreateSessionUseCase
func NewCreateSessionUseCase(sessionRepo repositories.SessionRepository) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute creates a new session
func (uc *CreateSessionUseCase) Execute(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	logger.Info().Str("name", req.Name).Msg("Creating new session")

	// Create new session entity
	session := entities.NewSession(req.Name)

	// Set optional fields
	if req.WebhookURL != "" || req.Events != "" {
		session.SetWebhook(req.WebhookURL, req.Events)
	}

	if req.ProxyConfig != nil {
		session.SetProxy(req.ProxyConfig)
	}

	// Validate session
	if err := session.Validate(); err != nil {
		logger.Error().Err(err).Msg("Session validation failed")
		return nil, err
	}

	// Save to repository
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		logger.Error().Err(err).Str("sessionId", session.ID).Msg("Failed to create session")
		return nil, err
	}

	logger.Info().Str("sessionId", session.ID).Msg("Session created successfully")

	// Convert to response DTO
	response := dto.ToSessionResponse(session)
	return &response, nil
}
