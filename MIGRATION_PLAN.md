# Plano de Migração para Bun ORM com Migrations Automáticas

## 📋 Resumo Executivo

Este documento detalha o plano para migrar o WazMeow para usar Bun ORM com migrations geradas automaticamente a partir dos models Go, eliminando migrations manuais em SQL e garantindo sincronização perfeita entre código e banco de dados.

## 🔍 Análise das Inconsistências Atuais

### 1. **Duplicação Crítica de Entidades Session**

**Problema**: Existem duas definições conflitantes da entidade Session:

#### `internal/domain/entities.go` (❌ Problemático)
```go
type Session struct {
    ID           string        `json:"id" bun:"id,pk"`           // TAGS BUN NO DOMAIN!
    Name         string        `json:"name" bun:"name,notnull"`
    Status       SessionStatus `json:"status" bun:"status,notnull"`
    Active       bool          `json:"active" bun:"active,notnull"`     // Campo extra
    Events       []string      `json:"events" bun:"events,array"`       // Array
    LastActivity *time.Time    `json:"last_activity" bun:"last_activity"` // Campo extra
    // ... outros campos
}
```

#### `internal/domain/entity/session.go` (✅ Clean Architecture)
```go
type Session struct {
    ID         string       `json:"id"`                    // SEM TAGS BUN
    Name       string       `json:"name"`
    Status     SessionStatus `json:"status"`
    Phone      string       `json:"phone,omitempty"`
    Events     string       `json:"events,omitempty"`     // String
    // ... campos diferentes
}
```

**Impacto**: 
- Viola princípios de Clean Architecture
- Confusão entre camadas de domínio e persistência
- Conversões inconsistentes no SessionModel
- Bugs potenciais na aplicação

### 2. **Inconsistências Estruturais**

| Campo | entities.go | entity/session.go | SessionModel | Migration SQL |
|-------|-------------|-------------------|--------------|---------------|
| `Active` | `bool` | ❌ Não existe | ❌ Não existe | ❌ Não existe |
| `Events` | `[]string` | `string` | `string` | `VARCHAR` |
| `LastActivity` | `*time.Time` | ❌ Não existe | ❌ Não existe | ❌ Não existe |
| `QRCode` | ❌ Não existe | ❌ Não existe | `string` | `TEXT` |
| `ProxyConfig` | ❌ Não existe | `*ProxyConfig` | Campos separados | Campos separados |

### 3. **Problemas na Arquitetura de Migrations**

**Atual (Problemático)**:
- Migration manual em SQL: [`20250109000001_create_sessions.go`](internal/infra/database/migrations/20250109000001_create_sessions.go)
- Schema hardcoded sem relação com models
- Alterações em models não geram migrations automáticas
- Risco de desincronização entre código e banco

## 🎯 Arquitetura Proposta

### **Princípios da Nova Arquitetura**

1. **Single Source of Truth**: Models Bun são a única fonte da verdade
2. **Migrations Automáticas**: Geradas a partir das mudanças nos models
3. **Versionamento Seguro**: Cada mudança gera nova migration versionada
4. **Clean Architecture**: Domínio livre de tags de persistência

### **Estrutura Proposta**

```
internal/
├── domain/
│   └── entity/
│       └── session.go           # ✅ Entidade clean (sem tags Bun)
├── infra/
│   ├── models/
│   │   └── session.go           # ✅ Model Bun (com todas as tags)
│   └── database/
│       ├── migrations/
│       │   ├── auto/            # 🆕 Migrations geradas automaticamente
│       │   └── manual/          # 🆕 Migrations manuais específicas
│       └── migrator/
│           ├── generator.go     # 🆕 Gerador de migrations
│           ├── differ.go        # 🆕 Detector de mudanças
│           └── schema.go        # 🆕 Análise de schema
```

## 🚀 Plano de Implementação

### **Fase 1: Limpeza e Padronização**

#### 1.1 Definir Entidade Canonical
- ✅ Usar [`internal/domain/entity/session.go`](internal/domain/entity/session.go) como base
- ❌ Remover [`internal/domain/entities.go`](internal/domain/entities.go) 
- 🔧 Atualizar todos os imports e referências

