package entities

import (
	"time"

	"go.mau.fi/whatsmeow"
)

// SessionStatus representa o status de uma sessão
type SessionStatus string

const (
	StatusDisconnected SessionStatus = "disconnected"
	StatusConnecting   SessionStatus = "connecting"
	StatusConnected    SessionStatus = "connected"
	StatusLoggedOut    SessionStatus = "logged_out"
)

// Session representa uma sessão do WhatsApp
type Session struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Status      SessionStatus     `json:"status"`
	Phone       string            `json:"phone,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Client      *whatsmeow.Client `json:"-"`
	ProxyConfig *ProxyConfig      `json:"proxy_config,omitempty"`
}

// ProxyConfig representa a configuração de proxy
type ProxyConfig struct {
	Type     string `bun:"proxy_type" json:"type"` // http, socks5
	Host     string `bun:"proxy_host" json:"host"`
	Port     int    `bun:"proxy_port" json:"port"`
	Username string `bun:"proxy_username" json:"username,omitempty"`
	Password string `bun:"proxy_password" json:"password,omitempty"`
}
