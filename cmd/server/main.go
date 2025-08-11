package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wazmeow/internal/config"
	"wazmeow/internal/infra/database"
	"wazmeow/internal/infra/http"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize WhatsApp store container
	waStore, err := whatsapp.NewStoreContainer(cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize WhatsApp store")
	}

	// Initialize HTTP server
	server, err := http.NewServer(cfg, db, waStore)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize HTTP server")
	}

	// Start server
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	logger.Info().
		Str("host", cfg.Server.Host).
		Int("port", cfg.Server.Port).
		Msg("Server started successfully")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited")
}
