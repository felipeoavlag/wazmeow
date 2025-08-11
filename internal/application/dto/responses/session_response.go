package responses

import "wazmeow/internal/domain/entity"

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
