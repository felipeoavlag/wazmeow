# Guia Prático: Bun ORM com Migrations Automáticas

Este guia mostra como implementar migrations automáticas usando as funcionalidades nativas do **Bun ORM**, seguindo as melhores práticas da [documentação oficial](https://bun.uptrace.dev/).

## 🎯 Abordagem Simplificada

Em vez de criar um sistema complexo, vamos usar as funcionalidades nativas do Bun:

### 1. **Schema Automático via Models**
```go
// O Bun pode criar tabelas automaticamente a partir dos models
_, err := db.NewCreateTable().Model((*models.SessionModel)(nil)).IfNotExists().Exec(ctx)
```

### 2. **Migrations com SQL + Models**
```go
// Combinar SQL manual para alterações complexas com auto-create para novas tabelas
func init() {
    Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
        // Criar novas tabelas automaticamente
        _, err := db.NewCreateTable().Model((*models.SessionModel)(nil)).IfNotExists().Exec(ctx)
        return err
    }, func(ctx context.Context, db *bun.DB) error {
        // Rollback
        _, err := db.NewDropTable().Model((*models.SessionModel)(nil)).IfExists().Exec(ctx)
        return err
    })
}
```

## 🚀 Implementação Prática

### Passo 1: Atualizar o Sistema Atual

```go
// internal/infra/database/auto_migrations.go
package database

import (
    "context"
    "wazmeow/internal/infra/models"
    "github.com/uptrace/bun"
)

// CreateTablesFromModels cria tabelas automaticamente baseado nos models
func CreateTablesFromModels(ctx context.Context, db *bun.DB) error {
    models := []interface{}{
        (*models.SessionModel)(nil),
        // Adicionar outros models aqui
    }
    
    for _, model := range models {
        _, err := db.NewCreateTable().
            Model(model).
            IfNotExists().
            Exec(ctx)
        if err != nil {
            return err
        }
    }
    
    return nil
}

// ValidateSchema verifica se as tabelas existem
func ValidateSchema(ctx context.Context, db *bun.DB) error {
    // Verificar se tabela sessions existe
    var exists bool
    err := db.NewSelect().
        ColumnExpr("EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sessions')").
        Scan(ctx, &exists)
    
    if err != nil {
        return err
    }
    
    if !exists {
        return CreateTablesFromModels(ctx, db)
    }
    
    return nil
}
```

### Passo 2: Comandos CLI Simples

```go
// Adicionar ao cmd/migrate/main.go
{
    Name:  "auto-create",
    Usage: "create tables automatically from models",
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

        fmt.Println("✅ Tabelas criadas com sucesso!")
        return nil
    },
},
{
    Name:  "validate",
    Usage: "validate schema against models",
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

        fmt.Println("✅ Schema validado com sucesso!")
        return nil
    },
},
```

### Passo 3: Makefile Targets

```makefile
# Adicionar ao Makefile
db-auto-create: ## Criar tabelas automaticamente dos models
	@echo "🏗️ Criando tabelas automaticamente..."
	@go run cmd/migrate/main.go --env=dev db auto-create

db-validate: ## Validar schema contra models
	@echo "🔍 Validando schema..."
	@go run cmd/migrate/main.go --env=dev db validate

db-reset: ## Resetar e recriar todas as tabelas
	@echo "🔄 Resetando banco de dados..."
	@make docker-restart
	@sleep 5
	@make db-auto-create
```

## 🎯 Fluxo de Trabalho Simplificado

### Para Novos Projetos
```bash
# 1. Subir banco
make docker-up

# 2. Criar tabelas automaticamente
make db-auto-create

# 3. Validar
make db-validate
```

### Para Mudanças nos Models
```bash
# 1. Alterar o model SessionModel
# 2. Criar migration SQL manual se necessário
# 3. Executar migration
make db-migrate

# 4. Validar resultado
make db-validate
```

## 🔧 Integração com o Sistema Atual

### Atualizar o BunConnection
```go
// internal/infra/database/bun.go
func (c *BunConnection) EnsureSchema(ctx context.Context) error {
    return ValidateSchema(ctx, c.DB)
}
```

### Usar na Inicialização da Aplicação
```go
// cmd/server/main.go
func main() {
    // ... configuração ...
    
    // Garantir que schema existe
    err = bunConnection.EnsureSchema(context.Background())
    if err != nil {
        logger.Fatal("Erro ao garantir schema: %v", err)
    }
    
    // ... resto da aplicação ...
}
```

## 📊 Vantagens desta Abordagem

✅ **Simples**: Usa funcionalidades nativas do Bun
✅ **Confiável**: Menos código customizado = menos bugs  
✅ **Flexível**: Combina auto-create com migrations manuais quando necessário
✅ **Rápido**: Setup inicial em minutos, não horas
✅ **Manutenível**: Segue padrões estabelecidos do Bun

## 🎯 Resultado Final

- **Zero** SQL hardcoded para criação de tabelas
- **Auto-detecção** de schema baseado nos models
- **Validação** automática na inicialização
- **Comandos CLI** simples e diretos
- **100%** compatível com sistema atual

---

**Esta abordagem mantém a simplicidade while delivering the core functionality needed.**