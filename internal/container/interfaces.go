package container

import (
	"context"

	"wazmeow/internal/application/usecase"
	"wazmeow/internal/config"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/domain/service"
	"wazmeow/internal/infra/database"
	"wazmeow/internal/infra/whatsapp"
)

// DIContainer define a interface para o container de injeção de dependências
// Esta interface facilita testes e permite diferentes implementações do container
type DIContainer interface {
	// ========================================
	// LIFECYCLE
	// ========================================
	IsInitialized() bool
	HealthCheck(ctx context.Context) error
	Close() error
	GetStatus() map[string]interface{}

	// ========================================
	// CONFIGURAÇÃO
	// ========================================
	GetConfig() *config.Config

	// ========================================
	// INFRAESTRUTURA
	// ========================================
	GetDB() *database.Connection
	GetBunDB() *database.BunConnection

	// ========================================
	// REPOSITORIES
	// ========================================
	GetSessionRepository() repository.SessionRepository

	// ========================================
	// DOMAIN SERVICES
	// ========================================
	GetSessionDomainService() *service.SessionDomainService

	// ========================================
	// SESSION MANAGER
	// ========================================
	GetSessionManager() *whatsapp.SessionManager

	// ========================================
	// USE CASES - SESSION MANAGEMENT
	// ========================================
	GetCreateSessionUseCase() *usecase.CreateSessionUseCase
	GetListSessionsUseCase() *usecase.ListSessionsUseCase
	GetDeleteSessionUseCase() *usecase.DeleteSessionUseCase
	GetGetSessionInfoUseCase() *usecase.GetSessionInfoUseCase

	// ========================================
	// USE CASES - SESSION CONNECTION
	// ========================================
	GetConnectSessionUseCase() *usecase.ConnectSessionUseCase
	GetLogoutSessionUseCase() *usecase.LogoutSessionUseCase
	GetGetQRCodeUseCase() *usecase.GetQRCodeUseCase

	// ========================================
	// USE CASES - SESSION SETUP
	// ========================================
	GetPairPhoneUseCase() *usecase.PairPhoneUseCase
	GetSetProxyUseCase() *usecase.SetProxyUseCase
}

// Verificar se Container implementa DIContainer em tempo de compilação
var _ DIContainer = (*Container)(nil)
