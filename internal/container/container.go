// Package container contém o sistema de injeção de dependências da aplicação
//
// Este pacote é responsável por:
// - Configurar e instanciar todas as dependências da aplicação
// - Seguir os princípios de Clean Architecture na organização das camadas
// - Implementar o padrão Dependency Injection Container
// - Fornecer uma interface clara para acesso às dependências
package container

import (
	"sync"

	"wazmeow/internal/application/usecase"
	"wazmeow/internal/config"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/domain/service"
	"wazmeow/internal/infra/database"
	"wazmeow/internal/infra/webhook"
	"wazmeow/internal/infra/whatsapp"
)

// Container representa o container de injeção de dependências
// Organizado seguindo as camadas da Clean Architecture:
// 1. Configuração
// 2. Infraestrutura (database, logger)
// 3. Repositories (implementações de persistência)
// 4. Domain Services (regras de negócio)
// 5. Use Cases (orquestração)
type Container struct {
	// ========================================
	// CONFIGURAÇÃO
	// ========================================
	config *config.Config // Configurações da aplicação carregadas do ambiente

	// ========================================
	// INFRAESTRUTURA
	// ========================================
	db             *database.Connection    // Conexão para WhatsApp (sqlstore)
	bunDB          *database.BunConnection // Conexão Bun ORM para sessões da aplicação
	sessionManager *whatsapp.SessionManager
	clientFactory  *whatsapp.ClientFactory
	webhookService *webhook.WebhookService // Serviço de webhooks

	// ========================================
	// REPOSITORIES
	// ========================================
	sessionRepo repository.SessionRepository

	// ========================================
	// DOMAIN SERVICES
	// ========================================
	sessionDomainService *service.SessionDomainService

	// ========================================
	// USE CASES - Organizados por categoria
	// ========================================
	sessionUseCases    *SessionUseCases
	messageUseCases    *MessageUseCases
	webhookUseCases    *WebhookUseCases
	userUseCases       *UserUseCases
	chatUseCases       *ChatUseCases
	groupUseCases      *GroupUseCases
	newsletterUseCases *NewsletterUseCases

	// ========================================
	// CONTROLE INTERNO
	// ========================================
	initialized bool
	mu          sync.RWMutex
}

// SessionUseCases agrupa todos os use cases relacionados a sessões
type SessionUseCases struct {
	// Gerenciamento básico de sessões
	Create  *usecase.CreateSessionUseCase
	List    *usecase.ListSessionsUseCase
	Delete  *usecase.DeleteSessionUseCase
	GetInfo *usecase.GetSessionInfoUseCase

	// Conectividade e autenticação
	Connect *usecase.ConnectSessionUseCase
	Logout  *usecase.LogoutSessionUseCase
	GetQR   *usecase.GetQRCodeUseCase

	// Setup e configuração
	PairPhone *usecase.PairPhoneUseCase
	SetProxy  *usecase.SetProxyUseCase
}

// MessageUseCases agrupa todos os use cases relacionados a mensagens
type MessageUseCases struct {
	// Envio de mensagens básicas
	SendText  *usecase.SendTextMessageUseCase
	SendMedia *usecase.SendMediaMessageUseCase

	// Envio de mensagens específicas
	SendImage    *usecase.SendImageMessageUseCase
	SendAudio    *usecase.SendAudioMessageUseCase
	SendDocument *usecase.SendDocumentMessageUseCase
	SendVideo    *usecase.SendVideoMessageUseCase
	SendSticker  *usecase.SendStickerMessageUseCase
	SendLocation *usecase.SendLocationMessageUseCase
	SendContact  *usecase.SendContactMessageUseCase
	SendButtons  *usecase.SendButtonsMessageUseCase
	SendList     *usecase.SendListMessageUseCase
	SendPoll     *usecase.SendPollMessageUseCase

	// Operações de mensagem
	SendEdit      *usecase.SendEditMessageUseCase
	DeleteMessage *usecase.DeleteMessageUseCase
	React         *usecase.ReactMessageUseCase
}

