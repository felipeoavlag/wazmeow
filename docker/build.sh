#!/bin/bash

# Script para build e push da imagem WazMeow para Docker Hub
# Uso: ./build.sh [tag] [docker-hub-username]

set -e

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configurações padrão
DEFAULT_TAG="latest"
IMAGE_NAME="wazmeow"

# Função para exibir ajuda
show_help() {
    echo -e "${GREEN}WazMeow Docker Build Script${NC}"
    echo ""
    echo "Uso: $0 [TAG] [DOCKER_HUB_USERNAME]"
    echo ""
    echo "Parâmetros:"
    echo "  TAG                 Tag da imagem (padrão: latest)"
    echo "  DOCKER_HUB_USERNAME Username do Docker Hub (obrigatório para push)"
    echo ""
    echo "Exemplos:"
    echo "  $0                           # Build local com tag 'latest'"
    echo "  $0 v1.0.0                    # Build local com tag 'v1.0.0'"
    echo "  $0 latest meuusername        # Build e push para meuusername/wazmeow:latest"
    echo "  $0 v1.0.0 meuusername        # Build e push para meuusername/wazmeow:v1.0.0"
    echo ""
}

# Verificar se foi solicitada ajuda
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_help
    exit 0
fi

# Obter parâmetros
TAG=${1:-$DEFAULT_TAG}
DOCKER_HUB_USERNAME=$2

# Definir nome completo da imagem
if [[ -n "$DOCKER_HUB_USERNAME" ]]; then
    FULL_IMAGE_NAME="$DOCKER_HUB_USERNAME/$IMAGE_NAME:$TAG"
    LOCAL_IMAGE_NAME="$IMAGE_NAME:$TAG"
else
    FULL_IMAGE_NAME="$IMAGE_NAME:$TAG"
    LOCAL_IMAGE_NAME="$FULL_IMAGE_NAME"
fi

echo -e "${GREEN}🐳 WazMeow Docker Build${NC}"
echo -e "${GREEN}========================${NC}"
echo -e "📦 Imagem: ${YELLOW}$FULL_IMAGE_NAME${NC}"
echo -e "🏷️  Tag: ${YELLOW}$TAG${NC}"

if [[ -n "$DOCKER_HUB_USERNAME" ]]; then
    echo -e "👤 Docker Hub: ${YELLOW}$DOCKER_HUB_USERNAME${NC}"
    echo -e "🚀 Push: ${GREEN}Sim${NC}"
else
    echo -e "🚀 Push: ${YELLOW}Não (apenas build local)${NC}"
fi

echo ""

# Verificar se Docker está rodando
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}❌ Docker não está rodando ou não está acessível${NC}"
    exit 1
fi

# Ir para o diretório raiz do projeto (um nível acima)
cd "$(dirname "$0")/.."

# Build da imagem usando o Dockerfile na pasta docker
echo -e "${GREEN}🔨 Iniciando build da imagem...${NC}"
docker build -f docker/Dockerfile -t "$LOCAL_IMAGE_NAME" .

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Build concluído com sucesso!${NC}"
else
    echo -e "${RED}❌ Erro durante o build${NC}"
    exit 1
fi

# Se username foi fornecido, fazer tag e push
if [[ -n "$DOCKER_HUB_USERNAME" ]]; then
    echo ""
    echo -e "${GREEN}🏷️ Criando tag para Docker Hub...${NC}"
    
    # Se a imagem local tem nome diferente, criar tag
    if [[ "$LOCAL_IMAGE_NAME" != "$FULL_IMAGE_NAME" ]]; then
        docker tag "$LOCAL_IMAGE_NAME" "$FULL_IMAGE_NAME"
    fi
    
    echo -e "${GREEN}🚀 Fazendo push para Docker Hub...${NC}"
    echo -e "${YELLOW}⚠️ Certifique-se de estar logado no Docker Hub (docker login)${NC}"
    
    docker push "$FULL_IMAGE_NAME"
    
    if [[ $? -eq 0 ]]; then
        echo ""
        echo -e "${GREEN}🎉 Imagem publicada com sucesso no Docker Hub!${NC}"
        echo -e "${GREEN}📦 Imagem disponível em: ${YELLOW}$FULL_IMAGE_NAME${NC}"
        echo ""
        echo -e "${GREEN}Para usar a imagem:${NC}"
        echo -e "  docker pull $FULL_IMAGE_NAME"
        echo -e "  docker run -p 8080:8080 -e DB_HOST=seu_db_host $FULL_IMAGE_NAME"
    else
        echo -e "${RED}❌ Erro durante o push${NC}"
        exit 1
    fi
else
    echo ""
    echo -e "${GREEN}✅ Imagem criada localmente: ${YELLOW}$LOCAL_IMAGE_NAME${NC}"
    echo ""
    echo -e "${GREEN}Para testar a imagem:${NC}"
    echo -e "  docker run -p 8080:8080 -e DB_HOST=seu_db_host $LOCAL_IMAGE_NAME"
    echo ""
    echo -e "${GREEN}Para fazer push para Docker Hub:${NC}"
    echo -e "  $0 $TAG SEU_USERNAME_DOCKERHUB"
fi

echo ""
echo -e "${GREEN}📋 Informações da imagem:${NC}"
docker images "$LOCAL_IMAGE_NAME" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
