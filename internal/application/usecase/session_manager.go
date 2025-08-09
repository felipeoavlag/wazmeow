// Package usecase contém os casos de uso da camada de aplicação
// Este arquivo (session_manager.go) contém os use cases para:
// - Orquestração do gerenciamento básico do ciclo de vida das sessões (CRUD)
// - Operações fundamentais: Create, List, Delete, GetInfo
// - Coordena entre domain services, repositories e infraestrutura
// - Não contém regras de negócio (delegadas para domain services)
package usecase

import (
	"fmt"
	"time"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/entity"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/domain/service"
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
	sessionRepo   repository.SessionRepository
	domainService *service.SessionDomainService
}

// NewCreateSessionUseCase cria uma nova instância do use case
func NewCreateSessionUseCase(sessionRepo repository.SessionRepository, domainService *service.SessionDomainService) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		sessionRepo:   sessionRepo,
		domainService: domainService,
	}
}

// Execute executa o caso de uso de criação de sessão
func (uc *CreateSessionUseCase) Execute(req *requests.CreateSessionRequest) (*entity.Session, error) {
	// Validar request usando domain service
	if err := uc.domainService.ValidateSessionName(req.Name); err != nil {
		return nil, err
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
	session := &entity.Session{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Status:    entity.StatusDisconnected,
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
	sessionRepo repository.SessionRepository
}

// NewListSessionsUseCase cria uma nova instância do use case
func NewListSessionsUseCase(sessionRepo repository.SessionRepository) *ListSessionsUseCase {
	return &ListSessionsUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de listagem de sessões
func (uc *ListSessionsUseCase) Execute() ([]*entity.Session, error) {
	sessions, err := uc.sessionRepo.List()
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// DeleteSessionUseCase representa o caso de uso para deletar sessões
type DeleteSessionUseCase struct {
	sessionRepo   repository.SessionRepository
	domainService *service.SessionDomainService
}

// NewDeleteSessionUseCase cria uma nova instância do use case
func NewDeleteSessionUseCase(sessionRepo repository.SessionRepository, domainService *service.SessionDomainService) *DeleteSessionUseCase {
	return &DeleteSessionUseCase{
		sessionRepo:   sessionRepo,
		domainService: domainService,
	}
}

// Execute executa o caso de uso de deleção de sessão
func (uc *DeleteSessionUseCase) Execute(sessionID string) error {
	// Buscar sessão para verificar se existe
	session, err := uc.findSession(sessionID)
	if err != nil {
		return err
	}

	// Verificar se pode deletar usando regras de negócio
	if err := uc.domainService.CanDelete(session); err != nil {
		return err
	}

	// TODO: Implementar desconexão usando SessionManager
	// Por enquanto, apenas log
	logger.Info("Deletando sessão '%s'", session.Name)

	// Deletar sessão do repositório
	if err := uc.sessionRepo.Delete(session.ID); err != nil {
		return fmt.Errorf("erro ao deletar sessão: %w", err)
	}

	logger.Info("Sessão '%s' deletada com sucesso", session.Name)
	return nil
}

// findSession busca uma sessão por ID ou nome
func (uc *DeleteSessionUseCase) findSession(identifier string) (*entity.Session, error) {
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
	sessionRepo repository.SessionRepository
}

// NewGetSessionInfoUseCase cria uma nova instância do use case
func NewGetSessionInfoUseCase(sessionRepo repository.SessionRepository) *GetSessionInfoUseCase {
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

	// TODO: Implementar verificação de status usando SessionManager
	// Por enquanto, definir como false
	sessionInfo.IsConnected = false
	sessionInfo.IsLoggedIn = false

	return sessionInfo, nil
}

// findSession busca uma sessão por ID ou nome
func (uc *GetSessionInfoUseCase) findSession(identifier string) (*entity.Session, error) {
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
