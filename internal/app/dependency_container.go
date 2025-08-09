// Package app contém o container de injeção de dependências da aplicação
//
// Este pacote é responsável por:
// - Configurar e instanciar todas as dependências da aplicação
// - Seguir os princípios de Clean Architecture na organização das camadas
// - Implementar o padrão Dependency Injection Container
// - Fornecer uma interface clara para acesso às dependências
//
// Estrutura das dependências:
// 1. Configuração (config)
// 2. Infraestrutura (database, whatsapp)
// 3. Repositories (implementações de persistência)
// 4. Domain Services (regras de negócio puras)
// 5. Use Cases (orquestração da aplicação)
package app

import (
	"context"
	"fmt"

	"wazmeow/internal/application/usecase"
	"wazmeow/internal/config"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/domain/service"
	"wazmeow/internal/infra/database"
	infraRepo "wazmeow/internal/infra/repository"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"
)

// DependencyContainer representa o container de injeção de dependências
//
// RESPONSABILIDADES:
// - Configurar e instanciar todas as dependências da aplicação
// - Garantir que as dependências sejam criadas na ordem correta
// - Implementar o padrão Dependency Injection Container
// - Fornecer acesso controlado às dependências através da interface DIContainer
//
// PRINCÍPIOS SEGUIDOS:
// - Clean Architecture: separação clara entre camadas
// - Dependency Inversion: dependências apontam para abstrações
// - Single Responsibility: cada dependência tem uma responsabilidade específica
// - Interface Segregation: interface DIContainer fornece acesso granular
//
// ORDEM DE INICIALIZAÇÃO:
// 1. Configuração (config.Load)
// 2. Infraestrutura (database, logger)
// 3. Repositories (implementações de persistência)
// 4. Domain Services (regras de negócio)
// 5. Use Cases (orquestração)
type DependencyContainer struct {
	// ========================================
	// CONFIGURAÇÃO
	// ========================================
	Config *config.Config // Configurações da aplicação carregadas do ambiente

	// ========================================
	// INFRAESTRUTURA
	// ========================================
	DB    *database.Connection    // Conexão para WhatsApp (sqlstore)
	BunDB *database.BunConnection // Conexão Bun ORM para sessões da aplicação

	// ========================================
	// REPOSITORIES (Camada de Persistência)
	// ========================================
	SessionRepo repository.SessionRepository // Interface para persistência de sessões

	// ========================================
	// INFRAESTRUTURA DE DOMÍNIO
	// ========================================
	SessionManager *whatsapp.SessionManager // Gerenciador de clientes WhatsApp

	// ========================================
	// USE CASES (Camada de Aplicação)
	// ========================================
	// Gerenciamento básico de sessões
	CreateSessionUC  *usecase.CreateSessionUseCase  // Criação de novas sessões
	ListSessionsUC   *usecase.ListSessionsUseCase   // Listagem de sessões
	DeleteSessionUC  *usecase.DeleteSessionUseCase  // Remoção de sessões
	GetSessionInfoUC *usecase.GetSessionInfoUseCase // Informações detalhadas

	// Conectividade WhatsApp
	ConnectSessionUC *usecase.ConnectSessionUseCase // Estabelecer conexão
	LogoutSessionUC  *usecase.LogoutSessionUseCase  // Desconectar e logout
	GetQRCodeUC      *usecase.GetQRCodeUseCase      // Gerar QR code

	// Configuração e setup
	PairPhoneUC *usecase.PairPhoneUseCase // Emparelhamento por telefone
	SetProxyUC  *usecase.SetProxyUseCase  // Configuração de proxy
}

