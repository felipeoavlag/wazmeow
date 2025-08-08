package repositories

import (
	"wazmeow/internal/domain/entities"
)

// SessionRepository define a interface para persistência de sessões
type SessionRepository interface {
	// Create cria uma nova sessão
	Create(session *entities.Session) error
	
	// GetByID busca uma sessão pelo ID
	GetByID(id string) (*entities.Session, error)
	
	// GetByName busca uma sessão pelo nome
	GetByName(name string) (*entities.Session, error)
	
	// List retorna todas as sessões
	List() ([]*entities.Session, error)
	
	// Update atualiza uma sessão existente
	Update(session *entities.Session) error
	
	// Delete remove uma sessão
	Delete(id string) error
	
	// ExistsByName verifica se existe uma sessão com o nome
	ExistsByName(name string) (bool, error)
}
