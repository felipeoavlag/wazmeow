// Package domain contains the core business entities and domain logic
// Following Clean Architecture principles, this layer is independent of external concerns
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Session representa uma sessão do WhatsApp
// Esta entidade encapsula todas as informações necessárias para gerenciar
// uma conexão ativa com o WhatsApp Web
type Session struct {
	// Identificação única da sessão
	ID string `json:"id" bun:"id,pk" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Nome amigável da sessão para identificação
	Name string `json:"name" bun:"name,notnull" example:"Minha Sessão WhatsApp"`

	// Número de telefone associado à sessão (formato internacional)
	Phone string `json:"phone" bun:"phone" example:"5511999999999"`

	// Status atual da sessão
	Status SessionStatus `json:"status" bun:"status,notnull" example:"connected"`

	// Indica se a sessão está ativa e pode receber/enviar mensagens
	Active bool `json:"active" bun:"active,notnull" example:"true"`

	// URL do webhook para receber eventos desta sessão
	WebhookURL string `json:"webhook_url" bun:"webhook_url" example:"https://api.exemplo.com/webhook"`

	// Lista de eventos que devem ser enviados para o webhook
	Events []string `json:"events" bun:"events,array" example:"message,message_ack,qr"`

	// Timestamps de controle
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`

	// Timestamp da última atividade registrada
	LastActivity *time.Time `json:"last_activity,omitempty" bun:"last_activity"`
}

// SessionStatus representa os possíveis estados de uma sessão
type SessionStatus string

const (
	// SessionStatusDisconnected indica que a sessão não está conectada
	SessionStatusDisconnected SessionStatus = "disconnected"

	// SessionStatusConnecting indica que a sessão está tentando conectar
	SessionStatusConnecting SessionStatus = "connecting"

	// SessionStatusConnected indica que a sessão está conectada e ativa
	SessionStatusConnected SessionStatus = "connected"

	// SessionStatusError indica que houve erro na conexão
	SessionStatusError SessionStatus = "error"

	// SessionStatusQR indica que está aguardando leitura do QR Code
	SessionStatusQR SessionStatus = "qr"
)

// IsValid verifica se o status é válido
func (s SessionStatus) IsValid() bool {
	switch s {
	case SessionStatusDisconnected, SessionStatusConnecting, SessionStatusConnected, SessionStatusError, SessionStatusQR:
		return true
	default:
		return false
	}
}

// String retorna a representação em string do status
func (s SessionStatus) String() string {
	return string(s)
}

// NewSession cria uma nova sessão com valores padrão
func NewSession(name string) *Session {
	return &Session{
		ID:        uuid.New().String(),
		Name:      name,
		Status:    SessionStatusDisconnected,
		Active:    false,
		Events:    []string{"message", "message_ack", "qr", "ready", "authenticated", "auth_failure", "disconnected"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdateStatus atualiza o status da sessão e o timestamp
func (s *Session) UpdateStatus(status SessionStatus) {
	s.Status = status
	s.UpdatedAt = time.Now()
	now := time.Now()
	s.LastActivity = &now
}

// SetActive define se a sessão está ativa
func (s *Session) SetActive(active bool) {
	s.Active = active
	s.UpdatedAt = time.Now()
}

// SetWebhook configura o webhook da sessão
func (s *Session) SetWebhook(url string, events []string) {
	s.WebhookURL = url
	s.Events = events
	s.UpdatedAt = time.Now()
}

// IsConnected verifica se a sessão está conectada
func (s *Session) IsConnected() bool {
	return s.Status == SessionStatusConnected && s.Active
}

// CanSendMessages verifica se a sessão pode enviar mensagens
func (s *Session) CanSendMessages() bool {
	return s.IsConnected()
}

// HasWebhook verifica se a sessão tem webhook configurado
func (s *Session) HasWebhook() bool {
	return s.WebhookURL != ""
}

// ShouldSendEvent verifica se um evento deve ser enviado para o webhook
func (s *Session) ShouldSendEvent(eventType string) bool {
	if !s.HasWebhook() {
		return false
	}

	for _, event := range s.Events {
		if event == eventType {
			return true
		}
	}
	return false
}

// TableName retorna o nome da tabela no banco de dados
func (s *Session) TableName() string {
	return "sessions"
}
