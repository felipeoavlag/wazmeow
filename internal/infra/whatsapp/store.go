package whatsapp

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun/driver/pgdriver"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"wazmeow/internal/config"
	"wazmeow/pkg/logger"
)

// NewStoreContainer creates a new WhatsApp store container using PostgreSQL
func NewStoreContainer(cfg config.DatabaseConfig) (*sqlstore.Container, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)

	// Create SQL database connection for WhatsApp store
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create WhatsApp store container
	container := sqlstore.NewWithDB(sqldb, "postgres", nil)

	// Upgrade store schema if needed
	ctx := context.Background()
	err := container.Upgrade(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to upgrade WhatsApp store schema")
		return nil, err
	}

	logger.Info().Msg("WhatsApp store container initialized successfully")
	return container, nil
}
