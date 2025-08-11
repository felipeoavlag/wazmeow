package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wazmeow/internal/http"
	"wazmeow/pkg/logger"

	_ "wazmeow/docs" // Import para registrar documentação Swagger

	"github.com/lib/pq"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

// @title WazMeow API
// @version 1.0
// @description API REST para gerenciamento de sessões WhatsApp usando whatsmeow
// @contact.name WazMeow Support
// @contact.email support@wazmeow.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https

// @tag.name sessions
// @tag.description Gerenciamento de sessões WhatsApp

// @tag.name messages
// @tag.description Envio e gerenciamento de mensagens

// @tag.name chats
// @tag.description Operações de chat e presença

// @tag.name groups
// @tag.description Gerenciamento de grupos WhatsApp

// @tag.name contacts
// @tag.description Gerenciamento de contatos

// @tag.name webhooks
// @tag.description Configuração de webhooks

// @tag.name newsletters
// @tag.description Gerenciamento de newsletters

// @tag.name health
// @tag.description Verificação de saúde da API

func main() {
	// Configurar PostgreSQL array wrapper
	sqlstore.PostgresArrayWrapper = pq.Array

	// Criar servidor
	server, err := http.NewServer()
	if err != nil {
		logger.Error("❌ Erro ao criar servidor: %v", err)
		return
	}

	// Canal para capturar sinais do sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Canal para erros do servidor
	serverErrors := make(chan error, 1)

	// Iniciar servidor em goroutine
	go func() {
		logger.Info("🚀 Iniciando servidor WazMeow...")
		if err := server.Start(); err != nil {
			serverErrors <- err
		}
	}()

	// Aguardar sinal de parada ou erro
	select {
	case err := <-serverErrors:
		logger.Fatal("Erro no servidor: %v", err)
	case sig := <-quit:
		logger.Info("Recebido sinal de parada: %v", sig)

		// Timeout para shutdown forçado
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		// Fazer shutdown direto e simples
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Erro durante shutdown: %v", err)
			os.Exit(1)
		}

		logger.Info("Servidor parado com sucesso")
	}
}
