# WazMeow API - Makefile
# =====================

# Variáveis
BINARY_NAME=wazmeow
SETUP_BINARY=setup
BUILD_DIR=bin
SERVER_CMD=cmd/server/main.go
SETUP_CMD=cmd/setup/main.go

# Comandos padrão
.PHONY: help build setup run clean test deps dev

# Ajuda
help: ## Mostra esta mensagem de ajuda
	@echo "WazMeow API - Comandos disponíveis:"
	@echo "=================================="
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build
build: ## Compila o servidor principal
	@echo "🔨 Compilando servidor..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SERVER_CMD)
	@echo "✅ Servidor compilado: $(BUILD_DIR)/$(BINARY_NAME)"

build-setup: ## Compila o aplicativo de configuração
	@echo "🔨 Compilando configurador..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(SETUP_BINARY) $(SETUP_CMD)
	@echo "✅ Configurador compilado: $(BUILD_DIR)/$(SETUP_BINARY)"

build-all: build build-setup ## Compila todos os binários

# Configuração
setup: build-setup ## Executa o configurador interativo
	@echo "⚙️  Iniciando configurador..."
	@./$(BUILD_DIR)/$(SETUP_BINARY)

# Execução
run: build ## Compila e executa o servidor
	@echo "🚀 Iniciando servidor..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Executa em modo desenvolvimento
	@echo "🔧 Iniciando em modo desenvolvimento..."
	@go run $(SERVER_CMD)

# Dependências
deps: ## Instala/atualiza dependências
	@echo "📦 Instalando dependências..."
	@go mod tidy
	@go mod download
	@echo "✅ Dependências atualizadas"

# Testes
test: ## Executa os testes
	@echo "🧪 Executando testes..."
	@go test -v ./...

test-coverage: ## Executa testes com cobertura
	@echo "🧪 Executando testes com cobertura..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Relatório de cobertura: coverage.html"

# Limpeza
clean: ## Remove arquivos compilados
	@echo "🧹 Limpando arquivos..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Limpeza concluída"

# Docker
docker-build: ## Constrói imagem Docker
	@echo "🐳 Construindo imagem Docker..."
	@docker build -t wazmeow:latest .

docker-run: ## Executa container Docker
	@echo "🐳 Executando container..."
	@docker run -p 8080:8080 --env-file .env wazmeow:latest

docker-up: ## Inicia serviços de desenvolvimento (PostgreSQL, Redis, DBGate)
	@echo "🐳 Iniciando serviços de desenvolvimento..."
	@docker-compose up -d postgres redis dbgate

docker-down: ## Para todos os serviços
	@echo "🐳 Parando serviços..."
	@docker-compose down

docker-logs: ## Mostra logs dos serviços
	@echo "📋 Logs dos serviços..."
	@docker-compose logs -f

docker-clean: ## Remove containers, volumes e imagens
	@echo "🧹 Limpando Docker..."
	@docker-compose down -v --rmi all
	@docker system prune -f

docker-restart: docker-down docker-up ## Reinicia os serviços

docker-status: ## Mostra status dos containers
	@echo "📊 Status dos containers..."
	@docker-compose ps



# Utilitários
fmt: ## Formata o código
	@echo "🎨 Formatando código..."
	@go fmt ./...
	@echo "✅ Código formatado"

lint: ## Executa linter
	@echo "🔍 Executando linter..."
	@golangci-lint run
	@echo "✅ Linting concluído"

mod-update: ## Atualiza módulos Go
	@echo "📦 Atualizando módulos..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Módulos atualizados"

# Instalação
install: build ## Instala o binário no sistema
	@echo "📥 Instalando binário..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Instalado em /usr/local/bin/$(BINARY_NAME)"

uninstall: ## Remove o binário do sistema
	@echo "🗑️  Removendo binário..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Binário removido"

# Informações
version: ## Mostra informações de versão
	@echo "WazMeow API Server"
	@echo "=================="
	@echo "Go version: $(shell go version)"
	@echo "Build date: $(shell date)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

# Banco de dados (utilitários)
db-create: ## Cria o banco de dados PostgreSQL
	@echo "🗄️  Criando banco de dados..."
	@createdb wazmeow || echo "Banco já existe ou erro na criação"

db-drop: ## Remove o banco de dados PostgreSQL
	@echo "🗑️  Removendo banco de dados..."
	@dropdb wazmeow || echo "Banco não existe ou erro na remoção"

db-reset: db-drop db-create ## Recria o banco de dados

# Desenvolvimento
watch: ## Executa com hot reload (requer air)
	@echo "👀 Iniciando hot reload..."
	@air -c .air.toml

# Release
release: clean deps test build-all ## Prepara release completo
	@echo "🎉 Release preparado!"
	@echo "Binários disponíveis em $(BUILD_DIR)/"
