package responses

import "wazmeow/internal/domain/entities"

// APIResponse representa uma resposta padrão da API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SessionInfo representa informações detalhadas de uma sessão
type SessionInfo struct {
	*entities.Session
	IsConnected bool `json:"is_connected"`
	IsLoggedIn  bool `json:"is_logged_in"`
}

// QRResponse representa a resposta do QR code
type QRResponse struct {
	QRCode string `json:"qr_code,omitempty"`
	Status string `json:"status"`
}

// PairCodeResponse representa a resposta do código de emparelhamento
type PairCodeResponse struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}