#### 1.2 Normalizar SessionModel
- 🔧 Atualizar [`internal/infra/models/session.go`](internal/infra/models/session.go) para refletir exatamente a entidade
- 🔧 Implementar conversões precisas `ToDomain()` e `FromDomain()`
- ✅ Manter todas as tags Bun necessárias no model

### **Fase 2: Sistema de Migrations Automáticas**

#### 2.1 Detector de Mudanças
```go
// internal/infra/database/migrator/differ.go
type SchemaDiffer struct {
    db *bun.DB
}

func (d *SchemaDiffer) DetectChanges(models []interface{}) (*MigrationDiff, error) {
    currentSchema := d.getCurrentSchema()
    expectedSchema := d.generateSchemaFromModels(models)
    return d.compareSchemas(currentSchema, expectedSchema)
}
```

#### 2.2 Gerador de Migrations
```go
// internal/infra/database/migrator/generator.go
type MigrationGenerator struct {
    diff *MigrationDiff
}

func (g *MigrationGenerator) GenerateMigration(name string) (*Migration, error) {
    timestamp := time.Now().Format("20060102150405")
    filename := fmt.Sprintf("%s_%s.go", timestamp, name)
    
    upSQL := g.generateUpSQL()
    downSQL := g.generateDownSQL()
    
    return &Migration{
        Name:     filename,
        UpSQL:    upSQL,
        DownSQL:  downSQL,
    }, nil
}
```

#### 2.3 Comandos CLI Automáticos
```bash
# Gerar migration baseada em mudanças dos models
go run cmd/migrate/main.go db generate --name="add_user_table"

# Detectar diferenças sem gerar migration
go run cmd/migrate/main.go db diff

# Aplicar migrations pendentes
go run cmd/migrate/main.go db migrate

# Status detalhado
go run cmd/migrate/main.go db status
```

### **Fase 3: Migração das Migrations Existentes**

#### 3.1 Converter Migration Manual
- 🔧 Analisar [`20250109000001_create_sessions.go`](internal/infra/database/migrations/20250109000001_create_sessions.go)
- 🆕 Gerar migration automática equivalente baseada no model atualizado
- ✅ Manter compatibilidade com dados existentes

#### 3.2 Validar Integridade
```go
// Comando para validar sincronização
func ValidateSchemaSync(db *bun.DB, models []interface{}) error {
    differ := NewSchemaDiffer(db)
    diff, err := differ.DetectChanges(models)
    if err != nil {
        return err
    }
    
    if !diff.IsEmpty() {
        return fmt.Errorf("schema out of sync: %+v", diff)
    }
    
    return nil
}
```

### **Fase 4: Sistema de Versionamento**

#### 4.1 Controle de Versão Automático
```go
type MigrationVersion struct {
    Timestamp time.Time
    Hash      string      // Hash dos models
    Models    []string    // Lista de models incluídos
    Changes   []Change    // Mudanças específicas
}
```

#### 4.2 Rollback Inteligente
- 🔧 Gerar automaticamente comandos DOWN baseados nas mudanças
- ✅ Validar segurança do rollback
- ⚠️ Alertar sobre possível perda de dados

## 🔧 Implementação Técnica Detalhada

### **Estrutura do Modelo Bun com Convenção de Nomenclatura**

