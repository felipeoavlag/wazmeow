package port

import (
	"context"
	"wazmeow/internal/domain/entity"
)

// WhatsAppClient define a interface para o cliente WhatsApp
type WhatsAppClient interface {
	// Connect estabelece conexão com o WhatsApp
	Connect() error
	
	// Disconnect desconecta do WhatsApp
	Disconnect()
	
	// IsConnected verifica se está conectado
	IsConnected() bool
	
	// IsLoggedIn verifica se está logado
	IsLoggedIn() bool
	
	// Logout faz logout da sessão
	Logout(ctx context.Context) error
	
	// PairPhone emparelha um telefone
	PairPhone(ctx context.Context, phone string, showPushNotification bool, clientType string, clientName string) (string, error)
	
	// AddEventHandler adiciona um handler de eventos
	AddEventHandler(handler func(interface{})) string
	
	// RemoveEventHandler remove um handler de eventos
	RemoveEventHandler(handlerID string)
}

// WhatsAppClientFactory define a interface para criar clientes WhatsApp
type WhatsAppClientFactory interface {
	// CreateClient cria um novo cliente WhatsApp para uma sessão
	CreateClient(session *entity.Session) (WhatsAppClient, error)
}

// QRCodeGenerator define a interface para gerar QR codes
type QRCodeGenerator interface {
	// GenerateQR gera um QR code para autenticação
	GenerateQR(session *entity.Session) (string, error)
}

// SessionManager define a interface para gerenciar sessões no nível de infraestrutura
type SessionManager interface {
	// GetClient retorna o cliente WhatsApp para uma sessão
	GetClient(sessionID string) (WhatsAppClient, error)
	
	// SetClient define o cliente WhatsApp para uma sessão
	SetClient(sessionID string, client WhatsAppClient) error
	
	// RemoveClient remove o cliente WhatsApp de uma sessão
	RemoveClient(sessionID string) error
	
	// IsSessionConnected verifica se uma sessão está conectada
	IsSessionConnected(sessionID string) bool
}
