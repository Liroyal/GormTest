# =============================================================================
# PostgreSQL Docker Configuration Constants
# =============================================================================
# These constants are used throughout the Makefile and should be referenced
# in the README.md file. Update these values here to change the configuration
# across all targets.
#
# Docker image:     postgres:latest
# Container name:   pg-local  
# Port mapping:     host 5439 → container 5432
# Username:         admin
# Password:         1234
# Database name:    appdb (optional)
# =============================================================================

DOCKER_IMAGE = postgres:latest
CONTAINER_NAME = pg-local
CONTAINER = $(CONTAINER_NAME)
USER = admin
HOST_PORT = 5439
CONTAINER_PORT = 5432
DB_USERNAME = admin
DB_PASSWORD = 1234
DB_NAME = appdb

# Example targets using the constants above
.PHONY: help start stop logs connect db-up db-down db-logs psql

help: ## Show this help message
	@echo "PostgreSQL Docker Management"
	@echo "============================"
	@echo "Configuration:"
	@echo "  Image: $(DOCKER_IMAGE)"
	@echo "  Container: $(CONTAINER_NAME)"
	@echo "  Port: $(HOST_PORT):$(CONTAINER_PORT)"
	@echo "  Username: $(DB_USERNAME)"
	@echo "  Database: $(DB_NAME)"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

start: ## Start PostgreSQL container
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(HOST_PORT):$(CONTAINER_PORT) \
		-e POSTGRES_USER=$(DB_USERNAME) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		$(DOCKER_IMAGE)

stop: ## Stop and remove PostgreSQL container
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

logs: ## Show container logs
	docker logs -f $(CONTAINER_NAME)

connect: ## Connect to PostgreSQL using psql
	docker exec -it $(CONTAINER_NAME) psql -U $(DB_USERNAME) -d $(DB_NAME)

# =============================================================================
# Database Management Targets
# =============================================================================

db-up: ## Start container
	docker rm $(CONTAINER_NAME) 2>/dev/null || true
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(HOST_PORT):$(CONTAINER_PORT) \
		-e POSTGRES_USER=$(DB_USERNAME) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		$(DOCKER_IMAGE)

db-down: ## Stop + remove container
	docker stop $(CONTAINER) && docker rm $(CONTAINER)

db-logs: ## Follow container logs
	docker logs -f $(CONTAINER)

psql: ## Open interactive psql inside the container
	docker exec -it $(CONTAINER) psql -U $(USER) -d postgres