```go
// internal/infra/models/session.go
type SessionModel struct {
    bun.BaseModel `bun:"table:sessions"`

    // Campos principais (seguindo entity/session.go)
    // 📝 Convenção: camelCase em Go -> snake_case automático no PostgreSQL
    ID     string `bun:"id,pk" json:"id"`
    Name   string `bun:"name,unique,notnull" json:"name"`
    Status string `bun:"status,notnull,default:'disconnected'" json:"status"`
    Phone  string `bun:"phone" json:"phone"`

    // Campos WhatsApp
    // DeviceJID (Go) -> device_jid (PostgreSQL) - conversão automática via Bun
    DeviceJID  string `bun:"deviceJID,default:''" json:"deviceJID"`
    // WebhookURL (Go) -> webhook_url (PostgreSQL) - conversão automática via Bun
    WebhookURL string `bun:"webhookURL,default:''" json:"webhookURL"`
    Events     string `bun:"events,default:''" json:"events"`

    // Campos de proxy (desnormalizados para performance)
    // ProxyType (Go) -> proxy_type (PostgreSQL) - conversão automática via Bun
    ProxyType     *string `bun:"proxyType" json:"proxyType"`
    ProxyHost     *string `bun:"proxyHost" json:"proxyHost"`
    ProxyPort     *int    `bun:"proxyPort" json:"proxyPort"`
    ProxyUsername *string `bun:"proxyUsername" json:"proxyUsername"`
    ProxyPassword *string `bun:"proxyPassword" json:"proxyPassword"`

    // Auditoria
    // CreatedAt (Go) -> created_at (PostgreSQL) - conversão automática via Bun
    CreatedAt time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp" json:"createdAt"`
    UpdatedAt time.Time `bun:"updatedAt,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}
```

**🎯 Convenção de Nomenclatura Adotada:**
- **Go structs**: `camelCase` (DeviceJID, WebhookURL, CreatedAt, ProxyType)
- **PostgreSQL**: `snake_case` (device_jid, webhook_url, created_at, proxy_type)
- **Conversão**: Automática via `bun.NamingStrategy` e tags Bun
- **Benefícios**:
  - ✅ Código Go idiomático (camelCase)
  - ✅ SQL padrão PostgreSQL (snake_case)
  - ✅ Conversão transparente e automática
  - ✅ Zero configuração adicional necessária

### **Sistema de Geração de Migrations**

```go
// internal/infra/database/migrator/schema.go
type TableSchema struct {
    Name        string
    Columns     []ColumnSchema
    Indexes     []IndexSchema
    Constraints []ConstraintSchema
}

type ColumnSchema struct {
    Name         string
    Type         string
    Nullable     bool
    Default      *string
    IsPrimaryKey bool
    IsUnique     bool
}

func (g *MigrationGenerator) GenerateCreateTable(table TableSchema) string {
    var sql strings.Builder
    sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", table.Name))
    
    for i, col := range table.Columns {
        sql.WriteString(fmt.Sprintf("    %s %s", col.Name, col.Type))
        
        if !col.Nullable {
            sql.WriteString(" NOT NULL")
        }
        
        if col.Default != nil {
            sql.WriteString(fmt.Sprintf(" DEFAULT %s", *col.Default))
        }
        
        if i < len(table.Columns)-1 {
            sql.WriteString(",")
        }
        sql.WriteString("\n")
    }
    
    // Adicionar constraints
    for _, constraint := range table.Constraints {
        sql.WriteString(fmt.Sprintf("    CONSTRAINT %s %s,\n", constraint.Name, constraint.Definition))
    }
    
    sql.WriteString(");")
    return sql.String()
}
```

## 🔧 Configuração de Nomenclatura no Bun ORM

### **Configuração Automática da Conversão**
```go
// internal/infra/database/bun.go
func NewBunConnection(cfg Config) (*BunConnection, error) {
    // ... código existente ...
    
    // Criar instância Bun com naming strategy automática
    db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
    
    // 🎯 Configurar conversão automática camelCase -> snake_case
    // Esta configuração garante que todos os fields em camelCase sejam
    // automaticamente convertidos para snake_case no PostgreSQL
    db.RegisterModel((*models.SessionModel)(nil))
    
    // ... resto do código ...
}
```

### **Vantagens da Abordagem**
- ✅ **Código Idiomático**: Go usa camelCase naturalmente
- ✅ **SQL Padrão**: PostgreSQL usa snake_case como convenção
- ✅ **Zero Config**: Bun faz conversão automaticamente
- ✅ **Manutenibilidade**: Mudanças de nome refletem automaticamente no banco
- ✅ **Consistência**: Um padrão único para todo o projeto

## 🔄 Fluxo de Desenvolvimento Futuro

