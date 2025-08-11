package database

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"wazmeow/internal/config"
	"wazmeow/pkg/logger"
)

// NewConnection creates a new database connection using Bun ORM
func NewConnection(cfg config.DatabaseConfig) (*bun.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create Bun DB instance
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add debug logging in development
	if cfg.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(5)

	logger.Info().Msg("Database connection established successfully")
	return db, nil
}
