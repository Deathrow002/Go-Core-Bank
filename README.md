# Core Banking Microservices Platform

A modular microservice architecture for core banking operations.

## Project Structure

```
Go-Core-Bank/
├── Customer-Service/       # Customer microservice (standalone)
│   ├── cmd/               # Application entry points
│   ├── internal/          # Private application code
│   │   ├── customer/      # Customer domain
│   │   │   ├── controllers/ # HTTP controllers (MVC pattern)
│   │   │   ├── service/   # Business logic
│   │   │   ├── repository/ # Data access
│   │   │   └── models/    # Domain models
│   │   ├── config/        # Configuration
│   │   └── database/      # Database utilities
│   └── pkg/              # Public packages
├── docker-compose.yml    # Multi-service deployment
├── Makefile             # Build automation
└── README.md           # This file
```

## Quick Start

### Using Docker (Recommended)
```bash
# Run all services
make docker-run

# Stop all services
make docker-stop
```

### Building Individual Services
```bash
# Build Customer Service
make customer-service

# Or build all services
make build
```

## Services

### Customer Service
- **Location**: `./Customer-Service/`
- **Port**: 8080
- **Documentation**: See `./Customer-Service/README.md`
- **Architecture**: Clean Architecture with MVC pattern

## Architecture

Each microservice is completely standalone with its own:
- Go module (`go.mod`)
- Database schema
- Docker configuration
- Documentation
- Build system
- MVC architecture (Controllers, Services, Repositories)

## Development

To work on a specific service, navigate to its directory:

```bash
cd Customer-Service
# Follow the README in that directory
```

## Future Services

This architecture supports adding more banking microservices:
- Account-Service
- Transaction-Service
- Loan-Service
- Card-Service
- Notification-Service

Each service will follow the same standalone pattern with MVC architecture.