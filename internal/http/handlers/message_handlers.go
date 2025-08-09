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
	// Envio de mensagens básicas
	sendTextUseCase  *usecase.SendTextMessageUseCase
	sendMediaUseCase *usecase.SendMediaMessageUseCase

	// Envio de mensagens específicas
	sendImageUseCase    *usecase.SendImageMessageUseCase
	sendAudioUseCase    *usecase.SendAudioMessageUseCase
	sendDocumentUseCase *usecase.SendDocumentMessageUseCase
	sendVideoUseCase    *usecase.SendVideoMessageUseCase
	sendStickerUseCase  *usecase.SendStickerMessageUseCase
	sendLocationUseCase *usecase.SendLocationMessageUseCase
	sendContactUseCase  *usecase.SendContactMessageUseCase
	sendButtonsUseCase  *usecase.SendButtonsMessageUseCase
	sendListUseCase     *usecase.SendListMessageUseCase
	sendPollUseCase     *usecase.SendPollMessageUseCase

	// Operações de mensagem
	sendEditUseCase      *usecase.SendEditMessageUseCase
	deleteMessageUseCase *usecase.DeleteMessageUseCase
	reactUseCase         *usecase.ReactMessageUseCase
}

// NewMessageHandlers cria uma nova instância dos handlers de mensagem
func NewMessageHandlers(
	sendTextUseCase *usecase.SendTextMessageUseCase,
	sendMediaUseCase *usecase.SendMediaMessageUseCase,
	sendImageUseCase *usecase.SendImageMessageUseCase,
	sendAudioUseCase *usecase.SendAudioMessageUseCase,
	sendDocumentUseCase *usecase.SendDocumentMessageUseCase,
	sendVideoUseCase *usecase.SendVideoMessageUseCase,
	sendStickerUseCase *usecase.SendStickerMessageUseCase,
	sendLocationUseCase *usecase.SendLocationMessageUseCase,
	sendContactUseCase *usecase.SendContactMessageUseCase,
	sendButtonsUseCase *usecase.SendButtonsMessageUseCase,
	sendListUseCase *usecase.SendListMessageUseCase,
	sendPollUseCase *usecase.SendPollMessageUseCase,
	sendEditUseCase *usecase.SendEditMessageUseCase,
	deleteMessageUseCase *usecase.DeleteMessageUseCase,
	reactUseCase *usecase.ReactMessageUseCase,
) *MessageHandlers {
	return &MessageHandlers{
		sendTextUseCase:      sendTextUseCase,
		sendMediaUseCase:     sendMediaUseCase,
		sendImageUseCase:     sendImageUseCase,
		sendAudioUseCase:     sendAudioUseCase,
		sendDocumentUseCase:  sendDocumentUseCase,
		sendVideoUseCase:     sendVideoUseCase,
		sendStickerUseCase:   sendStickerUseCase,
		sendLocationUseCase:  sendLocationUseCase,
		sendContactUseCase:   sendContactUseCase,
		sendButtonsUseCase:   sendButtonsUseCase,
		sendListUseCase:      sendListUseCase,
		sendPollUseCase:      sendPollUseCase,
		sendEditUseCase:      sendEditUseCase,
		deleteMessageUseCase: deleteMessageUseCase,
		reactUseCase:         reactUseCase,
	}
}

// SendTextMessage envia uma mensagem de texto
// @Summary Envia mensagem de texto
// @Description Envia uma mensagem de texto via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendTextMessageRequest true "Dados da mensagem"
// @Success 200 {object} map[string]interface{} "Mensagem enviada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/text [post]
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
// @Summary Envia mensagem de mídia
// @Description Envia uma mensagem de mídia (imagem, vídeo, áudio, documento) via WhatsApp
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendMediaMessageRequest true "Dados da mídia"
// @Success 200 {object} map[string]interface{} "Mídia enviada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/media [post]
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

