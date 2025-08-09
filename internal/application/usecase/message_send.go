package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendTextMessageUseCase representa o caso de uso para envio de mensagens de texto
type SendTextMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewSendTextMessageUseCase cria uma nova instância do use case
func NewSendTextMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendTextMessageUseCase {
	return &SendTextMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa o envio de mensagem de texto
func (uc *SendTextMessageUseCase) Execute(sessionID string, req *requests.SendTextMessageRequest) (*responses.SendMessageResponse, error) {
	// Verificar se a sessão existe e está conectada
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	// Obter cliente WhatsApp da sessão
	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	// Validar e parsear o número de telefone
	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	// Gerar ID da mensagem se não fornecido
	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	// Criar mensagem de texto
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(req.Body),
		},
	}

	// Adicionar informações de contexto se fornecidas
	if req.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*req.ContextInfo.StanzaID),
			Participant:   proto.String(*req.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}
	if req.ContextInfo.MentionedJID != nil {
		if msg.ExtendedTextMessage.ContextInfo == nil {
			msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{}
		}
		msg.ExtendedTextMessage.ContextInfo.MentionedJID = req.ContextInfo.MentionedJID
	}

	// Enviar mensagem
	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar mensagem: %w", err)
	}

	logger.Info("Mensagem de texto enviada - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendMediaMessageUseCase representa o caso de uso para envio de mídia
type SendMediaMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewSendMediaMessageUseCase cria uma nova instância do use case
func NewSendMediaMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendMediaMessageUseCase {
	return &SendMediaMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa o envio de mídia
func (uc *SendMediaMessageUseCase) Execute(sessionID string, req *requests.SendMediaMessageRequest) (*responses.SendMessageResponse, error) {
	// Verificar se a sessão existe e está conectada
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	// Obter cliente WhatsApp da sessão
	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	// Validar e parsear o número de telefone
	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	// Gerar ID da mensagem se não fornecido
	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	// Decodificar dados da mídia (base64)
	mediaData, err := decodeMediaData(req.MediaData)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados da mídia: %w", err)
	}

	// Determinar tipo de mídia
	mediaType := determineMediaType(req.MediaData, req.MimeType)

	// Fazer upload da mídia
	uploaded, err := client.GetClient().Upload(context.Background(), mediaData, mediaType)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload da mídia: %w", err)
	}

	// Criar mensagem baseada no tipo de mídia
	var msg *waE2E.Message
	switch mediaType {
	case whatsmeow.MediaImage:
		msg = createImageMessage(uploaded, mediaData, req)
	case whatsmeow.MediaAudio:
		msg = createAudioMessage(uploaded, mediaData, req)
	case whatsmeow.MediaVideo:
		msg = createVideoMessage(uploaded, mediaData, req)
	case whatsmeow.MediaDocument:
		msg = createDocumentMessage(uploaded, mediaData, req)
	default:
		return nil, fmt.Errorf("tipo de mídia não suportado")
	}

	// Enviar mensagem
	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar mídia: %w", err)
	}

	logger.Info("Mídia enviada - ID: %s, Timestamp: %v, Tipo: %v", msgID, resp.Timestamp, mediaType)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// parseJID converte uma string em JID do WhatsApp
func parseJID(phone string) (types.JID, error) {
	if phone[0] == '+' {
		phone = phone[1:]
	}
	if !strings.ContainsRune(phone, '@') {
		return types.NewJID(phone, types.DefaultUserServer), nil
	} else {
		recipient, err := types.ParseJID(phone)
		if err != nil {
			return recipient, fmt.Errorf("JID inválido: %w", err)
		}
		if recipient.User == "" {
			return recipient, fmt.Errorf("JID inválido: servidor não especificado")
		}
		return recipient, nil
	}
}

// decodeMediaData decodifica dados de mídia em base64
func decodeMediaData(mediaData string) ([]byte, error) {
	// Verificar se é data URL
	if strings.HasPrefix(mediaData, "data:") {
		// Extrair apenas a parte base64
		parts := strings.Split(mediaData, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("formato de data URL inválido")
		}
		mediaData = parts[1]
	}

	// Decodificar base64
	data, err := base64.StdEncoding.DecodeString(mediaData)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar base64: %w", err)
	}

	return data, nil
}

// determineMediaType determina o tipo de mídia baseado no conteúdo
func determineMediaType(mediaData, mimeType string) whatsmeow.MediaType {
	if mimeType != "" {
		if strings.HasPrefix(mimeType, "image/") {
			return whatsmeow.MediaImage
		}
		if strings.HasPrefix(mimeType, "audio/") {
			return whatsmeow.MediaAudio
		}
		if strings.HasPrefix(mimeType, "video/") {
			return whatsmeow.MediaVideo
		}
		return whatsmeow.MediaDocument
	}

	// Detectar pelo data URL
	if strings.HasPrefix(mediaData, "data:image/") {
		return whatsmeow.MediaImage
	}
	if strings.HasPrefix(mediaData, "data:audio/") {
		return whatsmeow.MediaAudio
	}
	if strings.HasPrefix(mediaData, "data:video/") {
		return whatsmeow.MediaVideo
	}

	return whatsmeow.MediaDocument
}

// createImageMessage cria uma mensagem de imagem
func createImageMessage(uploaded whatsmeow.UploadResponse, mediaData []byte, req *requests.SendMediaMessageRequest) *waE2E.Message {
	return &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(req.MimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
		},
	}
}

// createAudioMessage cria uma mensagem de áudio
func createAudioMessage(uploaded whatsmeow.UploadResponse, mediaData []byte, req *requests.SendMediaMessageRequest) *waE2E.Message {
	return &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(req.MimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
		},
	}
}

// createVideoMessage cria uma mensagem de vídeo
func createVideoMessage(uploaded whatsmeow.UploadResponse, mediaData []byte, req *requests.SendMediaMessageRequest) *waE2E.Message {
	return &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(req.MimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
		},
	}
}

// createDocumentMessage cria uma mensagem de documento
func createDocumentMessage(uploaded whatsmeow.UploadResponse, mediaData []byte, req *requests.SendMediaMessageRequest) *waE2E.Message {
	return &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(req.MimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
		},
	}
}
