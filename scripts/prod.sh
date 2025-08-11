#!/bin/bash

# Script para ambiente de produção
# Este script usa o docker-compose.yml principal para rodar tudo containerizado

set -e

echo "🚀 Iniciando ambiente de produção..."

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para log colorido
log() {
    echo -e "${GREEN}[PROD]${NC} $1"
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

# Verificar se .env existe
if [ ! -f .env ]; then
    warn "Arquivo .env não encontrado. Criando com valores padrão..."
    cp .env.example .env
fi

# Parar containers existentes se estiverem rodando
log "Parando containers existentes..."
docker-compose down > /dev/null 2>&1 || true

# Build e start dos containers
log "Fazendo build e iniciando containers..."
docker-compose up --build -d

log "Aguardando serviços ficarem prontos..."
sleep 10

# Verificar se os serviços estão rodando
if docker-compose ps | grep -q "Up"; then
    log "✅ Ambiente de produção iniciado com sucesso!"
    echo ""
    info "Serviços disponíveis:"
    info "  🌐 API: http://localhost:8080"
    info "  🗄️  DBGate: http://localhost:3000"
    echo ""
    info "Para ver logs:"
    echo "   docker-compose logs -f"
    echo ""
    info "Para ver logs apenas da aplicação:"
    echo "   docker-compose logs -f wazmeow"
    echo ""
    warn "📱 Para ver QR codes do WhatsApp:"
    echo "   docker-compose logs -f wazmeow | grep -A 20 -B 5 'QR'"
    echo ""
    info "Para parar:"
    echo "   docker-compose down"
else
    error "❌ Falha ao iniciar alguns serviços"
    docker-compose ps
    exit 1
fi
