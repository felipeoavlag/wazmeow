package responses

// WebhookResponse representa a resposta de operações de webhook
type WebhookResponse struct {
	Webhook   string   `json:"webhook"`
	Events    []string `json:"events,omitempty"`
	Active    bool     `json:"active,omitempty"`
	Subscribe []string `json:"subscribe,omitempty"`
}

// WebhookDeleteResponse representa a resposta de exclusão de webhook
type WebhookDeleteResponse struct {
	Details string `json:"details"`
}
