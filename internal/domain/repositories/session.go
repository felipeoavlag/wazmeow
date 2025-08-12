package repositories

import (
	"context"

	"wazmeow/internal/domain/entities"
)

// SessionRepository defines the interface for session persistence
type SessionRepository interface {
	// Create creates a new session
	Create(ctx context.Context, session *entities.Session) error

	// GetByID retrieves a session by its ID
	GetByID(ctx context.Context, id string) (*entities.Session, error)

	// GetByDeviceJID retrieves a session by its device JID
	GetByDeviceJID(ctx context.Context, deviceJID string) (*entities.Session, error)

	// GetAll retrieves all sessions
	GetAll(ctx context.Context) ([]*entities.Session, error)

	// GetConnectedSessions retrieves all sessions with connected status
	GetConnectedSessions(ctx context.Context) ([]*entities.Session, error)

	// Update updates an existing session
	Update(ctx context.Context, session *entities.Session) error

	// Delete deletes a session by its ID
	Delete(ctx context.Context, id string) error

	// UpdateStatus updates only the status of a session
	UpdateStatus(ctx context.Context, id string, status entities.SessionStatus) error
}
