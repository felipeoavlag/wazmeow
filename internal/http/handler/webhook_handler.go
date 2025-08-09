package handler

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/infra/webhook"
	"wazmeow/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// WebhookHandler representa o handler para operações de webhook
type WebhookHandler struct {
	setWebhookUseCase    *usecase.SetWebhookUseCase
	getWebhookUseCase    *usecase.GetWebhookUseCase
	updateWebhookUseCase *usecase.UpdateWebhookUseCase
	deleteWebhookUseCase *usecase.DeleteWebhookUseCase
	eventFilter          *webhook.EventFilter
}

// NewWebhookHandler cria uma nova instância do handler
func NewWebhookHandler(
	setWebhookUseCase *usecase.SetWebhookUseCase,
	getWebhookUseCase *usecase.GetWebhookUseCase,
	updateWebhookUseCase *usecase.UpdateWebhookUseCase,
	deleteWebhookUseCase *usecase.DeleteWebhookUseCase,
) *WebhookHandler {
	return &WebhookHandler{
		setWebhookUseCase:    setWebhookUseCase,
		getWebhookUseCase:    getWebhookUseCase,
		updateWebhookUseCase: updateWebhookUseCase,
		deleteWebhookUseCase: deleteWebhookUseCase,
		eventFilter:          webhook.NewEventFilter(),
	}
}

// SetWebhook define o webhook para uma sessão
func (h *WebhookHandler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SetWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.setWebhookUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao definir webhook: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Webhook definido com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Webhook definido com sucesso - Session: %s", sessionID)
}

// GetWebhook obtém o webhook de uma sessão
func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.getWebhookUseCase.Execute(sessionID)
	if err != nil {
		logger.Error("Erro ao obter webhook: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Webhook obtido com sucesso - Session: %s", sessionID)
}

// UpdateWebhook atualiza o webhook de uma sessão
func (h *WebhookHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.UpdateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.updateWebhookUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao atualizar webhook: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Webhook atualizado com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Webhook atualizado com sucesso - Session: %s", sessionID)
}

// DeleteWebhook remove o webhook de uma sessão
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.deleteWebhookUseCase.Execute(sessionID)
	if err != nil {
		logger.Error("Erro ao remover webhook: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Webhook removido com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Webhook removido com sucesso - Session: %s", sessionID)
}

// GetSupportedEvents retorna lista de eventos suportados
func (h *WebhookHandler) GetSupportedEvents(w http.ResponseWriter, r *http.Request) {
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
		},
		"examples": map[string]interface{}{
			"all_events":         "*",
			"only_messages":      "messages",
			"messages_and_calls": []string{"messages", "calls"},
			"specific_events":    []string{"message", "receipt", "connected"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Lista de eventos suportados retornada com sucesso")
}

// TestWebhook testa a conectividade de um webhook
func (h *WebhookHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	// Obter configuração atual do webhook
	webhookResponse, err := h.getWebhookUseCase.Execute(sessionID)
	if err != nil {
		logger.Error("Erro ao obter webhook para teste: %v", err)
		http.Error(w, "Webhook não configurado: "+err.Error(), http.StatusNotFound)
		return
	}

	if webhookResponse.Webhook == "" {
		logger.Error("Webhook URL não configurada para sessão: %s", sessionID)
		http.Error(w, "Webhook URL não configurada", http.StatusBadRequest)
		return
	}

	// TODO: Implementar teste real do webhook quando o dispatcher estiver disponível
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Teste de webhook será implementado quando o sistema estiver completo",
		"data": map[string]interface{}{
			"url":    webhookResponse.Webhook,
			"status": "pending",
		},
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Teste de webhook solicitado - Session: %s, URL: %s", sessionID, webhookResponse.Webhook)
}