### **1. Desenvolvedor Altera Model**
```go
// Developer adiciona novo campo no SessionModel
type SessionModel struct {
    // ... campos existentes
    // 🎯 Usando camelCase em Go (conversão automática para snake_case no PostgreSQL)
    LastActivity *time.Time `bun:"lastActivity" json:"lastActivity"` // 🆕 NOVO CAMPO
}
```

### **2. Sistema Detecta Mudança Automaticamente**
```bash
$ make db-diff
🔍 Detectando mudanças no schema...
┌─────────────────────────────────────────┐
│ MUDANÇAS DETECTADAS                     │
├─────────────────────────────────────────┤
│ Tabela: sessions                        │
│ + Adicionar coluna: last_activity       │
│   Tipo: TIMESTAMP WITH TIME ZONE NULL  │
│   Origin: lastActivity (Go field)      │
└─────────────────────────────────────────┘
```

### **3. Gerar Migration Automaticamente**
```bash
$ make db-generate name="add_session_last_activity"
🚀 Gerando migration...
✅ Migration criada: 20250810230000_add_session_last_activity.go
📝 Conversão automática: lastActivity -> last_activity
```

### **4. Migration Gerada Automaticamente**
```go
// internal/infra/database/migrations/auto/20250810230000_add_session_last_activity.go
func init() {
    Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
        fmt.Print(" [UP] adding column last_activity to sessions...")
        
        // 📝 Note: lastActivity (Go) -> last_activity (PostgreSQL)
        _, err := db.ExecContext(ctx, `
            ALTER TABLE sessions
            ADD COLUMN last_activity TIMESTAMP WITH TIME ZONE NULL;
        `)
        
        return err
    }, func(ctx context.Context, db *bun.DB) error {
        fmt.Print(" [DOWN] removing column last_activity from sessions...")
        
        _, err := db.ExecContext(ctx, `
            ALTER TABLE sessions
            DROP COLUMN IF EXISTS last_activity;
        `)
        
        return err
    })
}
```

### **5. Aplicar Migration**
```bash
$ make db-migrate
🗄️ Aplicando migrations...
✅ Migration 20250810230000_add_session_last_activity aplicada
```

## 📊 Benefícios da Nova Arquitetura

### **🔒 Segurança**
- ✅ Impossível desincronização entre código e banco
- ✅ Rollbacks automáticos e seguros
- ✅ Validação de integridade contínua

### **🚀 Produtividade**  
- ✅ Zero migrations manuais
- ✅ Detecção automática de mudanças
- ✅ Fluxo de desenvolvimento simplificado

### **🔧 Manutenibilidade**
- ✅ Single source of truth (models)
- ✅ Histórico completo de mudanças
- ✅ Clean Architecture preservada

### **📈 Escalabilidade**
- ✅ Suporte para múltiples models
- ✅ Migrations complexas automáticas
- ✅ Ambiente de desenvolvimento alinhado com produção

## ⚠️ Considerações e Riscos

### **Riscos Identificados**
1. **Dados Existentes**: Precisa migrar dados da estrutura atual
2. **Downtime**: Algumas migrations podem exigir parada temporária
3. **Rollback Complexo**: Mudanças estruturais podem ser irreversíveis

### **Mitigações**
1. **Backup Automático**: Sempre antes de aplicar migrations
2. **Dry Run**: Testar migrations em ambiente de desenvolvimento
3. **Rollback Testing**: Validar todas as operações de rollback

## 🎯 Critérios de Sucesso

- [ ] ✅ Zero migrations SQL manuais
- [ ] ✅ 100% sincronização entre models e banco
- [ ] ✅ Fluxo de desenvolvimento < 30 segundos para mudanças simples
- [ ] ✅ Zero regressões na aplicação existente
- [ ] ✅ Cobertura de testes > 90% para o sistema de migrations
- [ ] ✅ Documentação completa para desenvolvedores

---

**Próximos Passos**: Iniciar implementação seguindo o plano detalhado acima, começando pela limpeza das inconsistências e padronização das entidades.