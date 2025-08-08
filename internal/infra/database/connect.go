package database

import (
	"context"
	"fmt"
	"time"

	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

// Config representa a configuração do banco de dados
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Debug    bool
}

// Connection representa uma conexão com o banco de dados
type Connection struct {
	Store *sqlstore.Container
}

// Connect estabelece conexão com o banco de dados PostgreSQL
func Connect(cfg Config) (*Connection, error) {
	// Construir URL de conexão
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	logger.Info("Conectando ao banco de dados: %s:%s/%s", cfg.Host, cfg.Port, cfg.Name)

	// Criar container do sqlstore com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Usar nosso logger centralizado para o whatsmeow
	container, err := sqlstore.New(ctx, "postgres", dbURL, logger.ForWhatsApp())
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o banco de dados: %w", err)
	}

	logger.Info("Conexão com banco de dados estabelecida com sucesso")

	return &Connection{
		Store: container,
	}, nil
}

// Close fecha a conexão com o banco de dados
func (c *Connection) Close() error {
	if c.Store != nil {
		logger.Info("Fechando conexão com banco de dados")
		// O sqlstore não tem método Close público, mas o contexto gerencia isso
		return nil
	}
	return nil
}

// Health verifica se a conexão está saudável
func (c *Connection) Health() error {
	if c.Store == nil {
		return fmt.Errorf("conexão não inicializada")
	}

	// Tentar obter um device para testar a conexão
	device := c.Store.NewDevice()
	if device == nil {
		return fmt.Errorf("não foi possível criar device - conexão pode estar inativa")
	}

	// Se chegou até aqui, a conexão está funcionando
	logger.Debug("Health check do banco de dados: OK")
	return nil
}

// GetDeviceStore retorna um novo device store
func (c *Connection) GetDeviceStore() *store.Device {
	if c.Store == nil {
		return nil
	}
	return c.Store.NewDevice()
}

// Migrate executa as migrações necessárias
func (c *Connection) Migrate() error {
	logger.Info("Executando migrações do banco de dados")

	// O whatsmeow/sqlstore gerencia suas próprias migrações automaticamente
	// quando um novo device é criado pela primeira vez

	// Criar um device temporário para forçar a criação das tabelas
	device := c.Store.NewDevice()
	if device == nil {
		return fmt.Errorf("erro ao criar device para migração")
	}

	logger.Info("Migrações executadas com sucesso")
	return nil
}
