# Sistema de Migrações - WazMeow

Este diretório contém o sistema de migrações do banco de dados seguindo exatamente as melhores práticas do [Bun ORM](https://bun.uptrace.dev/guide/migrations.html) usando **migrações Go**.

## Estrutura

```
migrations/
├── migrations.go                           # Coleção de migrações
├── 20250109000001_create_sessions.go       # Migração Go
└── README.md                              # Esta documentação
```

## Como usar

### 1. Inicializar o sistema de migrações

```bash
go run cmd/migrate/main.go -env=dev db init
```

### 2. Executar migrações

```bash
go run cmd/migrate/main.go -env=dev db migrate
```

### 3. Verificar status das migrações

```bash
go run cmd/migrate/main.go -env=dev db status
```

### 4. Fazer rollback da última migração

```bash
go run cmd/migrate/main.go -env=dev db rollback
```

## Convenções de Nomenclatura

As migrações seguem o padrão: `YYYYMMDDHHMMSS_description.go`

Exemplo:
- `20250109000001_create_sessions.go`
- `20250109000002_add_user_table.go`

## Estrutura de Migração Go

Cada migração Go deve usar `Migrations.MustRegister` com funções UP e DOWN:

```go
package migrations

import (
    "context"
    "fmt"
    "github.com/uptrace/bun"
)

func init() {
    Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
        fmt.Print(" [UP] creating table...")
        // Código da migração UP
        return nil
    }, func(ctx context.Context, db *bun.DB) error {
        fmt.Print(" [DOWN] dropping table...")
        // Código da migração DOWN
        return nil
    })
}
```

## Estrutura do Código

### migrations.go

```go
package migrations

import (
    "github.com/uptrace/bun/migrate"
)

// Migrations é a coleção de migrações do sistema
// Seguindo exatamente a documentação do Bun ORM para migrações Go
var Migrations = migrate.NewMigrations()
```

## Adicionando Novas Migrações

1. Crie um arquivo `.go` com timestamp único: `YYYYMMDDHHMMSS_description.go`
2. Use `Migrations.MustRegister` com funções UP e DOWN
3. Implemente tratamento de erros adequado
4. Teste a migração em ambiente de desenvolvimento
5. Execute `go run cmd/migrate/main.go -env=dev db migrate`

## Melhores Práticas

1. **Sempre crie migrações reversíveis** (funções UP e DOWN)
2. **Teste em ambiente de desenvolvimento primeiro**
3. **Use tratamento de erros adequado** com `fmt.Errorf`
4. **Documente migrações complexas** com comentários
5. **Mantenha migrações atômicas** (uma responsabilidade por migração)
6. **Nunca edite migrações já aplicadas em produção**
7. **Use mensagens de log** para acompanhar progresso

## Tabelas de Controle

O Bun cria automaticamente as seguintes tabelas:
- `bun_migrations`: Controla quais migrações foram aplicadas
- `bun_migration_locks`: Previne execução concorrente de migrações