// NewDependencyContainer cria e configura um novo container com todas as dependências
// Este é o ponto central de configuração da aplicação onde todas as dependências
// são instanciadas e conectadas seguindo os princípios de Clean Architecture
func NewDependencyContainer() (*DependencyContainer, error) {
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
	sessionRepo := infraRepo.NewBunSessionRepository(bunConnection.GetDB())

	// Instanciar session manager
	sessionManager := whatsapp.NewSessionManager()

	// Instanciar domain services
	sessionDomainService := service.NewSessionDomainService()

	// Instanciar use cases
	createSessionUC := usecase.NewCreateSessionUseCase(sessionRepo, sessionDomainService)
	connectSessionUC := usecase.NewConnectSessionUseCase(sessionRepo, dbConnection.Store, sessionManager, sessionDomainService)
	listSessionsUC := usecase.NewListSessionsUseCase(sessionRepo)
	getQRCodeUC := usecase.NewGetQRCodeUseCase(sessionRepo, sessionManager)
	deleteSessionUC := usecase.NewDeleteSessionUseCase(sessionRepo, sessionDomainService)
	logoutSessionUC := usecase.NewLogoutSessionUseCase(sessionRepo, sessionManager)
	pairPhoneUC := usecase.NewPairPhoneUseCase(sessionRepo, sessionDomainService)
	getSessionInfoUC := usecase.NewGetSessionInfoUseCase(sessionRepo)
	setProxyUC := usecase.NewSetProxyUseCase(sessionRepo, sessionDomainService)

	return &DependencyContainer{
		Config: cfg,
		DB:     dbConnection,  // WhatsApp sqlstore
		BunDB:  bunConnection, // Bun ORM para sessões

		// Repositories
		SessionRepo: sessionRepo,

		// Session Manager
		SessionManager: sessionManager,

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
func (c *DependencyContainer) Close() error {
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

// ========================================
// IMPLEMENTAÇÃO DA INTERFACE DIContainer
// ========================================

// GetConfig retorna a configuração da aplicação
func (c *DependencyContainer) GetConfig() *config.Config {
	return c.Config
}

// GetDB retorna a conexão do banco de dados para WhatsApp
func (c *DependencyContainer) GetDB() *database.Connection {
	return c.DB
}

// GetBunDB retorna a conexão Bun ORM
func (c *DependencyContainer) GetBunDB() *database.BunConnection {
	return c.BunDB
}

// GetSessionRepository retorna o repository de sessões
func (c *DependencyContainer) GetSessionRepository() repository.SessionRepository {
	return c.SessionRepo
}

// GetSessionDomainService retorna o domain service de sessões
func (c *DependencyContainer) GetSessionDomainService() *service.SessionDomainService {
	// Por enquanto, criar uma nova instância
	// Em implementações futuras, pode ser cached
	return service.NewSessionDomainService()
}

// GetSessionManager retorna o gerenciador de sessões
func (c *DependencyContainer) GetSessionManager() *whatsapp.SessionManager {
	return c.SessionManager
}

// GetCreateSessionUseCase retorna o use case de criação de sessões
func (c *DependencyContainer) GetCreateSessionUseCase() *usecase.CreateSessionUseCase {
	return c.CreateSessionUC
}

// GetListSessionsUseCase retorna o use case de listagem de sessões
func (c *DependencyContainer) GetListSessionsUseCase() *usecase.ListSessionsUseCase {
	return c.ListSessionsUC
}

// GetDeleteSessionUseCase retorna o use case de deleção de sessões
func (c *DependencyContainer) GetDeleteSessionUseCase() *usecase.DeleteSessionUseCase {
	return c.DeleteSessionUC
}

// GetGetSessionInfoUseCase retorna o use case de informações de sessão
func (c *DependencyContainer) GetGetSessionInfoUseCase() *usecase.GetSessionInfoUseCase {
	return c.GetSessionInfoUC
}

// GetConnectSessionUseCase retorna o use case de conexão de sessões
func (c *DependencyContainer) GetConnectSessionUseCase() *usecase.ConnectSessionUseCase {
	return c.ConnectSessionUC
}

// GetLogoutSessionUseCase retorna o use case de logout de sessões
func (c *DependencyContainer) GetLogoutSessionUseCase() *usecase.LogoutSessionUseCase {
	return c.LogoutSessionUC
}

// GetGetQRCodeUseCase retorna o use case de QR code
func (c *DependencyContainer) GetGetQRCodeUseCase() *usecase.GetQRCodeUseCase {
	return c.GetQRCodeUC
}

// GetPairPhoneUseCase retorna o use case de emparelhamento por telefone
func (c *DependencyContainer) GetPairPhoneUseCase() *usecase.PairPhoneUseCase {
	return c.PairPhoneUC
}

// GetSetProxyUseCase retorna o use case de configuração de proxy
func (c *DependencyContainer) GetSetProxyUseCase() *usecase.SetProxyUseCase {
	return c.SetProxyUC
}
