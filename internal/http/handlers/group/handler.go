package group

import (
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/http/handlers/base"
	"wazmeow/internal/http/handlers/middleware"
)

// Handler contém os handlers para operações de grupo refatorados
type Handler struct {
	*base.BaseHandler
	createGroupUseCase             *usecase.CreateGroupUseCase
	setGroupPhotoUseCase           *usecase.SetGroupPhotoUseCase
	updateGroupParticipantsUseCase *usecase.UpdateGroupParticipantsUseCase
	leaveGroupUseCase              *usecase.LeaveGroupUseCase
	joinGroupUseCase               *usecase.JoinGroupUseCase
	getGroupInfoUseCase            *usecase.GetGroupInfoUseCase
	listGroupsUseCase              *usecase.ListGroupsUseCase
	getGroupInviteLinkUseCase      *usecase.GetGroupInviteLinkUseCase
	revokeGroupInviteLinkUseCase   *usecase.RevokeGroupInviteLinkUseCase
	setGroupNameUseCase            *usecase.SetGroupNameUseCase
	setGroupTopicUseCase           *usecase.SetGroupTopicUseCase
}

// NewHandler cria uma nova instância dos handlers de grupo refatorados
func NewHandler(
	createGroupUseCase *usecase.CreateGroupUseCase,
	setGroupPhotoUseCase *usecase.SetGroupPhotoUseCase,
	updateGroupParticipantsUseCase *usecase.UpdateGroupParticipantsUseCase,
	leaveGroupUseCase *usecase.LeaveGroupUseCase,
	joinGroupUseCase *usecase.JoinGroupUseCase,
	getGroupInfoUseCase *usecase.GetGroupInfoUseCase,
	listGroupsUseCase *usecase.ListGroupsUseCase,
	getGroupInviteLinkUseCase *usecase.GetGroupInviteLinkUseCase,
	revokeGroupInviteLinkUseCase *usecase.RevokeGroupInviteLinkUseCase,
	setGroupNameUseCase *usecase.SetGroupNameUseCase,
	setGroupTopicUseCase *usecase.SetGroupTopicUseCase,
) *Handler {
	return &Handler{
		BaseHandler:                    base.NewBaseHandler(),
		createGroupUseCase:             createGroupUseCase,
		setGroupPhotoUseCase:           setGroupPhotoUseCase,
		updateGroupParticipantsUseCase: updateGroupParticipantsUseCase,
		leaveGroupUseCase:              leaveGroupUseCase,
		joinGroupUseCase:               joinGroupUseCase,
		getGroupInfoUseCase:            getGroupInfoUseCase,
		listGroupsUseCase:              listGroupsUseCase,
		getGroupInviteLinkUseCase:      getGroupInviteLinkUseCase,
		revokeGroupInviteLinkUseCase:   revokeGroupInviteLinkUseCase,
		setGroupNameUseCase:            setGroupNameUseCase,
		setGroupTopicUseCase:           setGroupTopicUseCase,
	}
}

// CreateGroup cria um novo grupo
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.CreateGroupRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"name":         req.Name,
		"participants": req.Participants,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "criar grupo", func() (interface{}, error) {
		return h.createGroupUseCase.Execute(sessionID, &req)
	}, "Grupo criado com sucesso")
}

// SetGroupPhoto define foto do grupo
func (h *Handler) SetGroupPhoto(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SetGroupPhotoRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"group_id": req.GroupID,
		"image":    req.Image,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "definir foto do grupo", func() (interface{}, error) {
		return h.setGroupPhotoUseCase.Execute(sessionID, &req)
	}, "Foto do grupo definida com sucesso")
}

// UpdateGroupParticipants atualiza participantes do grupo
func (h *Handler) UpdateGroupParticipants(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.UpdateGroupParticipantsRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"group_id":     req.GroupID,
		"participants": req.Participants,
		"action":       req.Action,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "atualizar participantes do grupo", func() (interface{}, error) {
		return h.updateGroupParticipantsUseCase.Execute(sessionID, &req)
	}, "Participantes do grupo atualizados com sucesso")
}

// LeaveGroup sai do grupo
func (h *Handler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.LeaveGroupRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"group_id": req.GroupID,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "sair do grupo", func() (interface{}, error) {
		return h.leaveGroupUseCase.Execute(sessionID, &req)
	}, "Saiu do grupo com sucesso")
}

// JoinGroup entra no grupo via link
func (h *Handler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.JoinGroupRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"code": req.Code,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "entrar no grupo", func() (interface{}, error) {
		return h.joinGroupUseCase.Execute(sessionID, &req)
	}, "Entrou no grupo com sucesso")
}

// GetGroupInfo obtém informações do grupo
func (h *Handler) GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.GetGroupInfoRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"group_id": req.GroupID,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "obter informações do grupo", func() (interface{}, error) {
		return h.getGroupInfoUseCase.Execute(sessionID, &req)
	}, "Informações do grupo obtidas com sucesso")
}

// ListGroups lista grupos
func (h *Handler) ListGroups(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "listar grupos", func() (interface{}, error) {
		return h.listGroupsUseCase.Execute(sessionID)
	}, "Grupos listados com sucesso")
}

// GetGroupInviteLink obtém link de convite do grupo
func (h *Handler) GetGroupInviteLink(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.GetGroupInviteLinkRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"group_id": req.GroupID,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "obter link de convite do grupo", func() (interface{}, error) {
		return h.getGroupInviteLinkUseCase.Execute(sessionID, &req)
	}, "Link de convite do grupo obtido com sucesso")
}

// RevokeGroupInviteLink revoga link de convite do grupo
func (h *Handler) RevokeGroupInviteLink(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.RevokeGroupInviteLinkRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"group_id": req.GroupID,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "revogar link de convite do grupo", func() (interface{}, error) {
		return h.revokeGroupInviteLinkUseCase.Execute(sessionID, &req)
	}, "Link de convite do grupo revogado com sucesso")
}

// SetGroupName define nome do grupo
func (h *Handler) SetGroupName(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SetGroupNameRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"group_id": req.GroupID,
		"name":     req.Name,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "definir nome do grupo", func() (interface{}, error) {
		return h.setGroupNameUseCase.Execute(sessionID, &req)
	}, "Nome do grupo definido com sucesso")
}

// SetGroupTopic define tópico do grupo
func (h *Handler) SetGroupTopic(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SetGroupTopicRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"group_id": req.GroupID,
		"topic":    req.Topic,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "definir tópico do grupo", func() (interface{}, error) {
		return h.setGroupTopicUseCase.Execute(sessionID, &req)
	}, "Tópico do grupo definido com sucesso")
}