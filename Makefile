# Makefile para WazMeow
# Comandos básicos para desenvolvimento e deploy

# Variáveis
APP_NAME = wazmeow
BINARY_NAME = wazmeow
MAIN_PATH = ./cmd/server
BUILD_DIR = ./bin
DOCKER_COMPOSE_FILE = docker-compose.yml

# Cores para output
GREEN = \033[0;32m
YELLOW = \033[0;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: help build run clean test deps docker-up docker-down docker-logs docker-restart dev install lint fmt vet check swagger-gen swagger-serve swagger-clean

# Comando padrão
help: ## Mostra esta ajuda
	@echo "$(GREEN)WazMeow - Comandos disponíveis:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

# Comandos de Build
build: ## Compila a aplicação
	@echo "$(GREEN)🔨 Compilando $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build concluído: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-linux: ## Compila para Linux (útil para Docker)
	@echo "$(GREEN)🔨 Compilando $(APP_NAME) para Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "$(GREEN)✅ Build Linux concluído: $(BUILD_DIR)/$(BINARY_NAME)-linux$(NC)"

# Comandos de Execução
run: build ## Compila e executa a aplicação
	@echo "$(GREEN)🚀 Executando $(APP_NAME)...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Executa em modo desenvolvimento (com go run)
	@echo "$(GREEN)🔥 Executando em modo desenvolvimento...$(NC)"
	@go run $(MAIN_PATH)

# Comandos de Teste e Qualidade
test: ## Executa todos os testes
	@echo "$(GREEN)🧪 Executando testes...$(NC)"
	@go test ./... -v

test-coverage: ## Executa testes com coverage
	@echo "$(GREEN)🧪 Executando testes com coverage...$(NC)"
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)📊 Coverage report: coverage.html$(NC)"

lint: ## Executa linter (golangci-lint)
	@echo "$(GREEN)🔍 Executando linter...$(NC)"
	@golangci-lint run

fmt: ## Formata o código
	@echo "$(GREEN)✨ Formatando código...$(NC)"
	@go fmt ./...

vet: ## Executa go vet
	@echo "$(GREEN)🔍 Executando go vet...$(NC)"
	@go vet ./...

check: fmt vet test ## Executa formatação, vet e testes

# Comandos de Dependências
deps: ## Baixa e organiza dependências
	@echo "$(GREEN)📦 Baixando dependências...$(NC)"
	@go mod download
	@go mod tidy

deps-update: ## Atualiza todas as dependências
	@echo "$(GREEN)📦 Atualizando dependências...$(NC)"
	@go get -u ./...
	@go mod tidy

# Comandos Docker Compose
docker-up: ## Inicia todos os serviços (PostgreSQL, Redis, DBGate, Webhook Tester)
	@echo "$(GREEN)🐳 Iniciando serviços Docker...$(NC)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "$(GREEN)✅ Serviços iniciados:$(NC)"
	@echo "  📊 DBGate (Admin DB): http://localhost:3000"
	@echo "  🔗 Webhook Tester: http://localhost:8090"
	@echo "  🐘 PostgreSQL: localhost:5432"
	@echo "  🔴 Redis: localhost:6379"

docker-down: ## Para todos os serviços
	@echo "$(YELLOW)🐳 Parando serviços Docker...$(NC)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-restart: ## Reinicia todos os serviços
	@echo "$(YELLOW)🐳 Reiniciando serviços Docker...$(NC)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) restart

docker-logs: ## Mostra logs dos serviços
	@echo "$(GREEN)📋 Logs dos serviços Docker:$(NC)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-clean: ## Remove containers, volumes e imagens
	@echo "$(RED)🧹 Limpando Docker (containers, volumes, imagens)...$(NC)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --rmi all

# Comandos de Limpeza
clean: ## Remove arquivos de build
	@echo "$(YELLOW)🧹 Limpando arquivos de build...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

clean-all: clean docker-clean ## Limpeza completa (build + docker)

# Comandos de Instalação
install: deps build ## Instala dependências e compila

install-tools: ## Instala ferramentas de desenvolvimento
	@echo "$(GREEN)🛠️ Instalando ferramentas de desenvolvimento...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✅ Ferramentas instaladas$(NC)"

# Comandos de Desenvolvimento Completo
setup: install-tools deps docker-up ## Setup completo para desenvolvimento
	@echo "$(GREEN)🎉 Setup de desenvolvimento concluído!$(NC)"
	@echo ""
	@echo "$(GREEN)Próximos passos:$(NC)"
	@echo "  1. Execute: $(YELLOW)make dev$(NC) para iniciar a aplicação"
	@echo "  2. Acesse: $(YELLOW)http://localhost:3000$(NC) para DBGate"
	@echo "  3. Acesse: $(YELLOW)http://localhost:8090$(NC) para Webhook Tester"

# Comandos de Status
status: ## Mostra status dos serviços Docker
	@echo "$(GREEN)📊 Status dos serviços:$(NC)"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) ps

# Comando para desenvolvimento rápido
quick: docker-up dev ## Inicia Docker e executa app em modo dev

# Comando para produção local
prod: build docker-up ## Build e inicia com Docker para simular produção
	@echo "$(GREEN)🚀 Executando em modo produção local...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME)

# Comandos de Documentação Swagger
swagger-gen: ## Gera documentação Swagger
	@echo "$(GREEN)📝 Gerando documentação Swagger...$(NC)"
	@if ! command -v swag &> /dev/null; then \
		echo "$(YELLOW)⚠️ swag não encontrado. Instalando...$(NC)"; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@mkdir -p docs
	@swag init -g cmd/server/main.go -o docs/ --parseDependency --parseInternal
	@echo "$(GREEN)✅ Documentação Swagger gerada com sucesso!$(NC)"
	@echo "$(GREEN)📁 Arquivos gerados: docs/docs.go, docs/swagger.json, docs/swagger.yaml$(NC)"

swagger-serve: swagger-gen dev ## Gera documentação e inicia servidor
	@echo "$(GREEN)🌐 Acesse a documentação em: http://localhost:8080/swagger/$(NC)"

swagger-clean: ## Remove arquivos de documentação gerados
	@echo "$(YELLOW)🧹 Removendo documentação Swagger...$(NC)"
	@rm -f docs/docs.go docs/swagger.json docs/swagger.yaml
	@echo "$(GREEN)✅ Documentação removida$(NC)"
