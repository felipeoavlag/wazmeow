#!/bin/bash

# Script para gerar documentação Swagger do WazMeow API
# Uso: ./scripts/generate-docs.sh

set -e

echo "🔄 Gerando documentação Swagger..."

# Verificar se swag está instalado
if ! command -v swag &> /dev/null; then
    echo "❌ swag não encontrado. Instalando..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Criar diretório docs se não existir
mkdir -p docs

# Gerar documentação
echo "📝 Executando swag init..."
swag init -g cmd/server/main.go -o docs/ --parseDependency --parseInternal

# Verificar se os arquivos foram gerados
if [ -f "docs/docs.go" ] && [ -f "docs/swagger.json" ] && [ -f "docs/swagger.yaml" ]; then
    echo "✅ Documentação Swagger gerada com sucesso!"
    echo "📁 Arquivos gerados:"
    echo "   - docs/docs.go"
    echo "   - docs/swagger.json"
    echo "   - docs/swagger.yaml"
    echo ""
    echo "🌐 Para visualizar a documentação:"
    echo "   1. Execute o servidor: go run cmd/server/main.go"
    echo "   2. Acesse: http://localhost:8080/swagger/"
else
    echo "❌ Erro ao gerar documentação"
    exit 1
fi
