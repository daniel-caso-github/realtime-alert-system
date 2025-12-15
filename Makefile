# ============================================================================
# VARIABLES
# ============================================================================
APP_NAME := realtime-alerting-system
BINARY_NAME := api
BINARY_PATH := bin/$(BINARY_NAME)
GO := go
GOFLAGS := -v
MAIN_PATH := ./cmd/api

# Docker
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := latest

# ============================================================================
# COLORS
# ============================================================================
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# ============================================================================
# HELP
# ============================================================================
.PHONY: help
help: ## Show this help message
	@echo '$(BLUE)Usage:$(NC)'
	@echo '  make $(GREEN)<target>$(NC)'
	@echo ''
	@echo '$(BLUE)Targets:$(NC)'
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ============================================================================
# DEVELOPMENT
# ============================================================================
.PHONY: run
run: ## Run the application
	@echo "$(BLUE)Starting application...$(NC)"
	$(GO) run $(MAIN_PATH)

.PHONY: build
build: ## Build the application binary
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p bin
	$(GO) build $(GOFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Binary built: $(BINARY_PATH)$(NC)"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf bin/
	@rm -rf tmp/
	$(GO) clean
	@echo "$(GREEN)Clean complete$(NC)"

# ============================================================================
# TESTING
# ============================================================================
.PHONY: test
test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	$(GO) test -v -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

.PHONY: test-short
test-short: ## Run short tests only
	$(GO) test -v -short ./...

## test-integration: Run integration tests
test-integration:
	@echo "$(BLUE)Running integration tests...$(NC)"
	$(GO) test -v -cover ./test/integration/...

## test-integration-pretty: Run integration tests with pretty output
test-pretty:
	@echo "$(BLUE)Running integration tests...$(NC)"
	gotestsum --format testname ./test/...

## coverage: Generate coverage report
coverage:
	@echo "$(BLUE)Generating coverage report...$(NC)"
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	$(GO) tool cover -func=coverage.out
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

# ============================================================================
# CODE QUALITY
# ============================================================================
.PHONY: lint
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run linter and fix issues
	@echo "$(BLUE)Running linter with fix...$(NC)"
	golangci-lint run --fix ./...

.PHONY: fmt
fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GO) vet ./...

.PHONY: check
check: fmt vet lint ## Run all code quality checks

# ============================================================================
# DEPENDENCIES
# ============================================================================
.PHONY: deps
deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GO) mod download

.PHONY: tidy
tidy: ## Tidy dependencies
	@echo "$(BLUE)Tidying dependencies...$(NC)"
	$(GO) mod tidy

.PHONY: vendor
vendor: ## Vendor dependencies
	@echo "$(BLUE)Vendoring dependencies...$(NC)"
	$(GO) mod vendor

# ============================================================================
# DATABASE MIGRATIONS
# ============================================================================
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/alerting_db?sslmode=disable

.PHONY: migrate-up
migrate-up: ## Run all pending migrations
	@echo "$(BLUE)Running migrations up...$(NC)"
	migrate -path migrations -database "$(DATABASE_URL)" up
	@echo "$(GREEN)Migrations completed$(NC)"

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	migrate -path migrations -database "$(DATABASE_URL)" down 1

.PHONY: migrate-down-all
migrate-down-all: ## Rollback all migrations
	@echo "$(RED)Rolling back ALL migrations...$(NC)"
	migrate -path migrations -database "$(DATABASE_URL)" down -all

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@migrate -path migrations -database "$(DATABASE_URL)" version

.PHONY: migrate-force
migrate-force: ## Force set migration version (usage: make migrate-force VERSION=1)
	@echo "$(YELLOW)Forcing migration version to $(VERSION)...$(NC)"
	migrate -path migrations -database "$(DATABASE_URL)" force $(VERSION)

.PHONY: migrate-create
migrate-create: ## Create new migration (usage: make migrate-create NAME=add_users_table)
	@echo "$(BLUE)Creating migration: $(NAME)...$(NC)"
	migrate create -ext sql -dir migrations -seq $(NAME)
	@echo "$(GREEN)Migration files created$(NC)"

.PHONY: migrate-status
migrate-status: ## Show migration status
	@echo "$(BLUE)Migration files:$(NC)"
	@ls -la migrations/*.sql 2>/dev/null || echo "No migrations found"

# ============================================================================
# DOCKER
# ============================================================================
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "$(BLUE)Running Docker container...$(NC)"
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-up
docker-up: ## Start all services with docker-compose
	@echo "$(BLUE)Starting services...$(NC)"
	docker-compose up -d

.PHONY: docker-down
docker-down: ## Stop all services
	@echo "$(YELLOW)Stopping services...$(NC)"
	docker-compose down

.PHONY: docker-logs
docker-logs: ## View docker-compose logs
	docker-compose logs -f

# ============================================================================
# UTILITIES
# ============================================================================
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(BLUE)Installing tools...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)Tools installed$(NC)"

.PHONY: env
env: ## Create .env file from example
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(GREEN).env file created$(NC)"; \
	else \
		echo "$(YELLOW).env file already exists$(NC)"; \
	fi

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "$(BLUE)Generating Swagger docs...$(NC)"
	swag init -g cmd/api/main.go -o docs

# ============================================================================
# DEFAULT
# ============================================================================
.DEFAULT_GOAL := help


# ============================================================================
# DEVELOPMENT HELPERS
# ============================================================================
.PHONY: db-connect
db-connect: ## Connect to PostgreSQL database
	@./scripts/db-connect.sh

.PHONY: redis-connect
redis-connect: ## Connect to Redis
	@./scripts/redis-connect.sh

.PHONY: db-reset
db-reset: ## Reset database (WARNING: deletes all data)
	@./scripts/reset-db.sh

.PHONY: health
health: ## Check health of all services
	@./scripts/health-check.sh

.PHONY: logs
logs: ## View logs (usage: make logs or make logs SERVICE=api)
	@./scripts/logs.sh $(SERVICE)

.PHONY: dev
dev: docker-up logs ## Start development environment and show logs
