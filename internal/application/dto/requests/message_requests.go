package requests

// SendTextMessageRequest representa a requisição para envio de mensagem de texto
type SendTextMessageRequest struct {
	// Número de telefone do destinatário (formato: 5511999999999)
	Phone string `json:"phone" validate:"required" example:"5511999999999"`
	// Conteúdo da mensagem de texto
	Body string `json:"body" validate:"required" example:"Olá! Como você está?"`
	// ID personalizado da mensagem (opcional)
	ID string `json:"id,omitempty" example:"msg-123"`
	// Informações de contexto para reply ou menções (opcional)
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// SendMediaMessageRequest representa a requisição para envio de mídia
type SendMediaMessageRequest struct {
	// Número de telefone do destinatário (formato: 5511999999999)
	Phone string `json:"phone" validate:"required" example:"5511999999999"`
	// Dados da mídia codificados em Base64
	MediaData string `json:"media_data" validate:"required" example:"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="`
	// Legenda da mídia (opcional)
	Caption string `json:"caption,omitempty" example:"Minha imagem"`
	// Tipo MIME da mídia (opcional, detectado automaticamente se não fornecido)
	MimeType string `json:"mime_type,omitempty" example:"image/png"`
	// ID personalizado da mensagem (opcional)
	ID string `json:"id,omitempty" example:"media-123"`
	// Informações de contexto para reply ou menções (opcional)
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// SendImageMessageRequest representa a requisição para envio de imagem
type SendImageMessageRequest struct {
	Phone       string      `json:"phone" validate:"required"`
	Image       string      `json:"image" validate:"required"` // Base64 encoded data
	Caption     string      `json:"caption,omitempty"`
	MimeType    string      `json:"mime_type,omitempty"`
	ID          string      `json:"id,omitempty"`
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// SendAudioMessageRequest representa a requisição para envio de áudio
type SendAudioMessageRequest struct {
	Phone       string      `json:"phone" validate:"required"`
	Audio       string      `json:"audio" validate:"required"` // Base64 encoded data
	Caption     string      `json:"caption,omitempty"`
	ID          string      `json:"id,omitempty"`
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// SendDocumentMessageRequest representa a requisição para envio de documento
type SendDocumentMessageRequest struct {
	Phone       string      `json:"phone" validate:"required"`
	Document    string      `json:"document" validate:"required"` // Base64 encoded data
	FileName    string      `json:"filename" validate:"required"`
	Caption     string      `json:"caption,omitempty"`
	MimeType    string      `json:"mime_type,omitempty"`
	ID          string      `json:"id,omitempty"`
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// SendVideoMessageRequest representa a requisição para envio de vídeo
type SendVideoMessageRequest struct {
	Phone         string      `json:"phone" validate:"required"`
	Video         string      `json:"video" validate:"required"` // Base64 encoded data
	Caption       string      `json:"caption,omitempty"`
	MimeType      string      `json:"mime_type,omitempty"`
	JPEGThumbnail []byte      `json:"jpeg_thumbnail,omitempty"`
	ID            string      `json:"id,omitempty"`
	ContextInfo   ContextInfo `json:"context_info,omitempty"`
}

// SendStickerMessageRequest representa a requisição para envio de sticker
type SendStickerMessageRequest struct {
	Phone        string      `json:"phone" validate:"required"`
	Sticker      string      `json:"sticker" validate:"required"` // Base64 encoded data
	MimeType     string      `json:"mime_type,omitempty"`
	PngThumbnail []byte      `json:"png_thumbnail,omitempty"`
	ID           string      `json:"id,omitempty"`
	ContextInfo  ContextInfo `json:"context_info,omitempty"`
}

// SendLocationMessageRequest representa a requisição para envio de localização
type SendLocationMessageRequest struct {
	Phone       string      `json:"phone" validate:"required"`
	Name        string      `json:"name,omitempty"`
	Latitude    float64     `json:"latitude" validate:"required"`
	Longitude   float64     `json:"longitude" validate:"required"`
	ID          string      `json:"id,omitempty"`
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// SendContactMessageRequest representa a requisição para envio de contato
type SendContactMessageRequest struct {
	Phone       string      `json:"phone" validate:"required"`
	Name        string      `json:"name" validate:"required"`
	Vcard       string      `json:"vcard" validate:"required"`
	ID          string      `json:"id,omitempty"`
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// ButtonStruct representa um botão para mensagens interativas
type ButtonStruct struct {
	ButtonID   string `json:"button_id" validate:"required"`
	ButtonText string `json:"button_text" validate:"required"`
}

// SendButtonsMessageRequest representa a requisição para envio de botões
type SendButtonsMessageRequest struct {
	Phone       string         `json:"phone" validate:"required"`
	Title       string         `json:"title" validate:"required"`
	Buttons     []ButtonStruct `json:"buttons" validate:"required,min=1,max=3"`
	ID          string         `json:"id,omitempty"`
	ContextInfo ContextInfo    `json:"context_info,omitempty"`
}

// ListItem representa um item de lista
type ListItem struct {
	Title string `json:"title" validate:"required"`
	Desc  string `json:"desc,omitempty"`
	RowID string `json:"row_id" validate:"required"`
}

// ListSection representa uma seção de lista
type ListSection struct {
	Title string     `json:"title" validate:"required"`
	Rows  []ListItem `json:"rows" validate:"required,min=1"`
}

// SendListMessageRequest representa a requisição para envio de lista
type SendListMessageRequest struct {
	Phone      string        `json:"phone" validate:"required"`
	ButtonText string        `json:"button_text" validate:"required"`
	Desc       string        `json:"desc" validate:"required"`
	TopText    string        `json:"top_text" validate:"required"`
	Sections   []ListSection `json:"sections" validate:"required,min=1"`
	FooterText string        `json:"footer_text,omitempty"`
	ID         string        `json:"id,omitempty"`
}

// SendPollMessageRequest representa a requisição para envio de enquete
type SendPollMessageRequest struct {
	Phone   string   `json:"phone" validate:"required"`
	Header  string   `json:"header" validate:"required"`
	Options []string `json:"options" validate:"required,min=2"`
	ID      string   `json:"id,omitempty"`
}

// SendEditMessageRequest representa a requisição para edição de mensagem
type SendEditMessageRequest struct {
	Phone       string      `json:"phone" validate:"required"`
	Body        string      `json:"body" validate:"required"`
	ID          string      `json:"id" validate:"required"`
	ContextInfo ContextInfo `json:"context_info,omitempty"`
}

// DeleteMessageRequest representa a requisição para deletar mensagem
type DeleteMessageRequest struct {
	Phone string `json:"phone" validate:"required"`
	ID    string `json:"id" validate:"required"`
}

// ReactMessageRequest representa a requisição para reagir a mensagem
type ReactMessageRequest struct {
	Phone string `json:"phone" validate:"required"`
	Body  string `json:"body" validate:"required"` // Emoji ou "remove" para remover reação
	ID    string `json:"id" validate:"required"`
}

// ContextInfo representa informações de contexto para mensagens (reply, mentions)
type ContextInfo struct {
	// ID da mensagem original para reply (opcional)
	StanzaID *string `json:"stanza_id,omitempty" example:"3EB0C431C26A1916E07A"`
	// JID do participante da mensagem original (opcional)
	Participant *string `json:"participant,omitempty" example:"5511999999999@s.whatsapp.net"`
	// Lista de JIDs mencionados na mensagem (opcional)
	MentionedJID []string `json:"mentioned_jid,omitempty" example:"5511999999999@s.whatsapp.net,5511888888888@s.whatsapp.net"`
}
