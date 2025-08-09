package migrations

import (
	"github.com/uptrace/bun/migrate"
)

// Migrations é a coleção de migrações do sistema
// Seguindo exatamente a documentação do Bun ORM para migrações Go
var Migrations = migrate.NewMigrations()
