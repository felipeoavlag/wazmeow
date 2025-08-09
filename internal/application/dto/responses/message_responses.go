package responses

// SendMessageResponse representa a resposta do envio de mensagem
type SendMessageResponse struct {
	// Detalhes sobre o envio da mensagem
	Details string `json:"details" example:"Mensagem enviada com sucesso"`
	// Timestamp Unix do envio
	Timestamp int64 `json:"timestamp" example:"1692454800"`
	// ID único da mensagem enviada
	ID string `json:"id" example:"3EB0C431C26A1916E07A"`
}
