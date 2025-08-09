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
	sessionUseCases *SessionUseCases
	messageUseCases *MessageUseCases

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
	// Envio de mensagens
	SendText  *usecase.SendTextMessageUseCase
	SendMedia *usecase.SendMediaMessageUseCase
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
		c.clientFactory = whatsapp.NewClientFactory(c.db.Store, c.GetSessionRepository())
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
