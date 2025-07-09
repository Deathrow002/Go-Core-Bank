# Customer Service - Core Banking Microservice

A dedicated microservice for managing bank customers with full CRUD operations.

## Features

- **Create** customer profiles with validation
- **Read** customer information by ID
- **Update** customer details with conflict resolution
- **Delete** customer records (soft delete)
- **List** customers with pagination
- **Search** customers by multiple criteria

## Technology Stack

- **Language**: Go 1.23
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Environment**: Docker support

## Project Structure

```
customer-service/
├── cmd/                    # Application entry points
│   ├── main.go            # Customer service main
│   └── migrate/           # Database migration utility
│       └── main.go        # Migration runner
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── customer/         # Customer domain module
│   │   ├── handlers/     # HTTP request handlers
│   │   ├── models/       # Domain models and DTOs
│   │   ├── repository/   # Data access layer
│   │   └── service/      # Business logic layer
│   └── database/         # Database utilities
├── pkg/                  # Public packages
│   └── middleware/       # HTTP middlewares
├── .env                  # Environment variables
├── .env.example          # Environment template
├── docker-compose.yml    # Docker compose setup
├── Dockerfile           # Container image
├── go.mod               # Go module
├── go.sum               # Module checksums
├── Makefile             # Development tasks
└── README.md            # This file
```

## Getting Started

### Prerequisites
- Go 1.23 or higher
- PostgreSQL 13 or higher
- Docker (optional)

### Installation

1. Navigate to the customer service directory:
```bash
cd Customer-Service
```

2. Install dependencies:
```bash
make deps
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database configuration
```

4. Run database migrations:
```bash
make migrate
```

5. Start the service:
```bash
make run
```

### Using Docker

1. Start all services:
```bash
make docker-run
```

This will start:
- PostgreSQL database on port 5432
- Customer service on port 8080
- PgAdmin on port 5050

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/api/v1/customers` | Create a new customer |
| GET    | `/api/v1/customers/:id` | Get customer by ID |
| PUT    | `/api/v1/customers/:id` | Update customer |
| DELETE | `/api/v1/customers/:id` | Delete customer |
| GET    | `/api/v1/customers` | List customers with pagination |
| GET    | `/api/v1/customers/search` | Search customers |
| GET    | `/health` | Health check |

### Example Usage

```bash
# Create a customer
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "1234567890",
    "address": {
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "USA"
    }
  }'

# Get a customer
curl http://localhost:8080/api/v1/customers/{customer-id}

# List customers
curl "http://localhost:8080/api/v1/customers?page=1&page_size=10"

# Search customers
curl "http://localhost:8080/api/v1/customers/search?query=John&page=1&page_size=10"
```

## Development

### Available Make Commands

- `make build` - Build the application
- `make run` - Run the service
- `make docker-run` - Start with Docker
- `make docker-stop` - Stop Docker containers
- `make migrate` - Run database migrations
- `make clean` - Clean build artifacts
- `make fmt` - Format code
- `make mod-tidy` - Tidy Go modules

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | |
| DB_NAME | Database name | core_bank |
| DB_SSL_MODE | SSL mode | disable |
| SERVER_PORT | Server port | 8080 |
| APP_ENV | Environment | development |

## Architecture

The service follows Clean Architecture principles:

- **Handlers**: HTTP request/response handling
- **Service**: Business logic and validation
- **Repository**: Data access layer
- **Models**: Domain entities and DTOs

## Contributing

1. Make your changes
2. Format code: `make fmt`
3. Build: `make build`
4. Test locally: `make run`

## License

MIT License
