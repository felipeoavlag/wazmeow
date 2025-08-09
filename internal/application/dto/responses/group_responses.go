package responses

import "go.mau.fi/whatsmeow/types"

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
type GroupJoinResponse struct {
	Details string `json:"details"`
}

// GroupInviteInfoResponse representa a resposta de informações do convite
type GroupInviteInfoResponse struct {
	GroupInfo types.GroupInfo `json:"groupInfo"`
}

// UpdateGroupParticipantsResponse representa a resposta de atualização de participantes
type UpdateGroupParticipantsResponse struct {
	Details string `json:"details"`
}

// NewsletterListResponse representa a resposta de listagem de newsletters
type NewsletterListResponse struct {
	Newsletters []types.NewsletterMetadata `json:"newsletters"`
}
