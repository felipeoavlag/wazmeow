package handlers

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// GroupHandlers contém os handlers para operações de grupo
type GroupHandlers struct {
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

// NewGroupHandlers cria uma nova instância dos handlers de grupo
func NewGroupHandlers(
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
) *GroupHandlers {
	return &GroupHandlers{
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
func (h *GroupHandlers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.createGroupUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao criar grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Grupo criado com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Grupo criado com sucesso - Session: %s", sessionID)
}

// SetGroupPhoto define foto do grupo
func (h *GroupHandlers) SetGroupPhoto(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SetGroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.setGroupPhotoUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao definir foto do grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Foto do grupo definida com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Foto do grupo definida com sucesso - Session: %s", sessionID)
}

// UpdateGroupParticipants atualiza participantes do grupo
func (h *GroupHandlers) UpdateGroupParticipants(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.UpdateGroupParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.updateGroupParticipantsUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao atualizar participantes do grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Participantes do grupo atualizados com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Participantes do grupo atualizados com sucesso - Session: %s", sessionID)
}

// LeaveGroup sai do grupo
func (h *GroupHandlers) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.LeaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.leaveGroupUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao sair do grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Saiu do grupo com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Saiu do grupo com sucesso - Session: %s", sessionID)
}

// JoinGroup entra no grupo via link
func (h *GroupHandlers) JoinGroup(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.joinGroupUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao entrar no grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Entrou no grupo com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Entrou no grupo com sucesso - Session: %s", sessionID)
}

// GetGroupInfo obtém informações do grupo
func (h *GroupHandlers) GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.GetGroupInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.getGroupInfoUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao obter informações do grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Informações do grupo obtidas com sucesso - Session: %s", sessionID)
}

// ListGroups lista grupos
func (h *GroupHandlers) ListGroups(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.listGroupsUseCase.Execute(sessionID)
	if err != nil {
		logger.Error("Erro ao listar grupos: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Grupos listados com sucesso - Session: %s", sessionID)
}

// GetGroupInviteLink obtém link de convite do grupo
func (h *GroupHandlers) GetGroupInviteLink(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.GetGroupInviteLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.getGroupInviteLinkUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao obter link de convite do grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Link de convite do grupo obtido com sucesso - Session: %s", sessionID)
}

// RevokeGroupInviteLink revoga link de convite do grupo
func (h *GroupHandlers) RevokeGroupInviteLink(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.RevokeGroupInviteLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.revokeGroupInviteLinkUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao revogar link de convite do grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Link de convite do grupo revogado com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Link de convite do grupo revogado com sucesso - Session: %s", sessionID)
}

// SetGroupName define nome do grupo
func (h *GroupHandlers) SetGroupName(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SetGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.setGroupNameUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao definir nome do grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Nome do grupo definido com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Nome do grupo definido com sucesso - Session: %s", sessionID)
}

// SetGroupTopic define tópico do grupo
func (h *GroupHandlers) SetGroupTopic(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SetGroupTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.setGroupTopicUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao definir tópico do grupo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tópico do grupo definido com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Tópico do grupo definido com sucesso - Session: %s", sessionID)
}
