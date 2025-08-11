#!/bin/bash

# Script para desenvolvimento local
# Este script inicia o banco de dados e roda a aplicação localmente
# para que você possa ver os QR codes diretamente no terminal

set -e

echo "🚀 Iniciando ambiente de desenvolvimento local..."

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para log colorido
log() {
    echo -e "${GREEN}[DEV]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    error "Docker não está rodando. Por favor, inicie o Docker primeiro."
    exit 1
fi

# Parar containers existentes se estiverem rodando
log "Parando containers existentes..."
docker-compose -f docker-compose.dev.yml down > /dev/null 2>&1 || true

# Iniciar banco de dados
log "Iniciando PostgreSQL e DBGate..."
docker-compose -f docker-compose.dev.yml up -d

# Aguardar banco estar pronto
log "Aguardando PostgreSQL ficar pronto..."
timeout=30
counter=0
while ! docker exec wazmeow-postgres-dev pg_isready -U wazmeow -d wazmeow > /dev/null 2>&1; do
    if [ $counter -ge $timeout ]; then
        error "Timeout aguardando PostgreSQL ficar pronto"
        exit 1
    fi
    sleep 1
    counter=$((counter + 1))
    echo -n "."
done
echo ""

log "PostgreSQL está pronto!"

# Configurar variáveis de ambiente para desenvolvimento local
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=wazmeow
export DB_PASSWORD=wazmeow123
export DB_NAME=wazmeow
export DB_SSLMODE=disable
export DB_DEBUG=false
export SERVER_PORT=8080
export LOG_LEVEL=info

info "Variáveis de ambiente configuradas:"
info "  DB_HOST=$DB_HOST"
info "  DB_PORT=$DB_PORT"
info "  DB_NAME=$DB_NAME"
info "  SERVER_PORT=$SERVER_PORT"

echo ""
log "🎯 Ambiente pronto! Agora você pode:"
echo ""
info "1. Rodar a aplicação localmente:"
echo "   go run cmd/server/main.go"
echo ""
info "2. Acessar DBGate em: http://localhost:3000"
echo ""
info "3. API estará disponível em: http://localhost:8080"
echo ""
warn "📱 QR codes do WhatsApp aparecerão diretamente no terminal!"
echo ""
info "Para parar o ambiente:"
echo "   docker-compose -f docker-compose.dev.yml down"
echo ""

# Se foi passado o parâmetro --run, roda a aplicação automaticamente
if [ "$1" = "--run" ]; then
    log "Iniciando aplicação..."
    echo ""
    go run cmd/server/main.go
fi
