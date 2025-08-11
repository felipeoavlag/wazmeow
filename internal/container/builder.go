package container

import (
	"context"
	"fmt"

	"wazmeow/internal/application/usecase"
	"wazmeow/internal/config"
	"wazmeow/internal/domain/service"
	"wazmeow/internal/infra/database"
	infraRepo "wazmeow/internal/infra/repository"
	"wazmeow/internal/infra/webhook"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"
)

// Builder implementa o padrão Builder para construção do container
type Builder struct {
	config *config.Config
}

// NewBuilder cria um novo builder para o container
func NewBuilder() *Builder {
	return &Builder{}
}

// WithConfig define uma configuração específica para o builder
func (b *Builder) WithConfig(cfg *config.Config) *Builder {
	b.config = cfg
	return b
}

// Build constrói o container com todas as dependências
func (b *Builder) Build() (*Container, error) {
	container := &Container{}

	// Etapas de inicialização em ordem específica
	steps := []struct {
		name string
		fn   func(*Container) error
	}{
		{"configuração", b.setupConfig},
		{"infraestrutura", b.setupInfrastructure},
		{"repositories", b.setupRepositories},
		{"domain services", b.setupDomainServices},
		{"use cases", b.setupUseCases},
	}

	logger.Info("🚀 Iniciando inicialização do container...")

	for _, step := range steps {
		logger.Debug("📦 Configurando %s...", step.name)
		if err := step.fn(container); err != nil {
			return nil, fmt.Errorf("erro ao configurar %s: %w", step.name, err)
		}
		logger.Debug("✅ %s configurado com sucesso", step.name)
	}

	container.setInitialized(true)
	logger.Info("🎉 Container inicializado com sucesso!")
	return container, nil
}

// setupConfig carrega e valida a configuração
func (b *Builder) setupConfig(container *Container) error {
	var cfg *config.Config
	var err error

	if b.config != nil {
		cfg = b.config
		logger.Debug("Usando configuração fornecida")
	} else {
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("erro ao carregar configuração: %w", err)
		}
		logger.Debug("Configuração carregada das variáveis de ambiente")
	}

	// Validar configuração
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuração inválida: %w", err)
	}

	// Inicializar logger global com a configuração
	logger.InitGlobalLogger(cfg.Log.Level)

	container.config = cfg
	return nil
}

// setupInfrastructure configura as conexões de banco e infraestrutura
func (b *Builder) setupInfrastructure(container *Container) error {
	// Configuração do banco de dados
	dbConfig := database.Config{
		Host:     container.config.Database.Host,
		Port:     container.config.Database.Port,
		User:     container.config.Database.User,
		Password: container.config.Database.Password,
		Name:     container.config.Database.Name,
		SSLMode:  container.config.Database.SSLMode,
		Debug:    container.config.Database.Debug,
	}

	// Conectar ao banco de dados para WhatsApp (sqlstore)
	dbConnection, err := database.Connect(dbConfig)
	if err != nil {
		return fmt.Errorf("erro ao conectar com banco para WhatsApp: %w", err)
	}

	// Executar migrações do WhatsApp
	if err := dbConnection.Migrate(); err != nil {
		return fmt.Errorf("erro ao executar migrações do WhatsApp: %w", err)
	}

	// Conectar ao banco de dados com Bun ORM
	bunConnection, err := database.NewBunConnection(dbConfig)
	if err != nil {
		return fmt.Errorf("erro ao conectar com Bun: %w", err)
	}

	// Garantir que o schema esteja atualizado
	ctx := context.Background()
	if err := bunConnection.EnsureSchema(ctx); err != nil {
		return fmt.Errorf("erro ao garantir schema atualizado: %w", err)
	}

	// Instanciar session manager
	sessionManager := whatsapp.NewSessionManager()

	// Instanciar webhook service
	webhookService := webhook.NewWebhookService(&container.config.Webhook)

	container.db = dbConnection
	container.bunDB = bunConnection
	container.sessionManager = sessionManager
	container.webhookService = webhookService

	return nil
}

// setupRepositories configura todos os repositories
func (b *Builder) setupRepositories(container *Container) error {
	// Repository de sessões usando Bun ORM
	sessionRepo := infraRepo.NewBunSessionRepository(container.bunDB.GetDB())

	container.sessionRepo = sessionRepo
	return nil
}

// setupDomainServices configura os domain services
func (b *Builder) setupDomainServices(container *Container) error {
	// Domain service de sessões
	sessionDomainService := service.NewSessionDomainService()

	container.sessionDomainService = sessionDomainService
	return nil
}

