package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uptrace/bun"

	"wazmeow/internal/application/handlers"
	"wazmeow/internal/application/usecases/session"
	"wazmeow/internal/config"
	"wazmeow/internal/infra/database/repositories"
	"wazmeow/internal/infra/http/routes"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"
)

// Server represents the HTTP server
type Server struct {
	config     *config.Config
	httpServer *http.Server
	router     chi.Router
}

// NewServer creates a new HTTP server instance
func NewServer(cfg *config.Config, db *bun.DB, waStore interface{}) (*Server, error) {
	// Initialize repositories
	sessionRepo := repositories.NewSessionRepository(db)

	// Initialize WhatsApp store and service
	whatsappStore, err := whatsapp.NewStoreContainer(cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize WhatsApp store")
	}
	whatsappService := whatsapp.NewService(sessionRepo, whatsappStore)

	// Initialize use cases
	createSessionUC := session.NewCreateSessionUseCase(sessionRepo)
	listSessionsUC := session.NewListSessionsUseCase(sessionRepo)
	connectSessionUC := session.NewConnectSessionUseCase(sessionRepo, whatsappService)

	// Initialize handlers
	sessionHandler := handlers.NewSessionHandler(createSessionUC, listSessionsUC, connectSessionUC, whatsappService)

	// Create router
	router := chi.NewRouter()

	// Setup middleware
	setupMiddleware(router)

	// Setup routes
	routes.SetupRoutes(router, sessionHandler)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	server := &Server{
		config:     cfg,
		httpServer: httpServer,
		router:     router,
	}

	logger.Info().
		Str("address", addr).
		Msg("HTTP server configured")

	return server, nil
}

// setupMiddleware configures all middleware for the router
func setupMiddleware(router chi.Router) {
	// Basic middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(handlers.RecoveryMiddleware())
	router.Use(handlers.LoggingMiddleware())
	router.Use(handlers.CORSMiddleware())
	router.Use(handlers.ContentTypeMiddleware())

	// Timeout middleware
	router.Use(middleware.Timeout(60 * time.Second))

	// Compress responses
	router.Use(middleware.Compress(5))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	logger.Info().
		Str("address", s.httpServer.Addr).
		Msg("Starting HTTP server")

	if s.config.Server.SSLCertFile != "" && s.config.Server.SSLKeyFile != "" {
		logger.Info().Msg("Starting HTTPS server")
		return s.httpServer.ListenAndServeTLS(s.config.Server.SSLCertFile, s.config.Server.SSLKeyFile)
	}

	logger.Info().Msg("Starting HTTP server")
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info().Msg("Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}
