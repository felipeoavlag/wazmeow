package container

import (
	"context"
	"fmt"
	"time"

	"wazmeow/pkg/logger"
)

// Close fecha todas as conexões e recursos do container
func (c *Container) Close() error {
	if !c.IsInitialized() {
		return fmt.Errorf("container não foi inicializado")
	}

	logger.Info("🔄 Fechando container...")

	var errors []error

	// Fechar conexões em ordem reversa da inicialização
	if c.bunDB != nil {
		logger.Debug("Fechando conexão Bun...")
		if err := c.bunDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("erro ao fechar Bun DB: %w", err))
		}
	}

	if c.db != nil {
		logger.Debug("Fechando conexão WhatsApp...")
		if err := c.db.Close(); err != nil {
			errors = append(errors, fmt.Errorf("erro ao fechar WhatsApp DB: %w", err))
		}
	}

	// Marcar como não inicializado
	c.setInitialized(false)

	if len(errors) > 0 {
		logger.Error("❌ Erros ao fechar container: %v", errors)
		return fmt.Errorf("erros ao fechar container: %v", errors)
	}

	logger.Info("✅ Container fechado com sucesso")
	return nil
}

// HealthCheck verifica a saúde de todas as dependências críticas
func (c *Container) HealthCheck(ctx context.Context) error {
	if !c.IsInitialized() {
		return fmt.Errorf("container não inicializado")
	}

	// Timeout para health check
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Verificar banco de dados Bun
	if c.bunDB != nil {
		if err := c.bunDB.Health(ctx); err != nil {
			return fmt.Errorf("banco de dados Bun não saudável: %w", err)
		}
	}

	// Verificar repositories críticos
	if c.sessionRepo == nil {
		return fmt.Errorf("session repository não inicializado")
	}

	// Verificar use cases críticos
	if c.sessionUseCases == nil {
		return fmt.Errorf("session use cases não inicializados")
	}

	if c.sessionUseCases.Create == nil {
		return fmt.Errorf("create session use case não inicializado")
	}

	// Verificar domain services
	if c.sessionDomainService == nil {
		return fmt.Errorf("session domain service não inicializado")
	}

	return nil
}

// Restart reinicia o container (fecha e inicializa novamente)
func (c *Container) Restart() error {
	logger.Info("🔄 Reiniciando container...")

	// Salvar configuração atual
	currentConfig := c.config

	// Fechar container atual
	if err := c.Close(); err != nil {
		return fmt.Errorf("erro ao fechar container para restart: %w", err)
	}

	// Recriar com a mesma configuração
	newContainer, err := NewWithConfig(currentConfig)
	if err != nil {
		return fmt.Errorf("erro ao recriar container: %w", err)
	}

	// Copiar estado do novo container
	*c = *newContainer

	logger.Info("✅ Container reiniciado com sucesso")
	return nil
}

// GetStatus retorna informações sobre o status do container
func (c *Container) GetStatus() map[string]interface{} {
	status := map[string]interface{}{
		"initialized": c.IsInitialized(),
		"components":  map[string]bool{},
	}

	if c.IsInitialized() {
		components := status["components"].(map[string]bool)
		components["config"] = c.config != nil
		components["database"] = c.db != nil
		components["bun_database"] = c.bunDB != nil
		components["session_manager"] = c.sessionManager != nil
		components["session_repository"] = c.sessionRepo != nil
		components["session_domain_service"] = c.sessionDomainService != nil
		components["session_use_cases"] = c.sessionUseCases != nil
	}

	return status
}
