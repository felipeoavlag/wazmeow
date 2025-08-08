# WazMeow - API REST para WhatsApp

Uma API REST completa para gerenciar sessões do WhatsApp usando Go, Chi Router e PostgreSQL com sistema de configuração avançado.

## 🚀 Funcionalidades

- ✅ Criar e gerenciar múltiplas sessões do WhatsApp
- ✅ Autenticação via QR Code
- ✅ Emparelhamento por telefone
- ✅ Configuração de proxy
- ✅ Conexão e desconexão de sessões
- ✅ Logout de sessões
- ✅ Listagem e informações detalhadas das sessões
- ✅ Banco de dados PostgreSQL
- ✅ Sistema de configuração avançado
- ✅ Aplicativo de configuração interativo
- ✅ Makefile com comandos úteis
- ✅ Suporte a CORS configurável
- ✅ Rate limiting
- ✅ Logs estruturados

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
- PostgreSQL 12 ou superior
- Make (opcional, mas recomendado)

### Instalação Rápida

1. **Clone o repositório:**
```bash
git clone <repository-url>
cd wazmeow
```

2. **Configure o projeto:**
```bash
make setup
```
Este comando irá executar um configurador interativo que criará o arquivo `.env` com todas as configurações necessárias.

3. **Instale as dependências:**
```bash
make deps
```

4. **Execute o servidor:**
```bash
make run
```

### Instalação Manual

1. **Clone e configure dependências:**
```bash
git clone <repository-url>
cd wazmeow
go mod tidy
```

2. **Configure o banco de dados PostgreSQL:**
```sql
CREATE DATABASE wazmeow;
CREATE USER wazmeow_user WITH PASSWORD 'sua_senha';
GRANT ALL PRIVILEGES ON DATABASE wazmeow TO wazmeow_user;
```

3. **Configure as variáveis de ambiente:**
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

4. **Compile e execute:**
```bash
go build -o bin/wazmeow cmd/server/main.go
./bin/wazmeow
```

## ⚙️ Configuração

### Configurador Interativo

Execute o configurador para uma configuração guiada:
```bash
make setup
# ou
./bin/setup
```

### Variáveis de Ambiente

O sistema suporta as seguintes variáveis de ambiente:

#### Banco de Dados
- `DB_HOST`: Host do PostgreSQL (padrão: localhost)
- `DB_PORT`: Porta do PostgreSQL (padrão: 5432)
- `DB_USER`: Usuário do PostgreSQL (padrão: postgres)
- `DB_PASSWORD`: Senha do PostgreSQL (padrão: password)
- `DB_NAME`: Nome do banco (padrão: wazmeow)
- `DB_SSLMODE`: Modo SSL (padrão: disable)

#### Servidor HTTP
- `SERVER_HOST`: Host do servidor (padrão: 0.0.0.0)
- `SERVER_PORT`: Porta do servidor (padrão: 8080)
- `SERVER_READ_TIMEOUT`: Timeout de leitura (padrão: 30s)
- `SERVER_WRITE_TIMEOUT`: Timeout de escrita (padrão: 30s)

#### Logs
- `LOG_LEVEL`: Nível de log (DEBUG, INFO, WARN, ERROR - padrão: INFO)
- `LOG_FORMAT`: Formato do log (json, text - padrão: json)

#### CORS
- `CORS_ALLOWED_ORIGINS`: Origens permitidas (padrão: *)
- `CORS_ALLOWED_METHODS`: Métodos permitidos
- `CORS_ALLOWED_HEADERS`: Headers permitidos

#### Sessões
- `MAX_SESSIONS`: Máximo de sessões simultâneas (padrão: 100)
- `SESSION_TIMEOUT`: Timeout das sessões (padrão: 3600s)

#### Segurança
- `API_KEY`: Chave da API (opcional)
- `RATE_LIMIT_REQUESTS`: Limite de requisições (padrão: 100)
- `RATE_LIMIT_WINDOW`: Janela de tempo (padrão: 1m)

#### Aplicação
- `DEBUG`: Modo debug (true/false - padrão: false)
- `ENVIRONMENT`: Ambiente (development/production - padrão: production)

## 🔧 Comandos Make

```bash
make help          # Mostra todos os comandos disponíveis
make build          # Compila o servidor
make build-setup    # Compila o configurador
make build-all      # Compila todos os binários
make setup          # Executa o configurador interativo
make run            # Compila e executa o servidor
make dev            # Executa em modo desenvolvimento
make deps           # Instala/atualiza dependências
make test           # Executa os testes
make clean          # Remove arquivos compilados
make fmt            # Formata o código
make install        # Instala o binário no sistema
```

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

## 📁 Estrutura do Projeto

```
wazmeow/
├── cmd/
│   ├── server/main.go           # Servidor principal
│   └── setup/main.go            # Configurador interativo
├── internal/
│   ├── config/config.go         # Sistema de configuração
│   ├── http/
│   │   ├── handlers/            # Handlers HTTP
│   │   └── router.go            # Configuração das rotas
│   ├── models/                  # Modelos de dados
│   └── services/                # Lógica de negócio
├── pkg/
│   └── logger/                  # Logger personalizado
├── bin/                         # Binários compilados
├── .env.example                 # Exemplo de configuração
├── Makefile                     # Comandos de automação
└── README.md                    # Documentação
```

## 🔒 Segurança

- Validação de entrada em todos os endpoints
- Tratamento seguro de erros
- Logs estruturados para auditoria
- Suporte a CORS configurável
- Rate limiting configurável
- Conexão segura com PostgreSQL
- Suporte a API Key opcional

## 🚀 Produção

Para ambiente de produção:

1. Configure `ENVIRONMENT=production`
2. Use `DEBUG=false`
3. Configure um `API_KEY` forte
4. Configure SSL no PostgreSQL
5. Use um proxy reverso (nginx/traefik)
6. Configure logs estruturados
7. Monitore com health checks

## 📞 Suporte

Para dúvidas ou problemas, abra uma issue no repositório do projeto.

## 📄 Licença

Este projeto está sob a licença MIT.
