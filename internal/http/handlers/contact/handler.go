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
// @Summary Obtém informações do usuário
// @Description Retorna informações detalhadas sobre um usuário WhatsApp
// @Tags contacts
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.GetUserInfoRequest true "Dados do usuário"
// @Success 200 {object} base.APIResponse "Informações do usuário"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /contact/{sessionID}/info [post]
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
// @Summary Verifica se usuário existe no WhatsApp
// @Description Verifica se um número de telefone está registrado no WhatsApp
// @Tags contacts
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.CheckUserRequest true "Número do usuário"
// @Success 200 {object} base.APIResponse "Usuário verificado com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /contact/{sessionID}/check [post]
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
// @Summary Obtém avatar do usuário
// @Description Baixa a foto de perfil de um usuário WhatsApp
// @Tags contacts
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.GetAvatarRequest true "Dados do usuário"
// @Success 200 {object} base.APIResponse "Avatar obtido com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /contact/{sessionID}/avatar [post]
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
// @Summary Obtém lista de contatos
// @Description Retorna lista de todos os contatos da sessão WhatsApp
// @Tags contacts
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "Lista de contatos"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /contact/{sessionID}/list [get]
func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "obter contatos", func() (interface{}, error) {
		return h.getContactsUseCase.Execute(sessionID)
	}, "Contatos obtidos com sucesso")
}
