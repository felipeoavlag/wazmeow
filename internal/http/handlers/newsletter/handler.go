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
func (h *Handler) ListNewsletter(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "listar newsletters", func() (interface{}, error) {
		return h.listNewsletterUseCase.Execute(sessionID)
	}, "Newsletters listadas com sucesso")
}