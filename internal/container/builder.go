package container

import (
	"context"
	"fmt"

	"wazmeow/internal/application/usecase"
	"wazmeow/internal/config"
	"wazmeow/internal/domain/service"
	"wazmeow/internal/infra/database"
	infraRepo "wazmeow/internal/infra/repository"
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

	// Executar auto-migrações do Bun
	ctx := context.Background()
	if err := bunConnection.AutoMigrate(ctx); err != nil {
		return fmt.Errorf("erro nas auto-migrações do Bun: %w", err)
	}

	// Instanciar session manager
	sessionManager := whatsapp.NewSessionManager()

	container.db = dbConnection
	container.bunDB = bunConnection
	container.sessionManager = sessionManager

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
		Create:  usecase.NewCreateSessionUseCase(container.sessionRepo, container.sessionDomainService),
		List:    usecase.NewListSessionsUseCase(container.sessionRepo),
		Delete:  usecase.NewDeleteSessionUseCase(container.sessionRepo, container.sessionDomainService),
		GetInfo: usecase.NewGetSessionInfoUseCase(container.sessionRepo),

		// Conectividade
		Connect: usecase.NewConnectSessionUseCase(container.sessionRepo, container.db.Store, container.sessionManager, container.sessionDomainService),
		Logout:  usecase.NewLogoutSessionUseCase(container.sessionRepo, container.sessionManager),
		GetQR:   usecase.NewGetQRCodeUseCase(container.sessionRepo, container.sessionManager),

		// Setup e configuração
		PairPhone: usecase.NewPairPhoneUseCase(container.sessionRepo, container.sessionDomainService),
		SetProxy:  usecase.NewSetProxyUseCase(container.sessionRepo, container.sessionDomainService),
	}

	container.sessionUseCases = sessionUseCases
	return nil
}
