package message

import (
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/http/handlers/base"
	"wazmeow/internal/http/handlers/middleware"
)

// Handler contém os handlers para operações de mensagens refatorados
type Handler struct {
	*base.BaseHandler
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

// NewHandler cria uma nova instância dos handlers de mensagem refatorados
func NewHandler(
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
) *Handler {
	return &Handler{
		BaseHandler:          base.NewBaseHandler(),
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
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SendTextMessageRequest true "Dados da mensagem"
// @Success 200 {object} base.APIResponse "Mensagem enviada com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /message/{sessionID}/send/text [post]
func (h *Handler) SendTextMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendTextMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	// Validações específicas
	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
		"body":  req.Body,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar mensagem de texto", func() (interface{}, error) {
		return h.sendTextUseCase.Execute(sessionID, &req)
	}, "Mensagem enviada com sucesso")
}

// SendMediaMessage envia uma mensagem de mídia
// @Summary Envia mensagem de mídia
// @Description Envia uma mensagem de mídia (imagem, vídeo, áudio, documento) via WhatsApp
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SendMediaMessageRequest true "Dados da mídia"
// @Success 200 {object} base.APIResponse "Mídia enviada com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /message/{sessionID}/send/media [post]
func (h *Handler) SendMediaMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendMediaMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	// Validações específicas
	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone":      req.Phone,
		"media_data": req.MediaData,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar mídia", func() (interface{}, error) {
		return h.sendMediaUseCase.Execute(sessionID, &req)
	}, "Mídia enviada com sucesso")
}

// SendImageMessage envia uma mensagem de imagem
// @Summary Envia mensagem de imagem
// @Description Envia uma imagem via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SendImageMessageRequest true "Dados da imagem"
// @Success 200 {object} base.APIResponse "Imagem enviada com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /message/{sessionID}/send/image [post]
func (h *Handler) SendImageMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendImageMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
		"image": req.Image,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar imagem", func() (interface{}, error) {
		return h.sendImageUseCase.Execute(sessionID, &req)
	}, "Imagem enviada com sucesso")
}

// SendButtonsMessage envia uma mensagem com botões
// @Summary Envia mensagem com botões interativos
// @Description Envia uma mensagem com botões interativos via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SendButtonsMessageRequest true "Dados da mensagem com botões"
// @Success 200 {object} base.APIResponse "Mensagem com botões enviada com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /message/{sessionID}/send/buttons [post]
func (h *Handler) SendButtonsMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendButtonsMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	// Validações específicas para botões
	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone":   req.Phone,
		"title":   req.Title,
		"buttons": req.Buttons,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	// Validação específica de botões
	if err := h.GetValidator().ValidateSliceLength("buttons", req.Buttons, 1, 3); err != nil {
		h.SendError(w, err, base.GetHTTPStatus(err))
		return
	}

	h.HandleUseCaseExecution(w, "enviar botões", func() (interface{}, error) {
		return h.sendButtonsUseCase.Execute(sessionID, &req)
	}, "Botões enviados com sucesso")
}

// SendPollMessage envia uma mensagem de enquete
// @Summary Envia mensagem de enquete
// @Description Envia uma enquete via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SendPollMessageRequest true "Dados da enquete"
// @Success 200 {object} base.APIResponse "Enquete enviada com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /message/{sessionID}/send/poll [post]
func (h *Handler) SendPollMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendPollMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone":   req.Phone,
		"header":  req.Header,
		"options": req.Options,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	// Validação específica de enquete
	if err := h.GetValidator().ValidateSliceLength("options", req.Options, 2, 12); err != nil {
		h.SendError(w, err, base.GetHTTPStatus(err))
		return
	}

	h.HandleUseCaseExecution(w, "enviar enquete", func() (interface{}, error) {
		return h.sendPollUseCase.Execute(sessionID, &req)
	}, "Enquete enviada com sucesso")
}

