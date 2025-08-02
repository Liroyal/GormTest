# =============================================================================
# Configuration Constants
# =============================================================================
# These constants are used throughout the Makefile for consistency.
# Update these values here to change the configuration across all targets.
# =============================================================================

# Database configuration
DOCKER_IMAGE = postgres:latest
CONTAINER_NAME = pg-local
DB_HOST = localhost
DB_PORT = 5439
DB_USER = admin
DB_PASSWORD = 1234
DB_NAME = appdb
CONTAINER_PORT = 5432

# API configuration
SERVER_HOST = localhost
SERVER_PORT = 8080
API_BINARY = bin/api
API_PID_FILE = .api.pid
API_LOG_FILE = logs/api.log

# Environment variables for API
ENV_VARS = SERVER_HOST=$(SERVER_HOST) SERVER_PORT=$(SERVER_PORT) DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME)

# Migration configuration
MIGRATE_PATH = migrations
MIGRATE_DB_URL = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATE = migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_DB_URL)"

# =============================================================================
# Target Declarations
# =============================================================================

.PHONY: help api-build api-down api-up db-down db-logs db-rebuild db-up migrate-create migrate-down migrate-up

# =============================================================================
# Meta Targets
# =============================================================================

help: ## Show this help message with all available targets
	@echo "Development Environment Management"
	@echo "=================================="
	@echo "Configuration:"
	@echo "  Database: $(DB_USER)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)"
	@echo "  API Server: $(SERVER_HOST):$(SERVER_PORT)"
	@echo "  Container: $(CONTAINER_NAME)"
	@echo "  API Binary: $(API_BINARY)"
	@echo "  API Logs: $(API_LOG_FILE)"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# =============================================================================
# Database Targets
# =============================================================================

db-up: ## Start PostgreSQL container
	docker rm $(CONTAINER_NAME) 2>/dev/null || true
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(DB_PORT):$(CONTAINER_PORT) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		$(DOCKER_IMAGE)

db-down: ## Stop and remove PostgreSQL container
	@docker stop $(CONTAINER_NAME) 2>/dev/null || echo "Container $(CONTAINER_NAME) not running"
	@docker rm $(CONTAINER_NAME) 2>/dev/null || echo "Container $(CONTAINER_NAME) not found"

db-rebuild: ## Stop container, remove it, start fresh container, run migrate-up
	@echo "Rebuilding database container and running migrations..."
	$(MAKE) db-down
	$(MAKE) db-up
	@echo "Waiting for database to be ready..."
	@sleep 3
	$(MAKE) migrate-up

db-logs: ## Follow PostgreSQL container logs
	docker logs -f $(CONTAINER_NAME)

# =============================================================================
# Migration Targets
# =============================================================================

migrate-up: ## Apply all pending database migrations
	$(MIGRATE) up

migrate-down: ## Roll back the last database migration
	$(MIGRATE) down 1

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME parameter is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(NAME)

# =============================================================================
# API Targets
# =============================================================================

api-build: ## Build the API binary
	@echo "Building API binary..."
	mkdir -p bin
	go build -o $(API_BINARY) main.go
	@echo "API binary built at $(API_BINARY)"

api-up: ## Start the API server in the background
	@echo "Starting API server..."
	@if [ -f $(API_PID_FILE) ]; then \
		echo "API server may already be running (PID file exists)"; \
		echo "Use 'make api-down' to stop it first"; \
		exit 1; \
	fi
	@if [ ! -f $(API_BINARY) ]; then \
		echo "API binary not found, building first..."; \
		$(MAKE) api-build; \
	fi
	@mkdir -p logs
	@nohup env $(ENV_VARS) $(API_BINARY) >> $(API_LOG_FILE) 2>&1 & echo $$! > $(API_PID_FILE)
	@echo "API server started at http://$(SERVER_HOST):$(SERVER_PORT)"
	@echo "PID: $$(cat $(API_PID_FILE))"
	@echo "Logs: $(API_LOG_FILE)"

api-down: ## Stop the API server running in the background
	@if [ ! -f $(API_PID_FILE) ]; then \
		echo "No API server PID file found, nothing to stop"; \
		exit 0; \
	fi
	@PID=$$(cat $(API_PID_FILE)); \
	if kill -0 $$PID 2>/dev/null; then \
		echo "Stopping API server (PID: $$PID)"; \
		kill -TERM $$PID 2>/dev/null || true; \
		sleep 1; \
		if kill -0 $$PID 2>/dev/null; then \
			echo "Graceful shutdown failed, forcing kill"; \
			kill -KILL $$PID 2>/dev/null || true; \
		fi; \
		rm -f $(API_PID_FILE); \
		echo "API server stopped"; \
	else \
		echo "API server process not found, cleaning up PID file"; \
		rm -f $(API_PID_FILE); \
	fi
