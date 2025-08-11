package database

import (
	"context"
	"fmt"

	"wazmeow/internal/infra/models"
	"wazmeow/pkg/logger"

	"github.com/uptrace/bun"
)

// CreateTablesFromModels cria tabelas automaticamente baseado nos models usando Bun nativo
func CreateTablesFromModels(ctx context.Context, db *bun.DB) error {
	logger.Info("🏗️ Criando tabelas automaticamente baseado nos models...")

	// Lista de todos os models registrados
	models := []interface{}{
		(*models.SessionModel)(nil),
		// TODO: Adicionar outros models aqui conforme necessário
	}

	for _, model := range models {
		modelName := fmt.Sprintf("%T", model)
		logger.Info("📋 Criando tabela para: %s", modelName)

		// Usar funcionalidade nativa do Bun para criar tabela
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("erro ao criar tabela para %s: %w", modelName, err)
		}

		logger.Info("✅ Tabela criada: %s", modelName)
	}

	logger.Info("🎉 Todas as tabelas criadas com sucesso!")
	return nil
}

// ValidateSchema verifica se as tabelas existem e cria se necessário
func ValidateSchema(ctx context.Context, db *bun.DB) error {
	logger.Info("🔍 Validando schema contra models...")

	// Verificar se tabela sessions existe
	var exists bool
	err := db.NewSelect().
		ColumnExpr("EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sessions')").
		Scan(ctx, &exists)

	if err != nil {
		return fmt.Errorf("erro ao verificar tabela sessions: %w", err)
	}

	if !exists {
		logger.Info("⚠️ Tabela sessions não existe, criando automaticamente...")
		return CreateTablesFromModels(ctx, db)
	}

	logger.Info("✅ Schema validado - todas as tabelas existem!")
	return nil
}

// DropAllTables remove todas as tabelas (usar com cuidado!)
func DropAllTables(ctx context.Context, db *bun.DB) error {
	logger.Info("🗑️ Removendo todas as tabelas...")

	models := []interface{}{
		(*models.SessionModel)(nil),
	}

	// Reverter ordem para evitar problemas de dependência
	for i := len(models) - 1; i >= 0; i-- {
		model := models[i]
		modelName := fmt.Sprintf("%T", model)
		logger.Info("🗑️ Removendo: %s", modelName)

		_, err := db.NewDropTable().
			Model(model).
			IfExists().
			Cascade().
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("erro ao remover tabela %s: %w", modelName, err)
		}

		logger.Info("✅ Tabela removida: %s", modelName)
	}

	logger.Info("🎉 Todas as tabelas removidas!")
	return nil
}

// RecreateAllTables remove e recria todas as tabelas
func RecreateAllTables(ctx context.Context, db *bun.DB) error {
	logger.Info("🔄 Recriando todas as tabelas...")

	// Remover todas
	err := DropAllTables(ctx, db)
	if err != nil {
		return fmt.Errorf("erro ao remover tabelas: %w", err)
	}

	// Criar novamente
	err = CreateTablesFromModels(ctx, db)
	if err != nil {
		return fmt.Errorf("erro ao recriar tabelas: %w", err)
	}

	logger.Info("✅ Todas as tabelas recriadas com sucesso!")
	return nil
}

// GetSchemaStatus retorna informações sobre o status do schema
func GetSchemaStatus(ctx context.Context, db *bun.DB) (*SchemaStatus, error) {
	status := &SchemaStatus{
		Tables: make(map[string]bool),
	}

	// Lista das tabelas esperadas
	expectedTables := []string{"sessions"} // Adicionar outras conforme necessário

	for _, tableName := range expectedTables {
		var exists bool
		err := db.NewSelect().
			ColumnExpr("EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = ?)", tableName).
			Scan(ctx, &exists)

		if err != nil {
			return nil, fmt.Errorf("erro ao verificar tabela %s: %w", tableName, err)
		}

		status.Tables[tableName] = exists
		status.TotalTables++
		if exists {
			status.ExistingTables++
		} else {
			status.MissingTables++
		}
	}

	status.IsFullySynced = status.MissingTables == 0

	return status, nil
}

// SchemaStatus representa o status atual do schema
type SchemaStatus struct {
	TotalTables     int            `json:"total_tables"`
	ExistingTables  int            `json:"existing_tables"`
	MissingTables   int            `json:"missing_tables"`
	IsFullySynced   bool           `json:"is_fully_synced"`
	Tables          map[string]bool `json:"tables"`
}