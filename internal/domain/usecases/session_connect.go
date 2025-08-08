// Package usecases contém os casos de uso da aplicação
// Este arquivo (session_connect.go) contém os use cases para:
// - Conectividade e autenticação com WhatsApp
// - Gerenciamento do ciclo de vida da conexão
// - Operações: Connect, Logout, GetQR (autenticação)
// - Integração direta com a biblioteca whatsmeow
package usecases

import (
	"context"
	"fmt"
	"time"

	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/domain/responses"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
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
	sessionRepo repositories.SessionRepository
	deviceStore *sqlstore.Container
}

// NewConnectSessionUseCase cria uma nova instância do use case
func NewConnectSessionUseCase(sessionRepo repositories.SessionRepository, deviceStore *sqlstore.Container) *ConnectSessionUseCase {
	return &ConnectSessionUseCase{
		sessionRepo: sessionRepo,
		deviceStore: deviceStore,
	}
}

// Execute executa o caso de uso de conexão de sessão
func (uc *ConnectSessionUseCase) Execute(sessionID string) error {
	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return err
	}

	// Verificar se já está conectada
	if session.Client != nil && session.Client.IsConnected() {
		return fmt.Errorf("sessão '%s' já está conectada", session.Name)
	}

	// Atualizar status para conectando
	session.Status = entities.StatusConnecting
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(session); err != nil {
		return fmt.Errorf("erro ao atualizar status da sessão: %w", err)
	}

	// Criar device store para a sessão
	deviceStore := uc.deviceStore.NewDevice()

	// Criar cliente WhatsApp
	client := whatsmeow.NewClient(deviceStore, logger.ForWhatsApp())

	// Configurar handlers de eventos
	uc.setupEventHandlers(client, session)

	// Conectar cliente
	if err := client.Connect(); err != nil {
		session.Status = entities.StatusDisconnected
		session.UpdatedAt = time.Now()
		uc.sessionRepo.Update(session)
		return fmt.Errorf("erro ao conectar cliente: %w", err)
	}

	// Atualizar sessão com cliente
	session.Client = client
	session.Status = entities.StatusConnected
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(session); err != nil {
		logger.Error("Erro ao atualizar sessão após conexão: %v", err)
	}

	logger.Info("Sessão '%s' conectada com sucesso", session.Name)
	return nil
}

// findSession busca uma sessão por ID ou nome
func (uc *ConnectSessionUseCase) findSession(identifier string) (*entities.Session, error) {
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

// setupEventHandlers configura os handlers de eventos do WhatsApp
func (uc *ConnectSessionUseCase) setupEventHandlers(client *whatsmeow.Client, session *entities.Session) {
	client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.Connected:
			session.Status = entities.StatusConnected
			session.UpdatedAt = time.Now()
			uc.sessionRepo.Update(session)
			logger.Info("Sessão '%s' conectada ao WhatsApp", session.Name)

		case *events.Disconnected:
			session.Status = entities.StatusDisconnected
			session.UpdatedAt = time.Now()
			uc.sessionRepo.Update(session)
			logger.Info("Sessão '%s' desconectada do WhatsApp", session.Name)

		case *events.LoggedOut:
			session.Status = entities.StatusLoggedOut
			session.UpdatedAt = time.Now()
			uc.sessionRepo.Update(session)
			logger.Info("Sessão '%s' fez logout do WhatsApp", session.Name)
		}
	})
}

// LogoutSessionUseCase representa o caso de uso para fazer logout de sessões
type LogoutSessionUseCase struct {
	sessionRepo repositories.SessionRepository
}

// NewLogoutSessionUseCase cria uma nova instância do use case
func NewLogoutSessionUseCase(sessionRepo repositories.SessionRepository) *LogoutSessionUseCase {
	return &LogoutSessionUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de logout de sessão
func (uc *LogoutSessionUseCase) Execute(sessionID string) error {
	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return err
	}

	// Verificar se o cliente existe
	if session.Client == nil {
		return fmt.Errorf("sessão '%s' não possui cliente inicializado", session.Name)
	}

	// Verificar se está conectado
	if !session.Client.IsConnected() {
		return fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	// Fazer logout
	if err := session.Client.Logout(context.Background()); err != nil {
		return fmt.Errorf("erro ao fazer logout: %w", err)
	}

	// Atualizar status da sessão
	session.Status = entities.StatusLoggedOut
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(session); err != nil {
		logger.Error("Erro ao atualizar sessão após logout: %v", err)
	}

	logger.Info("Logout realizado com sucesso para sessão '%s'", session.Name)
	return nil
}

// findSession busca uma sessão por ID ou nome
func (uc *LogoutSessionUseCase) findSession(identifier string) (*entities.Session, error) {
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
	sessionRepo repositories.SessionRepository
}

// NewGetQRCodeUseCase cria uma nova instância do use case
func NewGetQRCodeUseCase(sessionRepo repositories.SessionRepository) *GetQRCodeUseCase {
	return &GetQRCodeUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de obtenção de QR code
func (uc *GetQRCodeUseCase) Execute(sessionID string) (*responses.QRResponse, error) {
	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return nil, err
	}

	// Verificar se o cliente existe
	if session.Client == nil {
		return nil, fmt.Errorf("sessão '%s' não está inicializada", session.Name)
	}

	// Verificar se já está logado
	if session.Client.IsLoggedIn() {
		return &responses.QRResponse{
			Status: "already_logged_in",
		}, nil
	}

	// Criar canal para receber QR code
	qrChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	// Configurar handler para QR code
	eventHandler := func(evt interface{}) {
		switch e := evt.(type) {
		case *events.QR:
			select {
			case qrChan <- e.Codes[0]:
			default:
			}
		case *events.PairSuccess:
			logger.Info("Emparelhamento bem-sucedido para sessão '%s'", session.Name)
		}
	}

	// Adicionar handler temporário
	handlerID := session.Client.AddEventHandler(eventHandler)
	defer session.Client.RemoveEventHandler(handlerID)

	// Solicitar QR code
	if err := session.Client.Connect(); err != nil {
		return nil, fmt.Errorf("erro ao conectar para obter QR: %w", err)
	}

	// Aguardar QR code com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	select {
	case qrCode := <-qrChan:
		return &responses.QRResponse{
			QRCode: qrCode,
			Status: "qr_generated",
		}, nil
	case err := <-errorChan:
		return nil, fmt.Errorf("erro ao gerar QR: %w", err)
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout ao aguardar QR code")
	}
}

// findSession busca uma sessão por ID ou nome
func (uc *GetQRCodeUseCase) findSession(identifier string) (*entities.Session, error) {
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
