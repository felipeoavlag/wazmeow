package handlers

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/infra/webhook"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// WebhookHandler gerencia as rotas relacionadas aos webhooks
type WebhookHandler struct {
	setWebhookUC *usecase.SetWebhookUseCase
	getWebhookUC *usecase.GetWebhookUseCase
	eventFilter  *webhook.EventFilter
}

// NewWebhookHandler cria uma nova instância do handler de webhooks
func NewWebhookHandler(
	setWebhookUC *usecase.SetWebhookUseCase,
	getWebhookUC *usecase.GetWebhookUseCase,
) *WebhookHandler {
	return &WebhookHandler{
		setWebhookUC: setWebhookUC,
		getWebhookUC: getWebhookUC,
		eventFilter:  webhook.NewEventFilter(),
	}
}

// SetWebhook configura ou remove webhook para uma sessão
// @Summary Configurar/Remover webhook
// @Description Configura URL e eventos de webhook para uma sessão específica. Para remover, envie URL vazia ou enabled=false
// @Tags webhooks
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body requests.SetWebhookRequest true "Dados do webhook"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /webhook/{sessionId}/set [post]
func (wh *WebhookHandler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Session ID é obrigatório"})
		return
	}

	var req requests.SetWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "JSON inválido: " + err.Error()})
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
		if err := wh.eventFilter.ValidateEvents(eventsStr); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Eventos inválidos: " + err.Error()})
			return
		}
	}

	response, err := wh.setWebhookUC.Execute(sessionID, &req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	// Determinar mensagem baseada na ação
	message := "Webhook configurado com sucesso"
	if req.WebhookURL == "" || (req.Enabled != nil && !*req.Enabled) {
		message = "Webhook removido com sucesso"
	}

	render.JSON(w, r, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    response,
	})
}

// FindWebhook obtém configuração de webhook de uma sessão
// @Summary Obter webhook
// @Description Obtém a configuração atual de webhook de uma sessão
// @Tags webhooks
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /webhook/{sessionId}/find [get]
func (wh *WebhookHandler) FindWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Session ID é obrigatório"})
		return
	}

	response, err := wh.getWebhookUC.Execute(sessionID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// GetSupportedEvents retorna lista de eventos suportados
// @Summary Listar eventos suportados
// @Description Retorna todos os eventos disponíveis para configuração de webhook
// @Tags webhooks
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /webhook/events [get]
func (wh *WebhookHandler) GetSupportedEvents(w http.ResponseWriter, r *http.Request) {
	events := wh.eventFilter.GetSupportedEvents()
	groups := wh.eventFilter.GetEventGroups()

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

	render.JSON(w, r, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}
