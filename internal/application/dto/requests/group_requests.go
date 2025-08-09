package requests

import "time"

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