// setupUseCases configura todos os use cases
func (b *Builder) setupUseCases(container *Container) error {
	// Instanciar use cases organizados por categoria
	sessionUseCases := &SessionUseCases{
		// Gerenciamento básico
		Create:  usecase.NewCreateSessionUseCase(container.sessionRepo, container.sessionDomainService, container.config),
		List:    usecase.NewListSessionsUseCase(container.sessionRepo),
		Delete:  usecase.NewDeleteSessionUseCase(container.sessionRepo, container.sessionDomainService),
		GetInfo: usecase.NewGetSessionInfoUseCase(container.sessionRepo),

		// Conectividade
		Connect: usecase.NewConnectSessionUseCase(container.sessionRepo, container.GetClientFactory(), container.sessionManager, container.sessionDomainService),
		Logout:  usecase.NewLogoutSessionUseCase(container.sessionRepo, container.sessionManager),
		GetQR:   usecase.NewGetQRCodeUseCase(container.sessionRepo, container.sessionManager),

		// Setup e configuração
		PairPhone: usecase.NewPairPhoneUseCase(container.sessionRepo, container.sessionDomainService),
		SetProxy:  usecase.NewSetProxyUseCase(container.sessionRepo, container.sessionDomainService),
	}

	container.sessionUseCases = sessionUseCases

	// Instanciar use cases de mensagem
	messageUseCases := &MessageUseCases{
		// Envio de mensagens básicas
		SendText:  usecase.NewSendTextMessageUseCase(container.sessionRepo, container.sessionManager),
		SendMedia: usecase.NewSendMediaMessageUseCase(container.sessionRepo, container.sessionManager),

		// Envio de mensagens específicas
		SendImage:    usecase.NewSendImageMessageUseCase(container.sessionRepo, container.sessionManager),
		SendAudio:    usecase.NewSendAudioMessageUseCase(container.sessionRepo, container.sessionManager),
		SendDocument: usecase.NewSendDocumentMessageUseCase(container.sessionRepo, container.sessionManager),
		SendVideo:    usecase.NewSendVideoMessageUseCase(container.sessionRepo, container.sessionManager),
		SendSticker:  usecase.NewSendStickerMessageUseCase(container.sessionRepo, container.sessionManager),
		SendLocation: usecase.NewSendLocationMessageUseCase(container.sessionRepo, container.sessionManager),
		SendContact:  usecase.NewSendContactMessageUseCase(container.sessionRepo, container.sessionManager),
		SendButtons:  usecase.NewSendButtonsMessageUseCase(container.sessionRepo, container.sessionManager),
		SendList:     usecase.NewSendListMessageUseCase(container.sessionRepo, container.sessionManager),
		SendPoll:     usecase.NewSendPollMessageUseCase(container.sessionRepo, container.sessionManager),

		// Operações de mensagem
		SendEdit:      usecase.NewSendEditMessageUseCase(container.sessionRepo, container.sessionManager),
		DeleteMessage: usecase.NewDeleteMessageUseCase(container.sessionRepo, container.sessionManager),
		React:         usecase.NewReactMessageUseCase(container.sessionRepo, container.sessionManager),
	}

	container.messageUseCases = messageUseCases

	// Instanciar use cases de webhook
	webhookUseCases := &WebhookUseCases{
		SetWebhook:    usecase.NewSetWebhookUseCase(container.sessionRepo),
		GetWebhook:    usecase.NewGetWebhookUseCase(container.sessionRepo),
		UpdateWebhook: usecase.NewUpdateWebhookUseCase(container.sessionRepo),
		DeleteWebhook: usecase.NewDeleteWebhookUseCase(container.sessionRepo),
	}

	container.webhookUseCases = webhookUseCases

	// Instanciar use cases de usuário
	userUseCases := &UserUseCases{
		GetUserInfo: usecase.NewGetUserInfoUseCase(container.sessionRepo, container.sessionManager),
		CheckUser:   usecase.NewCheckUserUseCase(container.sessionRepo, container.sessionManager),
		GetAvatar:   usecase.NewGetAvatarUseCase(container.sessionRepo, container.sessionManager),
		GetContacts: usecase.NewGetContactsUseCase(container.sessionRepo, container.sessionManager),
	}

	container.userUseCases = userUseCases

	// Instanciar use cases de chat
	chatUseCases := &ChatUseCases{
		SendPresence:     usecase.NewSendPresenceUseCase(container.sessionRepo, container.sessionManager),
		ChatPresence:     usecase.NewChatPresenceUseCase(container.sessionRepo, container.sessionManager),
		MarkRead:         usecase.NewMarkReadUseCase(container.sessionRepo, container.sessionManager),
		DownloadImage:    usecase.NewDownloadImageUseCase(container.sessionRepo, container.sessionManager),
		DownloadVideo:    usecase.NewDownloadVideoUseCase(container.sessionRepo, container.sessionManager),
		DownloadAudio:    usecase.NewDownloadAudioUseCase(container.sessionRepo, container.sessionManager),
		DownloadDocument: usecase.NewDownloadDocumentUseCase(container.sessionRepo, container.sessionManager),
	}

	container.chatUseCases = chatUseCases

	// Instanciar use cases de grupo
	groupUseCases := &GroupUseCases{
		CreateGroup:             usecase.NewCreateGroupUseCase(container.sessionRepo, container.sessionManager),
		SetGroupPhoto:           usecase.NewSetGroupPhotoUseCase(container.sessionRepo, container.sessionManager),
		UpdateGroupParticipants: usecase.NewUpdateGroupParticipantsUseCase(container.sessionRepo, container.sessionManager),
		LeaveGroup:              usecase.NewLeaveGroupUseCase(container.sessionRepo, container.sessionManager),
		JoinGroup:               usecase.NewJoinGroupUseCase(container.sessionRepo, container.sessionManager),
		GetGroupInfo:            usecase.NewGetGroupInfoUseCase(container.sessionRepo, container.sessionManager),
		ListGroups:              usecase.NewListGroupsUseCase(container.sessionRepo, container.sessionManager),
		GetGroupInviteLink:      usecase.NewGetGroupInviteLinkUseCase(container.sessionRepo, container.sessionManager),
		RevokeGroupInviteLink:   usecase.NewRevokeGroupInviteLinkUseCase(container.sessionRepo, container.sessionManager),
		SetGroupName:            usecase.NewSetGroupNameUseCase(container.sessionRepo, container.sessionManager),
		SetGroupTopic:           usecase.NewSetGroupTopicUseCase(container.sessionRepo, container.sessionManager),
	}

	container.groupUseCases = groupUseCases

	// Instanciar use cases de newsletter
	newsletterUseCases := &NewsletterUseCases{
		ListNewsletter: usecase.NewListNewsletterUseCase(container.sessionRepo, container.sessionManager),
	}

	container.newsletterUseCases = newsletterUseCases

	return nil
}
