# WazMeow - API REST para WhatsApp

Uma API REST completa para gerenciar sessões do WhatsApp usando Go e a biblioteca whatsmeow.

## 🚀 Funcionalidades

- ✅ Criar e gerenciar múltiplas sessões do WhatsApp
- ✅ Autenticação via QR Code
- ✅ Emparelhamento por telefone
- ✅ Configuração de proxy
- ✅ Conexão e desconexão de sessões
- ✅ Logout de sessões
- ✅ Listagem e informações detalhadas das sessões
- ✅ Event handlers completos para mensagens, presença, confirmações de leitura
- ✅ Reconexão automática de sessões na inicialização
- ✅ Gerenciamento de mídia (imagens, áudios, vídeos, documentos)
- ✅ Sistema de webhooks para eventos
- ✅ Graceful shutdown com desconexão de todas as sessões

## 🏗️ Tecnologias

- **Backend**: Go 1.23+ com Clean Architecture
- **ORM**: Bun ORM com auto-migrações
- **Banco de dados**: PostgreSQL 15+
- **WhatsApp**: whatsmeow library com event handlers completos
- **HTTP Router**: Chi v5
- **Containerização**: Docker & Docker Compose

## 📱 Implementação WhatsApp

A implementação do WhatsApp foi baseada no arquivo de referência `@reference/wuzapi/wmiau.go` e inclui:

### 🔧 Componentes Principais

- **WhatsAppClient**: Wrapper completo do cliente whatsmeow com event handlers
- **SessionManager**: Gerenciador de sessões ativas com thread-safety
- **ClientFactory**: Factory para criação e configuração de clientes
- **Event Handlers**: Tratamento completo de eventos do WhatsApp

### 📨 Eventos Suportados

- **Conexão**: Connected, Disconnected, LoggedOut
- **Autenticação**: QR Code generation, PairSuccess
- **Mensagens**: Recebimento de mensagens de texto e mídia
- **Confirmações**: Read receipts, delivery confirmations
- **Presença**: Online/offline status, chat presence
- **Mídia**: Processamento de imagens, áudios, vídeos e documentos

### 🔄 Funcionalidades Avançadas

- **Reconexão Automática**: Sessões conectadas são automaticamente reconectadas na inicialização
- **Graceful Shutdown**: Desconexão limpa de todas as sessões ao parar o servidor
- **Thread Safety**: Operações thread-safe em todos os componentes
- **Error Handling**: Tratamento robusto de erros e recuperação de falhas

## 📋 Endpoints da API

| Método | Endpoint                                      | Descrição                                                                 |
|--------|-----------------------------------------------|--------------------------------------------------------------------------|
| POST   | `/sessions/add`                               | Cria uma nova sessão do WhatsApp                                        |
| GET    | `/sessions/list`                              | Lista todas as sessões ativas e registradas no sistema                  |
| GET    | `/sessions/{sessionID}/info`                  | Retorna as informações detalhadas de uma sessão específica              |
| DELETE | `/sessions/{sessionID}`                       | Remove permanentemente uma sessão existente do sistema                  |
| POST   | `/sessions/{sessionID}/connect`               | Estabelece a conexão da sessão com o WhatsApp                           |
| POST   | `/sessions/{sessionID}/logout`                | Faz logout da sessão do WhatsApp, encerrando a comunicação              |
| GET    | `/sessions/{sessionID}/qr`                    | Gera e retorna o QR Code necessário para autenticar a sessão            |
| POST   | `/sessions/{sessionID}/pairphone`             | Emparelha um telefone com a sessão                                      |
| POST   | `/sessions/{sessionID}/proxy/set`             | Configura proxy para a sessão                                           |
| GET    | `/health`                                     | Health check da API                                                      |

## 🛠️ Instalação e Execução

### Pré-requisitos

- Go 1.23 ou superior
- PostgreSQL 15+
- Docker (opcional, para desenvolvimento)

### Instalação

1. Clone o repositório:
```bash
git clone <repository-url>
cd wazmeow
```

2. Instale as dependências:
```bash
go mod tidy
```

3. Execute o servidor:
```bash
go run cmd/server/main.go
```

O servidor será iniciado na porta 8080 por padrão.

### Variáveis de Ambiente

- `PORT`: Porta do servidor (padrão: 8080)
- `LOG_LEVEL`: Nível de log (DEBUG, INFO, WARN, ERROR - padrão: INFO)
- `DATA_DIR`: Diretório para armazenar dados (padrão: ./data)

## 📖 Exemplos de Uso

### Criar uma nova sessão

```bash
curl -X POST http://localhost:8080/sessions/add \
  -H "Content-Type: application/json" \
  -d '{"name": "Minha Sessão"}'
```

### Listar todas as sessões

```bash
curl http://localhost:8080/sessions/list
```

### Obter QR Code para autenticação

```bash
curl http://localhost:8080/sessions/{sessionID}/qr
```

### Conectar uma sessão

```bash
curl -X POST http://localhost:8080/sessions/{sessionID}/connect
```

### Emparelhar telefone

```bash
curl -X POST http://localhost:8080/sessions/{sessionID}/pairphone \
  -H "Content-Type: application/json" \
  -d '{"phone": "+5511999999999"}'
```

### Configurar proxy

```bash
curl -X POST http://localhost:8080/sessions/{sessionID}/proxy/set \
  -H "Content-Type: application/json" \
  -d '{
    "type": "http",
    "host": "proxy.example.com",
    "port": 8080,
    "username": "user",
    "password": "pass"
  }'
```

### Fazer logout

```bash
curl -X POST http://localhost:8080/sessions/{sessionID}/logout
```

### Remover sessão

```bash
curl -X DELETE http://localhost:8080/sessions/{sessionID}
```

## 📁 Estrutura do Projeto

```
wazmeow/
├── cmd/
│   └── server/
│       └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── http/
│   │   ├── handlers/
│   │   │   └── session_handler.go  # Handlers HTTP
│   │   └── router.go            # Configuração das rotas
│   ├── models/
│   │   └── session.go           # Modelos de dados
│   └── services/
│       └── session_service.go   # Lógica de negócio
├── pkg/
│   └── logger/
│       └── logger.go            # Logger personalizado
├── go.mod                       # Dependências do Go
└── README.md                    # Documentação
```

## 🔧 Tecnologias Utilizadas

- **Go 1.21**: Linguagem de programação
- **Chi Router**: Roteador HTTP leve e rápido
- **WhatsApp Web Multi-Device**: Biblioteca whatsmeow
- **SQLite**: Banco de dados para armazenar sessões
- **UUID**: Geração de identificadores únicos

## 📝 Formato das Respostas

Todas as respostas da API seguem o formato padrão:

```json
{
  "success": true,
  "message": "Mensagem descritiva",
  "data": {}, // Dados da resposta (opcional)
  "error": "" // Mensagem de erro (opcional)
}
```

## 🚨 Tratamento de Erros

A API retorna códigos de status HTTP apropriados:

- `200`: Sucesso
- `400`: Erro de validação ou dados inválidos
- `404`: Recurso não encontrado
- `500`: Erro interno do servidor

## 🔒 Segurança

- Validação de entrada em todos os endpoints
- Tratamento seguro de erros
- Logs estruturados para auditoria
- Suporte a CORS configurável

## 📞 Suporte

Para dúvidas ou problemas, abra uma issue no repositório do projeto.

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.
