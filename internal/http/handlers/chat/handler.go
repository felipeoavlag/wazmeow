package chat

import (
	"net/http"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/usecase"
	"wazmeow/internal/http/handlers/base"
	"wazmeow/internal/http/handlers/middleware"
)

// Handler contém os handlers para operações de chat refatorados
type Handler struct {
	*base.BaseHandler
	sendPresenceUseCase     *usecase.SendPresenceUseCase
	chatPresenceUseCase     *usecase.ChatPresenceUseCase
	markReadUseCase         *usecase.MarkReadUseCase
	downloadImageUseCase    *usecase.DownloadImageUseCase
	downloadVideoUseCase    *usecase.DownloadVideoUseCase
	downloadAudioUseCase    *usecase.DownloadAudioUseCase
	downloadDocumentUseCase *usecase.DownloadDocumentUseCase
}

// NewHandler cria uma nova instância dos handlers de chat refatorados
func NewHandler(
	sendPresenceUseCase *usecase.SendPresenceUseCase,
	chatPresenceUseCase *usecase.ChatPresenceUseCase,
	markReadUseCase *usecase.MarkReadUseCase,
	downloadImageUseCase *usecase.DownloadImageUseCase,
	downloadVideoUseCase *usecase.DownloadVideoUseCase,
	downloadAudioUseCase *usecase.DownloadAudioUseCase,
	downloadDocumentUseCase *usecase.DownloadDocumentUseCase,
) *Handler {
	return &Handler{
		BaseHandler:             base.NewBaseHandler(),
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
// @Summary Define presença do usuário
// @Description Define o status de presença do usuário (disponível, ocupado, etc.)
// @Tags chats
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.SendPresenceRequest true "Dados da presença"
// @Success 200 {object} base.APIResponse "Presença definida com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /chat/{sessionID}/presence [post]
func (h *Handler) SendPresence(w http.ResponseWriter, r *http.Request) {
	// Extrai sessionID com validação automática
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	// Decodifica JSON com validação automática
	var req requests.SendPresenceRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	// Executa use case com tratamento automático de erros
	h.HandleUseCaseExecution(w, "definir presença", func() (interface{}, error) {
		return h.sendPresenceUseCase.Execute(sessionID, &req)
	}, "Presença definida com sucesso")
}

// ChatPresence define presença no chat
// @Summary Define presença em um chat específico
// @Description Define o status de presença em um chat específico (digitando, gravando áudio, etc.)
// @Tags chats
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.ChatPresenceRequest true "Dados da presença no chat"
// @Success 200 {object} base.APIResponse "Presença no chat definida com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /chat/{sessionID}/chatpresence [post]
func (h *Handler) ChatPresence(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.ChatPresenceRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	h.HandleUseCaseExecution(w, "definir presença no chat", func() (interface{}, error) {
		return h.chatPresenceUseCase.Execute(sessionID, &req)
	}, "Presença no chat definida com sucesso")
}

// MarkRead marca mensagens como lidas
// @Summary Marca mensagens como lidas
// @Description Marca uma ou mais mensagens como lidas em um chat específico
// @Tags chats
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.MarkReadRequest true "Dados das mensagens para marcar como lidas"
// @Success 200 {object} base.APIResponse "Mensagens marcadas como lidas com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /chat/{sessionID}/markread [post]
func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.MarkReadRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	h.HandleUseCaseExecution(w, "marcar mensagens como lidas", func() (interface{}, error) {
		return h.markReadUseCase.Execute(sessionID, &req)
	}, "Mensagens marcadas como lidas com sucesso")
}

// DownloadImage faz download de imagem
// @Summary Faz download de imagem
// @Description Baixa uma imagem recebida via WhatsApp e retorna os dados em base64
// @Tags chats
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.DownloadImageRequest true "Dados da imagem para download"
// @Success 200 {object} base.APIResponse "Download da imagem concluído com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /chat/{sessionID}/download/image [post]
func (h *Handler) DownloadImage(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.DownloadImageRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	h.HandleUseCaseExecution(w, "fazer download da imagem", func() (interface{}, error) {
		return h.downloadImageUseCase.Execute(sessionID, &req)
	}, "Download da imagem concluído com sucesso")
}

// DownloadVideo faz download de vídeo
// @Summary Faz download de vídeo
// @Description Baixa um vídeo recebido via WhatsApp e retorna os dados em base64
// @Tags chats
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.DownloadVideoRequest true "Dados do vídeo para download"
// @Success 200 {object} base.APIResponse "Download do vídeo concluído com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /chat/{sessionID}/download/video [post]
func (h *Handler) DownloadVideo(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.DownloadVideoRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	h.HandleUseCaseExecution(w, "fazer download do vídeo", func() (interface{}, error) {
		return h.downloadVideoUseCase.Execute(sessionID, &req)
	}, "Download do vídeo concluído com sucesso")
}

// DownloadAudio faz download de áudio
// @Summary Faz download de áudio
// @Description Baixa um áudio recebido via WhatsApp e retorna os dados em base64
// @Tags chats
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.DownloadAudioRequest true "Dados do áudio para download"
// @Success 200 {object} base.APIResponse "Download do áudio concluído com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /chat/{sessionID}/download/audio [post]
func (h *Handler) DownloadAudio(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.DownloadAudioRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	h.HandleUseCaseExecution(w, "fazer download do áudio", func() (interface{}, error) {
		return h.downloadAudioUseCase.Execute(sessionID, &req)
	}, "Download do áudio concluído com sucesso")
}

// DownloadDocument faz download de documento
// @Summary Faz download de documento
// @Description Baixa um documento recebido via WhatsApp e retorna os dados em base64
// @Tags chats
// @Accept json
// @Produce json
// @Param sessionID path string true "ID ou nome da sessão"
// @Param request body requests.DownloadDocumentRequest true "Dados do documento para download"
// @Success 200 {object} base.APIResponse "Download do documento concluído com sucesso"
// @Failure 400 {object} base.APIResponse "Dados inválidos"
// @Failure 500 {object} base.APIResponse "Erro interno do servidor"
// @Router /chat/{sessionID}/download/document [post]
func (h *Handler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := middleware.RequireSessionID(w, r)
	if !ok {
		return
	}

	var req requests.DownloadDocumentRequest
	if !h.DecodeJSONOrError(w, r, &req) {
		return
	}

	h.HandleUseCaseExecution(w, "fazer download do documento", func() (interface{}, error) {
		return h.downloadDocumentUseCase.Execute(sessionID, &req)
	}, "Download do documento concluído com sucesso")
}
