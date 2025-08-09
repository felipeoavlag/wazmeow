package usecase

import (
	"fmt"

	"wazmeow/internal/domain/entity"
	"wazmeow/internal/domain/repository"
)

// SessionFinder é um helper para buscar sessões por ID ou nome
type SessionFinder struct {
	sessionRepo repository.SessionRepository
}

// NewSessionFinder cria uma nova instância do SessionFinder
func NewSessionFinder(sessionRepo repository.SessionRepository) *SessionFinder {
	return &SessionFinder{
		sessionRepo: sessionRepo,
	}
}

// FindSession busca uma sessão por ID ou nome
func (sf *SessionFinder) FindSession(identifier string) (*entity.Session, error) {
	// Tentar buscar por ID primeiro
	session, err := sf.sessionRepo.GetByID(identifier)
	if err == nil {
		return session, nil
	}

	// Se não encontrou por ID, tentar por nome
	session, err = sf.sessionRepo.GetByName(identifier)
	if err != nil {
		return nil, fmt.Errorf("sessão '%s' não encontrada", identifier)
	}

	return session, nil
}