// DeleteMessage deleta uma mensagem
// @Summary Deleta uma mensagem enviada
// @Description Remove uma mensagem já enviada via WhatsApp
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.DeleteMessageRequest true "Dados da mensagem a deletar"
// @Success 200 {object} base.APIResponse "Mensagem deletada com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /message/{sessionID}/delete [post]
func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.DeleteMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
		"id":    req.ID,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "deletar mensagem", func() (interface{}, error) {
		return h.deleteMessageUseCase.Execute(sessionID, &req)
	}, "Mensagem deletada com sucesso")
}

// SendAudioMessage envia uma mensagem de áudio
// @Summary Envia mensagem de áudio
// @Description Envia uma mensagem de áudio via WhatsApp para um número específico
// @Tags messages
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SendAudioMessageRequest true "Dados do áudio"
// @Success 200 {object} base.APIResponse "Áudio enviado com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /message/{sessionID}/send/audio [post]
func (h *Handler) SendAudioMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendAudioMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
		"audio": req.Audio,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar áudio", func() (interface{}, error) {
		return h.sendAudioUseCase.Execute(sessionID, &req)
	}, "Áudio enviado com sucesso")
}

// SendDocumentMessage envia uma mensagem de documento
func (h *Handler) SendDocumentMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendDocumentMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone":    req.Phone,
		"document": req.Document,
		"filename": req.FileName,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar documento", func() (interface{}, error) {
		return h.sendDocumentUseCase.Execute(sessionID, &req)
	}, "Documento enviado com sucesso")
}

// SendVideoMessage envia uma mensagem de vídeo
func (h *Handler) SendVideoMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendVideoMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
		"video": req.Video,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar vídeo", func() (interface{}, error) {
		return h.sendVideoUseCase.Execute(sessionID, &req)
	}, "Vídeo enviado com sucesso")
}

// SendStickerMessage envia uma mensagem de sticker
func (h *Handler) SendStickerMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendStickerMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone":   req.Phone,
		"sticker": req.Sticker,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar sticker", func() (interface{}, error) {
		return h.sendStickerUseCase.Execute(sessionID, &req)
	}, "Sticker enviado com sucesso")
}

// SendLocationMessage envia uma mensagem de localização
func (h *Handler) SendLocationMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendLocationMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone":     req.Phone,
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar localização", func() (interface{}, error) {
		return h.sendLocationUseCase.Execute(sessionID, &req)
	}, "Localização enviada com sucesso")
}

// SendContactMessage envia uma mensagem de contato
func (h *Handler) SendContactMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendContactMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
		"name":  req.Name,
		"vcard": req.Vcard,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar contato", func() (interface{}, error) {
		return h.sendContactUseCase.Execute(sessionID, &req)
	}, "Contato enviado com sucesso")
}

// SendListMessage envia uma mensagem de lista
func (h *Handler) SendListMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendListMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone":       req.Phone,
		"button_text": req.ButtonText,
		"desc":        req.Desc,
		"top_text":    req.TopText,
		"sections":    req.Sections,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "enviar lista", func() (interface{}, error) {
		return h.sendListUseCase.Execute(sessionID, &req)
	}, "Lista enviada com sucesso")
}

// SendEditMessage edita uma mensagem
func (h *Handler) SendEditMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.SendEditMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
		"body":  req.Body,
		"id":    req.ID,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "editar mensagem", func() (interface{}, error) {
		return h.sendEditUseCase.Execute(sessionID, &req)
	}, "Mensagem editada com sucesso")
}

// ReactMessage reage a uma mensagem
func (h *Handler) ReactMessage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.ReactMessageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	if !h.ValidateRequiredOrError(w, map[string]interface{}{
		"phone": req.Phone,
		"body":  req.Body,
		"id":    req.ID,
	}) {
		return
	}

	if !h.ValidatePhoneOrError(w, req.Phone) {
		return
	}

	h.HandleUseCaseExecution(w, "reagir à mensagem", func() (interface{}, error) {
		return h.reactUseCase.Execute(sessionID, &req)
	}, "Reação enviada com sucesso")
}