// SendImageMessage envia uma mensagem de imagem
// @Summary Envia mensagem de imagem
// @Description Envia uma imagem via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendImageMessageRequest true "Dados da imagem"
// @Success 200 {object} map[string]interface{} "Imagem enviada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/image [post]
func (h *MessageHandlers) SendImageMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendImageMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Image == "" {
		logger.Error("Campo 'image' é obrigatório")
		http.Error(w, "Campo 'image' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendImageUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar imagem: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Imagem enviada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Imagem enviada com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendAudioMessage envia uma mensagem de áudio
// @Summary Envia mensagem de áudio
// @Description Envia um áudio via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendAudioMessageRequest true "Dados do áudio"
// @Success 200 {object} map[string]interface{} "Áudio enviado com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/audio [post]
func (h *MessageHandlers) SendAudioMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendAudioMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Audio == "" {
		logger.Error("Campo 'audio' é obrigatório")
		http.Error(w, "Campo 'audio' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendAudioUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar áudio: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Áudio enviado com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Áudio enviado com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendDocumentMessage envia uma mensagem de documento
// @Summary Envia mensagem de documento
// @Description Envia um documento via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendDocumentMessageRequest true "Dados do documento"
// @Success 200 {object} map[string]interface{} "Documento enviado com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/document [post]
func (h *MessageHandlers) SendDocumentMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendDocumentMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Document == "" {
		logger.Error("Campo 'document' é obrigatório")
		http.Error(w, "Campo 'document' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.FileName == "" {
		logger.Error("Campo 'filename' é obrigatório")
		http.Error(w, "Campo 'filename' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendDocumentUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar documento: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Documento enviado com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Documento enviado com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendVideoMessage envia uma mensagem de vídeo
// @Summary Envia mensagem de vídeo
// @Description Envia um vídeo via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendVideoMessageRequest true "Dados do vídeo"
// @Success 200 {object} map[string]interface{} "Vídeo enviado com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/video [post]
func (h *MessageHandlers) SendVideoMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendVideoMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Video == "" {
		logger.Error("Campo 'video' é obrigatório")
		http.Error(w, "Campo 'video' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendVideoUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar vídeo: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Vídeo enviado com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Vídeo enviado com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendStickerMessage envia uma mensagem de sticker
// @Summary Envia mensagem de sticker
// @Description Envia um sticker via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendStickerMessageRequest true "Dados do sticker"
// @Success 200 {object} map[string]interface{} "Sticker enviado com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/sticker [post]
func (h *MessageHandlers) SendStickerMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendStickerMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Sticker == "" {
		logger.Error("Campo 'sticker' é obrigatório")
		http.Error(w, "Campo 'sticker' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendStickerUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar sticker: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Sticker enviado com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Sticker enviado com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendLocationMessage envia uma mensagem de localização
// @Summary Envia mensagem de localização
// @Description Envia uma localização via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendLocationMessageRequest true "Dados da localização"
// @Success 200 {object} map[string]interface{} "Localização enviada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/location [post]
func (h *MessageHandlers) SendLocationMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendLocationMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Latitude == 0 {
		logger.Error("Campo 'latitude' é obrigatório")
		http.Error(w, "Campo 'latitude' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Longitude == 0 {
		logger.Error("Campo 'longitude' é obrigatório")
		http.Error(w, "Campo 'longitude' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendLocationUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar localização: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Localização enviada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Localização enviada com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendContactMessage envia uma mensagem de contato
// @Summary Envia mensagem de contato
// @Description Envia um contato via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendContactMessageRequest true "Dados do contato"
// @Success 200 {object} map[string]interface{} "Contato enviado com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/contact [post]
func (h *MessageHandlers) SendContactMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendContactMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		logger.Error("Campo 'name' é obrigatório")
		http.Error(w, "Campo 'name' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Vcard == "" {
		logger.Error("Campo 'vcard' é obrigatório")
		http.Error(w, "Campo 'vcard' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendContactUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar contato: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Contato enviado com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Contato enviado com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendButtonsMessage envia uma mensagem com botões
// @Summary Envia mensagem com botões interativos
// @Description Envia uma mensagem com botões interativos via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendButtonsMessageRequest true "Dados da mensagem com botões"
// @Success 200 {object} map[string]interface{} "Mensagem com botões enviada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/buttons [post]
func (h *MessageHandlers) SendButtonsMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendButtonsMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		logger.Error("Campo 'title' é obrigatório")
		http.Error(w, "Campo 'title' é obrigatório", http.StatusBadRequest)
		return
	}

	if len(req.Buttons) == 0 {
		logger.Error("Campo 'buttons' é obrigatório")
		http.Error(w, "Campo 'buttons' é obrigatório", http.StatusBadRequest)
		return
	}

	if len(req.Buttons) > 3 {
		logger.Error("Máximo de 3 botões permitidos")
		http.Error(w, "Máximo de 3 botões permitidos", http.StatusBadRequest)
		return
	}

	response, err := h.sendButtonsUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar botões: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Botões enviados com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Botões enviados com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendListMessage envia uma mensagem de lista
// @Summary Envia mensagem de lista interativa
// @Description Envia uma mensagem de lista interativa via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendListMessageRequest true "Dados da mensagem de lista"
// @Success 200 {object} map[string]interface{} "Lista enviada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/list [post]
func (h *MessageHandlers) SendListMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendListMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.ButtonText == "" {
		logger.Error("Campo 'button_text' é obrigatório")
		http.Error(w, "Campo 'button_text' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Desc == "" {
		logger.Error("Campo 'desc' é obrigatório")
		http.Error(w, "Campo 'desc' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.TopText == "" {
		logger.Error("Campo 'top_text' é obrigatório")
		http.Error(w, "Campo 'top_text' é obrigatório", http.StatusBadRequest)
		return
	}

	if len(req.Sections) == 0 {
		logger.Error("Campo 'sections' é obrigatório")
		http.Error(w, "Campo 'sections' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendListUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar lista: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Lista enviada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Lista enviada com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendPollMessage envia uma mensagem de enquete
// @Summary Envia mensagem de enquete
// @Description Envia uma enquete via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendPollMessageRequest true "Dados da enquete"
// @Success 200 {object} map[string]interface{} "Enquete enviada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/poll [post]
func (h *MessageHandlers) SendPollMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendPollMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.Header == "" {
		logger.Error("Campo 'header' é obrigatório")
		http.Error(w, "Campo 'header' é obrigatório", http.StatusBadRequest)
		return
	}

	if len(req.Options) < 2 {
		logger.Error("Pelo menos 2 opções são obrigatórias")
		http.Error(w, "Pelo menos 2 opções são obrigatórias", http.StatusBadRequest)
		return
	}

	response, err := h.sendPollUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao enviar enquete: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Enquete enviada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Enquete enviada com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// SendEditMessage edita uma mensagem
// @Summary Edita uma mensagem enviada
// @Description Edita o conteúdo de uma mensagem já enviada via WhatsApp
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.SendEditMessageRequest true "Dados da edição"
// @Success 200 {object} map[string]interface{} "Mensagem editada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/send/edit [post]
func (h *MessageHandlers) SendEditMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.SendEditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

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

	if req.ID == "" {
		logger.Error("Campo 'id' é obrigatório")
		http.Error(w, "Campo 'id' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.sendEditUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao editar mensagem: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Mensagem editada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Mensagem editada com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// DeleteMessage deleta uma mensagem
// @Summary Deleta uma mensagem enviada
// @Description Remove uma mensagem já enviada via WhatsApp
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.DeleteMessageRequest true "Dados da mensagem a deletar"
// @Success 200 {object} map[string]interface{} "Mensagem deletada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/delete [post]
func (h *MessageHandlers) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.DeleteMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		logger.Error("Campo 'phone' é obrigatório")
		http.Error(w, "Campo 'phone' é obrigatório", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		logger.Error("Campo 'id' é obrigatório")
		http.Error(w, "Campo 'id' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.deleteMessageUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao deletar mensagem: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Mensagem deletada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Mensagem deletada com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}

// ReactMessage reage a uma mensagem
// @Summary Reage a uma mensagem
// @Description Adiciona ou remove uma reação (emoji) a uma mensagem via WhatsApp
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID da sessão"
// @Param request body requests.ReactMessageRequest true "Dados da reação"
// @Success 200 {object} map[string]interface{} "Reação adicionada com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Failure 500 {string} string "Erro interno do servidor"
// @Router /message/{sessionID}/react [post]
func (h *MessageHandlers) ReactMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if sessionID == "" {
		logger.Error("Session ID não fornecido")
		http.Error(w, "Session ID é obrigatório", http.StatusBadRequest)
		return
	}

	var req requests.ReactMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar payload: %v", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

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

	if req.ID == "" {
		logger.Error("Campo 'id' é obrigatório")
		http.Error(w, "Campo 'id' é obrigatório", http.StatusBadRequest)
		return
	}

	response, err := h.reactUseCase.Execute(sessionID, &req)
	if err != nil {
		logger.Error("Erro ao reagir à mensagem: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Reação enviada com sucesso",
		"data":    response,
	}); err != nil {
		logger.Error("Erro ao codificar resposta: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	logger.Info("Reação enviada com sucesso - Session: %s, Phone: %s, ID: %s",
		sessionID, req.Phone, response.ID)
}
