.PHONY: help db-up db-down db-drop db-reset db-seed db-seed-force db-status migrate-up migrate-down migrate-create seed test-db clean

# Default DATABASE_URL (override with environment variable)
DATABASE_URL ?= postgres://code_grader:dev_password@localhost:5432/code_auto_grader?sslmode=disable
TEST_DATABASE_URL ?= postgres://code_grader:test_password@localhost:5433/code_auto_grader_test?sslmode=disable

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Database Migration Commands

migrate-up: ## Run all pending migrations
	@echo "Running migrations..."
	migrate -path migrations -database "$(DATABASE_URL)" up
	@echo "✓ Migrations complete"

migrate-down: ## Rollback the last migration
	@echo "Rolling back last migration..."
	migrate -path migrations -database "$(DATABASE_URL)" down 1
	@echo "✓ Rollback complete"

migrate-drop: ## Drop all migrations (WARNING: destroys data)
	@echo "⚠️  WARNING: This will drop all tables and data"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		migrate -path migrations -database "$(DATABASE_URL)" drop -f; \
		echo "✓ All migrations dropped"; \
	else \
		echo "Cancelled"; \
	fi

migrate-version: ## Show current migration version
	@migrate -path migrations -database "$(DATABASE_URL)" version

migrate-force: ## Force migration to specific version (usage: make migrate-force VERSION=1)
	@echo "Forcing migration to version $(VERSION)..."
	migrate -path migrations -database "$(DATABASE_URL)" force $(VERSION)
	@echo "✓ Forced to version $(VERSION)"

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=add_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=add_users_table"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir migrations -seq $(NAME)
	@echo "✓ Created migration files for $(NAME)"

# Database Seeding Commands

db-seed: ## Seed database with content (users, tech trees, nodes, exercises)
	@echo "Seeding database..."
	@cd cmd/seed && go run main.go -all
	@echo "✓ Database seeded"

db-seed-users: ## Seed only users
	@echo "Seeding users..."
	@cd cmd/seed && go run main.go -users
	@echo "✓ Users seeded"

db-seed-trees: ## Seed only tech trees
	@echo "Seeding tech trees..."
	@cd cmd/seed && go run main.go -trees
	@echo "✓ Tech trees seeded"

db-seed-nodes: ## Seed only nodes
	@echo "Seeding nodes..."
	@cd cmd/seed && go run main.go -nodes
	@echo "✓ Nodes seeded"

db-seed-exercises: ## Seed only exercises
	@echo "Seeding exercises..."
	@cd cmd/seed && go run main.go -exercises
	@echo "✓ Exercises seeded"

db-seed-force: ## Force re-seed (deletes and re-inserts data)
	@echo "Force seeding database..."
	@cd cmd/seed && go run main.go -all -force
	@echo "✓ Database force seeded"

# Database Setup Commands

db-reset: migrate-drop migrate-up db-seed ## Reset database (drop, migrate, seed)
	@echo "✓ Database reset complete"

db-fresh: migrate-up db-seed ## Fresh start (migrate and seed without dropping)
	@echo "✓ Fresh database ready"

# Testing Commands

test-db-setup: ## Set up test database
	@echo "Setting up test database..."
	@DATABASE_URL="$(TEST_DATABASE_URL)" migrate -path migrations -database "$(TEST_DATABASE_URL)" up
	@echo "✓ Test database ready"

test-db-reset: ## Reset test database
	@echo "Resetting test database..."
	@DATABASE_URL="$(TEST_DATABASE_URL)" migrate -path migrations -database "$(TEST_DATABASE_URL)" drop -f || true
	@DATABASE_URL="$(TEST_DATABASE_URL)" migrate -path migrations -database "$(TEST_DATABASE_URL)" up
	@echo "✓ Test database reset"

test: ## Run all tests
	@echo "Running tests..."
	@TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -v ./...
	@echo "✓ Tests complete"

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Development Commands

run: ## Run the application
	@echo "Starting application..."
	@go run main.go

dev: ## Run with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
	@echo "Starting development server with auto-reload..."
	@air

build: ## Build the application binary
	@echo "Building application..."
	@go build -o code-auto-grader main.go
	@echo "✓ Built: code-auto-grader"

clean: ## Remove build artifacts and coverage files
	@echo "Cleaning build artifacts..."
	@rm -f code-auto-grader
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Docker Commands (if using Docker for Postgres)

docker-db-up: ## Start PostgreSQL in Docker
	@echo "Starting PostgreSQL container..."
	@docker run -d \
		--name code-grader-postgres \
		-e POSTGRES_USER=code_grader \
		-e POSTGRES_PASSWORD=dev_password \
		-e POSTGRES_DB=code_auto_grader \
		-p 5432:5432 \
		postgres:16-alpine
	@echo "✓ PostgreSQL started on localhost:5432"

docker-db-down: ## Stop PostgreSQL container
	@echo "Stopping PostgreSQL container..."
	@docker stop code-grader-postgres
	@docker rm code-grader-postgres
	@echo "✓ PostgreSQL stopped"

docker-db-logs: ## Show PostgreSQL logs
	@docker logs -f code-grader-postgres

# Utility Commands

db-console: ## Open psql console
	@psql "$(DATABASE_URL)"

db-dump: ## Dump database to file (usage: make db-dump FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		FILE="backup_$$(date +%Y%m%d_%H%M%S).sql"; \
	fi
	@echo "Dumping database to $$FILE..."
	@pg_dump "$(DATABASE_URL)" > $$FILE
	@echo "✓ Database dumped to $$FILE"

db-restore: ## Restore database from file (usage: make db-restore FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "Error: FILE is required. Usage: make db-restore FILE=backup.sql"; \
		exit 1; \
	fi
	@echo "Restoring database from $(FILE)..."
	@psql "$(DATABASE_URL)" < $(FILE)
	@echo "✓ Database restored from $(FILE)"

install-tools: ## Install required development tools
	@echo "Installing development tools..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/cosmtrek/air@latest
	@echo "✓ Tools installed (migrate, air)"