.PHONY: help build up down logs clean restart

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build all services"
	@echo "  up        - Start all services"
	@echo "  down      - Stop all services"
	@echo "  logs      - View logs"
	@echo "  clean     - Clean up containers and volumes"
	@echo "  restart   - Restart all services"

# Build all services
build:
	docker-compose build --no-cache

# Start all services
up:
	docker-compose up -d

# Start with logs
up-logs:
	docker-compose up

# Stop all services
down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# View specific service logs
logs-customer:
	docker-compose logs -f customer-service

logs-account:
	docker-compose logs -f account-service

# Clean up
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

# Restart services
restart:
	docker-compose restart

# Health check all services
health:
	@echo "Checking service health..."
	@curl -f http://localhost/health || echo "API Gateway: DOWN"
	@curl -f http://localhost:8080/health || echo "Customer Service: DOWN"
	@curl -f http://localhost:8081/health || echo "Account Service: DOWN"
