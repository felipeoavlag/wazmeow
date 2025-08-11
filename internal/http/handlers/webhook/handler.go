package webhook

import (
	"net/http"
	"time"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/http/handlers/base"
	"wazmeow/internal/http/handlers/middleware"
	"wazmeow/internal/infra/webhook"
	"wazmeow/pkg/logger"
)

// Handler contém os handlers para operações de webhook
type Handler struct {
	*base.BaseHandler
	setWebhookUC *usecase.SetWebhookUseCase
	getWebhookUC *usecase.GetWebhookUseCase
	eventFilter  *webhook.EventFilter
}

// NewHandler cria uma nova instância dos handlers de webhook
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
// @Summary Set - Configura webhook para sessão
// @Description Define ou remove URL de webhook para receber eventos de uma sessão
// @Tags webhooks
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SetWebhookRequest true "Configurações do webhook"
// @Success 200 {object} base.APIResponse "Webhook configurado com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID}/webhook/set [post]
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
// @Summary Find - Obtém configuração de webhook
// @Description Retorna a configuração atual de webhook de uma sessão
// @Tags webhooks
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} base.APIResponse "Configuração do webhook"
// @Failure 400 {object} base.APIResponse "Sessão não encontrada"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /sessions/{sessionID}/webhook [get]
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
// @Summary Events - Lista eventos suportados para webhooks
// @Description Retorna todos os tipos de eventos que podem ser configurados em webhooks
// @Tags webhooks
// @Produce json
// @Success 200 {object} base.APIResponse "Lista de eventos suportados"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /webhook/events [get]
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

// ReceiveWebhookEvent recebe eventos do webhook
// @Summary Receive - Recebe eventos do webhook
// @Description Endpoint para receber eventos do WhatsApp via webhook
// @Tags webhooks
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param event body object true "Dados do evento"
// @Success 200 {object} base.APIResponse "Evento processado com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /webhook/{sessionID} [post]
func (h *Handler) ReceiveWebhookEvent(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	// Decodificar o payload do evento
	var eventPayload map[string]interface{}
	if !h.DecodeJSONOrError(w, r, &eventPayload) {
		return
	}

	// Logar o evento recebido
	logger.Info("Evento recebido via webhook - Session: %s, Event: %+v", sessionID, eventPayload)

	// Responder com sucesso
	h.SendSuccess(w, map[string]interface{}{
		"sessionID": sessionID,
		"processed": true,
		"timestamp": time.Now().Unix(),
	}, "Evento processado com sucesso")
}
