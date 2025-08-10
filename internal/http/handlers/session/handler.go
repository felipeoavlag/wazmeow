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
func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	h.HandleUseCaseExecution(w, "listar sessões", func() (interface{}, error) {
		return h.listSessionsUC.Execute()
	}, "Sessões listadas com sucesso")
}

// ConnectSession conecta uma sessão
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