package handlers

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// MessageHandlers contém os handlers para operações de mensagens
type MessageHandlers struct {
	sendTextUseCase  *usecase.SendTextMessageUseCase
	sendMediaUseCase *usecase.SendMediaMessageUseCase
}

// NewMessageHandlers cria uma nova instância dos handlers de mensagem
func NewMessageHandlers(
	sendTextUseCase *usecase.SendTextMessageUseCase,
	sendMediaUseCase *usecase.SendMediaMessageUseCase,
) *MessageHandlers {
	return &MessageHandlers{
		sendTextUseCase:  sendTextUseCase,
		sendMediaUseCase: sendMediaUseCase,
	}
}

// SendTextMessage envia uma mensagem de texto
func (h *MessageHandlers) SendTextMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendTextMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Body == "" {
		logger.Error("Campo 'body' é obrigatório")
		http.Error(w, "Campo 'body' é obrigatório", http.StatusBadRequest)
		return
	}

	// Executar use case
	response, err := h.sendTextUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar mensagem de texto: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder com sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Mensagem enviada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Mensagem de texto enviada com sucesso - Session: %s, Phone: %s, ID: %s", 
		sessionID, req.Phone, response.ID)
}

// SendMediaMessage envia uma mensagem de mídia
func (h *MessageHandlers) SendMediaMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendMediaMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.MediaData == "" {
		logger.Error("Campo 'media_data' é obrigatório")
		http.Error(w, "Campo 'media_data' é obrigatório", http.StatusBadRequest)
		return
	}

	// Executar use case
	response, err := h.sendMediaUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar mídia: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder com sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Mídia enviada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Mídia enviada com sucesso - Session: %s, Phone: %s, ID: %s", 
		sessionID, req.Phone, response.ID)
}
