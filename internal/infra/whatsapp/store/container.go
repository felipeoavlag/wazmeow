package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun/driver/pgdriver"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"wazmeow/internal/config"
	"wazmeow/pkg/logger"
)

// NewContainer creates a new WhatsApp store container optimized for multi-session
func NewContainer(cfg config.DatabaseConfig) (*sqlstore.Container, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)

	// Create SQL database connection for WhatsApp store with optimized settings
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(connector)

	// Configure connection pool for multiple sessions (otimizado)
	sqldb.SetMaxOpenConns(100)                 // Máximo de conexões abertas (aumentado)
	sqldb.SetMaxIdleConns(25)                  // Máximo de conexões idle (aumentado)
	sqldb.SetConnMaxLifetime(30 * time.Minute) // Tempo de vida otimizado
	sqldb.SetConnMaxIdleTime(5 * time.Minute)  // Tempo idle otimizado

	// Create WhatsApp store container with logging
	storeLog := logger.NewWALogger("SQLStore")
	container := sqlstore.NewWithDB(sqldb, "postgres", storeLog)

	// Upgrade store schema if needed
	ctx := context.Background()
	err := container.Upgrade(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to upgrade WhatsApp store schema")
		return nil, err
	}

	logger.Info().Msg("WhatsApp store container initialized successfully for multi-session")
	return container, nil
}

// ValidateStoreHealth verifica a saúde do store
func ValidateStoreHealth(container *sqlstore.Container) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tentar buscar devices para validar conectividade
	_, err := container.GetAllDevices(ctx)
	if err != nil {
		return fmt.Errorf("store health check failed: %w", err)
	}

	logger.Info().Msg("WhatsApp store health check passed")
	return nil
}
