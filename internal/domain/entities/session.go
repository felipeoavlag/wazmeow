package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// SessionStatus represents the status of a WhatsApp session
type SessionStatus string

const (
	StatusDisconnected SessionStatus = "disconnected" // Estado inicial e após falha/timeout
	StatusConnecting   SessionStatus = "connecting"   // Durante processo de QR code
	StatusConnected    SessionStatus = "connected"    // Conectado e autenticado
)

// ProxyConfig holds proxy configuration for a session
type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	ProxyURL string `json:"proxyURL,omitempty"`
}

// Session representa uma sessão do WhatsApp
type Session struct {
	// ID único da sessão
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Nome da sessão
	Name string `json:"name" example:"minha-sessao"`
	// Status atual da sessão
	Status SessionStatus `json:"status" example:"connected"`
	// Número de telefone associado (opcional)
	Phone string `json:"phone,omitempty" example:"+5511999999999"`

	// JID do dispositivo WhatsApp (opcional)
	DeviceJID string `json:"deviceJID,omitempty" example:"5511999999999.0:1@s.whatsapp.net"`

	// QR Code para autenticação (base64 PNG)
	QRCode string `json:"qrCode,omitempty"`

	// Configuração de proxy (opcional)
	ProxyConfig *ProxyConfig `json:"proxyConfig,omitempty"`

	// URL do webhook para receber eventos (opcional)
	WebhookURL string `json:"webhookURL,omitempty" example:"https://example.com/webhook"`
	// Eventos subscritos separados por vírgula (opcional)
	Events string `json:"events,omitempty" example:"message,status"`

	// Data de criação da sessão
	CreatedAt time.Time `json:"createdAt" bun:"createdAt,nullzero,notnull,default:current_timestamp" example:"2023-08-19T10:30:00Z"`
	// Data da última atualização
	UpdatedAt time.Time `json:"updatedAt" bun:"updatedAt,nullzero,notnull,default:current_timestamp" example:"2023-08-19T10:30:00Z"`
}

// NewSession creates a new session with generated ID
func NewSession(name string) *Session {
	now := time.Now()
	return &Session{
		ID:        uuid.New().String(),
		Name:      name,
		Status:    StatusDisconnected,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate validates the session data
func (s *Session) Validate() error {
	if s.Name == "" {
		return errors.New("session name is required")
	}
	if s.ID == "" {
		return errors.New("session ID is required")
	}
	return nil
}

// UpdateStatus updates the session status and timestamp
func (s *Session) UpdateStatus(status SessionStatus) {
	s.Status = status
	s.UpdatedAt = time.Now()
}

// SetDeviceJID sets the device JID and updates timestamp
func (s *Session) SetDeviceJID(jid string) {
	s.DeviceJID = jid
	s.UpdatedAt = time.Now()
}

// SetPhone sets the phone number and updates timestamp
func (s *Session) SetPhone(phone string) {
	s.Phone = phone
	s.UpdatedAt = time.Now()
}

// SetWebhook sets the webhook URL and events
func (s *Session) SetWebhook(url, events string) {
	s.WebhookURL = url
	s.Events = events
	s.UpdatedAt = time.Now()
}

// SetProxy sets the proxy configuration
func (s *Session) SetProxy(config *ProxyConfig) {
	s.ProxyConfig = config
	s.UpdatedAt = time.Now()
}

// IsConnected returns true if the session is connected
func (s *Session) IsConnected() bool {
	return s.Status == StatusConnected
}

// IsConnecting returns true if the session is connecting
func (s *Session) IsConnecting() bool {
	return s.Status == StatusConnecting
}

// IsDisconnected returns true if the session is disconnected
func (s *Session) IsDisconnected() bool {
	return s.Status == StatusDisconnected
}
