package models

import (
	"time"

	"wazmeow/internal/domain/entity"
)

// SessionModel representa o modelo de persistência para sessões usando Bun ORM
// Este modelo contém as tags específicas do banco de dados e é usado apenas na camada de infraestrutura
type SessionModel struct {
	ID        string    `bun:"id,pk" json:"id"`
	Name      string    `bun:"name,unique,notnull" json:"name"`
	Status    string    `bun:"status,notnull,default:'disconnected'" json:"status"`
	Phone     string    `bun:"phone" json:"phone"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	// Campos de proxy (embedded da ProxyConfig)
	ProxyType     string `bun:"proxy_type" json:"proxy_type"`
	ProxyHost     string `bun:"proxy_host" json:"proxy_host"`
	ProxyPort     int    `bun:"proxy_port" json:"proxy_port"`
	ProxyUsername string `bun:"proxy_username" json:"proxy_username"`
	ProxyPassword string `bun:"proxy_password" json:"proxy_password"`
}

// TableName define o nome da tabela no banco de dados
func (SessionModel) TableName() string {
	return "sessions"
}

// ToDomain converte o modelo de persistência para a entidade de domínio
func (m *SessionModel) ToDomain() *entity.Session {
	session := &entity.Session{
		ID:        m.ID,
		Name:      m.Name,
		Status:    entity.SessionStatus(m.Status),
		Phone:     m.Phone,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	// Converter configuração de proxy se existir
	if m.ProxyType != "" {
		session.ProxyConfig = &entity.ProxyConfig{
			Type:     m.ProxyType,
			Host:     m.ProxyHost,
			Port:     m.ProxyPort,
			Username: m.ProxyUsername,
			Password: m.ProxyPassword,
		}
	}

	return session
}

// FromDomain converte a entidade de domínio para o modelo de persistência
func FromDomain(s *entity.Session) *SessionModel {
	model := &SessionModel{
		ID:        s.ID,
		Name:      s.Name,
		Status:    string(s.Status),
		Phone:     s.Phone,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}

	// Converter configuração de proxy se existir
	if s.ProxyConfig != nil {
		model.ProxyType = s.ProxyConfig.Type
		model.ProxyHost = s.ProxyConfig.Host
		model.ProxyPort = s.ProxyConfig.Port
		model.ProxyUsername = s.ProxyConfig.Username
		model.ProxyPassword = s.ProxyConfig.Password
	}

	return model
}

// ToDomainList converte uma lista de modelos para entidades de domínio
func ToDomainList(models []*SessionModel) []*entity.Session {
	sessions := make([]*entity.Session, len(models))
	for i, model := range models {
		sessions[i] = model.ToDomain()
	}
	return sessions
}
