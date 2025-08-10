package contact

import (
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/http/handlers/base"
	"wazmeow/internal/http/handlers/middleware"
)

// Handler contém os handlers para operações de contato refatorados
type Handler struct {
	*base.BaseHandler
	getUserInfoUseCase *usecase.GetUserInfoUseCase
	checkUserUseCase   *usecase.CheckUserUseCase
	getAvatarUseCase   *usecase.GetAvatarUseCase
	getContactsUseCase *usecase.GetContactsUseCase
}

// NewHandler cria uma nova instância dos handlers de contato refatorados
func NewHandler(
	getUserInfoUseCase *usecase.GetUserInfoUseCase,
	checkUserUseCase *usecase.CheckUserUseCase,
	getAvatarUseCase *usecase.GetAvatarUseCase,
	getContactsUseCase *usecase.GetContactsUseCase,
) *Handler {
	return &Handler{
		BaseHandler:        base.NewBaseHandler(),
		getUserInfoUseCase: getUserInfoUseCase,
		checkUserUseCase:   checkUserUseCase,
		getAvatarUseCase:   getAvatarUseCase,
		getContactsUseCase: getContactsUseCase,
	}
}

// GetUserInfo obtém informações do usuário
func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.GetUserInfoRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
	}) {
		return
	}

	// Validar cada telefone no slice
	for _, phone := range req.Phone {
		if !h.ValidatePhoneOrError(w, phone) {
			return
		}
	}

	h.HandleUseCaseExecution(w, "obter informações do usuário", func() (interface{}, error) {
		return h.getUserInfoUseCase.Execute(sessionID, &req)
	}, "Informações do usuário obtidas com sucesso")
}

// CheckUser verifica se usuário existe
func (h *Handler) CheckUser(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.CheckUserRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
	}) {
		return
	}

	// Validar cada telefone no slice
	for _, phone := range req.Phone {
		if !h.ValidatePhoneOrError(w, phone) {
			return
		}
	}

	h.HandleUseCaseExecution(w, "verificar usuário", func() (interface{}, error) {
		return h.checkUserUseCase.Execute(sessionID, &req)
	}, "Usuário verificado com sucesso")
}

// GetAvatar obtém avatar do usuário
func (h *Handler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.GetAvatarRequest
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

	h.HandleUseCaseExecution(w, "obter avatar", func() (interface{}, error) {
		return h.getAvatarUseCase.Execute(sessionID, &req)
	}, "Avatar obtido com sucesso")
}

// GetContacts obtém lista de contatos
func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "obter contatos", func() (interface{}, error) {
		return h.getContactsUseCase.Execute(sessionID)
	}, "Contatos obtidos com sucesso")
}