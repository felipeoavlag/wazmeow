package main

import (
	"fmt"
	"log"
	"os"

	"wazmeow/internal/config"
	"wazmeow/internal/infra/database"
	"wazmeow/pkg/logger"

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
			// ===================================
			// COMANDOS SIMPLIFICADOS - BUNO NATIVO
			// ===================================
			{
				Name:  "auto-create",
				Usage: "create tables automatically from models using Bun native functions",
				Action: func(c *cli.Context) error {
					db, err := connectDB(c.String("env"))
					if err != nil {
						return err
					}
					defer db.Close()

					err = database.CreateTablesFromModels(c.Context, db.GetDB())
					if err != nil {
						return fmt.Errorf("erro ao criar tabelas: %w", err)
					}

					fmt.Println("✅ Tabelas criadas automaticamente com sucesso!")
					return nil
				},
			},
			{
				Name:  "auto-validate",
				Usage: "validate schema against models and create missing tables",
				Action: func(c *cli.Context) error {
					db, err := connectDB(c.String("env"))
					if err != nil {
						return err
					}
					defer db.Close()

					err = database.ValidateSchema(c.Context, db.GetDB())
					if err != nil {
						return fmt.Errorf("erro na validação: %w", err)
					}

					fmt.Println("✅ Schema validado e sincronizado!")
					return nil
				},
			},
			{
				Name:  "auto-status",
				Usage: "show schema status compared to models",
				Action: func(c *cli.Context) error {
					db, err := connectDB(c.String("env"))
					if err != nil {
						return err
					}
					defer db.Close()

					status, err := database.GetSchemaStatus(c.Context, db.GetDB())
					if err != nil {
						return fmt.Errorf("erro ao obter status: %w", err)
					}

					fmt.Println("📊 Schema Status:")
					fmt.Printf("  📋 Total de tabelas esperadas: %d\n", status.TotalTables)
					fmt.Printf("  ✅ Tabelas existentes: %d\n", status.ExistingTables)
					fmt.Printf("  ❌ Tabelas faltando: %d\n", status.MissingTables)
					fmt.Printf("  🎯 Sincronizado: %t\n", status.IsFullySynced)
					
					if len(status.Tables) > 0 {
						fmt.Println("  📋 Detalhes por tabela:")
						for table, exists := range status.Tables {
							symbol := "❌"
							if exists {
								symbol = "✅"
							}
							fmt.Printf("    %s %s\n", symbol, table)
						}
					}

					if !status.IsFullySynced {
						fmt.Println("\n💡 Execute 'db auto-create' para criar tabelas faltando")
					}

					return nil
				},
			},
			{
				Name:  "recreate",
				Usage: "drop and recreate all tables (CAUTION: destroys data!)",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "confirm",
						Usage: "confirm that you want to destroy all data",
					},
				},
				Action: func(c *cli.Context) error {
					if !c.Bool("confirm") {
						fmt.Println("⚠️  Este comando irá DESTRUIR todos os dados!")
						fmt.Println("Para confirmar, use: --confirm")
						return nil
					}

					db, err := connectDB(c.String("env"))
					if err != nil {
						return err
					}
					defer db.Close()

					err = database.RecreateAllTables(c.Context, db.GetDB())
					if err != nil {
						return fmt.Errorf("erro ao recriar tabelas: %w", err)
					}

					fmt.Println("✅ Todas as tabelas foram recriadas!")
					return nil
				},
			},
		},
	}
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
