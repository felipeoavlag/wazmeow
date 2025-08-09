package responses

import "wazmeow/internal/domain/entity"

// SessionResponse representa uma resposta padrão para operações de sessão
type SessionResponse struct {
	// Indica se a operação foi bem-sucedida
	Success bool `json:"success" example:"true"`
	// Mensagem descritiva da operação (opcional)
	Message string `json:"message,omitempty" example:"Operação realizada com sucesso"`
	// Mensagem de erro caso a operação falhe (opcional)
	Error string `json:"error,omitempty" example:"Sessão não encontrada"`
	// Dados retornados pela operação (opcional)
	Data interface{} `json:"data,omitempty"`
}

// APIResponse representa uma resposta padrão da API
type APIResponse struct {
	// Indica se a operação foi bem-sucedida
	Success bool `json:"success" example:"true"`
	// Mensagem descritiva da operação (opcional)
	Message string `json:"message,omitempty" example:"Operação realizada com sucesso"`
	// Mensagem de erro caso a operação falhe (opcional)
	Error string `json:"error,omitempty" example:"Dados inválidos"`
	// Dados retornados pela operação (opcional)
	Data interface{} `json:"data,omitempty"`
}

// SessionInfo representa informações detalhadas de uma sessão
type SessionInfo struct {
	*entity.Session
	// Indica se a sessão está conectada ao WhatsApp
	IsConnected bool `json:"is_connected" example:"true"`
	// Indica se a sessão está autenticada no WhatsApp
	IsLoggedIn bool `json:"is_logged_in" example:"true"`
}

// QRResponse representa a resposta do QR code
type QRResponse struct {
	// Código QR para autenticação (opcional)
	QRCode string `json:"qr_code,omitempty" example:"2@BQcAEAYQAg==,f/9u+vz6zJTzOD0VGOEkjrU=,wU/DdpXJ0tPalzxUr6SQBlMAAAAAElFTkSuQmCC"`
	// Status do QR code
	Status string `json:"status" example:"qr_generated"`
}

// PairCodeResponse representa a resposta do código de emparelhamento
type PairCodeResponse struct {
	// Código de emparelhamento gerado
	Code string `json:"code" example:"ABCD-EFGH"`
	// Status do emparelhamento
	Status string `json:"status" example:"code_generated"`
}
