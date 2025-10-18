.PHONY: help build run test clean migrate-up migrate-down docker-build docker-run

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run in Docker"

build:
	@echo "Building Vignette user service..."
	@go build -o bin/user-service cmd/api/main.go cmd/api/routes.go
	@echo "Build complete!"

run:
	@echo "Starting Vignette user service..."
	@go run cmd/api/main.go cmd/api/routes.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete!"

migrate-up:
	@echo "Running database migrations..."
	@psql $(DB_DSN) -f migrations/001_users_table.sql
	@psql $(DB_DSN) -f migrations/002_sessions_table.sql
	@psql $(DB_DSN) -f migrations/003_password_reset_tokens.sql
	@psql $(DB_DSN) -f migrations/004_cross_platform_links.sql
	@echo "Migrations complete!"

migrate-down:
	@echo "Rolling back migrations..."
	@psql $(DB_DSN) -c "DROP TABLE IF EXISTS cross_platform_links CASCADE;"
	@psql $(DB_DSN) -c "DROP TABLE IF EXISTS password_reset_tokens CASCADE;"
	@psql $(DB_DSN) -c "DROP TABLE IF EXISTS sessions CASCADE;"
	@psql $(DB_DSN) -c "DROP TABLE IF EXISTS users CASCADE;"
	@echo "Rollback complete!"

docker-build:
	@echo "Building Docker image..."
	@docker build -t vignette-user-service:latest .
	@echo "Docker build complete!"

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8002:8002 --env-file .env vignette-user-service:latest

dev:
	@echo "Starting development server with hot reload..."
	@air -c .air.toml || go run cmd/api/main.go cmd/api/routes.go

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed!"

format:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete!"

lint:
	@echo "Running linter..."
	@golangci-lint run || go vet ./...
	@echo "Lint complete!"
