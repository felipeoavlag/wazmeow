// Package usecases contém os casos de uso da aplicação
// Este arquivo (session_setup.go) contém os use cases para:
// - Configuração e setup inicial das sessões
// - Preparação de sessões para uso (proxy, emparelhamento)
// - Operações: PairPhone (emparelhamento), SetProxy (configuração de rede)
// - Configurações técnicas que preparam a sessão para conectar
package usecases

import (
	"context"
	"fmt"
	"time"

	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow"
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
	sessionRepo repositories.SessionRepository
}

// NewPairPhoneUseCase cria uma nova instância do use case
func NewPairPhoneUseCase(sessionRepo repositories.SessionRepository) *PairPhoneUseCase {
	return &PairPhoneUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de emparelhamento por telefone
func (uc *PairPhoneUseCase) Execute(sessionID, phone string) (string, error) {
	// Validar número de telefone
	if phone == "" {
		return "", fmt.Errorf("número de telefone é obrigatório")
	}

	// Buscar sessão
	session, err := uc.findSession(sessionID)
	if err != nil {
		return "", err
	}

	// Verificar se o cliente existe
	if session.Client == nil {
		return "", fmt.Errorf("sessão '%s' não possui cliente inicializado", session.Name)
	}

	// Conectar se não estiver conectado
	if !session.Client.IsConnected() {
		if err := session.Client.Connect(); err != nil {
			return "", fmt.Errorf("erro ao conectar cliente: %w", err)
		}
	}

	// Emparelhar telefone
	code, err := session.Client.PairPhone(
		context.Background(),
		phone,
		true,
		whatsmeow.PairClientChrome,
		"Chrome (Linux)",
	)
	if err != nil {
		return "", fmt.Errorf("erro ao emparelhar telefone: %w", err)
	}

	// Atualizar sessão com o telefone
	session.Phone = phone
	session.UpdatedAt = time.Now()

	if err := uc.sessionRepo.Update(session); err != nil {
		logger.Error("Erro ao atualizar sessão após emparelhamento: %v", err)
	}

	logger.Info("Código de emparelhamento gerado para %s na sessão '%s': %s", phone, session.Name, code)
	return code, nil
}

// findSession busca uma sessão por ID ou nome
func (uc *PairPhoneUseCase) findSession(identifier string) (*entities.Session, error) {
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
	sessionRepo repositories.SessionRepository
}

// NewSetProxyUseCase cria uma nova instância do use case
func NewSetProxyUseCase(sessionRepo repositories.SessionRepository) *SetProxyUseCase {
	return &SetProxyUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute executa o caso de uso de configuração de proxy
func (uc *SetProxyUseCase) Execute(sessionID string, proxyConfig *entities.ProxyConfig) error {
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
func (uc *SetProxyUseCase) findSession(identifier string) (*entities.Session, error) {
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
