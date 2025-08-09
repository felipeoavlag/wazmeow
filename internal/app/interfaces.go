package app

import (
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
	// Configuração
	GetConfig() *config.Config

	// Infraestrutura
	GetDB() *database.Connection
	GetBunDB() *database.BunConnection

	// Repositories
	GetSessionRepository() repository.SessionRepository

	// Domain Services
	GetSessionDomainService() *service.SessionDomainService

	// Session Manager
	GetSessionManager() *whatsapp.SessionManager

	// Use Cases - Session Management
	GetCreateSessionUseCase() *usecase.CreateSessionUseCase
	GetListSessionsUseCase() *usecase.ListSessionsUseCase
	GetDeleteSessionUseCase() *usecase.DeleteSessionUseCase
	GetGetSessionInfoUseCase() *usecase.GetSessionInfoUseCase

	// Use Cases - Session Connection
	GetConnectSessionUseCase() *usecase.ConnectSessionUseCase
	GetLogoutSessionUseCase() *usecase.LogoutSessionUseCase
	GetGetQRCodeUseCase() *usecase.GetQRCodeUseCase

	// Use Cases - Session Setup
	GetPairPhoneUseCase() *usecase.PairPhoneUseCase
	GetSetProxyUseCase() *usecase.SetProxyUseCase

	// Lifecycle
	Close() error
}

// Verificar se DependencyContainer implementa DIContainer
var _ DIContainer = (*DependencyContainer)(nil)
