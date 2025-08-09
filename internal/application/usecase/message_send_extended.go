package usecase

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif" // Import GIF decoder
	"image/jpeg"
	_ "image/png" // Import PNG decoder
	"net/http"
	"os"
	"strings"
	"time"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"

	"github.com/nfnt/resize"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendImageMessageUseCase representa o caso de uso para envio de imagens
type SendImageMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendImageMessageUseCase cria uma nova instância do use case
func NewSendImageMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendImageMessageUseCase {
	return &SendImageMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de imagem
func (uc *SendImageMessageUseCase) Execute(sessionID string, req *requests.SendImageMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	mediaData, err := decodeMediaData(req.Image)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados da imagem: %w", err)
	}

	uploaded, err := client.GetClient().Upload(context.Background(), mediaData, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload da imagem: %w", err)
	}

	// Generate thumbnail for image
	var thumbnailBytes []byte
	reader := bytes.NewReader(mediaData)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar imagem para thumbnail: %w", err)
	}

	// Resize to width 72 using Lanczos resampling and preserve aspect ratio
	m := resize.Thumbnail(72, 72, img, resize.Lanczos3)

	tmpFile, err := os.CreateTemp("", "resized-*.jpg")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo temporário para thumbnail: %w", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Write new image to file
	if err := jpeg.Encode(tmpFile, m, nil); err != nil {
		return nil, fmt.Errorf("erro ao codificar JPEG: %w", err)
	}

	thumbnailBytes, err = os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de thumbnail: %w", err)
	}

	// Determine mimetype if not provided
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(mediaData)
	}

	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
			JPEGThumbnail: thumbnailBytes,
		},
	}

	if req.ContextInfo.StanzaID != nil {
		msg.ImageMessage.ContextInfo = &waE2E.ContextInfo{
			StanzaID:      proto.String(*req.ContextInfo.StanzaID),
			Participant:   proto.String(*req.ContextInfo.Participant),
			QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
		}
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar imagem: %w", err)
	}

	logger.Info("Imagem enviada - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendAudioMessageUseCase representa o caso de uso para envio de áudio
type SendAudioMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendAudioMessageUseCase cria uma nova instância do use case
func NewSendAudioMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendAudioMessageUseCase {
	return &SendAudioMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de áudio
func (uc *SendAudioMessageUseCase) Execute(sessionID string, req *requests.SendAudioMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	mediaData, err := decodeMediaData(req.Audio)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados do áudio: %w", err)
	}

	uploaded, err := client.GetClient().Upload(context.Background(), mediaData, whatsmeow.MediaAudio)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload do áudio: %w", err)
	}

	ptt := true
	mime := "audio/ogg; codecs=opus"

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      &mime,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
			PTT:           &ptt,
		},
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar áudio: %w", err)
	}

	logger.Info("Áudio enviado - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendDocumentMessageUseCase representa o caso de uso para envio de documentos
type SendDocumentMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendDocumentMessageUseCase cria uma nova instância do use case
func NewSendDocumentMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendDocumentMessageUseCase {
	return &SendDocumentMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de documento
func (uc *SendDocumentMessageUseCase) Execute(sessionID string, req *requests.SendDocumentMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	mediaData, err := decodeMediaData(req.Document)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados do documento: %w", err)
	}

	uploaded, err := client.GetClient().Upload(context.Background(), mediaData, whatsmeow.MediaDocument)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload do documento: %w", err)
	}

	// Determine mimetype if not provided
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(mediaData)
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           proto.String(uploaded.URL),
			FileName:      &req.FileName,
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
			Caption:       proto.String(req.Caption),
		},
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar documento: %w", err)
	}

	logger.Info("Documento enviado - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendVideoMessageUseCase representa o caso de uso para envio de vídeos
type SendVideoMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendVideoMessageUseCase cria uma nova instância do use case
func NewSendVideoMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendVideoMessageUseCase {
	return &SendVideoMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de vídeo
func (uc *SendVideoMessageUseCase) Execute(sessionID string, req *requests.SendVideoMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	mediaData, err := decodeMediaData(req.Video)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados do vídeo: %w", err)
	}

	uploaded, err := client.GetClient().Upload(context.Background(), mediaData, whatsmeow.MediaVideo)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload do vídeo: %w", err)
	}

	// Determine mimetype if not provided
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(mediaData)
	}

	msg := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			Caption:       proto.String(req.Caption),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
			JPEGThumbnail: req.JPEGThumbnail,
		},
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar vídeo: %w", err)
	}

	logger.Info("Vídeo enviado - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendStickerMessageUseCase representa o caso de uso para envio de stickers
type SendStickerMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendStickerMessageUseCase cria uma nova instância do use case
func NewSendStickerMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendStickerMessageUseCase {
	return &SendStickerMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de sticker
func (uc *SendStickerMessageUseCase) Execute(sessionID string, req *requests.SendStickerMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	mediaData, err := decodeMediaData(req.Sticker)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados do sticker: %w", err)
	}

	uploaded, err := client.GetClient().Upload(context.Background(), mediaData, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer upload do sticker: %w", err)
	}

	// Determine mimetype if not provided
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = http.DetectContentType(mediaData)
	}

	msg := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(mediaData))),
			PngThumbnail:  req.PngThumbnail,
		},
	}

	// Add ContextInfo if provided (following wuzapi pattern)
	if req.ContextInfo.StanzaID != nil {
		msg.ExtendedTextMessage = &waE2E.ExtendedTextMessage{
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String(*req.ContextInfo.StanzaID),
				Participant:   proto.String(*req.ContextInfo.Participant),
				QuotedMessage: &waE2E.Message{Conversation: proto.String("")},
			},
		}
	}
	if req.ContextInfo.MentionedJID != nil {
		if msg.ExtendedTextMessage == nil {
			msg.ExtendedTextMessage = &waE2E.ExtendedTextMessage{}
		}
		if msg.ExtendedTextMessage.ContextInfo == nil {
			msg.ExtendedTextMessage.ContextInfo = &waE2E.ContextInfo{}
		}
		msg.ExtendedTextMessage.ContextInfo.MentionedJID = req.ContextInfo.MentionedJID
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar sticker: %w", err)
	}

	logger.Info("Sticker enviado - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendLocationMessageUseCase representa o caso de uso para envio de localização
type SendLocationMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendLocationMessageUseCase cria uma nova instância do use case
func NewSendLocationMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendLocationMessageUseCase {
	return &SendLocationMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de localização
func (uc *SendLocationMessageUseCase) Execute(sessionID string, req *requests.SendLocationMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  &req.Latitude,
			DegreesLongitude: &req.Longitude,
			Name:             &req.Name,
		},
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar localização: %w", err)
	}

	logger.Info("Localização enviada - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendContactMessageUseCase representa o caso de uso para envio de contato
type SendContactMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendContactMessageUseCase cria uma nova instância do use case
func NewSendContactMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendContactMessageUseCase {
	return &SendContactMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de contato
func (uc *SendContactMessageUseCase) Execute(sessionID string, req *requests.SendContactMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	msg := &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: &req.Name,
			Vcard:       &req.Vcard,
		},
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar contato: %w", err)
	}

	logger.Info("Contato enviado - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendButtonsMessageUseCase representa o caso de uso para envio de botões
type SendButtonsMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendButtonsMessageUseCase cria uma nova instância do use case
func NewSendButtonsMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendButtonsMessageUseCase {
	return &SendButtonsMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de botões
func (uc *SendButtonsMessageUseCase) Execute(sessionID string, req *requests.SendButtonsMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	var buttons []*waE2E.ButtonsMessage_Button
	for _, item := range req.Buttons {
		buttons = append(buttons, &waE2E.ButtonsMessage_Button{
			ButtonID:       proto.String(item.ButtonID),
			ButtonText:     &waE2E.ButtonsMessage_Button_ButtonText{DisplayText: proto.String(item.ButtonText)},
			Type:           waE2E.ButtonsMessage_Button_RESPONSE.Enum(),
			NativeFlowInfo: &waE2E.ButtonsMessage_Button_NativeFlowInfo{},
		})
	}

	msg2 := &waE2E.ButtonsMessage{
		ContentText: proto.String(req.Title),
		HeaderType:  waE2E.ButtonsMessage_EMPTY.Enum(),
		Buttons:     buttons,
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, &waE2E.Message{ViewOnceMessage: &waE2E.FutureProofMessage{
		Message: &waE2E.Message{
			ButtonsMessage: msg2,
		},
	}}, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar botões: %w", err)
	}

	logger.Info("Botões enviados - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendListMessageUseCase representa o caso de uso para envio de lista
type SendListMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendListMessageUseCase cria uma nova instância do use case
func NewSendListMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendListMessageUseCase {
	return &SendListMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de lista
func (uc *SendListMessageUseCase) Execute(sessionID string, req *requests.SendListMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	var sections []*waE2E.ListMessage_Section
	for _, sec := range req.Sections {
		var rows []*waE2E.ListMessage_Row
		for _, item := range sec.Rows {
			rows = append(rows, &waE2E.ListMessage_Row{
				RowID:       proto.String(item.RowID),
				Title:       proto.String(item.Title),
				Description: proto.String(item.Desc),
			})
		}
		sections = append(sections, &waE2E.ListMessage_Section{
			Title: proto.String(sec.Title),
			Rows:  rows,
		})
	}

	listMsg := &waE2E.ListMessage{
		Title:       proto.String(req.TopText),
		Description: proto.String(req.Desc),
		ButtonText:  proto.String(req.ButtonText),
		ListType:    waE2E.ListMessage_SINGLE_SELECT.Enum(),
		Sections:    sections,
	}

	if req.FooterText != "" {
		listMsg.FooterText = proto.String(req.FooterText)
	}

	msg := &waE2E.Message{
		ViewOnceMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				ListMessage: listMsg,
			},
		},
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar lista: %w", err)
	}

	logger.Info("Lista enviada - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendPollMessageUseCase representa o caso de uso para envio de enquete
type SendPollMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendPollMessageUseCase cria uma nova instância do use case
func NewSendPollMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendPollMessageUseCase {
	return &SendPollMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa o envio de enquete
func (uc *SendPollMessageUseCase) Execute(sessionID string, req *requests.SendPollMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	if msgID == "" {
		msgID = client.GetClient().GenerateMessageID()
	}

	pollMessage := client.GetClient().BuildPollCreation(req.Header, req.Options, 1)
	resp, err := client.GetClient().SendMessage(context.Background(), recipient, pollMessage, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar enquete: %w", err)
	}

	logger.Info("Enquete enviada - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}

// SendEditMessageUseCase representa o caso de uso para edição de mensagem
type SendEditMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewSendEditMessageUseCase cria uma nova instância do use case
func NewSendEditMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendEditMessageUseCase {
	return &SendEditMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa a edição de mensagem
func (uc *SendEditMessageUseCase) Execute(sessionID string, req *requests.SendEditMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &req.Body,
		},
	}

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

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, client.GetClient().BuildEdit(recipient, req.ID, msg))
	if err != nil {
		return nil, fmt.Errorf("erro ao editar mensagem: %w", err)
	}

	logger.Info("Mensagem editada - ID: %s, Timestamp: %v", req.ID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        req.ID,
	}, nil
}

// DeleteMessageUseCase representa o caso de uso para deletar mensagem
type DeleteMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewDeleteMessageUseCase cria uma nova instância do use case
func NewDeleteMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *DeleteMessageUseCase {
	return &DeleteMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa a exclusão de mensagem
func (uc *DeleteMessageUseCase) Execute(sessionID string, req *requests.DeleteMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, client.GetClient().BuildRevoke(recipient, types.EmptyJID, req.ID))
	if err != nil {
		return nil, fmt.Errorf("erro ao deletar mensagem: %w", err)
	}

	logger.Info("Mensagem deletada - ID: %s, Timestamp: %v", req.ID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Deleted",
		Timestamp: resp.Timestamp.Unix(),
		ID:        req.ID,
	}, nil
}

// ReactMessageUseCase representa o caso de uso para reagir a mensagem
type ReactMessageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewReactMessageUseCase cria uma nova instância do use case
func NewReactMessageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *ReactMessageUseCase {
	return &ReactMessageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa a reação a mensagem
func (uc *ReactMessageUseCase) Execute(sessionID string, req *requests.ReactMessageRequest) (*responses.SendMessageResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	recipient, err := parseJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	msgID := req.ID
	fromMe := false
	if strings.HasPrefix(msgID, "me:") {
		fromMe = true
		msgID = msgID[len("me:"):]
	}

	reaction := req.Body
	if reaction == "remove" {
		reaction = ""
	}

	msg := &waE2E.Message{
		ReactionMessage: &waE2E.ReactionMessage{
			Key: &waCommon.MessageKey{
				RemoteJID: proto.String(recipient.String()),
				FromMe:    proto.Bool(fromMe),
				ID:        proto.String(msgID),
			},
			Text:              proto.String(reaction),
			GroupingKey:       proto.String(reaction),
			SenderTimestampMS: proto.Int64(time.Now().UnixMilli()),
		},
	}

	resp, err := client.GetClient().SendMessage(context.Background(), recipient, msg, whatsmeow.SendRequestExtra{ID: msgID})
	if err != nil {
		return nil, fmt.Errorf("erro ao reagir à mensagem: %w", err)
	}

	logger.Info("Reação enviada - ID: %s, Timestamp: %v", msgID, resp.Timestamp)

	return &responses.SendMessageResponse{
		Details:   "Sent",
		Timestamp: resp.Timestamp.Unix(),
		ID:        msgID,
	}, nil
}
