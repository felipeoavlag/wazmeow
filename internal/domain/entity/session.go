package entity

import (
	"time"
)

// SessionStatus representa o status de uma sessão
type SessionStatus string

const (
	StatusDisconnected SessionStatus = "disconnected" // Estado inicial e após falha/timeout
	StatusConnecting   SessionStatus = "connecting"   // Durante processo de QR code
	StatusConnected    SessionStatus = "connected"    // Conectado e autenticado
)

// Session representa uma sessão do WhatsApp
type Session struct {
	// Campos principais da sessão
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Status SessionStatus `json:"status"`
	Phone  string        `json:"phone,omitempty"`

	// Campos WhatsApp (conexão e autenticação)
	DeviceJID string `json:"device_jid,omitempty"`

	// Configuração de proxy
	ProxyConfig *ProxyConfig `json:"proxy_config,omitempty"`

	// Campos de auditoria (sempre no final)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProxyConfig representa a configuração de proxy
type ProxyConfig struct {
	Type     string `json:"type"` // http, socks5
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
