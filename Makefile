.PHONY: help build clean run dev test lint install-deps sqlc-generate \
        migrate-up migrate-down migrate-status migrate-reset seed \
        deploy deploy-staging \
        build-all build-linux build-darwin-amd64 build-darwin-arm64 build-windows

# Variables
BINARY_NAME=dopa
DB_PATH=dopadone.db
MIGRATIONS_DIR=migrations
GO=go
GOOSE=goose

# Version info (overridable via environment)
VERSION?=dev
GIT_COMMIT?=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-s -w \
	-X github.com/example/dopadone/internal/version.Version=$(VERSION) \
	-X github.com/example/dopadone/internal/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/example/dopadone/internal/version.BuildDate=$(BUILD_DATE)

# Default target
.DEFAULT_GOAL := help

## Help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Build
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "Build complete: bin/$(BINARY_NAME)"

build-versioned: ## Build with explicit version (VERSION=v1.0.0 make build-versioned)
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "Build complete: bin/$(BINARY_NAME)"

build-all: build-linux build-darwin-amd64 build-darwin-arm64 build-windows ## Build for all platforms

build-linux: ## Build for Linux amd64
	@echo "Building for Linux amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/$(BINARY_NAME)
	@echo "Build complete: bin/$(BINARY_NAME)-linux-amd64"

build-darwin-amd64: ## Build for macOS amd64
	@echo "Building for macOS amd64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/$(BINARY_NAME)
	@echo "Build complete: bin/$(BINARY_NAME)-darwin-amd64"

build-darwin-arm64: ## Build for macOS arm64 (M1/M2)
	@echo "Building for macOS arm64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/$(BINARY_NAME)
	@echo "Build complete: bin/$(BINARY_NAME)-darwin-arm64"

build-windows: ## Build for Windows amd64
	@echo "Building for Windows amd64..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/$(BINARY_NAME)
	@echo "Build complete: bin/$(BINARY_NAME)-windows-amd64.exe"

dist: build-all ## Create distribution archives
	@echo "Creating distribution archives..."
	cd bin && tar -czvf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd bin && tar -czvf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd bin && tar -czvf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	cd bin && zip -r $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Distribution archives created in bin/"

clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f $(BINARY_NAME)
	@echo "Clean complete"

## Development
run: build ## Build and run the application
	./bin/$(BINARY_NAME)

dev: ## Run in development mode (with hot reload if available)
	@echo "Running in development mode..."
	$(GO) run ./cmd/$(BINARY_NAME)

test: ## Run all tests
	@echo "Running tests..."
	$(GO) test ./... -v

lint: ## Run linter
	@echo "Running linter..."
	$(GO) vet ./...
	@echo "Lint complete"

install-deps: ## Install project dependencies
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies installed"

sqlc-generate: ## Generate sqlc code
	@echo "Generating sqlc code..."
	sqlc generate
	@echo "Sqlc code generated"

## Database (requires goose)
migrate-up: ## Apply all migrations
	@echo "Running migrations up..."
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) up

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) down

migrate-status: ## Check migration status
	@echo "Checking migration status..."
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) status

migrate-reset: ## Reset database (down + up)
	@echo "Resetting database..."
	$(GOOSE) -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) reset
	@echo "Database reset complete"

seed: ## Seed database with test data (optional DB_PATH=/path/to/db)
	@echo "Seeding database with test data..."
	@echo "Tip: Use ./dev.sh seed for contextual tasks"
	./scripts/seed-test-data.sh $(DB_PATH)

tui: ## Launch terminal UI
	@echo "Starting TUI..."
	$(GO) run ./cmd/$(BINARY_NAME) tui --db $(DB_PATH)

## Deployment (placeholders - customize as needed)
deploy: build ## Deploy to production
	@echo "DEPLOYMENT NOT CONFIGURED"
	@echo "Please configure deploy target in Makefile"
	@exit 1

deploy-staging: build ## Deploy to staging
	@echo "STAGING DEPLOYMENT NOT CONFIGURED"
	@echo "Please configure deploy-staging target in Makefile"
	@exit 1
