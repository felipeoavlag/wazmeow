package handlers

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// UserHandlers contém os handlers para operações de usuário
type UserHandlers struct {
	getUserInfoUseCase *usecase.GetUserInfoUseCase
	checkUserUseCase   *usecase.CheckUserUseCase
	getAvatarUseCase   *usecase.GetAvatarUseCase
	getContactsUseCase *usecase.GetContactsUseCase
}

// NewUserHandlers cria uma nova instância dos handlers de usuário
func NewUserHandlers(
	getUserInfoUseCase *usecase.GetUserInfoUseCase,
	checkUserUseCase *usecase.CheckUserUseCase,
	getAvatarUseCase *usecase.GetAvatarUseCase,
	getContactsUseCase *usecase.GetContactsUseCase,
) *UserHandlers {
	return &UserHandlers{
		getUserInfoUseCase: getUserInfoUseCase,
		checkUserUseCase:   checkUserUseCase,
		getAvatarUseCase:   getAvatarUseCase,
		getContactsUseCase: getContactsUseCase,
	}
}

// GetUserInfo obtém informações do usuário
func (h *UserHandlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.GetUserInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.getUserInfoUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao obter informações do usuário: %v", err)
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

	logger.Info("Informações do usuário obtidas com sucesso - Session: %s", sessionID)
}

// CheckUser verifica se usuário existe
func (h *UserHandlers) CheckUser(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.CheckUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.checkUserUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao verificar usuário: %v", err)
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

	logger.Info("Usuário verificado com sucesso - Session: %s", sessionID)
}

// GetAvatar obtém avatar do usuário
func (h *UserHandlers) GetAvatar(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.GetAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.getAvatarUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao obter avatar: %v", err)
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

	logger.Info("Avatar obtido com sucesso - Session: %s", sessionID)
}

// GetContacts obtém lista de contatos
func (h *UserHandlers) GetContacts(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.getContactsUseCase.Execute(sessionID)
	if err != nil {
		logger.Error("Erro ao obter contatos: %v", err)
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

	logger.Info("Contatos obtidos com sucesso - Session: %s", sessionID)
}