// WebhookUseCases agrupa todos os use cases relacionados a webhooks
type WebhookUseCases struct {
	SetWebhook    *usecase.SetWebhookUseCase
	GetWebhook    *usecase.GetWebhookUseCase
	UpdateWebhook *usecase.UpdateWebhookUseCase
	DeleteWebhook *usecase.DeleteWebhookUseCase
}

// UserUseCases agrupa todos os use cases relacionados a usuários
type UserUseCases struct {
	GetUserInfo *usecase.GetUserInfoUseCase
	CheckUser   *usecase.CheckUserUseCase
	GetAvatar   *usecase.GetAvatarUseCase
	GetContacts *usecase.GetContactsUseCase
}

// ChatUseCases agrupa todos os use cases relacionados a chat
type ChatUseCases struct {
	SendPresence     *usecase.SendPresenceUseCase
	ChatPresence     *usecase.ChatPresenceUseCase
	MarkRead         *usecase.MarkReadUseCase
	DownloadImage    *usecase.DownloadImageUseCase
	DownloadVideo    *usecase.DownloadVideoUseCase
	DownloadAudio    *usecase.DownloadAudioUseCase
	DownloadDocument *usecase.DownloadDocumentUseCase
}

// GroupUseCases agrupa todos os use cases relacionados a grupos
type GroupUseCases struct {
	CreateGroup             *usecase.CreateGroupUseCase
	SetGroupPhoto           *usecase.SetGroupPhotoUseCase
	UpdateGroupParticipants *usecase.UpdateGroupParticipantsUseCase
	LeaveGroup              *usecase.LeaveGroupUseCase
	JoinGroup               *usecase.JoinGroupUseCase
	GetGroupInfo            *usecase.GetGroupInfoUseCase
	ListGroups              *usecase.ListGroupsUseCase
	GetGroupInviteLink      *usecase.GetGroupInviteLinkUseCase
	RevokeGroupInviteLink   *usecase.RevokeGroupInviteLinkUseCase
	SetGroupName            *usecase.SetGroupNameUseCase
	SetGroupTopic           *usecase.SetGroupTopicUseCase
}

// NewsletterUseCases agrupa todos os use cases relacionados a newsletters
type NewsletterUseCases struct {
	ListNewsletter *usecase.ListNewsletterUseCase
}

// New cria um novo container com todas as dependências configuradas
// Este é o ponto de entrada principal para inicialização da aplicação
func New() (*Container, error) {
	builder := NewBuilder()
	return builder.Build()
}

// NewWithConfig cria um container com configuração específica
// Útil para testes e casos especiais
func NewWithConfig(cfg *config.Config) (*Container, error) {
	builder := NewBuilder().WithConfig(cfg)
	return builder.Build()
}

// IsInitialized verifica se o container foi completamente inicializado
func (c *Container) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}

// setInitialized marca o container como inicializado (thread-safe)
func (c *Container) setInitialized(value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.initialized = value
}

// ========================================
// GETTERS PARA DEPENDÊNCIAS
// ========================================

// GetConfig retorna a configuração da aplicação
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// GetDB retorna a conexão do banco de dados para WhatsApp
func (c *Container) GetDB() *database.Connection {
	return c.db
}

// GetBunDB retorna a conexão Bun ORM
func (c *Container) GetBunDB() *database.BunConnection {
	return c.bunDB
}

// GetSessionRepository retorna o repository de sessões
func (c *Container) GetSessionRepository() repository.SessionRepository {
	return c.sessionRepo
}

// GetSessionDomainService retorna o domain service de sessões
func (c *Container) GetSessionDomainService() *service.SessionDomainService {
	return c.sessionDomainService
}

// GetSessionManager retorna o gerenciador de sessões WhatsApp
func (c *Container) GetSessionManager() *whatsapp.SessionManager {
	return c.sessionManager
}

// GetClientFactory retorna o factory de clientes WhatsApp
func (c *Container) GetClientFactory() *whatsapp.ClientFactory {
	if c.clientFactory == nil {
		c.clientFactory = whatsapp.NewClientFactory(c.db.Store, c.GetSessionRepository(), c.GetSessionDomainService(), c.webhookService)
	}
	return c.clientFactory
}

