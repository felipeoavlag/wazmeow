package responses

import "go.mau.fi/whatsmeow/types"

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
