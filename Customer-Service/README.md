# Customer Service - Core Banking Microservice

A standalone microservice for managing customer data in a core banking system.

## Architecture Overview

This service follows a clean architecture pattern with the following structure:

```
Customer-Service/
├── cmd/                   # Application entry points
│   ├── main.go           # Service entry point
│   └── migrate/          # Database migration utility
│       └── main.go
├── internal/             # Private application code
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── customer/         # Customer domain
│   │   ├── controllers/  # HTTP controllers (was handlers)
│   │   │   └── customer_controller.go
│   │   ├── models/       # Domain models
│   │   │   └── customer.go
│   │   ├── repository/   # Data access layer
│   │   │   └── customer_repository.go
│   │   └── service/      # Business logic layer
│   │       └── customer_service.go
│   └── database/         # Database utilities
│       └── database.go
├── pkg/                  # Public packages
│   └── middleware/       # HTTP middlewares
│       └── middleware.go
├── .env.example         # Environment template
├── .gitignore          # Git ignore rules
├── docker-compose.yml  # Docker Compose config
├── Dockerfile         # Docker image config
├── go.mod            # Go dependencies
├── go.sum           # Dependency checksums
├── Makefile        # Build automation
└── README.md      # This documentation
```

## Technology Stack

- **Language**: Go 1.23+
- **Web Framework**: Gin (HTTP router and middleware)
- **Database**: PostgreSQL 15 with GORM ORM
- **Containerization**: Docker & Docker Compose
- **Configuration**: Environment-based with godotenv
- **Architecture**: Clean Architecture with MVC pattern

## Features

### Customer Management (CRUD)
- ✅ **Create** customer profiles with validation
- ✅ **Read** customer information by ID
- ✅ **Update** customer details with conflict resolution
- ✅ **Delete** customer records (soft delete)
- ✅ **List** customers with pagination
- ✅ **Search** customers by multiple criteria

### Technical Features
- ✅ RESTful API design with proper HTTP status codes
- ✅ Database migrations and connection management
- ✅ Comprehensive error handling and validation
- ✅ Middleware for CORS, logging, and recovery
- ✅ Docker containerization for easy deployment
- ✅ Environment-based configuration
- ✅ Health monitoring endpoint

## Quick Start with Docker

### Prerequisites
- Docker and Docker Compose installed

### Run the Service
```bash
# Clone and navigate to the service
git clone <repository-url>
cd Customer-Service

# Copy environment configuration
cp .env.example .env
# Edit .env with your preferred settings

# Start all services (customer-service, PostgreSQL, pgAdmin)
docker-compose up --build

# Service will be available at:
# - Customer API: http://localhost:8080
# - pgAdmin: http://localhost:5050 (admin@admin.com / admin)
```

### Verify Service Health
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "customer-service",
  "version": "1.0.0"
}
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/health` | Health check |
| POST   | `/api/v1/customers` | Create new customer |
| GET    | `/api/v1/customers/{id}` | Get customer by ID |
| PUT    | `/api/v1/customers/{id}` | Update customer |
| DELETE | `/api/v1/customers/{id}` | Delete customer |
| GET    | `/api/v1/customers` | List customers (paginated) |
| GET    | `/api/v1/customers/search` | Search customers |

### Example API Usage

#### Create a Customer
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "+1-555-0123",
    "date_of_birth": "1990-01-15",
    "address": "123 Main St, Anytown, USA"
  }'
```

#### Get Customer by ID
```bash
curl http://localhost:8080/api/v1/customers/{customer-id}
```

#### List Customers with Pagination
```bash
curl "http://localhost:8080/api/v1/customers?page=1&page_size=10"
```

#### Search Customers
```bash
curl "http://localhost:8080/api/v1/customers/search?query=john&status=active"
```

## Local Development

### Prerequisites
- Go 1.23+
- PostgreSQL 15+

### Setup
```bash
# Install dependencies
go mod download

# Set up environment
cp .env.example .env
# Edit .env with your database configuration

# Run database migrations
make migrate

# Build and run the service
make run
```

### Available Make Commands
```bash
make help          # Show available commands
make build         # Build the customer service
make run           # Build and run locally
make clean         # Clean build artifacts
make dev-setup     # Set up development environment
make migrate       # Run database migrations
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose
make docker-stop   # Stop Docker containers
```

## Docker Services

When running with `docker-compose up`, the following services are started:

| Service | Description | Port | Access |
|---------|-------------|------|--------|
| **customer-service** | Main API service | 8080 | http://localhost:8080 |
| **postgres** | PostgreSQL database | 5432 | Internal |
| **pgadmin** | Database admin UI | 5050 | http://localhost:5050 |

### Database Access
- **Database**: `core_bank`
- **Username**: `postgres`
- **Password**: `postgres`
- **pgAdmin Login**: `admin@admin.com` / `admin`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (development/production) | `development` |
| `SERVER_HOST` | Server bind address | `0.0.0.0` |
| `SERVER_PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `core_bank` |
| `DB_SSL_MODE` | SSL mode for database | `disable` |

## Development

### Project Philosophy
This Customer Service follows **Clean Architecture** principles:

1. **Controllers Layer**: HTTP request handling and routing
2. **Service Layer**: Business logic and validation
3. **Repository Layer**: Data access and persistence
4. **Models Layer**: Domain entities and data structures

### Database Features
- **GORM ORM**: Type-safe database operations
- **Auto-migrations**: Automatic schema management
- **Soft deletes**: Preserve data integrity
- **Connection pooling**: Optimized database performance
- **Health checks**: Database connectivity monitoring

### Future Enhancements
- Authentication and authorization
- API rate limiting
- Metrics and monitoring
- Swagger documentation
- Event sourcing integration

## Core Banking Architecture

This Customer Service is designed as part of a larger microservices architecture:

```
Core Banking Platform
├── Customer-Service (this service)
├── Account-Service (future)
├── Transaction-Service (future)
├── Loan-Service (future)
├── Card-Service (future)
└── Notification-Service (future)
```

Each service is:
- **Self-contained**: Own database and business logic
- **API-driven**: RESTful interfaces
- **Docker-ready**: Containerized deployment
- **Independently scalable**: Can be scaled based on demand

## Contributing

1. Follow Go coding standards
2. Ensure all endpoints have proper error handling
3. Update documentation for API changes
4. Use conventional commit messages
5. Test with Docker Compose before submitting

## License

This project is licensed under the MIT License.
