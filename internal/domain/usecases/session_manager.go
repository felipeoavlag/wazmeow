// Package usecases contém os casos de uso da aplicação
// Este arquivo (session_manager.go) contém os use cases para:
// - Gerenciamento básico do ciclo de vida das sessões (CRUD)
// - Operações fundamentais: Create, List, Delete, GetInfo
// - Não contém lógica de conexão WhatsApp ou configurações avançadas
package usecases

import (
	"fmt"
	"time"

	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/domain/requests"
	"wazmeow/internal/domain/responses"
	"wazmeow/pkg/logger"

	"github.com/google/uuid"
)

// ========================================
// SESSION MANAGER USE CASES
// ========================================
// Este arquivo agrupa os casos de uso para gerenciamento básico de sessões:
// 1. CreateSessionUseCase - Criar novas sessões
// 2. ListSessionsUseCase - Listar todas as sessões
// 3. DeleteSessionUseCase - Remover sessões
// 4. GetSessionInfoUseCase - Obter informações detalhadas
// ========================================

// CreateSessionUseCase representa o caso de uso para criar sessões
type CreateSessionUseCase struct {
	sessionRepo repositories.SessionRepository
}

// NewCreateSessionUseCase cria uma nova instância do use case
func NewCreateSessionUseCase(sessionRepo repositories.SessionRepository) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de criação de sessão
func (uc *CreateSessionUseCase) Execute(req *requests.CreateSessionRequest) (*entities.Session, error) {
	// Validar request
	if req.Name == "" {
		return nil, fmt.Errorf("nome da sessão é obrigatório")
	}

	if !req.IsValidURLName() {
		return nil, fmt.Errorf("nome da sessão deve ter entre 3-50 caracteres, conter apenas letras, números, hífens e underscores, e não pode começar ou terminar com hífen ou underscore")
	}

	// Verificar se já existe uma sessão com esse nome
	exists, err := uc.sessionRepo.ExistsByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar existência da sessão: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("já existe uma sessão com o nome '%s'", req.Name)
	}

	// Criar nova sessão
	session := &entities.Session{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Status:    entities.StatusDisconnected,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Persistir sessão
	if err := uc.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("erro ao criar sessão: %w", err)
	}

	logger.Info("Sessão '%s' criada com sucesso (ID: %s)", session.Name, session.ID)
	return session, nil
}

// ListSessionsUseCase representa o caso de uso para listar sessões
type ListSessionsUseCase struct {
	sessionRepo repositories.SessionRepository
}

// NewListSessionsUseCase cria uma nova instância do use case
func NewListSessionsUseCase(sessionRepo repositories.SessionRepository) *ListSessionsUseCase {
	return &ListSessionsUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de listagem de sessões
func (uc *ListSessionsUseCase) Execute() ([]*entities.Session, error) {
	sessions, err := uc.sessionRepo.List()
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// DeleteSessionUseCase representa o caso de uso para deletar sessões
type DeleteSessionUseCase struct {
	sessionRepo repositories.SessionRepository
}

// NewDeleteSessionUseCase cria uma nova instância do use case
func NewDeleteSessionUseCase(sessionRepo repositories.SessionRepository) *DeleteSessionUseCase {
	return &DeleteSessionUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de deleção de sessão
func (uc *DeleteSessionUseCase) Execute(sessionID string) error {
	// Buscar sessão para verificar se existe
	session, err := uc.findSession(sessionID)
	if err != nil {
		return err
	}

	// Desconectar cliente se estiver conectado
	if session.Client != nil && session.Client.IsConnected() {
		session.Client.Disconnect()
		logger.Info("Cliente da sessão '%s' desconectado antes da deleção", session.Name)
	}

	// Deletar sessão do repositório
	if err := uc.sessionRepo.Delete(session.ID); err != nil {
		return fmt.Errorf("erro ao deletar sessão: %w", err)
	}

	logger.Info("Sessão '%s' deletada com sucesso", session.Name)
	return nil
}

// findSession busca uma sessão por ID ou nome
func (uc *DeleteSessionUseCase) findSession(identifier string) (*entities.Session, error) {
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

// GetSessionInfoUseCase representa o caso de uso para obter informações de sessão
type GetSessionInfoUseCase struct {
	sessionRepo repositories.SessionRepository
}

// NewGetSessionInfoUseCase cria uma nova instância do use case
func NewGetSessionInfoUseCase(sessionRepo repositories.SessionRepository) *GetSessionInfoUseCase {
	return &GetSessionInfoUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de obtenção de informações da sessão
func (uc *GetSessionInfoUseCase) Execute(sessionID string) (*responses.SessionInfo, error) {
	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return nil, err
	}

	// Criar resposta com informações detalhadas
	sessionInfo := &responses.SessionInfo{
		Session:     session,
		IsConnected: false,
		IsLoggedIn:  false,
	}

	// Verificar status de conexão e login se cliente existe
	if session.Client != nil {
		sessionInfo.IsConnected = session.Client.IsConnected()
		sessionInfo.IsLoggedIn = session.Client.IsLoggedIn()
	}

	return sessionInfo, nil
}

// findSession busca uma sessão por ID ou nome
func (uc *GetSessionInfoUseCase) findSession(identifier string) (*entities.Session, error) {
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
