# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the mikrocloud binary (backend only)"
	@echo "  build-web    - Build the frontend assets"
	@echo "  build-full   - Build frontend assets and backend with embedded frontend"
	@echo "  run          - Run the mikrocloud server"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Download dependencies"
	@echo "  migrate      - Run database migrations"
	@echo "  migrate-up   - Apply all pending migrations"
	@echo "  migrate-down - Rollback the last migration"
	@echo "  migrate-status - Show migration status"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"

# Build the binary (backend only)
build: deps
	go build -o bin/mikrocloud ./main.go

# Build the frontend assets
build-web:
	@echo "Building frontend assets..."
	cd web && pnpm install
	cd web && pnpm run build
	@echo "✅ Frontend built successfully at web/dist/"

# Build everything: frontend + backend with embedded assets
build-full: build-web deps
	@echo "Building backend with embedded frontend..."
	go build -o bin/mikrocloud ./main.go
	@echo "✅ Full build complete! Frontend is embedded in the binary."

# Run the server (builds full if binary doesn't exist)
run: bin/mikrocloud
	./bin/mikrocloud serve

# Ensure the binary exists
bin/mikrocloud:
	$(MAKE) build-full

# Run in development mode with auto-reload (requires air)
dev:
	air

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf web/dist/
	rm -rf web/.svelte-kit/
	rm -rf web/node_modules/
	go clean

# Run tests
test:
	go test -v ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Download frontend dependencies
deps-web:
	@echo "Installing frontend dependencies..."
	cd web && npm install

# Download all dependencies
deps-all: deps deps-web

# Database migrations using goose
migrate: migrate-up

migrate-up:
	@mkdir -p $(HOME)/.local/share/mikrocloud
	goose -dir migrations sqlite3 "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations sqlite3 "$(DATABASE_URL)" down

migrate-status:
	goose -dir migrations sqlite3 "$(DATABASE_URL)" status

# Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

# Docker targets
docker-build:
	docker build -t mikrocloud:latest .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

# Development database - SQLite (no external dependencies needed)
db-init:
	@mkdir -p $(HOME)/.local/share/mikrocloud
	@echo "SQLite database will be created automatically"

db-clean:
	rm -f $(HOME)/.local/share/mikrocloud/mikrocloud.db*

# Install tools
install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/cosmtrek/air@latest
	npm install -g pnpm

# Development mode - frontend dev server + backend
dev-full:
	@echo "Starting development mode..."
	@echo "Frontend will be available at http://localhost:5173"
	@echo "Backend will be available at http://localhost:3000"
	cd web && pnpm run dev &
	$(MAKE) dev

# Development mode with frontend built and embedded
dev-embedded: build-full
	air

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Run security checks (requires gosec)
security:
	gosec ./...

# Generate code coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Default environment variables
export DATABASE_URL ?= $(HOME)/.local/share/mikrocloud/mikrocloud.db
export PORT ?= 3000
export LOG_LEVEL ?= info
