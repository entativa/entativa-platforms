.PHONY: help build up down restart logs ps clean test migrate proto

# Default target
help:
	@echo "Socialink Backend - Available Commands:"
	@echo ""
	@echo "  make up           - Start all services"
	@echo "  make down         - Stop all services"
	@echo "  make restart      - Restart all services"
	@echo "  make build        - Build all Docker images"
	@echo "  make logs         - View logs from all services"
	@echo "  make ps           - List running services"
	@echo "  make clean        - Remove all containers, volumes, and images"
	@echo "  make test         - Run all tests"
	@echo "  make migrate      - Run database migrations"
	@echo "  make proto        - Generate protobuf code"
	@echo "  make lint         - Run linters on all services"
	@echo ""

# Docker Compose commands
up:
	@echo "🚀 Starting Socialink services..."
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "📊 API Gateway: http://localhost:8081"
	@echo "📊 Grafana: http://localhost:3000 (admin/admin)"
	@echo "📊 Prometheus: http://localhost:9090"

down:
	@echo "🛑 Stopping Socialink services..."
	docker-compose down
	@echo "✅ Services stopped!"

restart:
	@echo "🔄 Restarting Socialink services..."
	docker-compose restart
	@echo "✅ Services restarted!"

build:
	@echo "🏗️  Building Docker images..."
	docker-compose build --parallel
	@echo "✅ Build complete!"

logs:
	docker-compose logs -f

ps:
	docker-compose ps

clean:
	@echo "🧹 Cleaning up..."
	docker-compose down -v --remove-orphans
	docker system prune -f
	@echo "✅ Cleanup complete!"

# Development commands
test:
	@echo "🧪 Running tests..."
	@cd services/user-service && go test ./...
	@cd services/posting-service && go test ./...
	@cd services/community-service && go test ./...
	@cd services/settings-service && go test ./...
	@cd services/search-service && go test ./...
	@cd services/live-streaming-service && go test ./...
	@cd services/creator-service && go test ./...
	@echo "✅ All tests passed!"

migrate:
	@echo "📦 Running migrations..."
	@./scripts/run-migrations.sh
	@echo "✅ Migrations complete!"

proto:
	@echo "📦 Generating protobuf code..."
	@cd proto && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		*.proto
	@echo "✅ Protobuf code generated!"

lint:
	@echo "🔍 Running linters..."
	@cd services/user-service && golangci-lint run
	@cd services/posting-service && golangci-lint run
	@cd services/community-service && golangci-lint run
	@echo "✅ Linting complete!"

# Database commands
db-shell:
	docker exec -it socialink-postgres psql -U postgres -d socialink

redis-cli:
	docker exec -it socialink-redis redis-cli

mongo-shell:
	docker exec -it socialink-mongodb mongosh -u mongo -p mongo

# Individual service commands
up-gateway:
	docker-compose up -d api-gateway

up-user:
	docker-compose up -d user-service

up-post:
	docker-compose up -d post-service

up-messaging:
	docker-compose up -d messaging-service

# Monitoring
grafana:
	@echo "📊 Opening Grafana..."
	@open http://localhost:3000 || xdg-open http://localhost:3000

prometheus:
	@echo "📊 Opening Prometheus..."
	@open http://localhost:9090 || xdg-open http://localhost:9090

# Health checks
health:
	@echo "🏥 Checking service health..."
	@curl -s http://localhost:8081/health | jq .
