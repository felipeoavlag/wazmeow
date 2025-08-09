// Package usecase contém os casos de uso da camada de aplicação
// Este arquivo (session_setup.go) contém os use cases para:
// - Orquestração da configuração e setup inicial das sessões
// - Coordenação da preparação de sessões para uso (proxy, emparelhamento)
// - Operações: PairPhone (emparelhamento), SetProxy (configuração de rede)
// - Integração entre domain services e infraestrutura para configurações
package usecase

import (
	"fmt"
	"time"

	"wazmeow/internal/domain/entity"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/domain/service"
	"wazmeow/pkg/logger"
)

// ========================================
// SESSION SETUP USE CASES
// ========================================
// Este arquivo agrupa os casos de uso para configuração inicial:
// 1. PairPhoneUseCase - Emparelhamento via número de telefone
// 2. SetProxyUseCase - Configuração de proxy para conexão
//
// Responsabilidades:
// - Configuração de proxy para conexões
// - Emparelhamento alternativo via telefone
// - Setup técnico antes da conexão
// - Preparação de sessões para diferentes cenários de uso
// ========================================

// PairPhoneUseCase representa o caso de uso para emparelhar telefone
type PairPhoneUseCase struct {
	sessionRepo   repository.SessionRepository
	domainService *service.SessionDomainService
}

// NewPairPhoneUseCase cria uma nova instância do use case
func NewPairPhoneUseCase(sessionRepo repository.SessionRepository, domainService *service.SessionDomainService) *PairPhoneUseCase {
	return &PairPhoneUseCase{
		sessionRepo:   sessionRepo,
		domainService: domainService,
	}
}

// Execute executa o caso de uso de emparelhamento por telefone
func (uc *PairPhoneUseCase) Execute(sessionID, phone string) (string, error) {
	// Validar número de telefone usando domain service
	if err := uc.domainService.ValidatePhoneNumber(phone); err != nil {
		return "", err
	}

	// Buscar sessão
	_, err := uc.findSession(sessionID)
	if err != nil {
		return "", err
	}

	// TODO: Implementar usando SessionManager
	// Por enquanto, retornar erro
	return "", fmt.Errorf("funcionalidade de emparelhamento não implementada ainda")
}

// findSession busca uma sessão por ID ou nome
func (uc *PairPhoneUseCase) findSession(identifier string) (*entity.Session, error) {
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

// SetProxyUseCase representa o caso de uso para configurar proxy
type SetProxyUseCase struct {
	sessionRepo   repository.SessionRepository
	domainService *service.SessionDomainService
}

// NewSetProxyUseCase cria uma nova instância do use case
func NewSetProxyUseCase(sessionRepo repository.SessionRepository, domainService *service.SessionDomainService) *SetProxyUseCase {
	return &SetProxyUseCase{
		sessionRepo:   sessionRepo,
		domainService: domainService,
	}
}

// Execute executa o caso de uso de configuração de proxy
func (uc *SetProxyUseCase) Execute(sessionID string, proxyConfig *entity.ProxyConfig) error {
	// Validar configuração de proxy
	if proxyConfig == nil {
		return fmt.Errorf("configuração de proxy é obrigatória")
	}

	if proxyConfig.Host == "" {
		return fmt.Errorf("host do proxy é obrigatório")
	}

	if proxyConfig.Port <= 0 || proxyConfig.Port > 65535 {
		return fmt.Errorf("porta do proxy deve estar entre 1 e 65535")
	}

	if proxyConfig.Type != "http" && proxyConfig.Type != "socks5" {
		return fmt.Errorf("tipo de proxy deve ser 'http' ou 'socks5'")
	}

	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return err
	}

	// Configurar proxy na sessão
	session.ProxyConfig = proxyConfig
	session.UpdatedAt = time.Now()

	// Atualizar sessão no repositório
	if err := uc.sessionRepo.Update(session); err != nil {
		return fmt.Errorf("erro ao atualizar configuração de proxy: %w", err)
	}

	logger.Info("Proxy configurado para sessão '%s': %s://%s:%d",
		session.Name, proxyConfig.Type, proxyConfig.Host, proxyConfig.Port)

	return nil
}

// findSession busca uma sessão por ID ou nome
func (uc *SetProxyUseCase) findSession(identifier string) (*entity.Session, error) {
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
