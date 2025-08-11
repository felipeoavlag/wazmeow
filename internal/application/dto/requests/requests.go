package requests

import (
	"time"

	"go.mau.fi/whatsmeow/types"
)

// =============================================================================
// CONTACT REQUESTS
// =============================================================================

// GetUserInfoRequest representa a requisição para obter informações do usuário
type GetUserInfoRequest struct {
	Phone []string `json:"phone" validate:"required,min=1"`
}

// CheckUserRequest representa a requisição para verificar se usuário está no WhatsApp
type CheckUserRequest struct {
	Phone []string `json:"phone" validate:"required,min=1"`
}

// GetAvatarRequest representa a requisição para obter avatar do usuário
type GetAvatarRequest struct {
	Phone   string `json:"phone" validate:"required"`
	Preview bool   `json:"preview"`
}

// =============================================================================
// WEBHOOK REQUESTS
// =============================================================================

// SetWebhookRequest representa a requisição para definir webhook
type SetWebhookRequest struct {
	WebhookURL string   `json:"webhookurl" validate:"omitempty,url"`
	Events     []string `json:"events,omitempty"`
	Enabled    *bool    `json:"enabled,omitempty"`
}

// UpdateWebhookRequest representa a requisição para atualizar webhook
type UpdateWebhookRequest struct {
	WebhookURL string   `json:"webhook" validate:"required,url"`
	Events     []string `json:"events,omitempty"`
	Active     bool     `json:"active"`
}

// =============================================================================
// SESSION REQUESTS
// =============================================================================

// SessionRequest representa a requisição para operações de sessão
type SessionRequest struct {
	// Nome da sessão
	Name string `json:"name" example:"minha-sessao"`
}

// CreateSessionRequest representa a requisição para criar uma sessão
type CreateSessionRequest struct {
	// Nome único da sessão
	Name string `json:"name" example:"minha-sessao" validate:"required"`
	// URL do webhook (opcional)
	WebhookURL string `json:"webhookUrl,omitempty" example:"https://meusite.com/webhook"`
	// Configuração de proxy (opcional)
	Proxy *ProxyConfig `json:"proxy,omitempty"`
}

// ProxyConfig representa a configuração de proxy para criação de sessão
type ProxyConfig struct {
	// Tipo do proxy (http, socks5)
	Type string `json:"type" example:"http" validate:"required"`
	// Host do servidor proxy
	Host string `json:"host" example:"proxy.example.com" validate:"required"`
	// Porta do servidor proxy
	Port int `json:"port" example:"8080" validate:"required,min=1,max=65535"`
	// Nome de usuário para autenticação (opcional)
	Username string `json:"username,omitempty" example:"usuario"`
	// Senha para autenticação (opcional)
	Password string `json:"password,omitempty" example:"senha"`
}

// PairPhoneRequest representa a requisição para emparelhar telefone
type PairPhoneRequest struct {
	// Número de telefone com código do país (formato: +5511999999999)
	Phone string `json:"phone" example:"+5511999999999" validate:"required"`
}

// SetProxyRequest representa a requisição para configurar proxy
type SetProxyRequest struct {
	// Tipo do proxy (http, socks5)
	Type string `json:"type" example:"http" validate:"required"`
	// Host do servidor proxy
	Host string `json:"host" example:"proxy.example.com" validate:"required"`
	// Porta do servidor proxy
	Port int `json:"port" example:"8080" validate:"required,min=1,max=65535"`
	// Nome de usuário para autenticação (opcional)
	Username string `json:"username,omitempty" example:"usuario"`
	// Senha para autenticação (opcional)
	Password string `json:"password,omitempty" example:"senha"`
}

// ===============================
// CHAT REQUESTS
// ===============================

// SendPresenceRequest representa a requisição para definir presença global
type SendPresenceRequest struct {
	Type string `json:"type" validate:"required,oneof=available unavailable"`
}

// ChatPresenceRequest representa a requisição para definir presença no chat
type ChatPresenceRequest struct {
	Phone string                  `json:"phone" validate:"required"`
	State string                  `json:"state" validate:"required,oneof=typing paused recording"`
	Media types.ChatPresenceMedia `json:"media,omitempty"`
}

// MarkReadRequest representa a requisição para marcar mensagens como lidas
type MarkReadRequest struct {
	ID     []string  `json:"id" validate:"required,min=1"`
	Chat   types.JID `json:"chat" validate:"required"`
	Sender types.JID `json:"sender,omitempty"`
}

// DownloadImageRequest representa a requisição para download de imagem
type DownloadImageRequest struct {
	URL           string `json:"url" validate:"required"`
	DirectPath    string `json:"directPath" validate:"required"`
	MediaKey      []byte `json:"mediaKey" validate:"required"`
	Mimetype      string `json:"mimetype" validate:"required"`
	FileEncSHA256 []byte `json:"fileEncSHA256" validate:"required"`
	FileSHA256    []byte `json:"fileSHA256" validate:"required"`
	FileLength    uint64 `json:"fileLength" validate:"required"`
}