// ========================================
// GETTERS PARA USE CASES
// ========================================

// GetCreateSessionUseCase retorna o use case de criação de sessão
func (c *Container) GetCreateSessionUseCase() *usecase.CreateSessionUseCase {
	return c.sessionUseCases.Create
}

// GetListSessionsUseCase retorna o use case de listagem de sessões
func (c *Container) GetListSessionsUseCase() *usecase.ListSessionsUseCase {
	return c.sessionUseCases.List
}

// GetDeleteSessionUseCase retorna o use case de remoção de sessão
func (c *Container) GetDeleteSessionUseCase() *usecase.DeleteSessionUseCase {
	return c.sessionUseCases.Delete
}

// GetGetSessionInfoUseCase retorna o use case de informações da sessão
func (c *Container) GetGetSessionInfoUseCase() *usecase.GetSessionInfoUseCase {
	return c.sessionUseCases.GetInfo
}

// GetConnectSessionUseCase retorna o use case de conexão de sessão
func (c *Container) GetConnectSessionUseCase() *usecase.ConnectSessionUseCase {
	return c.sessionUseCases.Connect
}

// GetLogoutSessionUseCase retorna o use case de logout de sessão
func (c *Container) GetLogoutSessionUseCase() *usecase.LogoutSessionUseCase {
	return c.sessionUseCases.Logout
}

// GetGetQRCodeUseCase retorna o use case de obtenção de QR code
func (c *Container) GetGetQRCodeUseCase() *usecase.GetQRCodeUseCase {
	return c.sessionUseCases.GetQR
}

// GetPairPhoneUseCase retorna o use case de emparelhamento por telefone
func (c *Container) GetPairPhoneUseCase() *usecase.PairPhoneUseCase {
	return c.sessionUseCases.PairPhone
}

// GetSetProxyUseCase retorna o use case de configuração de proxy
func (c *Container) GetSetProxyUseCase() *usecase.SetProxyUseCase {
	return c.sessionUseCases.SetProxy
}

// ========================================
// GETTERS PARA MESSAGE USE CASES
// ========================================

// GetSendTextMessageUseCase retorna o use case de envio de mensagem de texto
func (c *Container) GetSendTextMessageUseCase() *usecase.SendTextMessageUseCase {
	return c.messageUseCases.SendText
}

// GetSendMediaMessageUseCase retorna o use case de envio de mídia
func (c *Container) GetSendMediaMessageUseCase() *usecase.SendMediaMessageUseCase {
	return c.messageUseCases.SendMedia
}

// GetSendImageMessageUseCase retorna o use case de envio de imagem
func (c *Container) GetSendImageMessageUseCase() *usecase.SendImageMessageUseCase {
	return c.messageUseCases.SendImage
}

// GetSendAudioMessageUseCase retorna o use case de envio de áudio
func (c *Container) GetSendAudioMessageUseCase() *usecase.SendAudioMessageUseCase {
	return c.messageUseCases.SendAudio
}

// GetSendDocumentMessageUseCase retorna o use case de envio de documento
func (c *Container) GetSendDocumentMessageUseCase() *usecase.SendDocumentMessageUseCase {
	return c.messageUseCases.SendDocument
}

// GetSendVideoMessageUseCase retorna o use case de envio de vídeo
func (c *Container) GetSendVideoMessageUseCase() *usecase.SendVideoMessageUseCase {
	return c.messageUseCases.SendVideo
}

// GetSendStickerMessageUseCase retorna o use case de envio de sticker
func (c *Container) GetSendStickerMessageUseCase() *usecase.SendStickerMessageUseCase {
	return c.messageUseCases.SendSticker
}

// GetSendLocationMessageUseCase retorna o use case de envio de localização
func (c *Container) GetSendLocationMessageUseCase() *usecase.SendLocationMessageUseCase {
	return c.messageUseCases.SendLocation
}

