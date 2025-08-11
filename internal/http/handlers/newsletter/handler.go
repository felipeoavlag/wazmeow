package newsletter

import (
	"net/http"

	"wazmeow/internal/application/usecase"
	"wazmeow/internal/http/handlers/base"
	"wazmeow/internal/http/handlers/middleware"
)

// Handler contém os handlers para operações de newsletter refatorados
type Handler struct {
	*base.BaseHandler
	listNewsletterUseCase *usecase.ListNewsletterUseCase
}

// NewHandler cria uma nova instância dos handlers de newsletter refatorados
func NewHandler(
	listNewsletterUseCase *usecase.ListNewsletterUseCase,
) *Handler {
	return &Handler{
		BaseHandler:           base.NewBaseHandler(),
		listNewsletterUseCase: listNewsletterUseCase,
	}
}

// ListNewsletter lista newsletters subscritas
// @Summary Lista newsletters subscritas
// @Description Retorna lista de todas as newsletters que a sessão está inscrita
// @Tags newsletters
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "Lista de newsletters"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /newsletter/{sessionID}/list [get]
func (h *Handler) ListNewsletter(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "listar newsletters", func() (interface{}, error) {
		return h.listNewsletterUseCase.Execute(sessionID)
	}, "Newsletters listadas com sucesso")
}
