# Skopidom project Makefile.
# Run `make help` to see available targets.
# Place this file in the project root (next to docker-compose.yml).

# ── Variables ─────────────────────────────────────────────────────────────────

BACKEND_DIR := ./backend
BINARY      := inventory
BUILD_DIR   := $(BACKEND_DIR)/bin
CMD         := ./cmd/server

DOCKER_COMPOSE := sudo docker compose
DB_CONTAINER   := skopidom-postgres-1
DB_USER        := skopidom
DB_NAME        := skopidom

# ── Default target ────────────────────────────────────────────────────────────

.DEFAULT_GOAL := help

# ── Development ───────────────────────────────────────────────────────────────

.PHONY: run
run: ## Run the backend server locally (requires .env and running postgres)
	cd $(BACKEND_DIR) && go run $(CMD)

.PHONY: build
build: ## Compile the binary to backend/bin/inventory
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o bin/$(BINARY) $(CMD)

.PHONY: generate
generate: ## Regenerate sqlc code from queries/*.sql and migrations
	cd $(BACKEND_DIR) && sqlc generate

.PHONY: tidy
tidy: ## Tidy go modules
	cd $(BACKEND_DIR) && go mod tidy

# ── Quality ───────────────────────────────────────────────────────────────────

.PHONY: lint
lint: ## Run golangci-lint (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	cd $(BACKEND_DIR) && golangci-lint run ./...

.PHONY: vet
vet: ## Run go vet
	cd $(BACKEND_DIR) && go vet ./...

.PHONY: check
check: vet ## Run all static checks (vet + build)
	cd $(BACKEND_DIR) && go build ./...

# ── Docker ────────────────────────────────────────────────────────────────────

.PHONY: up
up: ## Start all services (postgres + backend)
	$(DOCKER_COMPOSE) up -d

.PHONY: down
down: ## Stop all services
	$(DOCKER_COMPOSE) down

.PHONY: reset
reset: ## Stop all services and wipe volumes (full reset)
	$(DOCKER_COMPOSE) down -v

.PHONY: postgres-up
postgres-up: ## Start only postgres
	$(DOCKER_COMPOSE) up postgres -d

.PHONY: postgres-down
postgres-down: ## Stop only postgres
	$(DOCKER_COMPOSE) stop postgres

.PHONY: logs
logs: ## Follow logs of all services
	$(DOCKER_COMPOSE) logs -f

.PHONY: logs-backend
logs-backend: ## Follow backend logs only
	$(DOCKER_COMPOSE) logs -f backend

# ── Database ──────────────────────────────────────────────────────────────────

.PHONY: psql
psql: ## Open psql shell in the running postgres container
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: create-admin
create-admin: ## Insert a default admin user (password: password)
	docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) <<'SQL'
	INSERT INTO users (full_name, email, password_hash, role)
	VALUES (
	  'Администратор',
	  'admin@university.ru',
	  '$$2a$$12$$iB.8j9wRbmfza6qfuGTBn.l2dCkBc9ojVcIWnZi80nDM4bwO0RhEy',
	  'admin'
	) ON CONFLICT (email) DO NOTHING;
	SQL

# ── Help ──────────────────────────────────────────────────────────────────────

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
