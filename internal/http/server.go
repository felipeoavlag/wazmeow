package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wazmeow/internal/container"
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
	router := NewRouter(container)

	// Criar servidor HTTP
	httpServer := &http.Server{
		Addr:         container.GetConfig().GetServerAddress(),
		Handler:      router,
		ReadTimeout:  container.GetConfig().Server.ReadTimeout,
		WriteTimeout: container.GetConfig().Server.WriteTimeout,
	}

	return &Server{
		container:  container,
		httpServer: httpServer,
	}, nil
}

// Start inicia o servidor HTTP com graceful shutdown
func (s *Server) Start() error {
	// Exibir informações de inicialização
	s.printStartupInfo()

	// Inicializar conexões WhatsApp das sessões que estavam conectadas
	if err := s.initializeWhatsAppConnections(); err != nil {
		logger.Error("Erro ao inicializar conexões WhatsApp: %v", err)
		// Não retornar erro para não impedir a inicialização do servidor
	}

	// Canal para capturar sinais do sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Canal para erros do servidor
	serverErrors := make(chan error, 1)

	// Iniciar servidor em goroutine
	go func() {
		logger.Info("Iniciando servidor HTTP...")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("erro ao iniciar servidor: %w", err)
		}
	}()

	// Aguardar sinal de parada ou erro
	select {
	case err := <-serverErrors:
		return err
	case sig := <-quit:
		logger.Info("Recebido sinal de parada: %v", sig)
		return s.shutdown()
	}
}

// shutdown realiza o graceful shutdown do servidor
func (s *Server) shutdown() error {
	logger.Info("Iniciando graceful shutdown...")

	// Criar contexto com timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Desconectar todas as sessões WhatsApp
	logger.Info("Desconectando sessões WhatsApp...")
	sessionManager := s.container.GetSessionManager()
	sessionManager.DisconnectAll()

	// Parar de aceitar novas conexões e aguardar as existentes terminarem
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Erro durante shutdown do servidor HTTP: %v", err)
		return err
	}

	// Fechar recursos do container
	if err := s.container.Close(); err != nil {
		logger.Error("Erro ao fechar recursos do container: %v", err)
		return err
	}

	logger.Info("Servidor parado com sucesso")
	return nil
}

// printStartupInfo exibe informações sobre a inicialização do servidor
func (s *Server) printStartupInfo() {
	cfg := s.container.GetConfig()

	fmt.Printf("🚀 WazMeow API Server\n")
	fmt.Printf("=====================\n")
	fmt.Printf("🌐 Servidor: http://%s\n", cfg.GetServerAddress())
	fmt.Printf("🗄️  Banco: %s:%s/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	fmt.Printf("📋 Ambiente: %s\n", cfg.App.Environment)
	fmt.Printf("📊 Log Level: %s\n", cfg.Log.Level)
	fmt.Printf("📋 Health Check: http://%s/health\n", cfg.GetServerAddress())
	fmt.Printf("=====================\n")
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
