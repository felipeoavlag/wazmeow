package models

import (
	"time"

	"wazmeow/internal/domain/entity"

	"github.com/uptrace/bun"
)

// SessionModel representa o modelo de persistência para sessões usando Bun ORM
// Este modelo reflete EXATAMENTE a estrutura da entidade de domínio (internal/domain/entity/session.go)
// Usa convenção camelCase → snake_case automática via Bun ORM
type SessionModel struct {
	bun.BaseModel `bun:"table:sessions"`

	// Campos principais da sessão (seguem exatamente entity.Session)
	ID     string `bun:"id,pk" json:"id"`
	Name   string `bun:"name,unique,notnull" json:"name"`
	Status string `bun:"status,notnull,default:'disconnected'" json:"status"`
	Phone  string `bun:"phone" json:"phone"`

	// Campos WhatsApp (conversão automática camelCase → snake_case)
	// DeviceJID → device_jid no PostgreSQL
	DeviceJID  string `bun:"deviceJID,default:''" json:"deviceJID"`
	// WebhookURL → webhook_url no PostgreSQL
	WebhookURL string `bun:"webhookURL,default:''" json:"webhookURL"`
	Events     string `bun:"events,default:''" json:"events"`

	// Campos de proxy (desnormalizados para performance - *nullable para opcionais*)
	// ProxyType → proxy_type no PostgreSQL
	ProxyType     *string `bun:"proxyType" json:"proxyType,omitempty"`
	ProxyHost     *string `bun:"proxyHost" json:"proxyHost,omitempty"`
	ProxyPort     *int    `bun:"proxyPort" json:"proxyPort,omitempty"`
	ProxyUsername *string `bun:"proxyUsername" json:"proxyUsername,omitempty"`
	ProxyPassword *string `bun:"proxyPassword" json:"proxyPassword,omitempty"`

	// Campos de auditoria (conversão automática camelCase → snake_case)
	// CreatedAt → created_at, UpdatedAt → updated_at
	CreatedAt time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:"updatedAt,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

// TableName define o nome da tabela no banco de dados
func (*SessionModel) TableName() string {
	return "sessions"
}

// ToDomain converte o modelo de persistência para a entidade de domínio
func (m *SessionModel) ToDomain() *entity.Session {
	session := &entity.Session{
		ID:         m.ID,
		Name:       m.Name,
		Status:     entity.SessionStatus(m.Status),
		Phone:      m.Phone,
		DeviceJID:  m.DeviceJID,
		WebhookURL: m.WebhookURL,
		Events:     m.Events,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}

	// Converter configuração de proxy se existir (verificar se algum campo não é nil)
	if m.ProxyType != nil && *m.ProxyType != "" {
		session.ProxyConfig = &entity.ProxyConfig{
			Type: *m.ProxyType,
			Host: safeStringValue(m.ProxyHost),
			Port: safeIntValue(m.ProxyPort),
			Username: safeStringValue(m.ProxyUsername),
			Password: safeStringValue(m.ProxyPassword),
		}
	}

	return session
}

// FromDomain converte a entidade de domínio para o modelo de persistência
func FromDomain(s *entity.Session) *SessionModel {
	model := &SessionModel{
		ID:         s.ID,
		Name:       s.Name,
		Status:     string(s.Status),
		Phone:      s.Phone,
		DeviceJID:  s.DeviceJID,
		WebhookURL: s.WebhookURL,
		Events:     s.Events,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}

	// Converter configuração de proxy se existir
	if s.ProxyConfig != nil {
		model.ProxyType = stringPtr(s.ProxyConfig.Type)
		model.ProxyHost = stringPtr(s.ProxyConfig.Host)
		model.ProxyPort = intPtr(s.ProxyConfig.Port)
		model.ProxyUsername = stringPtr(s.ProxyConfig.Username)
		model.ProxyPassword = stringPtr(s.ProxyConfig.Password)
	}

	return model
}

// Funções auxiliares para conversão de ponteiros
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func safeStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeIntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// ToDomainList converte uma lista de modelos para entidades de domínio
func ToDomainList(models []*SessionModel) []*entity.Session {
	sessions := make([]*entity.Session, len(models))
	for i, model := range models {
		sessions[i] = model.ToDomain()
	}
	return sessions
}
