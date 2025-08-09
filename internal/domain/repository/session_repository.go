package repository

import (
	"wazmeow/internal/domain/entity"
)

// SessionRepository define a interface para persistência de sessões
type SessionRepository interface {
	// Create cria uma nova sessão
	Create(session *entity.Session) error

	// GetByID busca uma sessão pelo ID
	GetByID(id string) (*entity.Session, error)

	// GetByName busca uma sessão pelo nome
	GetByName(name string) (*entity.Session, error)

	// List retorna todas as sessões
	List() ([]*entity.Session, error)

	// Update atualiza uma sessão existente
	Update(session *entity.Session) error

	// Delete remove uma sessão
	Delete(id string) error

	// ExistsByName verifica se existe uma sessão com o nome
	ExistsByName(name string) (bool, error)
}
