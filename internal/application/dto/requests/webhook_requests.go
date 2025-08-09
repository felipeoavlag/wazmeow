package requests

// SetWebhookRequest representa a requisição para definir webhook
type SetWebhookRequest struct {
	WebhookURL string   `json:"webhookurl" validate:"required,url"`
	Events     []string `json:"events,omitempty"`
}

// UpdateWebhookRequest representa a requisição para atualizar webhook
type UpdateWebhookRequest struct {
	WebhookURL string   `json:"webhook" validate:"required,url"`
	Events     []string `json:"events,omitempty"`
	Active     bool     `json:"active"`
}
