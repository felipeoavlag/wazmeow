package http

import (
	"context"
	"fmt"
	"net/http"

	"wazmeow/internal/container"
	"wazmeow/internal/http/router"
	"wazmeow/pkg/logger"

	"github.com/lib/pq"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

// Server representa o servidor HTTP da aplicação
type Server struct {
	container  *container.Container
	httpServer *http.Server
}

// NewServer cria um novo servidor HTTP com todas as dependências configuradas
func NewServer() (*Server, error) {
	// Configurar PostgreSQL array wrapper
	sqlstore.PostgresArrayWrapper = pq.Array

	// Criar container com todas as dependências
	container, err := container.New()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar container: %w", err)
	}

	// Configurar roteador HTTP usando use cases
	httpRouter := router.NewRouter(container)

	// Criar servidor HTTP
	httpServer := &http.Server{
		Addr:         container.GetConfig().GetServerAddress(),
		Handler:      httpRouter,
		ReadTimeout:  container.GetConfig().Server.ReadTimeout,
		WriteTimeout: container.GetConfig().Server.WriteTimeout,
	}

	return &Server{
		container:  container,
		httpServer: httpServer,
	}, nil
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	// Exibir informações de inicialização
	s.printStartupInfo()

	// Iniciar serviço de webhooks
	logger.Info("Iniciando serviço de webhooks...")
	webhookService := s.container.GetWebhookService()
	if webhookService != nil {
		if err := webhookService.Start(); err != nil {
			logger.Error("Erro ao iniciar webhook service: %v", err)
			// Não retornar erro para não impedir a inicialização do servidor
		}
	}

	// Inicializar conexões WhatsApp das sessões que estavam conectadas
	if err := s.initializeWhatsAppConnections(); err != nil {
		logger.Error("Erro ao inicializar conexões WhatsApp: %v", err)
		// Não retornar erro para não impedir a inicialização do servidor
	}

	// Iniciar servidor HTTP
	logger.Info("🚀 Servidor HTTP iniciado em %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("erro ao iniciar servidor: %w", err)
	}

	return nil
}

// Shutdown realiza o graceful shutdown do servidor
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Iniciando graceful shutdown...")

	// 1. Parar servidor HTTP primeiro (para de aceitar novas conexões)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Erro ao parar servidor HTTP: %v", err)
	}

	// 2. Parar webhook service de forma simples
	logger.Info("Parando webhook service...")
	if webhookService := s.container.GetWebhookService(); webhookService != nil {
		webhookService.Stop()
	}

	// 3. Desconectar sessões WhatsApp
	logger.Info("Desconectando sessões WhatsApp...")
	if sessionManager := s.container.GetSessionManager(); sessionManager != nil {
		sessionManager.DisconnectAll()
	}

	// 4. Fechar container
	logger.Info("Fechando container...")
	if err := s.container.Close(); err != nil {
		logger.Error("Erro ao fechar container: %v", err)
	}

	logger.Info("Servidor parado com sucesso")
	return nil
}

// printStartupInfo exibe informações sobre a inicialização do servidor
func (s *Server) printStartupInfo() {
	cfg := s.container.GetConfig()

	logger.Info("🚀 WazMeow API Server")
	logger.Info("=====================")
	logger.Info("🌐 Servidor: http://%s", cfg.GetServerAddress())
	logger.Info("🗄️  Banco: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	logger.Info("📋 Ambiente: %s", cfg.App.Environment)
	logger.Info("📊 Log Level: %s", cfg.Log.Level)
	logger.Info("📋 Health Check: http://%s/health", cfg.GetServerAddress())
	logger.Info("=====================")
}

// GetContainer retorna o container de dependências (útil para testes)
func (s *Server) GetContainer() *container.Container {
	return s.container
}

// initializeWhatsAppConnections inicializa as conexões WhatsApp das sessões que estavam conectadas
func (s *Server) initializeWhatsAppConnections() error {
	logger.Info("Inicializando conexões WhatsApp...")

	clientFactory := s.container.GetClientFactory()
	sessionManager := s.container.GetSessionManager()

	if err := clientFactory.ConnectOnStartup(sessionManager); err != nil {
		return fmt.Errorf("erro ao conectar sessões na inicialização: %w", err)
	}

	logger.Info("Conexões WhatsApp inicializadas com sucesso")
	return nil
}
