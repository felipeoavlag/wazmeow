package services

import (
	"context"
	"time"

	"wazmeow/internal/domain/entities"
)

// WhatsAppService defines the interface for WhatsApp operations
type WhatsAppService interface {
	// StartSession starts a WhatsApp session
	StartSession(ctx context.Context, sessionID string) error

	// StopSession stops a WhatsApp session
	StopSession(ctx context.Context, sessionID string) error

	// GetQRCode gets the QR code for session authentication
	GetQRCode(ctx context.Context, sessionID string) (string, error)

	// PairPhone pairs a phone number with the session
	PairPhone(ctx context.Context, sessionID, phone string) (string, error)

	// Logout logs out from WhatsApp
	Logout(ctx context.Context, sessionID string) error

	// IsConnected checks if a session is connected
	IsConnected(sessionID string) bool

	// IsLoggedIn checks if a session is logged in
	IsLoggedIn(sessionID string) bool

	// SetProxy sets proxy configuration for a session
	SetProxy(sessionID string, config *entities.ProxyConfig) error

	// GetSessionInfo gets detailed session information
	GetSessionInfo(sessionID string) (*SessionInfo, error)

	// SetSubscriptions sets event subscriptions for a session
	SetSubscriptions(sessionID string, subscriptions []string) error

	// GetSubscriptions gets event subscriptions for a session
	GetSubscriptions(sessionID string) ([]string, error)

	// AddSubscription adds a single event subscription to a session
	AddSubscription(sessionID, eventType string) error

	// RemoveSubscription removes a single event subscription from a session
	RemoveSubscription(sessionID, eventType string) error

	// GetSupportedEventTypes returns list of supported event types
	GetSupportedEventTypes() []string

	// GetAllSessionsInfo returns information about all active sessions
	GetAllSessionsInfo() []map[string]interface{}
}

// SessionInfo holds detailed information about a WhatsApp session
type SessionInfo struct {
	SessionID     string   `json:"sessionId"`
	Connected     bool     `json:"connected"`
	LoggedIn      bool     `json:"loggedIn"`
	Phone         string   `json:"phone,omitempty"`
	DeviceJID     string   `json:"deviceJID,omitempty"`
	QRCode        string   `json:"qrCode,omitempty"`
	Subscriptions []string `json:"subscriptions,omitempty"`
	Webhook       string   `json:"webhook,omitempty"`
}

// QRCodeData represents QR code information
type QRCodeData struct {
	Code      string    `json:"code"`
	Base64PNG string    `json:"base64Png"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// PairData represents phone pairing information
type PairData struct {
	LinkingCode string    `json:"linkingCode"`
	Phone       string    `json:"phone"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// WhatsAppEvent represents events from WhatsApp client
type WhatsAppEvent struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"sessionId"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}
