package handlers

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/requests"
	"wazmeow/internal/domain/responses"
	"wazmeow/internal/domain/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// SessionHandler gerencia as rotas relacionadas às sessões usando use cases
type SessionHandler struct {
	createSessionUC  *usecases.CreateSessionUseCase
	connectSessionUC *usecases.ConnectSessionUseCase
	listSessionsUC   *usecases.ListSessionsUseCase
	getQRCodeUC      *usecases.GetQRCodeUseCase
	deleteSessionUC  *usecases.DeleteSessionUseCase
	logoutSessionUC  *usecases.LogoutSessionUseCase
	pairPhoneUC      *usecases.PairPhoneUseCase
	getSessionInfoUC *usecases.GetSessionInfoUseCase
	setProxyUC       *usecases.SetProxyUseCase
}

// NewSessionHandler cria uma nova instância do handler de sessões com use cases
func NewSessionHandler(
	createSessionUC *usecases.CreateSessionUseCase,
	connectSessionUC *usecases.ConnectSessionUseCase,
	listSessionsUC *usecases.ListSessionsUseCase,
	getQRCodeUC *usecases.GetQRCodeUseCase,
	deleteSessionUC *usecases.DeleteSessionUseCase,
	logoutSessionUC *usecases.LogoutSessionUseCase,
	pairPhoneUC *usecases.PairPhoneUseCase,
	getSessionInfoUC *usecases.GetSessionInfoUseCase,
	setProxyUC *usecases.SetProxyUseCase,
) *SessionHandler {
	return &SessionHandler{
		createSessionUC:  createSessionUC,
		connectSessionUC: connectSessionUC,
		listSessionsUC:   listSessionsUC,
		getQRCodeUC:      getQRCodeUC,
		deleteSessionUC:  deleteSessionUC,
		logoutSessionUC:  logoutSessionUC,
		pairPhoneUC:      pairPhoneUC,
		getSessionInfoUC: getSessionInfoUC,
		setProxyUC:       setProxyUC,
	}
}

// CreateSession cria uma nova sessão
// POST /sessions/add
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateSessionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Executar use case
	session, err := h.createSessionUC.Execute(&req)
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Message: "Sessão criada com sucesso",
		Data:    session,
	})
}

// ListSessions lista todas as sessões
// GET /sessions
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.listSessionsUC.Execute()
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Data:    sessions,
	})
}

// ConnectSession conecta uma sessão
// POST /sessions/{sessionId}/connect
func (h *SessionHandler) ConnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	err := h.connectSessionUC.Execute(sessionID)
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Message: "Sessão conectada com sucesso",
	})
}

// GetQRCode obtém o QR code para uma sessão
// GET /sessions/{sessionId}/qr
func (h *SessionHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	qrResponse, err := h.getQRCodeUC.Execute(sessionID)
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Data:    qrResponse,
	})
}

// DeleteSession deleta uma sessão
// DELETE /sessions/{sessionId}
func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	err := h.deleteSessionUC.Execute(sessionID)
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Message: "Sessão deletada com sucesso",
	})
}

// LogoutSession faz logout de uma sessão
// POST /sessions/{sessionId}/logout
func (h *SessionHandler) LogoutSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	err := h.logoutSessionUC.Execute(sessionID)
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Message: "Logout realizado com sucesso",
	})
}

// GetSessionInfo obtém informações de uma sessão
// GET /sessions/{sessionId}
func (h *SessionHandler) GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	sessionInfo, err := h.getSessionInfoUC.Execute(sessionID)
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Data:    sessionInfo,
	})
}

// PairPhone emparelha um telefone com uma sessão
// POST /sessions/{sessionId}/pair
func (h *SessionHandler) PairPhone(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	var req struct {
		Phone string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   "Dados inválidos: " + err.Error(),
		})
		return
	}

	code, err := h.pairPhoneUC.Execute(sessionID, req.Phone)
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Message: "Código de emparelhamento gerado",
		Data:    map[string]string{"code": code},
	})
}

// SetProxy configura proxy para uma sessão
// POST /sessions/{sessionId}/proxy
func (h *SessionHandler) SetProxy(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	var req struct {
		Type     string `json:"type"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Criar configuração de proxy
	proxyConfig := &entities.ProxyConfig{
		Type:     req.Type,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}

	err := h.setProxyUC.Execute(sessionID, proxyConfig)
	if err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	render.JSON(w, r, responses.APIResponse{
		Success: true,
		Message: "Proxy configurado com sucesso",
	})
}
