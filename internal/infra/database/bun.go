package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"wazmeow/internal/infra/models"
	"wazmeow/pkg/logger"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// BunConnection representa uma conexão com o banco de dados usando Bun ORM
type BunConnection struct {
	DB *bun.DB
}

// NewBunConnection cria uma nova conexão com PostgreSQL usando Bun ORM
func NewBunConnection(cfg Config) (*BunConnection, error) {
	// Construir DSN para PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	logger.Info("Conectando ao PostgreSQL com Bun ORM: %s:%s/%s", cfg.Host, cfg.Port, cfg.Name)

	// Criar conexão SQL usando pgdriver
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Configurar timeouts e limites de conexão
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(5 * time.Minute)

	// Criar instância Bun
	db := bun.NewDB(sqldb, pgdialect.New())

	// Adicionar debug em desenvolvimento
	if cfg.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
		logger.Info("Bun debug mode habilitado")
	}

	// Testar conexão com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("erro ao conectar com PostgreSQL via Bun: %w", err)
	}

	logger.Info("Conexão Bun estabelecida com sucesso")

	return &BunConnection{DB: db}, nil
}

// RegisterModels registra os modelos no Bun para uso com ORM
func (c *BunConnection) RegisterModels() {
	logger.Info("Registrando modelos Bun...")

	// Registrar modelos no Bun
	c.DB.RegisterModel((*models.SessionModel)(nil))

	logger.Info("Modelos Bun registrados com sucesso")
}

// Health verifica se a conexão está saudável
func (c *BunConnection) Health(ctx context.Context) error {
	if c.DB == nil {
		return fmt.Errorf("conexão Bun não inicializada")
	}

	// Usar ping nativo do Bun - mais direto e eficiente
	err := c.DB.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("health check falhou: %w", err)
	}

	logger.Debug("Health check Bun: OK")
	return nil
}

// Close fecha a conexão com o banco de dados
func (c *BunConnection) Close() error {
	if c.DB != nil {
		logger.Info("Fechando conexão Bun com PostgreSQL")
		return c.DB.Close()
	}
	return nil
}

// GetDB retorna a instância do Bun DB para uso direto se necessário
func (c *BunConnection) GetDB() *bun.DB {
	return c.DB
}

// BeginTx inicia uma transação
func (c *BunConnection) BeginTx(ctx context.Context) (bun.Tx, error) {
	return c.DB.BeginTx(ctx, nil)
}

// RunInTransaction executa uma função dentro de uma transação
func (c *BunConnection) RunInTransaction(ctx context.Context, fn func(tx bun.Tx) error) error {
	tx, err := c.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.Error("Erro ao fazer rollback: %v", rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit: %w", err)
	}

	return nil
}

// EnsureSchema garante que as tabelas existam baseado nos models
func (c *BunConnection) EnsureSchema(ctx context.Context) error {
	return ValidateSchema(ctx, c.DB)
}
