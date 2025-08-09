// Package usecase contém os casos de uso da camada de aplicação
// Este arquivo (session_connect.go) contém os use cases para:
// - Orquestração da conectividade e autenticação com WhatsApp
// - Coordenação do ciclo de vida da conexão
// - Operações: Connect, Logout, GetQR (autenticação)
// - Integração entre domain services e infraestrutura WhatsApp
package usecase

import (
	"context"
	"fmt"
	"time"

	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/entity"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/domain/service"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"
)

// ========================================
// SESSION CONNECT USE CASES
// ========================================
// Este arquivo agrupa os casos de uso para conectividade WhatsApp:
// 1. ConnectSessionUseCase - Estabelecer conexão com WhatsApp
// 2. LogoutSessionUseCase - Desconectar e fazer logout
// 3. GetQRCodeUseCase - Gerar QR code para autenticação
//
// Responsabilidades:
// - Gerenciar clientes WhatsApp (whatsmeow.Client)
// - Configurar event handlers para eventos de conexão
// - Controlar estados de conexão (connecting, connected, disconnected)
// - Autenticação via QR code
// ========================================

// ConnectSessionUseCase representa o caso de uso para conectar sessões
type ConnectSessionUseCase struct {
	sessionRepo    repository.SessionRepository
	clientFactory  *whatsapp.ClientFactory
	sessionManager *whatsapp.SessionManager
	domainService  *service.SessionDomainService
}

// NewConnectSessionUseCase cria uma nova instância do use case
func NewConnectSessionUseCase(sessionRepo repository.SessionRepository, clientFactory *whatsapp.ClientFactory, sessionManager *whatsapp.SessionManager, domainService *service.SessionDomainService) *ConnectSessionUseCase {
	return &ConnectSessionUseCase{
		sessionRepo:    sessionRepo,
		clientFactory:  clientFactory,
		sessionManager: sessionManager,
		domainService:  domainService,
	}
}

// Execute executa o caso de uso de conexão de sessão
func (uc *ConnectSessionUseCase) Execute(sessionID string) error {
	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return err
	}

	// Verificar se pode conectar usando regras de negócio
	if err := uc.domainService.CanConnect(session); err != nil {
		return err
	}

	// Verificar se já está conectada na infraestrutura
	if uc.sessionManager.IsConnected(session.ID) {
		return fmt.Errorf("sessão '%s' já está conectada", session.Name)
	}

	// Atualizar status para conectando
	session.Status = entity.StatusConnecting
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(session); err != nil {
		return fmt.Errorf("erro ao atualizar status da sessão: %w", err)
	}

	// Criar cliente usando o factory
	client, err := uc.clientFactory.CreateClient(session)
	if err != nil {
		session.Status = entity.StatusDisconnected
		session.UpdatedAt = time.Now()
		uc.sessionRepo.Update(session)
		return fmt.Errorf("erro ao criar cliente: %w", err)
	}

	// Conectar cliente usando ConnectWithQR
	// Não usar timeout para permitir que o loop QR continue ativo
	if err := client.ConnectWithQR(context.Background()); err != nil {
		session.Status = entity.StatusDisconnected
		session.UpdatedAt = time.Now()
		uc.sessionRepo.Update(session)
		return fmt.Errorf("erro ao conectar cliente: %w", err)
	}

	// Armazenar cliente no gerenciador de sessões
	uc.sessionManager.SetClient(session.ID, client)

	// NÃO atualizar status para connected aqui!
	// O status será atualizado para 'connected' apenas quando:
	// - handleConnected() for chamado (conexão estabelecida)
	// - handlePairSuccess() for chamado (pareamento bem-sucedido)
	// Por enquanto, manter status 'connecting' que foi definido anteriormente

	logger.Info("Sessão '%s' conectada com sucesso", session.Name)
	return nil
}

