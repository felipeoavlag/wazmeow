package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"wazmeow/internal/domain/entity"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/infra/models"

	"github.com/uptrace/bun"
)

// BunSessionRepository implementa o SessionRepository usando Bun ORM
type BunSessionRepository struct {
	db *bun.DB
}

// NewBunSessionRepository cria uma nova instância do repository usando Bun
func NewBunSessionRepository(db *bun.DB) repository.SessionRepository {
	return &BunSessionRepository{db: db}
}

// Create cria uma nova sessão no banco de dados
func (r *BunSessionRepository) Create(session *entity.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Definir timestamps se não estiverem definidos
	now := time.Now()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}
	if session.UpdatedAt.IsZero() {
		session.UpdatedAt = now
	}

	// Converter para modelo de persistência
	model := models.FromDomain(session)

	// Inserir no banco
	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	if err != nil {
		// Tratar erro de constraint unique
		if isUniqueConstraintError(err) {
			return fmt.Errorf("já existe uma sessão com o nome '%s'", session.Name)
		}
		return fmt.Errorf("erro ao criar sessão: %w", err)
	}

	return nil
}

// GetByID busca uma sessão pelo ID
func (r *BunSessionRepository) GetByID(id string) (*entity.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	model := new(models.SessionModel)
	err := r.db.NewSelect().
		Model(model).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sessão com ID '%s' não encontrada", id)
		}
		return nil, fmt.Errorf("erro ao buscar sessão por ID: %w", err)
	}

	return model.ToDomain(), nil
}

// GetByName busca uma sessão pelo nome
func (r *BunSessionRepository) GetByName(name string) (*entity.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	model := new(models.SessionModel)
	err := r.db.NewSelect().
		Model(model).
		Where("name = ?", name).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sessão com nome '%s' não encontrada", name)
		}
		return nil, fmt.Errorf("erro ao buscar sessão por nome: %w", err)
	}

	return model.ToDomain(), nil
}

// List retorna todas as sessões ordenadas por data de criação
func (r *BunSessionRepository) List() ([]*entity.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sessionModels []*models.SessionModel
	err := r.db.NewSelect().
		Model(&sessionModels).
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("erro ao listar sessões: %w", err)
	}

	return models.ToDomainList(sessionModels), nil
}

// Update atualiza uma sessão existente
func (r *BunSessionRepository) Update(session *entity.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Atualizar timestamp
	session.UpdatedAt = time.Now()

	// Converter para modelo de persistência
	model := models.FromDomain(session)

	// Atualizar no banco
	result, err := r.db.NewUpdate().
		Model(model).
		Where("id = ?", session.ID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("erro ao atualizar sessão: %w", err)
	}

	// Verificar se alguma linha foi afetada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sessão com ID '%s' não encontrada para atualização", session.ID)
	}

	return nil
}

// Delete remove uma sessão do banco de dados
func (r *BunSessionRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.db.NewDelete().
		Model((*models.SessionModel)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("erro ao deletar sessão: %w", err)
	}

	// Verificar se alguma linha foi afetada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sessão com ID '%s' não encontrada para deleção", id)
	}

	return nil
}

// ExistsByName verifica se existe uma sessão com o nome especificado
func (r *BunSessionRepository) ExistsByName(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := r.db.NewSelect().
		Model((*models.SessionModel)(nil)).
		Where("name = ?", name).
		Exists(ctx)

	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência da sessão: %w", err)
	}

	return exists, nil
}

// isUniqueConstraintError verifica se o erro é de violação de constraint unique
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	// Verificar se é erro de constraint unique do PostgreSQL
	errStr := err.Error()
	return contains(errStr, "unique constraint") ||
		contains(errStr, "duplicate key") ||
		contains(errStr, "UNIQUE constraint failed")
}

// contains verifica se uma string contém uma substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOfSubstring(s, substr) >= 0)))
}

// indexOfSubstring encontra a posição de uma substring
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
