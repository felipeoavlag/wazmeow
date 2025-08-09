package main

import (
	"fmt"
	"log"
	"os"

	"wazmeow/internal/config"
	"wazmeow/internal/infra/database"
	"wazmeow/internal/infra/database/migrations"
	"wazmeow/pkg/logger"

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

					migrator := createMigrator(db)

					// Verificar status das migrações
					ms, err := migrator.MigrationsWithStatus(c.Context)
					if err != nil {
						return fmt.Errorf("erro ao verificar status das migrações: %w", err)
					}

					fmt.Printf("Migrations status:\n")
					for _, m := range ms {
						status := "❌ PENDING"
						if m.GroupID != 0 {
							status = "✅ APPLIED"
						}
						fmt.Printf("  %s %s\n", status, m.Name)
					}

					// Verificar se tabela sessions existe
					var exists bool
					err = db.GetDB().NewSelect().
						ColumnExpr("EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sessions')").
						Scan(c.Context, &exists)

					if err != nil {
						return err
					}

					fmt.Printf("\nDatabase status:\n")
					if exists {
						fmt.Println("  ✅ Tabela sessions existe")
					} else {
						fmt.Println("  ❌ Tabela sessions não existe")
					}

					return nil
				},
			},
		},
	}
}

func createMigrator(db *database.BunConnection) *migrate.Migrator {
	// Usar a coleção de migrações do pacote migrations
	// Seguindo exatamente a documentação do Bun ORM
	return migrate.NewMigrator(
		db.GetDB(),
		migrations.Migrations,
		migrate.WithTableName("bun_migrations"),
		migrate.WithLocksTableName("bun_migration_locks"),
	)
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
