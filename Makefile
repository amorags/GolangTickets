.PHONY: help dev build up down test clean deploy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start development environment
	docker compose up --build

up: ## Start all services in detached mode
	docker compose up -d

down: ## Stop all services
	docker compose down

build: ## Build all Docker images
	docker compose build

rebuild: ## Rebuild all images without cache
	docker compose build --no-cache

test: ## Run all tests
	@echo "Running Go tests..."
	go test -v ./...
	@echo "\nRunning Web tests..."
	cd web && bun test

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

logs: ## Show logs from all services
	docker compose logs -f

logs-api: ## Show API logs
	docker compose logs -f api

logs-web: ## Show web logs
	docker compose logs -f web

logs-db: ## Show database logs
	docker compose logs -f db

clean: ## Clean up containers, volumes, and images
	docker compose down -v
	docker system prune -f

restart: ## Restart all services
	docker compose restart

restart-api: ## Restart only API service
	docker compose restart api

restart-web: ## Restart only web service
	docker compose restart web

shell-api: ## Open shell in API container
	docker compose exec api sh

shell-web: ## Open shell in web container
	docker compose exec web sh

shell-db: ## Open PostgreSQL shell
	docker compose exec db psql -U user -d ticket_db

# Production commands
prod-up: ## Start production environment
	docker compose -f docker-compose.prod.yml up -d

prod-down: ## Stop production environment
	docker compose -f docker-compose.prod.yml down

prod-logs: ## Show production logs
	docker compose -f docker-compose.prod.yml logs -f

prod-pull: ## Pull latest images from registry
	docker compose -f docker-compose.prod.yml pull

prod-deploy: prod-pull ## Deploy latest images
	docker compose -f docker-compose.prod.yml up -d
	docker image prune -af

# Development helpers
install: ## Install dependencies
	go mod download
	cd web && bun install

fmt: ## Format code
	go fmt ./...
	cd web && bun run format || echo "No format script defined"

lint: ## Lint code
	golangci-lint run ./...
	cd web && bun run lint || echo "No lint script defined"

migration-create: ## Create a new migration (use name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: Please specify migration name. Usage: make migration-create name=create_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir migrations -seq $(name)

migration-up: ## Run all migrations
	migrate -path migrations -database "postgres://user:password@localhost:5432/ticket_db?sslmode=disable" up

migration-down: ## Rollback last migration
	migrate -path migrations -database "postgres://user:password@localhost:5432/ticket_db?sslmode=disable" down 1
