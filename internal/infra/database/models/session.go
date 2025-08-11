package models

import (
	"time"

	"github.com/uptrace/bun"

	"wazmeow/internal/domain/entities"
)

// SessionModel represents the session table in the database
type SessionModel struct {
	bun.BaseModel `bun:"table:Sessions,alias:s"`

	ID           string    `bun:"id,pk" json:"id"`
	Name         string    `bun:"name,notnull" json:"name"`
	Status       string    `bun:"status,notnull,default:'disconnected'" json:"status"`
	Phone        *string   `bun:"phone" json:"phone,omitempty"`
	DeviceJID    *string   `bun:"deviceJID" json:"deviceJID,omitempty"`
	ProxyEnabled bool      `bun:"proxyEnabled,default:false" json:"proxyEnabled"`
	ProxyURL     *string   `bun:"proxyURL" json:"proxyURL,omitempty"`
	WebhookURL   *string   `bun:"webhookURL" json:"webhookURL,omitempty"`
	Events       *string   `bun:"events" json:"events,omitempty"`
	CreatedAt    time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt    time.Time `bun:"updatedAt,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

// ToEntity converts the database model to a domain entity
func (m *SessionModel) ToEntity() *entities.Session {
	session := &entities.Session{
		ID:        m.ID,
		Name:      m.Name,
		Status:    entities.SessionStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if m.Phone != nil {
		session.Phone = *m.Phone
	}

	if m.DeviceJID != nil {
		session.DeviceJID = *m.DeviceJID
	}

	if m.ProxyEnabled && m.ProxyURL != nil {
		session.ProxyConfig = &entities.ProxyConfig{
			Enabled:  true,
			ProxyURL: *m.ProxyURL,
		}
	}

	if m.WebhookURL != nil {
		session.WebhookURL = *m.WebhookURL
	}

	if m.Events != nil {
		session.Events = *m.Events
	}

	return session
}

// FromEntity converts a domain entity to a database model
func (m *SessionModel) FromEntity(session *entities.Session) {
	m.ID = session.ID
	m.Name = session.Name
	m.Status = string(session.Status)
	m.CreatedAt = session.CreatedAt
	m.UpdatedAt = session.UpdatedAt

	if session.Phone != "" {
		m.Phone = &session.Phone
	}

	if session.DeviceJID != "" {
		m.DeviceJID = &session.DeviceJID
	}

	if session.ProxyConfig != nil {
		m.ProxyEnabled = session.ProxyConfig.Enabled
		if session.ProxyConfig.ProxyURL != "" {
			m.ProxyURL = &session.ProxyConfig.ProxyURL
		}
	}

	if session.WebhookURL != "" {
		m.WebhookURL = &session.WebhookURL
	}

	if session.Events != "" {
		m.Events = &session.Events
	}
}

// NewSessionModel creates a new SessionModel from a domain entity
func NewSessionModel(session *entities.Session) *SessionModel {
	model := &SessionModel{}
	model.FromEntity(session)
	return model
}
