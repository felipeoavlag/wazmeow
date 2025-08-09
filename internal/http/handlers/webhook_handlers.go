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
	setWebhookUC    *usecase.SetWebhookUseCase
	getWebhookUC    *usecase.GetWebhookUseCase
	updateWebhookUC *usecase.UpdateWebhookUseCase
	deleteWebhookUC *usecase.DeleteWebhookUseCase
	eventFilter     *webhook.EventFilter
}

// NewWebhookHandler cria uma nova instância do handler de webhooks
func NewWebhookHandler(
	setWebhookUC *usecase.SetWebhookUseCase,
	getWebhookUC *usecase.GetWebhookUseCase,
	updateWebhookUC *usecase.UpdateWebhookUseCase,
	deleteWebhookUC *usecase.DeleteWebhookUseCase,
) *WebhookHandler {
	return &WebhookHandler{
		setWebhookUC:    setWebhookUC,
		getWebhookUC:    getWebhookUC,
		updateWebhookUC: updateWebhookUC,
		deleteWebhookUC: deleteWebhookUC,
		eventFilter:     webhook.NewEventFilter(),
	}
}

// SetWebhook configura webhook para uma sessão
// @Summary Configurar webhook
// @Description Configura URL e eventos de webhook para uma sessão específica
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body requests.SetWebhookRequest true "Dados do webhook"
// @Success 200 {object} responses.WebhookResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /sessions/{sessionId}/webhook [post]
func (wh *WebhookHandler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
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

	render.JSON(w, r, response)
}

// GetWebhook obtém configuração de webhook de uma sessão
// @Summary Obter webhook
// @Description Obtém a configuração atual de webhook de uma sessão
// @Tags Webhooks
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} responses.WebhookResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /sessions/{sessionId}/webhook [get]
func (wh *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
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

	render.JSON(w, r, response)
}

// UpdateWebhook atualiza configuração de webhook de uma sessão
// @Summary Atualizar webhook
// @Description Atualiza a configuração de webhook de uma sessão (ativar/desativar)
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Param request body requests.UpdateWebhookRequest true "Dados de atualização"
// @Success 200 {object} responses.WebhookResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /sessions/{sessionId}/webhook [put]
func (wh *WebhookHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Session ID é obrigatório"})
		return
	}

	var req requests.UpdateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "JSON inválido: " + err.Error()})
		return
	}

	// Validar eventos se fornecidos
	if req.Active && len(req.Events) > 0 {
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

	response, err := wh.updateWebhookUC.Execute(sessionID, &req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, response)
}

// DeleteWebhook remove configuração de webhook de uma sessão
// @Summary Remover webhook
// @Description Remove a configuração de webhook de uma sessão
// @Tags Webhooks
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} responses.WebhookDeleteResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /sessions/{sessionId}/webhook [delete]
func (wh *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Session ID é obrigatório"})
		return
	}

	response, err := wh.deleteWebhookUC.Execute(sessionID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, response)
}

// GetSupportedEvents retorna lista de eventos suportados
// @Summary Listar eventos suportados
// @Description Retorna a lista de todos os eventos suportados pelo sistema de webhooks
// @Tags Webhooks
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
		},
		"examples": map[string]interface{}{
			"all_events":        "*",
			"only_messages":     "messages",
			"messages_and_calls": []string{"messages", "calls"},
			"specific_events":   []string{"message", "receipt", "connected"},
		},
	}

	render.JSON(w, r, response)
}

// TestWebhook testa a conectividade de um webhook
// @Summary Testar webhook
// @Description Envia um webhook de teste para verificar conectividade
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param sessionId path string true "ID da sessão"
// @Success 200 {object} map[string]string
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /sessions/{sessionId}/webhook/test [post]
func (wh *WebhookHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionId")
	if sessionID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Session ID é obrigatório"})
		return
	}

	// Obter configuração atual do webhook
	webhookResponse, err := wh.getWebhookUC.Execute(sessionID)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "Webhook não configurado: " + err.Error()})
		return
	}

	if webhookResponse.Webhook == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Webhook URL não configurada"})
		return
	}

	// TODO: Implementar teste real do webhook quando o dispatcher estiver disponível
	// Por enquanto, apenas validar a URL
	render.JSON(w, r, map[string]string{
		"message": "Teste de webhook será implementado quando o sistema estiver completo",
		"url":     webhookResponse.Webhook,
		"status":  "pending",
	})
}
