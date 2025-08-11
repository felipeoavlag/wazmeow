# 🚀 Guia de Desenvolvimento - WazMeow

Este projeto oferece duas formas de executar a aplicação, otimizadas para diferentes cenários.

## 📋 Pré-requisitos

- Docker e Docker Compose
- Go 1.23+ (apenas para desenvolvimento local)
- Git

## 🛠️ Desenvolvimento Local (Recomendado para desenvolvimento)

**Ideal para:** Ver QR codes diretamente no terminal, debugging, desenvolvimento ativo.

### Início Rápido

```bash
# Inicia banco de dados e roda a aplicação automaticamente
./scripts/dev.sh --run
```

### Passo a Passo

```bash
# 1. Inicia apenas o banco de dados
./scripts/dev.sh

# 2. Em outro terminal, roda a aplicação Go
go run cmd/server/main.go
```

### Vantagens do Desenvolvimento Local

- ✅ **QR codes aparecem diretamente no terminal**
- ✅ Hot reload com ferramentas como `air`
- ✅ Debugging direto no IDE
- ✅ Logs coloridos e formatados
- ✅ Acesso direto ao banco (localhost:5432)

### Serviços Disponíveis

- 🌐 **API**: http://localhost:8080
- 🗄️ **DBGate**: http://localhost:3000
- 🐘 **PostgreSQL**: localhost:5432

## 🐳 Produção Containerizada

**Ideal para:** Deploy, testes de integração, ambiente similar à produção.

### Início Rápido

```bash
# Inicia tudo containerizado
./scripts/prod.sh
```

### Manual

```bash
# Build e start
docker-compose up --build -d

# Ver logs (incluindo QR codes)
docker-compose logs -f wazmeow
```

### Vantagens da Produção

- ✅ Ambiente isolado e reproduzível
- ✅ Configuração via variáveis de ambiente
- ✅ Restart automático dos serviços
- ✅ Logs centralizados

## 📱 Testando WhatsApp

### 1. Criar uma Sessão

```bash
curl -X POST http://localhost:8080/sessions/add \
  -H "Content-Type: application/json" \
  -d '{"name": "minha-sessao"}'
```

### 2. Conectar (Gerar QR Code)

```bash
# Substitua {sessionID} pelo ID retornado
curl -X POST http://localhost:8080/sessions/{sessionID}/connect
```

### 3. Ver QR Code

**Desenvolvimento Local:**
- QR code aparece automaticamente no terminal onde você rodou `go run`

**Produção:**
```bash
# Ver logs para encontrar o QR code
docker-compose logs -f wazmeow

# Ou buscar especificamente por QR codes
docker-compose logs wazmeow | grep -A 20 -B 5 "QR"
```

**Via API (ambos):**
```bash
curl http://localhost:8080/sessions/{sessionID}/qr
```

### 4. Verificar Status

```bash
curl http://localhost:8080/sessions/{sessionID}/info
```

## 🔧 Configuração

### Desenvolvimento Local

As variáveis são configuradas automaticamente pelo script `dev.sh`:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=wazmeow
DB_PASSWORD=wazmeow123
DB_NAME=wazmeow
SERVER_PORT=8080
```

### Produção

Configure o arquivo `.env` (copie de `.env.example`):

```bash
cp .env.example .env
# Edite .env conforme necessário
```

## 🛑 Parando os Serviços

### Desenvolvimento Local
```bash
docker-compose -f docker-compose.dev.yml down
```

### Produção
```bash
docker-compose down
```

## 🐛 Troubleshooting

### QR Code não aparece
- **Local**: Verifique se está rodando `go run` diretamente no terminal
- **Produção**: Use `docker-compose logs -f wazmeow`

### Banco não conecta
```bash
# Verificar se PostgreSQL está rodando
docker ps | grep postgres

# Verificar logs do banco
docker-compose logs postgres
```

### Porta já em uso
```bash
# Verificar o que está usando a porta
lsof -i :8080
lsof -i :5432
```

## 📚 Próximos Passos

1. **Hot Reload**: Instale `air` para reload automático
2. **Testes**: Execute `go test ./...`
3. **Lint**: Use `golangci-lint run`
4. **Docs**: Acesse http://localhost:8080/docs (quando implementado)

---

💡 **Dica**: Use desenvolvimento local para codificar e produção para testar deploys!
