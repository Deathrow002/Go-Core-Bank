# Core Banking Microservice - Customer Module

A modular microservice architecture for core banking operations, starting with Customer CRUD functionality.

## Architecture Overview

This project follows a clean architecture pattern with the following structure:

```
go-core-bank/
├── Customer-Service/       # Customer microservice
│   ├── cmd/               # Application entry points
│   │   └── main.go        # Service entry point
│   ├── internal/          # Private application code
│   │   ├── customer/      # Customer domain
│   │   │   ├── handlers/  # HTTP handlers
│   │   │   ├── repository/ # Data access layer
│   │   │   ├── service/   # Business logic layer
│   │   │   └── models/    # Domain models
│   │   ├── config/        # Configuration
│   │   └── database/      # Database utilities
│   ├── pkg/               # Public packages
│   │   ├── middleware/    # HTTP middlewares
│   │   └── utils/         # Utility functions
│   ├── Dockerfile         # Docker configuration
│   ├── docker-compose.yml # Local development
│   └── go.mod             # Go dependencies
├── docker-compose.yml     # Production deployment
└── Makefile              # Build automation
```

## Features

### Customer Service
- Create customer profiles
- Read customer information
- Update customer details
- Delete customer records
- List customers with pagination
- Search customers by various criteria

## Technology Stack

- **Language**: Go 1.23
- **Web Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Containerization**: Docker & Docker Compose
- **Health Monitoring**: Built-in health checks

## Getting Started

### Prerequisites
- Go 1.23 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose

### Quick Start with Docker

1. Clone the repository:
```bash
git clone <repository-url>
cd Go-Core-Bank
```

2. Start the services with Docker Compose:
```bash
docker-compose up -d
```

3. Verify the service is running:
```bash
curl http://localhost:8080/health
```

### Local Development

1. Navigate to Customer-Service directory:
```bash
cd Customer-Service
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Start PostgreSQL (if not using Docker):
```bash
# Using Docker for just the database
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15-alpine
```

5. Run the service:
```bash
go run cmd/main.go
```

## API Endpoints

### Health Check
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | /health  | Service health status |

### Customer Service (Port: 8080)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | /api/v1/customers | Create a new customer |
| GET    | /api/v1/customers/:id | Get customer by ID |
| PUT    | /api/v1/customers/:id | Update customer |
| DELETE | /api/v1/customers/:id | Delete customer |
| GET    | /api/v1/customers | List customers with pagination |
| GET    | /api/v1/customers/search | Search customers |

### Example Usage

Create a customer:
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "date_of_birth": "1990-01-15T00:00:00Z"
  }'
```

List customers:
```bash
curl http://localhost:8080/api/v1/customers
```

## Docker Services

The application runs with the following services:

- **customer-service**: Main API service (Port 8080)
- **postgres**: PostgreSQL database (Port 5432)
- **pgadmin**: Database management UI (Port 5050)

Access pgAdmin at http://localhost:5050 with:
- Email: admin@admin.com
- Password: admin

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | postgres (Docker) / localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | core_bank |
| DB_SSL_MODE | SSL mode | disable |
| SERVER_HOST | Server bind address | 0.0.0.0 |
| SERVER_PORT | Server port | 8080 |
| APP_ENV | Application environment | development |

## Development

### Building the Service

```bash
# Build locally
cd Customer-Service
go build -o customer-service cmd/main.go

# Build with Docker
docker build -t customer-service .
```

### Project Structure

This is a standalone Customer microservice that can be:
- Deployed independently
- Scaled horizontally
- Integrated with other microservices
- Monitored via health endpoints

### Database

The service uses PostgreSQL with GORM for:
- Automatic migrations
- Connection pooling
- Query optimization
- Health monitoring

## Future Modules

The architecture is designed to support additional microservices:
- Account Service
- Transaction Service  
- Loan Service
- Card Service
- Notification Service
- Auth Service

Each service will be self-contained with its own:
- Database schema
- Docker configuration
- API endpoints
- Health monitoring

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.