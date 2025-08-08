# 🐳 Docker Setup - WazMeow API

Este documento explica como usar o Docker para desenvolvimento e produção do WazMeow API.

## 📋 Serviços Disponíveis

O `docker-compose.yml` inclui os seguintes serviços:

### 🗄️ Banco de Dados
- **PostgreSQL 15**: Banco principal na porta `5432` (com Bun ORM e auto-migrações)
- **Redis 7**: Cache e sessões na porta `6379`

### 🔧 Ferramentas de Administração
- **DBGate**: Interface web unificada na porta `3000`

### 🚀 API (Opcional)
- **WazMeow API**: Aplicação principal na porta `8080`

## 🚀 Início Rápido

### 1. Iniciar Serviços de Desenvolvimento
```bash
# Inicia PostgreSQL, Redis e DBGate
make docker-up

# Ou usando docker-compose diretamente
docker-compose up -d postgres redis dbgate
```

### 2. Verificar Status
```bash
make docker-status
```

### 3. Acessar Interfaces Web

#### DBGate (Recomendado)
- **URL**: http://localhost:3000
- **Descrição**: Interface unificada para PostgreSQL e Redis
- **Login**: Não requer (pré-configurado)



### 4. Executar a API Localmente
```bash
# Com os serviços Docker rodando, execute a API localmente
make run
```

### 5. Ou Executar Tudo no Docker
```bash
# Inicia todos os serviços incluindo a API
make docker-up-all
```

## ⚙️ Configuração

### Variáveis de Ambiente

O arquivo `.env` contém as configurações. Para Docker, ajuste:

```bash
# Para desenvolvimento local (API fora do Docker)
DB_HOST=localhost
REDIS_HOST=localhost

# Para API dentro do Docker
DB_HOST=postgres
REDIS_HOST=redis
```

### Credenciais Padrão

#### PostgreSQL
- **Host**: localhost:5432
- **Database**: wazmeow
- **User**: postgres
- **Password**: password

#### Redis
- **Host**: localhost:6379
- **Password**: redispassword

## 🔧 Comandos Úteis

### Gerenciamento de Serviços
```bash
make docker-up          # Inicia serviços de desenvolvimento
make docker-up-all      # Inicia todos os serviços
make docker-down        # Para todos os serviços
make docker-restart     # Reinicia serviços
make docker-status      # Mostra status
make docker-logs        # Mostra logs
```

### Limpeza
```bash
make docker-clean       # Remove tudo (containers, volumes, imagens)
```

### Build e Deploy
```bash
make docker-build       # Constrói imagem da API
make docker-run         # Executa container da API
```

## 📊 Monitoramento

### Logs
```bash
# Todos os serviços
docker-compose logs -f

# Serviço específico
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f wazmeow-api
```

### Health Checks
Todos os serviços têm health checks configurados:
```bash
docker-compose ps
```

## 🗄️ Persistência de Dados

Os dados são persistidos em volumes Docker:
- `postgres_data`: Dados do PostgreSQL
- `redis_data`: Dados do Redis
- `pgadmin_data`: Configurações do PgAdmin

### Backup
```bash
# Backup PostgreSQL
docker-compose exec postgres pg_dump -U postgres wazmeow > backup.sql

# Restore PostgreSQL
docker-compose exec -T postgres psql -U postgres wazmeow < backup.sql
```

## 🔒 Segurança

### Produção
Para produção, altere as senhas padrão:

1. **PostgreSQL**: Altere `POSTGRES_PASSWORD`
2. **Redis**: Altere a senha no comando Redis
3. **PgAdmin**: Altere `PGADMIN_DEFAULT_PASSWORD`
4. **Redis Commander**: Altere `HTTP_PASSWORD`

### Rede
Os serviços estão isolados na rede `wazmeow-network`.

## 🐛 Troubleshooting

### Problemas Comuns

#### Porta já em uso
```bash
# Verificar o que está usando a porta
lsof -i :5432
lsof -i :6379
lsof -i :3000

# Parar serviços conflitantes ou alterar portas no docker-compose.yml
```

#### Problemas de permissão
```bash
# Limpar volumes e recriar
docker-compose down -v
docker-compose up -d
```

#### Logs de erro
```bash
# Verificar logs específicos
docker-compose logs postgres
docker-compose logs redis
```

### Reset Completo
```bash
# Para tudo e remove volumes
make docker-clean

# Recria tudo
make docker-up
```

## 📝 Desenvolvimento

### Workflow Recomendado

1. **Inicie os serviços de infraestrutura**:
   ```bash
   make docker-up
   ```

2. **Execute a API localmente** para desenvolvimento:
   ```bash
   make dev
   ```

3. **Use DBGate** para administrar bancos de dados:
   - Acesse http://localhost:3000

4. **Para testes de integração**, use a API no Docker:
   ```bash
   make docker-up-all
   ```

### Hot Reload
Para desenvolvimento com hot reload, use `air`:
```bash
# Instalar air
go install github.com/cosmtrek/air@latest

# Executar com hot reload
make watch
```

## 🚀 Deploy

### Build para Produção
```bash
# Build otimizado
make docker-build

# Tag para registry
docker tag wazmeow:latest your-registry/wazmeow:v1.0.0

# Push para registry
docker push your-registry/wazmeow:v1.0.0
```

### Docker Compose para Produção
Crie um `docker-compose.prod.yml` com:
- Senhas seguras
- Volumes externos
- Configurações de rede adequadas
- Health checks robustos
