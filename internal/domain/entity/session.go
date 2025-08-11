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

	// Configuração de proxy (opcional)
	ProxyConfig *ProxyConfig `json:"proxyConfig,omitempty"`

	// URL do webhook para receber eventos (opcional)
	WebhookURL string `json:"webhookURL,omitempty" example:"https://example.com/webhook"`
	// Eventos subscritos separados por vírgula (opcional)
	Events string `json:"events,omitempty" example:"message,status"`

	// Data de criação da sessão
	CreatedAt time.Time `json:"createdAt" example:"2023-08-19T10:30:00Z"`
	// Data da última atualização
	UpdatedAt time.Time `json:"updatedAt" example:"2023-08-19T10:30:00Z"`
}

// ProxyConfig representa a configuração de proxy
type ProxyConfig struct {
	// Tipo do proxy (http, socks5)
	Type string `json:"type" example:"http"`
	// Host do servidor proxy
	Host string `json:"host" example:"proxy.example.com"`
	// Porta do servidor proxy
	Port int `json:"port" example:"8080"`
	// Nome de usuário para autenticação (opcional)
	Username string `json:"username,omitempty" example:"usuario"`
	// Senha para autenticação (opcional)
	Password string `json:"password,omitempty" example:"senha"`
}
