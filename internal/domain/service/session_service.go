package service

import (
	"fmt"
	"time"

	"wazmeow/internal/domain/entity"
)

// SessionDomainService contém as regras de negócio puras para sessões
// Esta camada não deve ter dependências de infraestrutura ou frameworks externos
type SessionDomainService struct{}

// NewSessionDomainService cria uma nova instância do domain service
func NewSessionDomainService() *SessionDomainService {
	return &SessionDomainService{}
}

// ========================================
// VALIDAÇÕES DE NEGÓCIO
// ========================================

// ValidateSessionName valida se o nome da sessão atende às regras de negócio
func (s *SessionDomainService) ValidateSessionName(name string) error {
	if name == "" {
		return fmt.Errorf("nome da sessão é obrigatório")
	}

	if len(name) < 3 || len(name) > 50 {
		return fmt.Errorf("nome da sessão deve ter entre 3-50 caracteres")
	}

	// Permitir apenas letras, números, hífens e underscores
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return fmt.Errorf("nome da sessão deve conter apenas letras, números, hífens e underscores")
		}
	}

	// Não pode começar ou terminar com hífen ou underscore
	if name[0] == '-' || name[0] == '_' ||
		name[len(name)-1] == '-' || name[len(name)-1] == '_' {
		return fmt.Errorf("nome da sessão não pode começar ou terminar com hífen ou underscore")
	}

	return nil
}

// ValidatePhoneNumber valida se o número de telefone está no formato correto
func (s *SessionDomainService) ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("número de telefone é obrigatório")
	}

	// Validação básica - pode ser expandida conforme necessário
	if len(phone) < 10 || len(phone) > 15 {
		return fmt.Errorf("número de telefone deve ter entre 10-15 dígitos")
	}

	return nil
}

// ========================================
// REGRAS DE CONEXÃO
// ========================================

// CanConnect verifica se uma sessão pode ser conectada
func (s *SessionDomainService) CanConnect(session *entity.Session) error {
	if session == nil {
		return fmt.Errorf("sessão não pode ser nula")
	}

	switch session.Status {
	case entity.StatusConnected:
		return fmt.Errorf("sessão '%s' já está conectada", session.Name)
	case entity.StatusConnecting:
		return fmt.Errorf("sessão '%s' já está em processo de conexão", session.Name)
	case entity.StatusDisconnected:
		return nil // Pode conectar
	default:
		return fmt.Errorf("status da sessão '%s' é inválido: %s", session.Name, session.Status)
	}
}

// ShouldReconnect determina se uma sessão deve tentar reconectar automaticamente
func (s *SessionDomainService) ShouldReconnect(session *entity.Session, lastDisconnectTime time.Time) bool {
	if session == nil {
		return false
	}

	// Sempre tentar reconectar se estiver desconectado
	// (não há mais distinção entre logout manual e desconexão)

	// Reconectar apenas se desconectou há menos de 1 hora
	timeSinceDisconnect := time.Since(lastDisconnectTime)
	return timeSinceDisconnect < time.Hour
}

// CalculateRetryInterval calcula o intervalo para próxima tentativa de conexão
func (s *SessionDomainService) CalculateRetryInterval(attemptCount int) time.Duration {
	// Backoff exponencial com limite máximo
	baseInterval := 5 * time.Second
	maxInterval := 5 * time.Minute

	interval := time.Duration(attemptCount) * baseInterval
	if interval > maxInterval {
		interval = maxInterval
	}

	return interval
}

// ========================================
// REGRAS DE TIMEOUT E LIFECYCLE
// ========================================

// CalculateSessionTimeout calcula o timeout para uma sessão baseado em suas características
func (s *SessionDomainService) CalculateSessionTimeout(session *entity.Session) time.Duration {
	if session == nil {
		return 30 * time.Minute // Default timeout
	}

	// Sessões com proxy podem precisar de timeout maior
	if session.ProxyConfig != nil {
		return 45 * time.Minute
	}

	return 30 * time.Minute
}

// IsSessionExpired verifica se uma sessão expirou baseado em suas regras de negócio
func (s *SessionDomainService) IsSessionExpired(session *entity.Session, lastActivity time.Time) bool {
	if session == nil {
		return true
	}

	timeout := s.CalculateSessionTimeout(session)
	return time.Since(lastActivity) > timeout
}

// ========================================
// REGRAS DE PROXY
// ========================================

// ValidateProxyConfig valida a configuração de proxy
func (s *SessionDomainService) ValidateProxyConfig(config *entity.ProxyConfig) error {
	if config == nil {
		return fmt.Errorf("configuração de proxy não pode ser nula")
	}

	if config.Type != "http" && config.Type != "socks5" {
		return fmt.Errorf("tipo de proxy deve ser 'http' ou 'socks5'")
	}

	if config.Host == "" {
		return fmt.Errorf("host do proxy é obrigatório")
	}

	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("porta do proxy deve estar entre 1 e 65535")
	}

	return nil
}

// ========================================
// REGRAS DE NEGÓCIO GERAIS
// ========================================

// CanDelete verifica se uma sessão pode ser deletada
func (s *SessionDomainService) CanDelete(session *entity.Session) error {
	if session == nil {
		return fmt.Errorf("sessão não pode ser nula")
	}

	// Não permitir deletar sessões conectadas
	if session.Status == entity.StatusConnected || session.Status == entity.StatusConnecting {
		return fmt.Errorf("não é possível deletar sessão '%s' enquanto estiver conectada ou conectando", session.Name)
	}

	return nil
}

// GenerateSessionID gera um ID único para uma nova sessão
// Esta é uma regra de negócio sobre como IDs devem ser gerados
func (s *SessionDomainService) GenerateSessionID() string {
	// Por enquanto, delegar para UUID, mas poderia ter regras específicas
	// como prefixos, formatos específicos, etc.
	return "" // Será implementado com UUID no use case
}

// ========================================
// REGRAS DE STATUS
// ========================================

// GetNextValidStatus retorna os próximos status válidos para uma sessão
func (s *SessionDomainService) GetNextValidStatus(currentStatus entity.SessionStatus) []entity.SessionStatus {
	switch currentStatus {
	case entity.StatusDisconnected:
		return []entity.SessionStatus{entity.StatusConnecting}
	case entity.StatusConnecting:
		return []entity.SessionStatus{entity.StatusConnected, entity.StatusDisconnected}
	case entity.StatusConnected:
		return []entity.SessionStatus{entity.StatusDisconnected}
	default:
		return []entity.SessionStatus{}
	}
}

// CanTransitionTo verifica se é possível transicionar de um status para outro
func (s *SessionDomainService) CanTransitionTo(from, to entity.SessionStatus) bool {
	validStatuses := s.GetNextValidStatus(from)
	for _, status := range validStatuses {
		if status == to {
			return true
		}
	}
	return false
}
