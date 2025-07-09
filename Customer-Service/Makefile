.PHONY: help build run clean docker-build docker-run docker-stop migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the customer service binary"
	@echo "  run           - Run the customer service"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker containers"
	@echo "  migrate       - Run database migrations"
	@echo "  mod-tidy      - Run go mod tidy"

# Build the application
build:
	go build -o bin/customer-service ./cmd/main.go
	go build -o bin/migrate ./cmd/migrate/main.go

# Run the application
run:
	go run ./cmd/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	docker build -t customer-service .

# Run with Docker Compose
docker-run:
	docker-compose up --build

# Stop Docker containers
docker-stop:
	docker-compose down

# Run database migrations
migrate:
	go run ./cmd/migrate/main.go

# Run go mod tidy
mod-tidy:
	go mod tidy

# Install dependencies
deps:
	go mod download

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Development setup
dev-setup: mod-tidy
	@echo "Development environment setup complete"
	@echo "To start the application:"
	@echo "  1. Set up your .env file: cp .env.example .env"
	@echo "  2. Start PostgreSQL (or use docker-compose)"
	@echo "  3. Run migrations: make migrate"
	@echo "  4. Start the service: make run"
