package handlers

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/application/usecase"
	"wazmeow/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// NewsletterHandlers contém os handlers para operações de newsletter
type NewsletterHandlers struct {
	listNewsletterUseCase *usecase.ListNewsletterUseCase
}

// NewNewsletterHandlers cria uma nova instância dos handlers de newsletter
func NewNewsletterHandlers(
	listNewsletterUseCase *usecase.ListNewsletterUseCase,
) *NewsletterHandlers {
	return &NewsletterHandlers{
		listNewsletterUseCase: listNewsletterUseCase,
	}
}

// ListNewsletter lista newsletters subscritas
// @Summary Lista newsletters subscritas
// @Description Retorna a lista de newsletters/canais WhatsApp que o usuário está inscrito
// @Tags newsletters
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Success 200 {object} map[string]interface{} "Newsletters listadas com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /newsletter/{sessionID}/list [get]
func (h *NewsletterHandlers) ListNewsletter(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.listNewsletterUseCase.Execute(sessionID)
	if err != nil {
		logger.Error("Erro ao listar newsletters: %v", err)
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

	logger.Info("Newsletters listadas com sucesso - Session: %s", sessionID)
}
