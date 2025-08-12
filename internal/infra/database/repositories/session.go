package repositories

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/infra/database/models"
	"wazmeow/pkg/logger"
)

// sessionRepository implements the SessionRepository interface using Bun ORM
type sessionRepository struct {
	db *bun.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *bun.DB) repositories.SessionRepository {
	return &sessionRepository{db: db}
}

// Create creates a new session using Bun query builder
func (r *sessionRepository) Create(ctx context.Context, session *entities.Session) error {
	model := models.NewSessionModel(session)

	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	if err != nil {
		logger.Error().Err(err).Str("sessionId", session.ID).Msg("Failed to create session")
		return err
	}

	logger.Debug().Str("sessionId", session.ID).Msg("Session created successfully")
	return nil
}

// GetByID retrieves a session by its ID using Bun query builder
func (r *sessionRepository) GetByID(ctx context.Context, id string) (*entities.Session, error) {
	model := &models.SessionModel{}

	err := r.db.NewSelect().
		Model(model).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		logger.Error().Err(err).Str("sessionId", id).Msg("Failed to get session by ID")
		return nil, err
	}

	return model.ToEntity(), nil
}

// GetByDeviceJID retrieves a session by its device JID using Bun query builder
func (r *sessionRepository) GetByDeviceJID(ctx context.Context, deviceJID string) (*entities.Session, error) {
	model := &models.SessionModel{}

	err := r.db.NewSelect().
		Model(model).
		Where("deviceJID = ?", deviceJID).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		logger.Error().Err(err).Str("deviceJID", deviceJID).Msg("Failed to get session by deviceJID")
		return nil, err
	}

	return model.ToEntity(), nil
}

// GetAll retrieves all sessions using Bun query builder
func (r *sessionRepository) GetAll(ctx context.Context) ([]*entities.Session, error) {
	var models []*models.SessionModel

	err := r.db.NewSelect().
		Model(&models).
		Order("createdAt DESC").
		Scan(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get all sessions")
		return nil, err
	}

	sessions := make([]*entities.Session, len(models))
	for i, model := range models {
		sessions[i] = model.ToEntity()
	}

	return sessions, nil
}

// GetConnectedSessions retrieves all sessions with connected status using Bun query builder
func (r *sessionRepository) GetConnectedSessions(ctx context.Context) ([]*entities.Session, error) {
	var models []*models.SessionModel

	err := r.db.NewSelect().
		Model(&models).
		Where("status = ?", "connected").
		Scan(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get connected sessions")
		return nil, err
	}

	sessions := make([]*entities.Session, len(models))
	for i, model := range models {
		sessions[i] = model.ToEntity()
	}

	logger.Debug().Int("count", len(sessions)).Msg("Retrieved connected sessions")
	return sessions, nil
}

// Update updates an existing session using Bun query builder
func (r *sessionRepository) Update(ctx context.Context, session *entities.Session) error {
	model := models.NewSessionModel(session)
	model.UpdatedAt = time.Now()

	_, err := r.db.NewUpdate().
		Model(model).
		Where("id = ?", session.ID).
		Exec(ctx)

	if err != nil {
		logger.Error().Err(err).Str("sessionId", session.ID).Msg("Failed to update session")
		return err
	}

	logger.Debug().Str("sessionId", session.ID).Msg("Session updated successfully")
	return nil
}

// Delete deletes a session by its ID using Bun query builder
func (r *sessionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.SessionModel)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		logger.Error().Err(err).Str("sessionId", id).Msg("Failed to delete session")
		return err
	}

	logger.Debug().Str("sessionId", id).Msg("Session deleted successfully")
	return nil
}

// UpdateStatus updates only the status of a session using Bun query builder
func (r *sessionRepository) UpdateStatus(ctx context.Context, id string, status entities.SessionStatus) error {
	_, err := r.db.NewUpdate().
		Model((*models.SessionModel)(nil)).
		Set("status = ?", string(status)).
		Set(`"updatedAt" = ?`, time.Now()).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		logger.Error().Err(err).Str("sessionId", id).Str("status", string(status)).Msg("Failed to update session status")
		return err
	}

	logger.Debug().Str("sessionId", id).Str("status", string(status)).Msg("Session status updated successfully")
	return nil
}
