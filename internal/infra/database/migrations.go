package database

import (
	"context"

	"github.com/uptrace/bun"

	"wazmeow/internal/infra/database/models"
	"wazmeow/pkg/logger"
)

// RunMigrations runs all database migrations using Bun
func RunMigrations(db *bun.DB) error {
	ctx := context.Background()
	logger.Info().Msg("Running database migrations")

	// Create tables using Bun models
	tables := []interface{}{
		(*models.SessionModel)(nil),
	}

	for _, table := range tables {
		_, err := db.NewCreateTable().
			Model(table).
			IfNotExists().
			Exec(ctx)

		if err != nil {
			logger.Error().Err(err).Msg("Failed to create table")
			return err
		}
	}

	// Create indexes using Bun query builder (zero SQL)
	if err := createIndexes(ctx, db); err != nil {
		return err
	}

	logger.Info().Msg("Database migrations completed successfully")
	return nil
}

// createIndexes creates database indexes using Bun query builder
func createIndexes(ctx context.Context, db *bun.DB) error {
	// Create index on Sessions.status
	_, err := db.NewCreateIndex().
		Model((*models.SessionModel)(nil)).
		Index("idx_sessions_status").
		Column("status").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create status index")
		return err
	}

	// Create index on Sessions.createdAt
	_, err = db.NewCreateIndex().
		Model((*models.SessionModel)(nil)).
		Index("idx_sessions_created_at").
		Column("createdAt").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create createdAt index")
		return err
	}

	// Create index on Sessions.deviceJID for better performance on device lookups
	_, err = db.NewCreateIndex().
		Model((*models.SessionModel)(nil)).
		Index("idx_sessions_device_jid").
		Column("deviceJID").
		IfNotExists().
		Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create deviceJID index")
		return err
	}

	logger.Debug().Msg("Database indexes created successfully")
	return nil
}