// DownloadVideoRequest representa a requisição para download de vídeo
type DownloadVideoRequest struct {
	URL           string `json:"url" validate:"required"`
	DirectPath    string `json:"directPath" validate:"required"`
	MediaKey      []byte `json:"mediaKey" validate:"required"`
	Mimetype      string `json:"mimetype" validate:"required"`
	FileEncSHA256 []byte `json:"fileEncSHA256" validate:"required"`
	FileSHA256    []byte `json:"fileSHA256" validate:"required"`
	FileLength    uint64 `json:"fileLength" validate:"required"`
}

// DownloadAudioRequest representa a requisição para download de áudio
type DownloadAudioRequest struct {
	URL           string `json:"url" validate:"required"`
	DirectPath    string `json:"directPath" validate:"required"`
	MediaKey      []byte `json:"mediaKey" validate:"required"`
	Mimetype      string `json:"mimetype" validate:"required"`
	FileEncSHA256 []byte `json:"fileEncSHA256" validate:"required"`
	FileSHA256    []byte `json:"fileSHA256" validate:"required"`
	FileLength    uint64 `json:"fileLength" validate:"required"`
}

// DownloadDocumentRequest representa a requisição para download de documento
type DownloadDocumentRequest struct {
	URL           string `json:"url" validate:"required"`
	DirectPath    string `json:"directPath" validate:"required"`
	MediaKey      []byte `json:"mediaKey" validate:"required"`
	Mimetype      string `json:"mimetype" validate:"required"`
	FileEncSHA256 []byte `json:"fileEncSHA256" validate:"required"`
	FileSHA256    []byte `json:"fileSHA256" validate:"required"`
	FileLength    uint64 `json:"fileLength" validate:"required"`
}

// ===============================
// GROUP REQUESTS
// ===============================

// CreateGroupRequest representa a requisição para criar grupo
type CreateGroupRequest struct {
	Name         string   `json:"name" validate:"required"`
	Participants []string `json:"participants" validate:"required,min=1"`
}

// GetGroupInfoRequest representa a requisição para obter informações do grupo
type GetGroupInfoRequest struct {
	GroupID string `json:"groupjID" validate:"required"`
}

// GetGroupInviteLinkRequest representa a requisição para obter link de convite
type GetGroupInviteLinkRequest struct {
	GroupID string `json:"groupJID" validate:"required"`
	Reset   bool   `json:"reset"`
}

// SetGroupPhotoRequest representa a requisição para definir foto do grupo
type SetGroupPhotoRequest struct {
	GroupID string `json:"groupJID" validate:"required"`
	Image   string `json:"image" validate:"required"`
}

// RemoveGroupPhotoRequest representa a requisição para remover foto do grupo
type RemoveGroupPhotoRequest struct {
	GroupID string `json:"groupJID" validate:"required"`
}

// LeaveGroupRequest representa a requisição para sair do grupo
type LeaveGroupRequest struct {
	GroupID string `json:"groupJID" validate:"required"`
}

// JoinGroupRequest representa a requisição para entrar no grupo
type JoinGroupRequest struct {
	Code string `json:"code" validate:"required"`
}

// SetGroupNameRequest representa a requisição para definir nome do grupo
type SetGroupNameRequest struct {
	GroupID string `json:"groupJID" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

// SetGroupTopicRequest representa a requisição para definir descrição do grupo
type SetGroupTopicRequest struct {
	GroupID   string    `json:"groupJID" validate:"required"`
	Topic     string    `json:"topic" validate:"required"`
	TopicID   string    `json:"topicID,omitempty"`
	TopicTime time.Time `json:"topicTime,omitempty"`
}

// SetGroupAnnounceRequest representa a requisição para configurar anúncios
type SetGroupAnnounceRequest struct {
	GroupID  string `json:"groupJID" validate:"required"`
	Announce bool   `json:"announce"`
}

// SetGroupLockedRequest representa a requisição para bloquear grupo
type SetGroupLockedRequest struct {
	GroupID string `json:"groupJID" validate:"required"`
	Locked  bool   `json:"locked"`
}

// SetDisappearingTimerRequest representa a requisição para mensagens temporárias
type SetDisappearingTimerRequest struct {
	GroupID  string `json:"groupJID" validate:"required"`
	Duration string `json:"duration" validate:"required,oneof=24h 7d 90d off"`
}

// GetGroupInviteInfoRequest representa a requisição para informações do convite
type GetGroupInviteInfoRequest struct {
	Code string `json:"code" validate:"required"`
}

// UpdateGroupParticipantsRequest representa a requisição para atualizar participantes
type UpdateGroupParticipantsRequest struct {
	GroupID      string   `json:"groupJID" validate:"required"`
	Participants []string `json:"participants" validate:"required,min=1"`
	Action       string   `json:"action" validate:"required,oneof=add remove promote demote"`
}

// RevokeGroupInviteLinkRequest representa a requisição para revogar link de convite
type RevokeGroupInviteLinkRequest struct {
	GroupID string `json:"groupJID" validate:"required"`
}

// ===============================
// MESSAGE REQUESTS
// ===============================

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