// findSession busca uma sessão por ID ou nome
func (uc *ConnectSessionUseCase) findSession(identifier string) (*entity.Session, error) {
	// Tentar buscar por ID primeiro
	session, err := uc.sessionRepo.GetByID(identifier)
	if err == nil {
		return session, nil
	}

	// Se não encontrou por ID, tentar por nome
	session, err = uc.sessionRepo.GetByName(identifier)
	if err != nil {
		return nil, fmt.Errorf("sessão '%s' não encontrada", identifier)
	}

	return session, nil
}

// LogoutSessionUseCase representa o caso de uso para fazer logout de sessões
type LogoutSessionUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewLogoutSessionUseCase cria uma nova instância do use case
func NewLogoutSessionUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *LogoutSessionUseCase {
	return &LogoutSessionUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa o caso de uso de logout de sessão
func (uc *LogoutSessionUseCase) Execute(sessionID string) error {
	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return err
	}

	// Obter cliente do gerenciador
	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return fmt.Errorf("sessão '%s' não possui cliente inicializado", session.Name)
	}

	// Verificar se está conectado
	if !client.IsConnected() {
		return fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	// Fazer logout
	if err := client.Logout(context.Background()); err != nil {
		return fmt.Errorf("erro ao fazer logout: %w", err)
	}

	// Remover cliente do gerenciador
	uc.sessionManager.RemoveClient(session.ID)

	// Atualizar status da sessão para desconectado (logout = desconectado)
	session.Status = entity.StatusDisconnected
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(session); err != nil {
		logger.Error("Erro ao atualizar sessão após logout: %v", err)
	}

	logger.Info("Logout realizado com sucesso para sessão '%s'", session.Name)
	return nil
}

// findSession busca uma sessão por ID ou nome
func (uc *LogoutSessionUseCase) findSession(identifier string) (*entity.Session, error) {
	// Tentar buscar por ID primeiro
	session, err := uc.sessionRepo.GetByID(identifier)
	if err == nil {
		return session, nil
	}

	// Se não encontrou por ID, tentar por nome
	session, err = uc.sessionRepo.GetByName(identifier)
	if err != nil {
		return nil, fmt.Errorf("sessão '%s' não encontrada", identifier)
	}

	return session, nil
}

// GetQRCodeUseCase representa o caso de uso para obter QR code
type GetQRCodeUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewGetQRCodeUseCase cria uma nova instância do use case
func NewGetQRCodeUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *GetQRCodeUseCase {
	return &GetQRCodeUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa o caso de uso de obtenção de QR code
func (uc *GetQRCodeUseCase) Execute(sessionID string) (*responses.QRResponse, error) {
	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return nil, err
	}

	// Verificar se cliente existe
	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return &responses.QRResponse{
			Status: "session_not_connected",
		}, nil
	}

	// Verificar se já está logado
	if client.IsLoggedIn() {
		return &responses.QRResponse{
			Status: "already_logged_in",
		}, nil
	}

	// Verificar se QR está ativo
	if !client.IsQRActive() {
		return &responses.QRResponse{
			Status: "qr_not_active",
		}, nil
	}

	// Obter QR code atual do canal interno (compatibilidade)
	qrChannel := client.GetQRChannel()

	// Tentar obter QR code do canal (não bloqueante)
	select {
	case qrCode := <-qrChannel:
		return &responses.QRResponse{
			QRCode: qrCode,
			Status: "qr_active",
		}, nil
	default:
		// Nenhum QR disponível no momento
		return &responses.QRResponse{
			Status: "qr_pending",
		}, nil
	}
}

// findSession busca uma sessão por ID ou nome
func (uc *GetQRCodeUseCase) findSession(identifier string) (*entity.Session, error) {
	// Tentar buscar por ID primeiro
	session, err := uc.sessionRepo.GetByID(identifier)
	if err == nil {
		return session, nil
	}

	// Se não encontrou por ID, tentar por nome
	session, err = uc.sessionRepo.GetByName(identifier)
	if err != nil {
		return nil, fmt.Errorf("sessão '%s' não encontrada", identifier)
	}

	return session, nil
}
