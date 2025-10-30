# Makefile for Code Auto-Grader
# Usage: make <target>

.PHONY: help dev dev-watch db-up db-down db-reset db-dump db-restore migrate-up migrate-down migrate-create db-shell db-seed db-seed-users db-seed-exercises db-seed-force test test-coverage clean install-tools

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
NC     := \033[0m # No Color

# Configuration
MIGRATE_VERSION ?= v4.17.0

help: ## Show this help message
	@echo "$(GREEN)Code Auto-Grader - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# Development Commands

dev: db-up ## Start full development environment (database + migrations)
	@echo "$(GREEN)Starting development environment...$(NC)"
	@echo "Waiting for database to be ready..."
	@sleep 3
	@$(MAKE) migrate-up
	@echo "$(GREEN)Database ready!$(NC)"
	@echo ""
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  Terminal 1: cd app && go run main.go"
	@echo "  Terminal 2: cd web && npm run dev"

dev-watch: ## Run backend with auto-reload (requires air)
	@echo "$(GREEN)Starting backend with auto-reload...$(NC)"
	@cd app && air

# Database Commands

db-up: ## Start PostgreSQL containers (dev + test)
	@echo "$(GREEN)Starting PostgreSQL containers...$(NC)"
	@docker compose up -d postgres postgres_test
	@echo "$(GREEN)Waiting for databases to be healthy...$(NC)"
	@sleep 2
	@echo "$(GREEN)Databases are ready!$(NC)"

db-down: ## Stop PostgreSQL containers
	@echo "$(YELLOW)Stopping PostgreSQL containers...$(NC)"
	@docker compose down

db-reset: ## Reset database (WARNING: deletes all data)
	@echo "$(YELLOW)⚠️  WARNING: This will DELETE ALL DATA!$(NC)"
	@read -p "Type 'yes' to continue: " confirm && [ "$$confirm" = "yes" ] || (echo "Cancelled." && exit 1)
	@echo "$(YELLOW)Resetting database...$(NC)"
	@docker compose down -v
	@docker compose up -d postgres postgres_test
	@sleep 3
	@$(MAKE) migrate-up
	@echo "$(GREEN)Database reset complete!$(NC)"

db-shell: ## Open psql shell to development database
	@docker compose exec postgres psql -U code_grader -d code_auto_grader

db-shell-test: ## Open psql shell to test database
	@docker compose exec postgres_test psql -U code_grader -d code_auto_grader_test

db-dump: ## Dump database to file (Usage: make db-dump FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		FILE="backup_$$(date +%Y%m%d_%H%M%S).sql"; \
	else \
		FILE="$(FILE)"; \
	fi; \
	echo "$(GREEN)Dumping database to $$FILE...$(NC)"; \
	docker compose exec -T postgres pg_dump -U code_grader code_auto_grader > $$FILE; \
	echo "$(GREEN)✓ Database dumped to $$FILE$(NC)"

db-restore: ## Restore database from file (Usage: make db-restore FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "$(YELLOW)Error: FILE is required. Usage: make db-restore FILE=backup.sql$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Restoring database from $(FILE)...$(NC)"
	@docker compose exec -T postgres psql -U code_grader code_auto_grader < $(FILE)
	@echo "$(GREEN)✓ Database restored from $(FILE)$(NC)"

# Migration Commands (using golang-migrate)

install-tools: ## Install development tools (migrate, air)
	@echo "$(GREEN)Installing development tools...$(NC)"
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
	@go install github.com/cosmtrek/air@latest
	@echo "$(GREEN)✓ Tools installed: migrate, air$(NC)"
	@echo "$(YELLOW)Make sure ~/go/bin is in your PATH$(NC)"

migrate-up: ## Run all pending migrations
	@echo "$(GREEN)Running migrations...$(NC)"
	@migrate -database "postgres://code_grader:dev_password_change_in_production@localhost:5432/code_auto_grader?sslmode=disable" -path ./app/migrations up
	@echo "$(GREEN)Migrations complete!$(NC)"

migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	@migrate -database "postgres://code_grader:dev_password_change_in_production@localhost:5432/code_auto_grader?sslmode=disable" -path ./app/migrations down 1

migrate-version: ## Show current migration version
	@migrate -database "postgres://code_grader:dev_password_change_in_production@localhost:5432/code_auto_grader?sslmode=disable" -path ./app/migrations version

migrate-force: ## Force migration version (use when dirty) - Usage: make migrate-force VERSION=1
	@if [ -z "$(VERSION)" ]; then \
		echo "$(YELLOW)Error: VERSION not specified$(NC)"; \
		echo "Usage: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Forcing migration version to $(VERSION)...$(NC)"
	@migrate -database "postgres://code_grader:dev_password_change_in_production@localhost:5432/code_auto_grader?sslmode=disable" -path ./app/migrations force $(VERSION)

migrate-create: ## Create new migration - Usage: make migrate-create NAME=add_user_avatar
	@if [ -z "$(NAME)" ]; then \
		echo "$(YELLOW)Error: NAME not specified$(NC)"; \
		echo "Usage: make migrate-create NAME=add_user_avatar"; \
		exit 1; \
	fi
	@echo "$(GREEN)Creating new migration: $(NAME)$(NC)"
	@migrate create -ext sql -dir ./app/migrations -seq $(NAME)

# Seed Commands

db-seed: ## Seed database with all content (from content/ directory)
	@echo "$(GREEN)Seeding database from content/...$(NC)"
	@cd app && go run cmd/seed/main.go --all
	@echo "$(GREEN)✓ Seeding complete!$(NC)"

db-seed-users: ## Seed only users
	@echo "$(GREEN)Seeding users...$(NC)"
	@cd app && go run cmd/seed/main.go --users

db-seed-exercises: ## Seed only exercises
	@echo "$(GREEN)Seeding exercises...$(NC)"
	@cd app && go run cmd/seed/main.go --exercises

db-seed-force: ## Force reseed (overwrites existing data)
	@echo "$(YELLOW)⚠️  WARNING: This will overwrite existing seed data!$(NC)"
	@read -p "Type 'yes' to continue: " confirm && [ "$$confirm" = "yes" ] || (echo "Cancelled." && exit 1)
	@cd app && go run cmd/seed/main.go --all --force
	@echo "$(GREEN)✓ Force seeding complete!$(NC)"

# Test Commands

test: db-up ## Run tests with test database
	@echo "$(GREEN)Running tests...$(NC)"
	@$(MAKE) test-migrate
	@cd app && TEST_DATABASE_URL="postgres://code_grader:test_password@localhost:5433/code_auto_grader_test?sslmode=disable" go test -v ./...

test-coverage: db-up ## Run tests with coverage report
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@$(MAKE) test-migrate
	@cd app && TEST_DATABASE_URL="postgres://code_grader:test_password@localhost:5433/code_auto_grader_test?sslmode=disable" go test -v -coverprofile=coverage.out ./...
	@cd app && go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: app/coverage.html$(NC)"

test-migrate: ## Run migrations on test database
	@echo "$(GREEN)Migrating test database...$(NC)"
	@migrate -database "postgres://code_grader:test_password@localhost:5433/code_auto_grader_test?sslmode=disable" -path ./app/migrations up

# Build Commands

build: ## Build production binary
	@echo "$(GREEN)Building backend...$(NC)"
	@cd app && go build -o code-auto-grader main.go
	@echo "$(GREEN)Building frontend...$(NC)"
	@cd web && npm run build
	@echo "$(GREEN)Build complete!$(NC)"

# Cleanup Commands

clean: ## Clean up containers, volumes, and build artifacts
	@echo "$(YELLOW)Cleaning up...$(NC)"
	@docker compose down -v
	@rm -f app/code-auto-grader
	@rm -f app/coverage.out app/coverage.html
	@rm -rf web/dist
	@echo "$(GREEN)Cleanup complete!$(NC)"

# Quick Start (for first-time setup)

setup: ## First-time setup (installs tools, starts DB, runs migrations)
	@echo "$(GREEN)Setting up Code Auto-Grader for the first time...$(NC)"
	@echo ""
	@$(MAKE) install-tools
	@if [ ! -f app/.env ]; then \
		cp app/.env.example app/.env; \
		echo "$(GREEN)Created app/.env from template$(NC)"; \
		echo "$(YELLOW)⚠️  Please edit app/.env and add your GEMINI_API_KEY$(NC)"; \
		read -p "Press enter when ready..."; \
	fi
	@$(MAKE) db-up
	@sleep 3
	@$(MAKE) migrate-up
	@echo ""
	@echo "$(GREEN)✓ Setup complete! Run 'make dev' to start developing.$(NC)"