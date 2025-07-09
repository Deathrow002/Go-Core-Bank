.PHONY: help build run clean docker-build docker-run docker-stop customer-service

# Default target
help:
	@echo "Available commands:"
	@echo "  customer-service - Build and manage Customer Service"
	@echo "  build            - Build all services"
	@echo "  run              - Run all services with Docker Compose"
	@echo "  clean            - Clean build artifacts"
	@echo "  docker-build     - Build Docker images"
	@echo "  docker-run       - Run with Docker Compose"
	@echo "  docker-stop      - Stop Docker containers"

# Customer Service commands
customer-service:
	@echo "Building Customer Service..."
	cd Customer-Service && go build -o customer-service cmd/main.go

# Build all services
build: customer-service

# Run all services
run:
	docker-compose up --build

# Clean build artifacts
clean:
	cd Customer-Service && rm -f customer-service
	docker-compose down --volumes --remove-orphans

# Build Docker images
docker-build:
	docker-compose build

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker containers
docker-stop:
	docker-compose down

# Development utilities
fmt:
	cd Customer-Service && go fmt ./...

# Development setup
dev-setup:
	@echo "Development environment setup complete"
	@echo "To start the application:"
	@echo "  1. Run: make run"
	@echo "  2. Access API at: http://localhost:8080"
	@echo "  3. Access pgAdmin at: http://localhost:5050"
	@echo "  4. Check health: curl http://localhost:8080/health"
