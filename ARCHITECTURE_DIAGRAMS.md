# Diagramas da Arquitetura de Migrations Automáticas

## 🏗️ Arquitetura Atual vs Proposta

### Arquitetura Atual (Problemática)
```mermaid
graph TB
    subgraph "Problemas Atuais"
        A[Domain/entities.go<br/>❌ COM TAGS BUN] 
        B[Domain/entity/session.go<br/>✅ Clean]
        C[Models/session.go<br/>❌ Inconsistente]
        D[Migration Manual SQL<br/>❌ Desalinhada]
        
        A -.->|"Confusão"| C
        B -->|"Conversão"| C
        C -->|"Repository"| E[Database]
        D -->|"Schema Fixo"| E
        
        style A fill:#ffcccc
        style D fill:#ffcccc
    end
```

### Arquitetura Proposta (Solução)
```mermaid
graph TB
    subgraph "Nova Arquitetura"
        A[Domain/entity/session.go<br/>✅ Single Source]
        B[Models/session.go<br/>✅ Com Tags Bun]
        C[Migration Generator<br/>🆕 Automático]
        D[Schema Differ<br/>🆕 Detector]
        
        A -->|"Clean Conversion"| B
        B -->|"Analyze"| D
        D -->|"Generate"| C
        C -->|"Auto Migration"| E[Database]
        B -->|"Repository"| E
        
        style A fill:#ccffcc
        style B fill:#ccffcc
        style C fill:#cceeff
        style D fill:#cceeff
    end
```

## 🔄 Fluxo de Desenvolvimento com Migrations Automáticas

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant Model as SessionModel
    participant Differ as SchemaDiffer
    participant Gen as MigrationGenerator
    participant DB as Database
    
    Dev->>Model: Adiciona novo campo
    Dev->>Differ: make db-diff
    Differ->>DB: Analisa schema atual
    Differ->>Model: Analisa schema esperado
    Differ-->>Dev: Lista mudanças detectadas
    
    Dev->>Gen: make db-generate
    Gen->>Differ: Pega diferenças
    Gen-->>Dev: Gera migration automática
    
    Dev->>DB: make db-migrate
    DB-->>Dev: Schema atualizado
    
    Note over Dev,DB: ✅ Zero SQL manual!
```

## 📊 Comparação: Manual vs Automático

### Fluxo Manual Atual
```mermaid
graph LR
    A[Alterar Model] --> B[Escrever SQL Manual]
    B --> C[Criar Migration File]
    C --> D[Testar SQL]
    D --> E{SQL OK?}
    E -->|Não| B
    E -->|Sim| F[Aplicar Migration]
    F --> G[Verificar Schema]
    G --> H{Sincronizado?}
    H -->|Não| B
    H -->|Sim| I[✅ Concluído]
    
    style B fill:#ffeecc
    style D fill:#ffeecc
    style G fill:#ffeecc
```

### Fluxo Automático Proposto  
```mermaid
graph LR
    A[Alterar Model] --> B[make db-diff]
    B --> C[make db-generate]
    C --> D[make db-migrate]
    D --> E[✅ Concluído]
    
    style B fill:#ccffcc
    style C fill:#ccffcc
    style D fill:#ccffcc
```

## 🗄️ Estrutura de Dados da Migration

```mermaid
erDiagram
    Migration {
        string Name
        string Hash
        time Timestamp
        string UpSQL
        string DownSQL
        bool IsAuto
    }
    
    MigrationDiff {
        string TableName
        array AddedColumns
        array RemovedColumns
        array ModifiedColumns
        array AddedIndexes
        array RemovedIndexes
    }
    
    ColumnChange {
        string Name
        string OldType
        string NewType
        bool OldNullable
        bool NewNullable
        string OldDefault
        string NewDefault
    }
    
    Migration ||--|| MigrationDiff : generates
    MigrationDiff ||--o{ ColumnChange : contains
```

## 🔧 Componentes do Sistema

```mermaid
graph TB
    subgraph "Sistema de Migrations Automáticas"
        subgraph "Análise"
            A[SchemaDiffer]
            B[ModelAnalyzer]
            C[DatabaseIntrospector]
        end
        
        subgraph "Geração"
            D[MigrationGenerator]
            E[SQLGenerator]
            F[FileGenerator]
        end
        
        subgraph "Execução"
            G[MigrationRunner]
            H[RollbackManager]
            I[ValidationEngine]
        end
        
        subgraph "CLI Commands"
            J[db diff]
            K[db generate]
            L[db migrate]
            M[db rollback]
            N[db status]
        end
    end
    
    A --> D
    B --> A
    C --> A
    D --> E
    E --> F
    F --> G
    G --> I
    
    J --> A
    K --> D
    L --> G
    M --> H
    N --> I
```

## 📈 Timeline de Implementação

```mermaid
gantt
    title Cronograma de Implementação
    dateFormat X
    axisFormat %d
    
    section Fase 1: Limpeza
    Remover entities.go        :done, p1, 0, 1
    Padronizar entity/session  :done, p2, 1, 2
    Atualizar SessionModel     :active, p3, 2, 4
    
    section Fase 2: Sistema Base
    SchemaDiffer              :p4, 4, 6
    MigrationGenerator        :p5, 6, 8
    CLI Commands              :p6, 7, 9
    
    section Fase 3: Migração
    Converter migration atual  :p7, 8, 10
    Testes de validação       :p8, 9, 11
    
    section Fase 4: Finalização
    Documentação              :p9, 10, 12
    Validação completa        :p10, 11, 13
```

## 🎯 Estados do Sistema

```mermaid
stateDiagram-v2
    [*] --> SchemaAnalysis
    
    SchemaAnalysis --> NoChanges: Schema em sinc
    SchemaAnalysis --> ChangesDetected: Diferenças encontradas
    
    NoChanges --> [*]
    
    ChangesDetected --> MigrationGeneration: Gerar migration
    
    MigrationGeneration --> MigrationReady: Migration criada
    
    MigrationReady --> MigrationApplied: Aplicar
    MigrationReady --> MigrationDiscarded: Descartar
    
    MigrationApplied --> SchemaValidation: Validar resultado
    MigrationDiscarded --> [*]
    
    SchemaValidation --> Success: ✅ Validação OK
    SchemaValidation --> RollbackRequired: ❌ Erro detectado
    
    Success --> [*]
    
    RollbackRequired --> RollbackExecuted: Executar rollback
    RollbackExecuted --> SchemaAnalysis: Revalidar
```

## 🔍 Detalhes Técnicos do SchemaDiffer

```mermaid
flowchart TD
    A[Models Bun] --> B[Extract Schema Info]
    C[Current Database] --> D[Introspect Schema]
    
    B --> E[Expected Schema]
    D --> F[Current Schema]
    
    E --> G[Schema Comparison]
    F --> G
    
    G --> H{Changes Found?}
    
    H -->|Yes| I[Generate MigrationDiff]
    H -->|No| J[Schema Synchronized]
    
    I --> K[Categorize Changes]
    K --> L[Table Changes]
    K --> M[Column Changes]  
    K --> N[Index Changes]
    K --> O[Constraint Changes]
    
    L --> P[MigrationDiff Object]
    M --> P
    N --> P
    O --> P
    
    P --> Q[Ready for Migration Generation]
```

---

Estes diagramas fornecem uma visão visual completa da arquitetura proposta e dos fluxos de trabalho que implementaremos.