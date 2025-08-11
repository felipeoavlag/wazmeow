package responses

import (
	"wazmeow/internal/domain/entity"

	"go.mau.fi/whatsmeow/types"
)

// ===============================
// BASE RESPONSES
// ===============================

// SimpleResponse representa uma resposta simples com apenas detalhes
type SimpleResponse struct {
	Details string `json:"details"`
}

// TimestampedResponse representa uma resposta com detalhes e timestamp
type TimestampedResponse struct {
	Details   string `json:"details"`
	Timestamp int64  `json:"timestamp"`
}

// ===============================
// CHAT RESPONSES
// ===============================

// ChatPresenceResponse representa a resposta de presença no chat
type ChatPresenceResponse = SimpleResponse

// MarkReadResponse representa a resposta de marcar como lida
type MarkReadResponse = SimpleResponse

// DownloadResponse representa a resposta de download de mídia
type DownloadResponse struct {
	Mimetype string `json:"mimetype"`
	Data     string `json:"data"` // Base64 encoded data
}

// HistorySyncResponse representa a resposta de sincronização de histórico
type HistorySyncResponse = TimestampedResponse

// ===============================
// CONTACT RESPONSES
// ===============================

// PresenceResponse representa a resposta de definição de presença
type PresenceResponse = SimpleResponse

// UserInfo representa informações de um usuário
type UserInfo struct {
	Query        string `json:"query"`
	IsInWhatsapp bool   `json:"isInWhatsapp"`
	JID          string `json:"jid"`
	VerifiedName string `json:"verifiedName,omitempty"`
}

// UserInfoResponse representa a resposta de informações do usuário
type UserInfoResponse struct {
	Users map[types.JID]types.UserInfo `json:"users"`
}

// CheckUserResponse representa a resposta de verificação de usuário
type CheckUserResponse struct {
	Users []UserInfo `json:"users"`
}

// AvatarResponse representa a resposta de avatar do usuário
type AvatarResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// ContactsResponse representa a resposta de contatos
type ContactsResponse struct {
	Contacts map[types.JID]types.ContactInfo `json:"contacts"`
}

// ===============================
// GROUP RESPONSES
// ===============================

// GroupResponse representa a resposta básica de operações de grupo
type GroupResponse struct {
	GroupID      string   `json:"groupID,omitempty"`
	Name         string   `json:"name,omitempty"`
	Participants []string `json:"participants,omitempty"`
	Details      string   `json:"details"`
}

// CreateGroupResponse representa a resposta de criação de grupo
type CreateGroupResponse struct {
	GroupJID     string   `json:"groupJID"`
	GroupName    string   `json:"groupName"`
	Participants []string `json:"participants"`
}

// GroupListResponse representa a resposta de listagem de grupos
type GroupListResponse struct {
	Groups []GroupSummary `json:"groups"`
	Count  int            `json:"count"`
}

// GroupSummary representa um resumo de grupo
type GroupSummary struct {
	GroupID string `json:"groupID"`
	Name    string `json:"name"`
}

// GroupInfoResponse representa a resposta de informações do grupo
type GroupInfoResponse struct {
	GroupID      string             `json:"groupID"`
	Name         string             `json:"name"`
	Topic        string             `json:"topic"`
	Owner        string             `json:"owner"`
	CreatedAt    int64              `json:"createdAt"`
	Participants []GroupParticipant `json:"participants"`
	Size         int                `json:"size"`
}

// GroupParticipant representa um participante do grupo
type GroupParticipant struct {
	JID          string `json:"jid"`
	IsAdmin      bool   `json:"isAdmin"`
	IsSuperAdmin bool   `json:"isSuperAdmin"`
}

// GroupInviteLinkResponse representa a resposta de link de convite
type GroupInviteLinkResponse struct {
	GroupID string `json:"groupID"`
	Link    string `json:"link"`
}

// SetGroupPhotoResponse representa a resposta de definição de foto
type SetGroupPhotoResponse struct {
	Details   string `json:"details"`
	PictureID string `json:"pictureID,omitempty"`
}

// GroupJoinResponse representa a resposta de entrada no grupo
type GroupJoinResponse = SimpleResponse

// GroupInviteInfoResponse representa a resposta de informações do convite
type GroupInviteInfoResponse struct {
	GroupInfo types.GroupInfo `json:"groupInfo"`
}

// UpdateGroupParticipantsResponse representa a resposta de atualização de participantes
type UpdateGroupParticipantsResponse = SimpleResponse

// NewsletterListResponse representa a resposta de listagem de newsletters
type NewsletterListResponse struct {
	Newsletters []types.NewsletterMetadata `json:"newsletters"`
}

// ===============================
// MESSAGE RESPONSES
// ===============================

// SendMessageResponse representa a resposta do envio de mensagem
type SendMessageResponse struct {
	// Detalhes sobre o envio da mensagem
	Details string `json:"details" example:"Mensagem enviada com sucesso"`
	// Timestamp Unix do envio
	Timestamp int64 `json:"timestamp" example:"1692454800"`
	// ID único da mensagem enviada
	ID string `json:"id" example:"3EB0C431C26A1916E07A"`
}

// ===============================
// SESSION RESPONSES
// ===============================

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

// ===============================
// WEBHOOK RESPONSES
// ===============================

// WebhookResponse representa a resposta de operações de webhook
type WebhookResponse struct {
	Webhook   string   `json:"webhook"`
	Events    []string `json:"events,omitempty"`
	Active    bool     `json:"active,omitempty"`
	Subscribe []string `json:"subscribe,omitempty"`
}

// WebhookDeleteResponse representa a resposta de exclusão de webhook
type WebhookDeleteResponse = SimpleResponse
