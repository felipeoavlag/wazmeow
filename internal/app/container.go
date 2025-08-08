package app

import (
	"context"
	"fmt"

	"wazmeow/internal/config"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/domain/usecases"
	"wazmeow/internal/infra/database"
	"wazmeow/internal/infra/repository"
	"wazmeow/pkg/logger"
)

// Container representa o container de injeção de dependências
type Container struct {
	// Configuração
	Config *config.Config

	// Infraestrutura
	DB    *database.Connection    // Para WhatsApp (sqlstore)
	BunDB *database.BunConnection // Para sessões da aplicação (Bun ORM)

	// Repositories
	SessionRepo repositories.SessionRepository

	// Use Cases
	CreateSessionUC  *usecases.CreateSessionUseCase
	ConnectSessionUC *usecases.ConnectSessionUseCase
	ListSessionsUC   *usecases.ListSessionsUseCase
	GetQRCodeUC      *usecases.GetQRCodeUseCase
	DeleteSessionUC  *usecases.DeleteSessionUseCase
	LogoutSessionUC  *usecases.LogoutSessionUseCase
	PairPhoneUC      *usecases.PairPhoneUseCase
	GetSessionInfoUC *usecases.GetSessionInfoUseCase
	SetProxyUC       *usecases.SetProxyUseCase
}

// NewContainer cria e configura um novo container com todas as dependências
func NewContainer() (*Container, error) {
	// Carregar configuração
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Validar configuração
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuração inválida: %w", err)
	}

	// Inicializar logger global
	logger.InitGlobalLogger(cfg.Log.Level)

	// Criar configuração do banco de dados
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
		Debug:    cfg.Database.Debug,
	}

	// Conectar ao banco de dados para WhatsApp (sqlstore)
	dbConnection, err := database.Connect(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com banco para WhatsApp: %w", err)
	}

	// Executar migrações do WhatsApp
	if err := dbConnection.Migrate(); err != nil {
		return nil, fmt.Errorf("erro ao executar migrações do WhatsApp: %w", err)
	}

	// Conectar ao banco de dados com Bun ORM para sessões da aplicação
	bunConnection, err := database.NewBunConnection(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com Bun: %w", err)
	}

	// Executar auto-migrações do Bun
	ctx := context.Background()
	if err := bunConnection.AutoMigrate(ctx); err != nil {
		return nil, fmt.Errorf("erro nas auto-migrações do Bun: %w", err)
	}

	// Instanciar repositories usando Bun
	sessionRepo := repository.NewBunSessionRepository(bunConnection.GetDB())

	// Instanciar use cases
	createSessionUC := usecases.NewCreateSessionUseCase(sessionRepo)
	connectSessionUC := usecases.NewConnectSessionUseCase(sessionRepo, dbConnection.Store)
	listSessionsUC := usecases.NewListSessionsUseCase(sessionRepo)
	getQRCodeUC := usecases.NewGetQRCodeUseCase(sessionRepo)
	deleteSessionUC := usecases.NewDeleteSessionUseCase(sessionRepo)
	logoutSessionUC := usecases.NewLogoutSessionUseCase(sessionRepo)
	pairPhoneUC := usecases.NewPairPhoneUseCase(sessionRepo)
	getSessionInfoUC := usecases.NewGetSessionInfoUseCase(sessionRepo)
	setProxyUC := usecases.NewSetProxyUseCase(sessionRepo)

	return &Container{
		Config: cfg,
		DB:     dbConnection,  // WhatsApp sqlstore
		BunDB:  bunConnection, // Bun ORM para sessões

		// Repositories
		SessionRepo: sessionRepo,

		// Use Cases
		CreateSessionUC:  createSessionUC,
		ConnectSessionUC: connectSessionUC,
		ListSessionsUC:   listSessionsUC,
		GetQRCodeUC:      getQRCodeUC,
		DeleteSessionUC:  deleteSessionUC,
		LogoutSessionUC:  logoutSessionUC,
		PairPhoneUC:      pairPhoneUC,
		GetSessionInfoUC: getSessionInfoUC,
		SetProxyUC:       setProxyUC,
	}, nil
}

// Close fecha todas as conexões e recursos do container
func (c *Container) Close() error {
	var err error

	// Fechar conexão Bun
	if c.BunDB != nil {
		if closeErr := c.BunDB.Close(); closeErr != nil {
			err = closeErr
		}
	}

	// Fechar conexão WhatsApp
	if c.DB != nil {
		if closeErr := c.DB.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			}
		}
	}

	return err
}
