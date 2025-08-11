package session

import (
	"context"

	"wazmeow/internal/application/dto"
	"wazmeow/internal/domain/repositories"
	"wazmeow/pkg/logger"
)

// ListSessionsUseCase handles listing all sessions
type ListSessionsUseCase struct {
	sessionRepo repositories.SessionRepository
}

// NewListSessionsUseCase creates a new ListSessionsUseCase
func NewListSessionsUseCase(sessionRepo repositories.SessionRepository) *ListSessionsUseCase {
	return &ListSessionsUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute lists all sessions
func (uc *ListSessionsUseCase) Execute(ctx context.Context) (*dto.SessionListResponse, error) {
	logger.Debug().Msg("Listing all sessions")

	// Get all sessions from repository
	sessions, err := uc.sessionRepo.GetAll(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list sessions")
		return nil, err
	}

	logger.Info().Int("count", len(sessions)).Msg("Sessions retrieved successfully")

	// Convert to response DTO
	response := dto.ToSessionListResponse(sessions)
	return &response, nil
}
