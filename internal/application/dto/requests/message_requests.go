package requests

// SendTextMessageRequest representa a requisição para envio de mensagem de texto
type SendTextMessageRequest struct {
	Phone       string      `json:"phone" validate:"required"`
	Body        string      `json:"body" validate:"required"`
	ID          string      `json:"id,omitempty"`
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// SendMediaMessageRequest representa a requisição para envio de mídia
type SendMediaMessageRequest struct {
	Phone       string      `json:"phone" validate:"required"`
	MediaData   string      `json:"media_data" validate:"required"` // Base64 encoded data
	Caption     string      `json:"caption,omitempty"`
	MimeType    string      `json:"mime_type,omitempty"`
	ID          string      `json:"id,omitempty"`
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// ContextInfo representa informações de contexto para mensagens (reply, mentions)
type ContextInfo struct {
	StanzaID     *string   `json:"stanza_id,omitempty"`
	Participant  *string   `json:"participant,omitempty"`
	MentionedJID []string  `json:"mentioned_jid,omitempty"`
}
