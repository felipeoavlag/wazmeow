# WazMeow - WhatsApp Session Management API

Uma API REST limpa e focada exclusivamente no gerenciamento de múltiplas sessões do WhatsApp, construída com Go e seguindo padrões idiomáticos.

## 🎯 Características

- **Foco Exclusivo**: Gerenciamento de múltiplas sessões WhatsApp via API REST
- **Arquitetura Limpa**: Separação clara de responsabilidades (Domain, Application, Infrastructure)
- **Zero SQL**: Uso exclusivo do Bun ORM com query builder (sem SQL manual)
- **CamelCase**: Consistência em Go, PostgreSQL e JSON
- **Logging Estruturado**: Sistema centralizado com zerolog
- **Chi Router**: Router HTTP rápido e idiomático

## 🏗️ Arquitetura

```
wazmeow/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── config/                     # Configurações
│   ├── domain/
│   │   ├── entities/               # Entidades de domínio
│   │   ├── repositories/           # Interfaces de repositório
│   │   └── services/               # Interfaces de serviços
│   ├── application/
│   │   ├── dto/                    # Request/Response DTOs
│   │   ├── usecases/               # Use cases
│   │   └── handlers/               # HTTP handlers
│   └── infra/
│       ├── database/               # Bun ORM + PostgreSQL
│       ├── whatsapp/               # Cliente WhatsApp
│       └── http/                   # Servidor HTTP
├── pkg/
│   └── logger/                     # Logger centralizado
└── go.mod
```

## 📋 Entidade Session

```go
type Session struct {
    ID           string                 // UUID único
    Name         string                 // Nome da sessão
    Status       SessionStatus          // disconnected|connecting|connected
    Phone        string                 // Número de telefone (opcional)
    DeviceJID    string                 // JID do dispositivo WhatsApp
    ProxyConfig  *ProxyConfig          // Configuração de proxy
    WebhookURL   string                 // URL do webhook para eventos
    Events       string                 // Eventos subscritos
    CreatedAt    time.Time             // Data de criação
    UpdatedAt    time.Time             // Data de atualização
}
```

## 🛣️ Endpoints da API

| Método | Endpoint                                      | Descrição                                                                 |
|--------|-----------------------------------------------|--------------------------------------------------------------------------|
| POST   | `/api/v1/sessions/add`                        | Cria uma nova sessão do WhatsApp                                        |
| GET    | `/api/v1/sessions/list`                       | Lista todas as sessões registradas                                      |
| GET    | `/api/v1/sessions/{sessionID}/info`           | Retorna informações detalhadas de uma sessão                            |
| DELETE | `/api/v1/sessions/{sessionID}`                | Remove permanentemente uma sessão                                        |
| POST   | `/api/v1/sessions/{sessionID}/connect`        | Estabelece conexão da sessão com o WhatsApp                             |
| POST   | `/api/v1/sessions/{sessionID}/logout`         | Faz logout da sessão do WhatsApp                                        |
| GET    | `/api/v1/sessions/{sessionID}/qr`             | Gera e retorna o QR Code para autenticação                              |
| POST   | `/api/v1/sessions/{sessionID}/pairphone`      | Emparelha um telefone com a sessão                                      |
| POST   | `/api/v1/sessions/{sessionID}/proxy/set`      | Configura proxy para a sessão                                           |

## 🚀 Configuração

### Variáveis de Ambiente

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=wazmeow
DB_PASSWORD=password
DB_NAME=wazmeow
DB_SSLMODE=disable
DB_DEBUG=false

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# WhatsApp
WA_DEBUG=false
WA_OS_NAME=Mac OS 10

# Logging
LOG_LEVEL=info
LOG_FORMAT=console
```

### Banco de Dados

O sistema usa PostgreSQL com uma única tabela `Sessions` em camelCase:

```sql
CREATE TABLE Sessions (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'disconnected',
    phone VARCHAR(20),
    deviceJID VARCHAR(255),
    proxyEnabled BOOLEAN DEFAULT FALSE,
    proxyURL TEXT,
    webhookURL TEXT,
    events TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔧 Desenvolvimento

### Pré-requisitos

- Go 1.23+
- PostgreSQL 12+

### Executar

```bash
# Instalar dependências
go mod tidy

# Executar servidor
go run cmd/server/main.go
```

### Estrutura de Resposta

Todas as respostas seguem o padrão:

```json
{
  "success": true,
  "message": "Operação realizada com sucesso",
  "data": { ... },
  "error": null
}
```

## 🎨 Padrões Utilizados

- **Domain-Driven Design**: Separação clara entre domínio, aplicação e infraestrutura
- **Repository Pattern**: Abstração da camada de dados
- **Use Case Pattern**: Lógica de negócio encapsulada
- **Dependency Injection**: Inversão de dependências
- **Clean Architecture**: Arquitetura limpa e testável

## 📦 Tecnologias

- **Go 1.23**: Linguagem principal
- **Chi Router**: Router HTTP
- **Bun ORM**: ORM com query builder (zero SQL)
- **PostgreSQL**: Banco de dados
- **Zerolog**: Logging estruturado
- **WhatsApp Web API**: Integração WhatsApp

## 🔍 Health Check

```bash
curl http://localhost:8080/health
```

Resposta:
```json
{
  "status": "ok",
  "service": "wazmeow"
}
```

## 📝 Licença

Este projeto está sob a licença MIT.
