package handlers

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/domain/entity"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// SessionHandler gerencia as rotas relacionadas às sessões usando use cases
type SessionHandler struct {
	createSessionUC  *usecase.CreateSessionUseCase
	connectSessionUC *usecase.ConnectSessionUseCase
	listSessionsUC   *usecase.ListSessionsUseCase
	getQRCodeUC      *usecase.GetQRCodeUseCase
	deleteSessionUC  *usecase.DeleteSessionUseCase
	logoutSessionUC  *usecase.LogoutSessionUseCase
	pairPhoneUC      *usecase.PairPhoneUseCase
	getSessionInfoUC *usecase.GetSessionInfoUseCase
	setProxyUC       *usecase.SetProxyUseCase
}

// NewSessionHandler cria uma nova instância do handler de sessões com use cases
func NewSessionHandler(
	createSessionUC *usecase.CreateSessionUseCase,
	connectSessionUC *usecase.ConnectSessionUseCase,
	listSessionsUC *usecase.ListSessionsUseCase,
	getQRCodeUC *usecase.GetQRCodeUseCase,
	deleteSessionUC *usecase.DeleteSessionUseCase,
	logoutSessionUC *usecase.LogoutSessionUseCase,
	pairPhoneUC *usecase.PairPhoneUseCase,
	getSessionInfoUC *usecase.GetSessionInfoUseCase,
	setProxyUC *usecase.SetProxyUseCase,
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
// @Summary Cria uma nova sessão WhatsApp
// @Description Cria uma nova sessão WhatsApp com as configurações especificadas
// @Tags sessions
// @Accept json
// @Produce json
// @Param request body requests.CreateSessionRequest true "Dados da sessão"
// @Success 200 {object} responses.APIResponse{data=entity.Session} "Sessão criada com sucesso"
// @Failure 400 {object} responses.APIResponse "Dados inválidos"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions/add [post]
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
// @Summary Lista todas as sessões WhatsApp
// @Description Retorna uma lista com todas as sessões WhatsApp cadastradas no sistema
// @Tags sessions
// @Accept json
// @Produce json
// @Success 200 {object} responses.APIResponse{data=[]entity.Session} "Lista de sessões"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions [get]
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
// @Summary Conecta uma sessão WhatsApp
// @Description Inicia a conexão de uma sessão WhatsApp específica
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} responses.APIResponse "Sessão conectada com sucesso"
// @Failure 404 {object} responses.APIResponse "Sessão não encontrada"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionId}/connect [post]
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
// @Summary Obtém QR Code para autenticação
// @Description Retorna o QR Code para autenticação de uma sessão WhatsApp
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} responses.APIResponse{data=object} "QR Code gerado com sucesso"
// @Failure 404 {object} responses.APIResponse "Sessão não encontrada"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionId}/qr [get]
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
// @Summary Remove uma sessão WhatsApp
// @Description Remove permanentemente uma sessão WhatsApp do sistema
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} responses.APIResponse "Sessão deletada com sucesso"
// @Failure 404 {object} responses.APIResponse "Sessão não encontrada"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionId} [delete]
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
// @Summary Faz logout de uma sessão WhatsApp
// @Description Desconecta uma sessão WhatsApp específica
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} responses.APIResponse "Logout realizado com sucesso"
// @Failure 404 {object} responses.APIResponse "Sessão não encontrada"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionId}/logout [post]
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
// @Summary Obtém informações de uma sessão específica
// @Description Retorna informações detalhadas de uma sessão WhatsApp específica
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} responses.APIResponse{data=entity.Session} "Informações da sessão"
// @Failure 404 {object} responses.APIResponse "Sessão não encontrada"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionId} [get]
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
// @Summary Emparelha um telefone com a sessão
// @Description Gera código de emparelhamento para conectar um telefone à sessão WhatsApp
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body requests.PairPhoneRequest true "Dados do telefone"
// @Success 200 {object} responses.APIResponse{data=object} "Código de emparelhamento gerado"
// @Failure 400 {object} responses.APIResponse "Dados inválidos"
// @Failure 404 {object} responses.APIResponse "Sessão não encontrada"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionId}/pair [post]
func (h *SessionHandler) PairPhone(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	var req requests.PairPhoneRequest

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
// @Summary Configura proxy para a sessão
// @Description Define configurações de proxy para conexão da sessão WhatsApp
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body requests.SetProxyRequest true "Configurações do proxy"
// @Success 200 {object} responses.APIResponse "Proxy configurado com sucesso"
// @Failure 400 {object} responses.APIResponse "Dados inválidos"
// @Failure 404 {object} responses.APIResponse "Sessão não encontrada"
// @Failure 500 {object} responses.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionId}/proxy [post]
func (h *SessionHandler) SetProxy(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")

	var req requests.SetProxyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.JSON(w, r, responses.APIResponse{
			Success: false,
			Error:   "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Criar configuração de proxy
	proxyConfig := &entity.ProxyConfig{
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
