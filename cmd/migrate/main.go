package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"wazmeow/internal/config"
	"wazmeow/internal/infra/database"
	"wazmeow/pkg/logger"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "wazmeow-migrate",
		Usage: "WazMeow database migration tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment (dev, prod, test)",
			},
		},
		Commands: []*cli.Command{
			newDBCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newDBCommand() *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					db, err := connectDB(c.String("env"))
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := createMigrator(db)

					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					db, err := connectDB(c.String("env"))
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := createMigrator(db)

					group, err := migrator.Migrate(c.Context)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to run\n")
						return nil
					}

					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					db, err := connectDB(c.String("env"))
					if err != nil {
						return err
					}
					defer db.Close()

					migrator := createMigrator(db)

					group, err := migrator.Rollback(c.Context)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}

					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					db, err := connectDB(c.String("env"))
					if err != nil {
						return err
					}
					defer db.Close()

					// Verificar se tabela sessions existe
					var exists bool
					err = db.GetDB().NewSelect().
						ColumnExpr("EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sessions')").
						Scan(c.Context, &exists)

					if err != nil {
						return err
					}

					if exists {
						fmt.Println("✅ Tabela sessions existe")
					} else {
						fmt.Println("❌ Tabela sessions não existe")
					}

					return nil
				},
			},
		},
	}
}

func createMigrator(db *database.BunConnection) *migrate.Migrator {
	migrations := migrate.NewMigrations()

	// Registrar migração para criar tabela sessions
	migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [UP] creating sessions table...")

		// Criar tabela sessions completa
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS sessions (
				id VARCHAR NOT NULL,
				name VARCHAR NOT NULL,
				status VARCHAR NOT NULL DEFAULT 'disconnected',
				phone VARCHAR NULL,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

				-- Campos de proxy
				proxy_type VARCHAR NULL,
				proxy_host VARCHAR NULL,
				proxy_port BIGINT NULL,
				proxy_username VARCHAR NULL,
				proxy_password VARCHAR NULL,

				-- Campos WhatsApp
				webhook_url VARCHAR DEFAULT '',
				qrcode TEXT DEFAULT '',
				device_jid VARCHAR DEFAULT '',
				is_connected BOOLEAN DEFAULT FALSE,

				CONSTRAINT sessions_pkey PRIMARY KEY (id),
				CONSTRAINT sessions_name_key UNIQUE (name)
			)
		`)
		if err != nil {
			return fmt.Errorf("erro ao criar tabela sessions: %w", err)
		}

		// Criar índices para performance
		indexes := []string{
			`CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status)`,
			`CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at)`,
			`CREATE INDEX IF NOT EXISTS idx_sessions_phone ON sessions(phone) WHERE phone IS NOT NULL`,
			`CREATE INDEX IF NOT EXISTS idx_sessions_is_connected ON sessions(is_connected)`,
			`CREATE INDEX IF NOT EXISTS idx_sessions_device_jid ON sessions(device_jid) WHERE device_jid != ''`,
		}

		for _, indexSQL := range indexes {
			if _, err := db.ExecContext(ctx, indexSQL); err != nil {
				return fmt.Errorf("erro ao criar índice: %w", err)
			}
		}

		fmt.Println(" OK")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [DOWN] dropping sessions table...")

		// Remover tabela sessions
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS sessions CASCADE`)
		if err != nil {
			return fmt.Errorf("erro ao remover tabela sessions: %w", err)
		}

		fmt.Println(" OK")
		return nil
	})

	return migrate.NewMigrator(db.GetDB(), migrations, migrate.WithTableName("migrations"), migrate.WithLocksTableName("migration_locks"))
}

func connectDB(env string) (*database.BunConnection, error) {
	// Carregar configuração
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	// Inicializar logger
	logger.InitGlobalLogger(cfg.Log.Level)

	// Configuração do banco de dados
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
		Debug:    cfg.Database.Debug,
	}

	// Conectar ao banco de dados com Bun ORM
	bunConnection, err := database.NewBunConnection(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com Bun: %w", err)
	}

	return bunConnection, nil
}
