package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [UP] creating sessions table...")

		// Criar tabela sessions completa
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS sessions (
				-- Campos principais da sessão
				id VARCHAR NOT NULL,
				name VARCHAR NOT NULL,
				status VARCHAR NOT NULL DEFAULT 'disconnected',
				phone VARCHAR NULL,

				-- Campos WhatsApp (conexão e autenticação)
				device_jid VARCHAR DEFAULT '',
				qrcode TEXT DEFAULT '',
				webhook_url VARCHAR DEFAULT '',
				events VARCHAR DEFAULT '',

				-- Campos de proxy
				proxy_type VARCHAR NULL,
				proxy_host VARCHAR NULL,
				proxy_port BIGINT NULL,
				proxy_username VARCHAR NULL,
				proxy_password VARCHAR NULL,

				-- Campos de auditoria (sempre no final)
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

				CONSTRAINT sessions_pkey PRIMARY KEY (id),
				CONSTRAINT sessions_name_key UNIQUE (name)
			);
		`)
		if err != nil {
			return fmt.Errorf("erro ao criar tabela sessions: %w", err)
		}

		// Criar índices para performance
		indexes := []string{
			`CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status)`,
			`CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at)`,
			`CREATE INDEX IF NOT EXISTS idx_sessions_phone ON sessions(phone) WHERE phone IS NOT NULL`,
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

		// Remover índices primeiro
		indexes := []string{
			`DROP INDEX IF EXISTS idx_sessions_device_jid`,
			`DROP INDEX IF EXISTS idx_sessions_phone`,
			`DROP INDEX IF EXISTS idx_sessions_created_at`,
			`DROP INDEX IF EXISTS idx_sessions_status`,
		}

		for _, indexSQL := range indexes {
			if _, err := db.ExecContext(ctx, indexSQL); err != nil {
				return fmt.Errorf("erro ao remover índice: %w", err)
			}
		}

		// Remover tabela sessions
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS sessions CASCADE`)
		if err != nil {
			return fmt.Errorf("erro ao remover tabela sessions: %w", err)
		}

		fmt.Println(" OK")
		return nil
	})
}