// GetSendContactMessageUseCase retorna o use case de envio de contato
func (c *Container) GetSendContactMessageUseCase() *usecase.SendContactMessageUseCase {
	return c.messageUseCases.SendContact
}

// GetSendButtonsMessageUseCase retorna o use case de envio de botões
func (c *Container) GetSendButtonsMessageUseCase() *usecase.SendButtonsMessageUseCase {
	return c.messageUseCases.SendButtons
}

// GetSendListMessageUseCase retorna o use case de envio de lista
func (c *Container) GetSendListMessageUseCase() *usecase.SendListMessageUseCase {
	return c.messageUseCases.SendList
}

// GetSendPollMessageUseCase retorna o use case de envio de enquete
func (c *Container) GetSendPollMessageUseCase() *usecase.SendPollMessageUseCase {
	return c.messageUseCases.SendPoll
}

// GetSendEditMessageUseCase retorna o use case de edição de mensagem
func (c *Container) GetSendEditMessageUseCase() *usecase.SendEditMessageUseCase {
	return c.messageUseCases.SendEdit
}

// GetDeleteMessageUseCase retorna o use case de exclusão de mensagem
func (c *Container) GetDeleteMessageUseCase() *usecase.DeleteMessageUseCase {
	return c.messageUseCases.DeleteMessage
}

// GetReactMessageUseCase retorna o use case de reação a mensagem
func (c *Container) GetReactMessageUseCase() *usecase.ReactMessageUseCase {
	return c.messageUseCases.React
}

// ========================================
// GETTERS PARA WEBHOOK USE CASES
// ========================================

// GetSetWebhookUseCase retorna o use case de definição de webhook
func (c *Container) GetSetWebhookUseCase() *usecase.SetWebhookUseCase {
	return c.webhookUseCases.SetWebhook
}

// GetGetWebhookUseCase retorna o use case de obtenção de webhook
func (c *Container) GetGetWebhookUseCase() *usecase.GetWebhookUseCase {
	return c.webhookUseCases.GetWebhook
}

// GetUpdateWebhookUseCase retorna o use case de atualização de webhook
func (c *Container) GetUpdateWebhookUseCase() *usecase.UpdateWebhookUseCase {
	return c.webhookUseCases.UpdateWebhook
}

// GetDeleteWebhookUseCase retorna o use case de remoção de webhook
func (c *Container) GetDeleteWebhookUseCase() *usecase.DeleteWebhookUseCase {
	return c.webhookUseCases.DeleteWebhook
}

// GetWebhookService retorna o serviço de webhooks
func (c *Container) GetWebhookService() *webhook.WebhookService {
	return c.webhookService
}

// ========================================
// GETTERS PARA USER USE CASES
// ========================================

// GetSendPresenceUseCase retorna o use case de definição de presença
func (c *Container) GetSendPresenceUseCase() *usecase.SendPresenceUseCase {
	return c.chatUseCases.SendPresence
}

// GetGetUserInfoUseCase retorna o use case de obtenção de informações do usuário
func (c *Container) GetGetUserInfoUseCase() *usecase.GetUserInfoUseCase {
	return c.userUseCases.GetUserInfo
}

// GetCheckUserUseCase retorna o use case de verificação de usuário
func (c *Container) GetCheckUserUseCase() *usecase.CheckUserUseCase {
	return c.userUseCases.CheckUser
}

// GetGetAvatarUseCase retorna o use case de obtenção de avatar
func (c *Container) GetGetAvatarUseCase() *usecase.GetAvatarUseCase {
	return c.userUseCases.GetAvatar
}

// GetGetContactsUseCase retorna o use case de obtenção de contatos
func (c *Container) GetGetContactsUseCase() *usecase.GetContactsUseCase {
	return c.userUseCases.GetContacts
}

// ========================================
// GETTERS PARA CHAT USE CASES
// ========================================

// GetChatPresenceUseCase retorna o use case de presença no chat
func (c *Container) GetChatPresenceUseCase() *usecase.ChatPresenceUseCase {
	return c.chatUseCases.ChatPresence
}

