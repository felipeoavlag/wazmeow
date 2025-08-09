package entity

import (
	"time"
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
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Status      SessionStatus `json:"status"`
	Phone       string        `json:"phone,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	ProxyConfig *ProxyConfig  `json:"proxy_config,omitempty"`
}

// ProxyConfig representa a configuração de proxy
type ProxyConfig struct {
	Type     string `json:"type"` // http, socks5
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
