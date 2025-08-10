package webhook

import (
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/http/handlers/base"
	"wazmeow/internal/http/handlers/middleware"
	"wazmeow/internal/infra/webhook"
)

// Handler contém os handlers para operações de webhook refatorados
type Handler struct {
	*base.BaseHandler
	setWebhookUC *usecase.SetWebhookUseCase
	getWebhookUC *usecase.GetWebhookUseCase
	eventFilter  *webhook.EventFilter
}

// NewHandler cria uma nova instância dos handlers de webhook refatorados
func NewHandler(
	setWebhookUC *usecase.SetWebhookUseCase,
	getWebhookUC *usecase.GetWebhookUseCase,
) *Handler {
	return &Handler{
		BaseHandler:  base.NewBaseHandler(),
		setWebhookUC: setWebhookUC,
		getWebhookUC: getWebhookUC,
		eventFilter:  webhook.NewEventFilter(),
	}
}

// SetWebhook configura ou remove webhook para uma sessão
func (h *Handler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SetWebhookRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	// Validar eventos se fornecidos
	if len(req.Events) > 0 {
		eventsStr := ""
		for i, event := range req.Events {
			if i > 0 {
				eventsStr += ","
			}
			eventsStr += event
		}
		if err := h.eventFilter.ValidateEvents(eventsStr); err != nil {
			h.SendBadRequest(w, "Eventos inválidos: "+err.Error())
			return
		}
	}

	// Determinar mensagem baseada na ação
	message := "Webhook configurado com sucesso"
	if req.WebhookURL == "" || (req.Enabled != nil && !*req.Enabled) {
		message = "Webhook removido com sucesso"
	}

	h.HandleUseCaseExecution(w, "configurar webhook", func() (interface{}, error) {
		return h.setWebhookUC.Execute(sessionID, &req)
	}, message)
}

// FindWebhook obtém configuração de webhook de uma sessão
func (h *Handler) FindWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	h.HandleUseCaseExecution(w, "obter webhook", func() (interface{}, error) {
		return h.getWebhookUC.Execute(sessionID)
	}, "Webhook obtido com sucesso")
}

// GetSupportedEvents retorna lista de eventos suportados
func (h *Handler) GetSupportedEvents(w http.ResponseWriter, r *http.Request) {
	events := h.eventFilter.GetSupportedEvents()
	groups := h.eventFilter.GetEventGroups()

	response := map[string]interface{}{
		"events": events,
		"groups": groups,
		"wildcards": []string{
			"*",
			"all",
			"messages.*",
			"calls.*",
			"groups.*",
			"presence.*",
			"connection.*",
			"media.*",
			"newsletters.*",
			"contacts.*",
		},
		"examples": map[string]interface{}{
			"all_events":         "*",
			"only_messages":      "messages",
			"messages_and_calls": []string{"messages", "calls"},
			"specific_events":    []string{"message", "receipt", "connected"},
			"contacts_events":    "contacts",
		},
	}

	h.SendSuccess(w, response, "Eventos suportados obtidos com sucesso")
}