# WazMeow - WhatsApp Session Management API
# Makefile for development and deployment

.PHONY: help build run clean deps fmt lint up down restart status logs logs-app logs-db logs-dbgate db-shell db-admin db-reset cleanup dev build-prod

# Default target
help: ## Show this help message
	@echo "WazMeow - WhatsApp Session Management API"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
deps: ## Install dependencies
	go mod tidy
	go mod download

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

build: ## Build the application
	go build -o bin/wazmeow cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Docker Compose Commands
up: ## Start all services with Docker Compose
	@echo "🚀 Starting WazMeow services..."
	docker-compose up -d --build
	@echo "✅ Services started! API available at http://localhost:8080"

down: ## Stop all services
	@echo "🛑 Stopping WazMeow services..."
	docker-compose down
	@echo "✅ Services stopped!"

restart: ## Restart all services
	@echo "🔄 Restarting WazMeow services..."
	docker-compose down
	docker-compose up -d --build
	@echo "✅ Services restarted!"

status: ## Show service status
	@echo "📊 Service Status:"
	docker-compose ps

logs: ## Show logs for all services
	docker-compose logs -f

logs-app: ## Show logs for WazMeow app only
	docker-compose logs -f wazmeow

logs-db: ## Show logs for PostgreSQL only
	docker-compose logs -f postgres

logs-dbgate: ## Show logs for DBGate only
	docker-compose logs -f dbgate

db-shell: ## Connect to PostgreSQL database
	@echo "🗄️  Connecting to PostgreSQL..."
	docker-compose exec postgres psql -U postgres -d wazmeow

db-admin: ## Open DBGate database admin interface
	@echo "🗄️  Opening DBGate database admin..."
	@echo "DBGate is available at: http://localhost:3000"
	@command -v open >/dev/null 2>&1 && open http://localhost:3000 || \
	command -v xdg-open >/dev/null 2>&1 && xdg-open http://localhost:3000 || \
	echo "Please open http://localhost:3000 in your browser"

db-reset: ## Reset database (WARNING: destroys all data)
	@echo "⚠️  Resetting database (this will destroy all data)..."
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo ""; \
		docker-compose down -v; \
		docker-compose up -d postgres; \
		echo "✅ Database reset complete!"; \
	else \
		echo ""; \
		echo "❌ Database reset cancelled."; \
	fi

cleanup: ## Clean up Docker resources
	@echo "🧹 Cleaning up Docker resources..."
	docker-compose down --volumes --remove-orphans
	docker system prune -f
	@echo "✅ Cleanup complete!"

# Development with hot reload (requires air)
dev: ## Run with hot reload
	air

# Production
build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/wazmeow cmd/server/main.go
