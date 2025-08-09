package responses

// SendMessageResponse representa a resposta do envio de mensagem
type SendMessageResponse struct {
	Details   string `json:"details"`
	Timestamp int64  `json:"timestamp"`
	ID        string `json:"id"`
}
