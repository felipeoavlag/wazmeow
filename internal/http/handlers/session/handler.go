package session

import (
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/domain/entity"
	"wazmeow/internal/http/handlers/base"
	"wazmeow/internal/http/handlers/middleware"
)

// Handler contém os handlers para operações de sessão refatorados
type Handler struct {
	*base.BaseHandler
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

// NewHandler cria uma nova instância dos handlers de sessão refatorados
func NewHandler(
	createSessionUC *usecase.CreateSessionUseCase,
	connectSessionUC *usecase.ConnectSessionUseCase,
	listSessionsUC *usecase.ListSessionsUseCase,
	getQRCodeUC *usecase.GetQRCodeUseCase,
	deleteSessionUC *usecase.DeleteSessionUseCase,
	logoutSessionUC *usecase.LogoutSessionUseCase,
	pairPhoneUC *usecase.PairPhoneUseCase,
	getSessionInfoUC *usecase.GetSessionInfoUseCase,
	setProxyUC *usecase.SetProxyUseCase,
) *Handler {
	return &Handler{
		BaseHandler:      base.NewBaseHandler(),
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
// @Description Cria uma nova sessão WhatsApp com nome único
// @Tags sessions
// @Accept json
// @Produce json
// @Param request body requests.CreateSessionRequest true "Dados da sessão"
// @Success 200 {object} base.APIResponse "Sessão criada com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/add [post]
func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateSessionRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"name": req.Name,
	}) {
		return
	}

	h.HandleUseCaseExecution(w, "criar sessão", func() (interface{}, error) {
		return h.createSessionUC.Execute(&req)
	}, "Sessão criada com sucesso")
}

// ListSessions lista todas as sessões
// @Summary Lista todas as sessões WhatsApp
// @Description Retorna uma lista com todas as sessões WhatsApp criadas
// @Tags sessions
// @Produce json
// @Success 200 {object} base.APIResponse "Lista de sessões"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions [get]
func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	h.HandleUseCaseExecution(w, "listar sessões", func() (interface{}, error) {
		return h.listSessionsUC.Execute()
	}, "Sessões listadas com sucesso")
}

// ConnectSession conecta uma sessão
// @Summary Conecta uma sessão WhatsApp
// @Description Conecta uma sessão WhatsApp específica ao servidor
// @Tags sessions
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "Sessão conectada com sucesso"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID}/connect [post]
func (h *Handler) ConnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "conectar sessão", func() (interface{}, error) {
		err := h.connectSessionUC.Execute(sessionID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "connected"}, nil
	}, "Sessão conectada com sucesso")
}

// GetQRCode obtém o QR code para uma sessão
// @Summary Obtém QR Code para autenticação
// @Description Gera um QR Code para autenticação da sessão WhatsApp
// @Tags sessions
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "QR Code gerado com sucesso"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID}/qr [get]
func (h *Handler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "obter QR code", func() (interface{}, error) {
		return h.getQRCodeUC.Execute(sessionID)
	}, "QR Code gerado com sucesso")
}

// DeleteSession deleta uma sessão
// @Summary Deleta uma sessão WhatsApp
// @Description Remove uma sessão WhatsApp e todos os seus dados
// @Tags sessions
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "Sessão deletada com sucesso"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID} [delete]
func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "deletar sessão", func() (interface{}, error) {
		err := h.deleteSessionUC.Execute(sessionID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "deleted"}, nil
	}, "Sessão deletada com sucesso")
}

// LogoutSession faz logout de uma sessão
// @Summary Faz logout de uma sessão WhatsApp
// @Description Desconecta uma sessão WhatsApp sem deletá-la
// @Tags sessions
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "Logout realizado com sucesso"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID}/logout [post]
func (h *Handler) LogoutSession(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "fazer logout da sessão", func() (interface{}, error) {
		err := h.logoutSessionUC.Execute(sessionID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "logged_out"}, nil
	}, "Logout realizado com sucesso")
}

// GetSessionInfo obtém informações de uma sessão
// @Summary Obtém informações de uma sessão
// @Description Retorna informações detalhadas sobre uma sessão WhatsApp
// @Tags sessions
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "Informações da sessão"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID} [get]
func (h *Handler) GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "obter informações da sessão", func() (interface{}, error) {
		return h.getSessionInfoUC.Execute(sessionID)
	}, "Informações da sessão obtidas com sucesso")
}

// PairPhone emparelha um telefone com uma sessão
// @Summary Emparelha telefone com sessão
// @Description Gera código de emparelhamento para conectar telefone à sessão
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.PairPhoneRequest true "Dados do telefone"
// @Success 200 {object} base.APIResponse "Código de emparelhamento gerado"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID}/pair [post]
func (h *Handler) PairPhone(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.PairPhoneRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "emparelhar telefone", func() (interface{}, error) {
		code, err := h.pairPhoneUC.Execute(sessionID, req.Phone)
		if err != nil {
			return nil, err
		}
		return map[string]string{"code": code}, nil
	}, "Código de emparelhamento gerado")
}

// SetProxy configura proxy para uma sessão
// @Summary Configura proxy para sessão
// @Description Define configurações de proxy para uma sessão WhatsApp
// @Tags sessions
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SetProxyRequest true "Configurações do proxy"
// @Success 200 {object} base.APIResponse "Proxy configurado com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID}/proxy [post]
func (h *Handler) SetProxy(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SetProxyRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"type": req.Type,
		"host": req.Host,
		"port": req.Port,
	}) {
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

	h.HandleUseCaseExecution(w, "configurar proxy", func() (interface{}, error) {
		err := h.setProxyUC.Execute(sessionID, proxyConfig)
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "proxy_configured"}, nil
	}, "Proxy configurado com sucesso")
}
