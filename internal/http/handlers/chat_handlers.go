package handlers

import (
	"encoding/json"
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// ChatHandlers contém os handlers para operações de chat
type ChatHandlers struct {
	sendPresenceUseCase     *usecase.SendPresenceUseCase
	chatPresenceUseCase     *usecase.ChatPresenceUseCase
	markReadUseCase         *usecase.MarkReadUseCase
	downloadImageUseCase    *usecase.DownloadImageUseCase
	downloadVideoUseCase    *usecase.DownloadVideoUseCase
	downloadAudioUseCase    *usecase.DownloadAudioUseCase
	downloadDocumentUseCase *usecase.DownloadDocumentUseCase
}

// NewChatHandlers cria uma nova instância dos handlers de chat
func NewChatHandlers(
	sendPresenceUseCase *usecase.SendPresenceUseCase,
	chatPresenceUseCase *usecase.ChatPresenceUseCase,
	markReadUseCase *usecase.MarkReadUseCase,
	downloadImageUseCase *usecase.DownloadImageUseCase,
	downloadVideoUseCase *usecase.DownloadVideoUseCase,
	downloadAudioUseCase *usecase.DownloadAudioUseCase,
	downloadDocumentUseCase *usecase.DownloadDocumentUseCase,
) *ChatHandlers {
	return &ChatHandlers{
		sendPresenceUseCase:     sendPresenceUseCase,
		chatPresenceUseCase:     chatPresenceUseCase,
		markReadUseCase:         markReadUseCase,
		downloadImageUseCase:    downloadImageUseCase,
		downloadVideoUseCase:    downloadVideoUseCase,
		downloadAudioUseCase:    downloadAudioUseCase,
		downloadDocumentUseCase: downloadDocumentUseCase,
	}
}

// SendPresence define presença do usuário
func (h *ChatHandlers) SendPresence(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendPresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.sendPresenceUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao definir presença: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Presença definida com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Presença definida com sucesso - Session: %s", sessionID)
}

// ChatPresence define presença no chat
func (h *ChatHandlers) ChatPresence(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.ChatPresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.chatPresenceUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao definir presença no chat: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Presença no chat definida com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Presença no chat definida com sucesso - Session: %s", sessionID)
}

// MarkRead marca mensagens como lidas
func (h *ChatHandlers) MarkRead(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.MarkReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.markReadUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao marcar mensagens como lidas: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Mensagens marcadas como lidas com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Mensagens marcadas como lidas com sucesso - Session: %s", sessionID)
}

// DownloadImage faz download de imagem
func (h *ChatHandlers) DownloadImage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.DownloadImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.downloadImageUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao fazer download da imagem: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Download da imagem concluído com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Download da imagem concluído com sucesso - Session: %s", sessionID)
}

// DownloadVideo faz download de vídeo
func (h *ChatHandlers) DownloadVideo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.DownloadVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.downloadVideoUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao fazer download do vídeo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Download do vídeo concluído com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Download do vídeo concluído com sucesso - Session: %s", sessionID)
}

// DownloadAudio faz download de áudio
func (h *ChatHandlers) DownloadAudio(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.DownloadAudioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.downloadAudioUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao fazer download do áudio: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Download do áudio concluído com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Download do áudio concluído com sucesso - Session: %s", sessionID)
}

// DownloadDocument faz download de documento
func (h *ChatHandlers) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.DownloadDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	response, err := h.downloadDocumentUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao fazer download do documento: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Download do documento concluído com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Download do documento concluído com sucesso - Session: %s", sessionID)
}
