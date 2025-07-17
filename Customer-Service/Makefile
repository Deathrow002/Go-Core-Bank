.PHONY: help build run clean dev-setup migrate docker-build docker-run docker-stop

# Default target
help:
	@echo "Available commands:"
	@echo "  build            - Build the customer service"
	@echo "  run              - Run the customer service locally"
	@echo "  clean            - Clean build artifacts"
	@echo "  dev-setup        - Set up development environment"
	@echo "  migrate          - Run database migrations"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run with Docker Compose"
	@echo "  docker-stop      - Stop Docker containers"

# Build the application
build:
	go build -o customer-service ./cmd

# Run the application locally
run: build
	./customer-service

# Clean build artifacts
clean:
	rm -f customer-service
	go clean

# Set up development environment
dev-setup:
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file"; fi
	go mod download

# Run database migrations
migrate:
	go run ./cmd/migrate

# Build Docker image
docker-build:
	docker build -t customer-service .

# Run with Docker Compose
docker-run: docker-build
	docker-compose up

# Stop Docker containers
docker-stop:
	docker-compose down
