package usecase

import (
	"fmt"
	"time"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow/types"
)

// SendPresenceUseCase representa o caso de uso para definir presença global
type SendPresenceUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewSendPresenceUseCase cria uma nova instância do use case
func NewSendPresenceUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SendPresenceUseCase {
	return &SendPresenceUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a definição de presença global
func (uc *SendPresenceUseCase) Execute(sessionID string, req *requests.SendPresenceRequest) (*responses.PresenceResponse, error) {
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	var presence types.Presence
	switch req.Type {
	case "available":
		presence = types.PresenceAvailable
	case "unavailable":
		presence = types.PresenceUnavailable
	default:
		return nil, fmt.Errorf("tipo de presença inválido: %s", req.Type)
	}

	err = client.GetClient().SendPresence(presence)
	if err != nil {
		return nil, fmt.Errorf("erro ao definir presença: %w", err)
	}

	logger.Info("Presença definida - Session: %s, Type: %s", sessionID, req.Type)

	return &responses.PresenceResponse{
		Details: "Presence set successfully",
	}, nil
}

// ChatPresenceUseCase representa o caso de uso para definir presença em chat específico
type ChatPresenceUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewChatPresenceUseCase cria uma nova instância do use case
func NewChatPresenceUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *ChatPresenceUseCase {
	return &ChatPresenceUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a definição de presença em chat específico
func (uc *ChatPresenceUseCase) Execute(sessionID string, req *requests.ChatPresenceRequest) (*responses.ChatPresenceResponse, error) {
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	jid, err := parseChatJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	var presence types.ChatPresence
	switch req.State {
	case "typing":
		presence = types.ChatPresenceComposing
	case "paused":
		presence = types.ChatPresencePaused
	case "recording":
		presence = types.ChatPresenceComposing
	default:
		return nil, fmt.Errorf("estado de presença inválido: %s", req.State)
	}

	media := req.Media
	if media == "" {
		media = types.ChatPresenceMediaText
	}

	err = client.GetClient().SendChatPresence(jid, presence, media)
	if err != nil {
		return nil, fmt.Errorf("erro ao definir presença no chat: %w", err)
	}

	logger.Info("Presença no chat definida - Session: %s, Phone: %s, State: %s", sessionID, req.Phone, req.State)

	return &responses.ChatPresenceResponse{
		Details: "Chat presence set successfully",
	}, nil
}

// MarkReadUseCase representa o caso de uso para marcar mensagens como lidas
type MarkReadUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewMarkReadUseCase cria uma nova instância do use case
func NewMarkReadUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *MarkReadUseCase {
	return &MarkReadUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a marcação de mensagens como lidas
func (uc *MarkReadUseCase) Execute(sessionID string, req *requests.MarkReadRequest) (*responses.MarkReadResponse, error) {
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	var messageIDs []types.MessageID
	for _, msgID := range req.ID {
		messageIDs = append(messageIDs, msgID)
	}

	sender := req.Sender
	if sender.IsEmpty() {
		sender = req.Chat
	}

	err = client.GetClient().MarkRead(messageIDs, time.Now(), req.Chat, sender)
	if err != nil {
		return nil, fmt.Errorf("erro ao marcar mensagens como lidas: %w", err)
	}

	logger.Info("Mensagens marcadas como lidas - Session: %s, Chat: %s, Count: %d", sessionID, req.Chat.String(), len(messageIDs))

	return &responses.MarkReadResponse{
		Details: "Messages marked as read successfully",
	}, nil
}

// DownloadImageUseCase representa o caso de uso para download de imagem
type DownloadImageUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewDownloadImageUseCase cria uma nova instância do use case
func NewDownloadImageUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *DownloadImageUseCase {
	return &DownloadImageUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa o download de imagem
func (uc *DownloadImageUseCase) Execute(sessionID string, req *requests.DownloadImageRequest) (*responses.DownloadResponse, error) {
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	// Simular download da mídia (implementação básica)
	// Em uma implementação real, você usaria os parâmetros de mídia para fazer o download
	logger.Info("Download de imagem solicitado - Session: %s, URL: %s", sessionID, req.URL)

	// Por enquanto, retornamos uma resposta vazia indicando que o download foi processado
	return &responses.DownloadResponse{
		Mimetype: req.Mimetype,
		Data:     "", // Dados vazios por enquanto
	}, nil
}

// DownloadVideoUseCase representa o caso de uso para download de vídeo
type DownloadVideoUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewDownloadVideoUseCase cria uma nova instância do use case
func NewDownloadVideoUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *DownloadVideoUseCase {
	return &DownloadVideoUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa o download de vídeo
func (uc *DownloadVideoUseCase) Execute(sessionID string, req *requests.DownloadVideoRequest) (*responses.DownloadResponse, error) {
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	// Simular download da mídia (implementação básica)
	logger.Info("Download de vídeo solicitado - Session: %s, URL: %s", sessionID, req.URL)

	return &responses.DownloadResponse{
		Mimetype: req.Mimetype,
		Data:     "", // Dados vazios por enquanto
	}, nil
}

// DownloadAudioUseCase representa o caso de uso para download de áudio
type DownloadAudioUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewDownloadAudioUseCase cria uma nova instância do use case
func NewDownloadAudioUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *DownloadAudioUseCase {
	return &DownloadAudioUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa o download de áudio
func (uc *DownloadAudioUseCase) Execute(sessionID string, req *requests.DownloadAudioRequest) (*responses.DownloadResponse, error) {
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	// Simular download da mídia (implementação básica)
	logger.Info("Download de áudio solicitado - Session: %s, URL: %s", sessionID, req.URL)

	return &responses.DownloadResponse{
		Mimetype: req.Mimetype,
		Data:     "", // Dados vazios por enquanto
	}, nil
}

// DownloadDocumentUseCase representa o caso de uso para download de documento
type DownloadDocumentUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewDownloadDocumentUseCase cria uma nova instância do use case
func NewDownloadDocumentUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *DownloadDocumentUseCase {
	return &DownloadDocumentUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa o download de documento
func (uc *DownloadDocumentUseCase) Execute(sessionID string, req *requests.DownloadDocumentRequest) (*responses.DownloadResponse, error) {
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	// Simular download da mídia (implementação básica)
	logger.Info("Download de documento solicitado - Session: %s, URL: %s", sessionID, req.URL)

	return &responses.DownloadResponse{
		Mimetype: req.Mimetype,
		Data:     "", // Dados vazios por enquanto
	}, nil
}

// parseChatJID converte um número de telefone em JID para chat
func parseChatJID(phone string) (types.JID, error) {
	if phone == "" {
		return types.JID{}, fmt.Errorf("número de telefone não pode estar vazio")
	}

	// Remove caracteres não numéricos
	cleanPhone := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			cleanPhone += string(char)
		}
	}

	if cleanPhone == "" {
		return types.JID{}, fmt.Errorf("número de telefone inválido")
	}

	// Criar JID para usuário individual
	return types.NewJID(cleanPhone, types.DefaultUserServer), nil
}