// GetMarkReadUseCase retorna o use case de marcar como lida
func (c *Container) GetMarkReadUseCase() *usecase.MarkReadUseCase {
	return c.chatUseCases.MarkRead
}

// GetDownloadImageUseCase retorna o use case de download de imagem
func (c *Container) GetDownloadImageUseCase() *usecase.DownloadImageUseCase {
	return c.chatUseCases.DownloadImage
}

// GetDownloadVideoUseCase retorna o use case de download de vídeo
func (c *Container) GetDownloadVideoUseCase() *usecase.DownloadVideoUseCase {
	return c.chatUseCases.DownloadVideo
}

// GetDownloadAudioUseCase retorna o use case de download de áudio
func (c *Container) GetDownloadAudioUseCase() *usecase.DownloadAudioUseCase {
	return c.chatUseCases.DownloadAudio
}

// GetDownloadDocumentUseCase retorna o use case de download de documento
func (c *Container) GetDownloadDocumentUseCase() *usecase.DownloadDocumentUseCase {
	return c.chatUseCases.DownloadDocument
}

// ========================================
// GETTERS PARA GROUP USE CASES
// ========================================

// GetCreateGroupUseCase retorna o use case de criação de grupo
func (c *Container) GetCreateGroupUseCase() *usecase.CreateGroupUseCase {
	return c.groupUseCases.CreateGroup
}

// GetSetGroupPhotoUseCase retorna o use case de definição de foto do grupo
func (c *Container) GetSetGroupPhotoUseCase() *usecase.SetGroupPhotoUseCase {
	return c.groupUseCases.SetGroupPhoto
}

// GetUpdateGroupParticipantsUseCase retorna o use case de atualização de participantes
func (c *Container) GetUpdateGroupParticipantsUseCase() *usecase.UpdateGroupParticipantsUseCase {
	return c.groupUseCases.UpdateGroupParticipants
}

// GetLeaveGroupUseCase retorna o use case de saída do grupo
func (c *Container) GetLeaveGroupUseCase() *usecase.LeaveGroupUseCase {
	return c.groupUseCases.LeaveGroup
}

// GetJoinGroupUseCase retorna o use case de entrada no grupo
func (c *Container) GetJoinGroupUseCase() *usecase.JoinGroupUseCase {
	return c.groupUseCases.JoinGroup
}

// GetGetGroupInfoUseCase retorna o use case de obtenção de informações do grupo
func (c *Container) GetGetGroupInfoUseCase() *usecase.GetGroupInfoUseCase {
	return c.groupUseCases.GetGroupInfo
}

// GetListGroupsUseCase retorna o use case de listagem de grupos
func (c *Container) GetListGroupsUseCase() *usecase.ListGroupsUseCase {
	return c.groupUseCases.ListGroups
}

// GetGetGroupInviteLinkUseCase retorna o use case de obtenção de link de convite
func (c *Container) GetGetGroupInviteLinkUseCase() *usecase.GetGroupInviteLinkUseCase {
	return c.groupUseCases.GetGroupInviteLink
}

// GetRevokeGroupInviteLinkUseCase retorna o use case de revogação de link de convite
func (c *Container) GetRevokeGroupInviteLinkUseCase() *usecase.RevokeGroupInviteLinkUseCase {
	return c.groupUseCases.RevokeGroupInviteLink
}

// GetSetGroupNameUseCase retorna o use case de definição de nome do grupo
func (c *Container) GetSetGroupNameUseCase() *usecase.SetGroupNameUseCase {
	return c.groupUseCases.SetGroupName
}

// GetSetGroupTopicUseCase retorna o use case de definição de tópico do grupo
func (c *Container) GetSetGroupTopicUseCase() *usecase.SetGroupTopicUseCase {
	return c.groupUseCases.SetGroupTopic
}

// ========================================
// GETTERS PARA NEWSLETTER USE CASES
// ========================================

// GetListNewsletterUseCase retorna o use case de listagem de newsletters
func (c *Container) GetListNewsletterUseCase() *usecase.ListNewsletterUseCase {
	return c.newsletterUseCases.ListNewsletter
}
