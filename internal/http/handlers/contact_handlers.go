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
// @Summary Obtém informações de um usuário
// @Description Retorna informações detalhadas de um usuário WhatsApp específico
// @Tags contacts
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.GetUserInfoRequest true "Dados do usuário"
// @Success 200 {object} map[string]interface{} "Informações do usuário obtidas com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /contact/{sessionID}/info [post]
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
// @Summary Verifica se um usuário existe no WhatsApp
// @Description Verifica se um número de telefone possui conta WhatsApp ativa
// @Tags contacts
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.CheckUserRequest true "Dados do usuário para verificar"
// @Success 200 {object} map[string]interface{} "Usuário verificado com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /contact/{sessionID}/check [post]
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
// @Summary Obtém avatar de um usuário
// @Description Baixa a foto de perfil de um usuário WhatsApp e retorna em base64
// @Tags contacts
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.GetAvatarRequest true "Dados do usuário"
// @Success 200 {object} map[string]interface{} "Avatar obtido com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /contact/{sessionID}/avatar [post]
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
// @Summary Obtém lista de contatos
// @Description Retorna a lista de contatos do usuário WhatsApp
// @Tags contacts
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} map[string]interface{} "Contatos obtidos com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /contact/{sessionID}/list [get]
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
