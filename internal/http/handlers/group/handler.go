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
// @Summary Cria um novo grupo WhatsApp
// @Description Cria um novo grupo WhatsApp com nome e participantes especificados
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.CreateGroupRequest true "Dados do grupo"
// @Success 200 {object} base.APIResponse "Grupo criado com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/create [post]
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
// @Summary Define foto do grupo
// @Description Define ou atualiza a foto de perfil de um grupo WhatsApp
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SetGroupPhotoRequest true "Dados da foto"
// @Success 200 {object} base.APIResponse "Foto do grupo definida com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/photo [post]
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
// @Summary Atualiza participantes do grupo
// @Description Adiciona ou remove participantes de um grupo WhatsApp
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.UpdateGroupParticipantsRequest true "Dados dos participantes"
// @Success 200 {object} base.APIResponse "Participantes atualizados com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/participants [post]
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
// @Summary Sai de um grupo
// @Description Remove a sessão atual de um grupo WhatsApp
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.LeaveGroupRequest true "Dados do grupo"
// @Success 200 {object} base.APIResponse "Saiu do grupo com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/leave [post]
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
// @Summary Entra em um grupo via link
// @Description Entra em um grupo WhatsApp usando link de convite
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.JoinGroupRequest true "Link do grupo"
// @Success 200 {object} base.APIResponse "Entrou no grupo com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/join [post]
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
// @Summary Obtém informações do grupo
// @Description Retorna informações detalhadas sobre um grupo WhatsApp
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.GetGroupInfoRequest true "ID do grupo"
// @Success 200 {object} base.APIResponse "Informações do grupo"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/info [post]
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
// @Summary Lista grupos da sessão
// @Description Retorna lista de todos os grupos WhatsApp da sessão
// @Tags groups
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "Lista de grupos"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/list [get]
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
// @Summary Obtém link de convite do grupo
// @Description Gera ou obtém o link de convite de um grupo WhatsApp
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.GetGroupInviteLinkRequest true "ID do grupo"
// @Success 200 {object} base.APIResponse "Link de convite do grupo"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/invite-link [post]
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
// @Summary Revoga link de convite do grupo
// @Description Revoga o link de convite atual e gera um novo para o grupo
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.RevokeGroupInviteLinkRequest true "ID do grupo"
// @Success 200 {object} base.APIResponse "Link de convite revogado com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/revoke-invite-link [post]
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
// @Summary Define nome do grupo
// @Description Altera o nome de um grupo WhatsApp
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SetGroupNameRequest true "Novo nome do grupo"
// @Success 200 {object} base.APIResponse "Nome do grupo definido com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/name [post]
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
// @Summary Define tópico/descrição do grupo
// @Description Altera a descrição/tópico de um grupo WhatsApp
// @Tags groups
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SetGroupTopicRequest true "Novo tópico do grupo"
// @Success 200 {object} base.APIResponse "Tópico do grupo definido com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /group/{sessionID}/topic [post]
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
