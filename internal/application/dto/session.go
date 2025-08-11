package dto

import (
	"time"

	"wazmeow/internal/domain/entities"
)

// CreateSessionRequest represents the request to create a new session
type CreateSessionRequest struct {
	Name        string                  `json:"name" validate:"required"`
	WebhookURL  string                  `json:"webhookURL,omitempty"`
	Events      string                  `json:"events,omitempty"`
	ProxyConfig *entities.ProxyConfig   `json:"proxyConfig,omitempty"`
}

// SessionResponse represents a session in API responses
type SessionResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Status      entities.SessionStatus  `json:"status"`
	Phone       string                  `json:"phone,omitempty"`
	DeviceJID   string                  `json:"deviceJID,omitempty"`
	ProxyConfig *entities.ProxyConfig   `json:"proxyConfig,omitempty"`
	WebhookURL  string                  `json:"webhookURL,omitempty"`
	Events      string                  `json:"events,omitempty"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}

// SessionListResponse represents the response for listing sessions
type SessionListResponse struct {
	Sessions []SessionResponse `json:"sessions"`
	Total    int               `json:"total"`
}

// SessionInfoResponse represents detailed session information
type SessionInfoResponse struct {
	SessionResponse
	Connected bool   `json:"connected"`
	LoggedIn  bool   `json:"loggedIn"`
	QRCode    string `json:"qrCode,omitempty"`
}

// ConnectSessionRequest represents the request to connect a session
type ConnectSessionRequest struct {
	Events    string `json:"events,omitempty"`
	Immediate bool   `json:"immediate,omitempty"`
}

// PairPhoneRequest represents the request to pair a phone
type PairPhoneRequest struct {
	Phone string `json:"phone" validate:"required"`
}

// PairPhoneResponse represents the response for phone pairing
type PairPhoneResponse struct {
	LinkingCode string `json:"linkingCode"`
}

// SetProxyRequest represents the request to set proxy configuration
type SetProxyRequest struct {
	Enabled  bool   `json:"enabled"`
	ProxyURL string `json:"proxyURL,omitempty"`
}

// QRCodeResponse represents the QR code response
type QRCodeResponse struct {
	QRCode string `json:"qrCode"`
	Image  string `json:"image,omitempty"` // Base64 encoded PNG image
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ToSessionResponse converts a domain session to a response DTO
func ToSessionResponse(session *entities.Session) SessionResponse {
	return SessionResponse{
		ID:          session.ID,
		Name:        session.Name,
		Status:      session.Status,
		Phone:       session.Phone,
		DeviceJID:   session.DeviceJID,
		ProxyConfig: session.ProxyConfig,
		WebhookURL:  session.WebhookURL,
		Events:      session.Events,
		CreatedAt:   session.CreatedAt,
		UpdatedAt:   session.UpdatedAt,
	}
}

// ToSessionListResponse converts a slice of domain sessions to a list response DTO
func ToSessionListResponse(sessions []*entities.Session) SessionListResponse {
	responses := make([]SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = ToSessionResponse(session)
	}
	return SessionListResponse{
		Sessions: responses,
		Total:    len(sessions),
	}
}
