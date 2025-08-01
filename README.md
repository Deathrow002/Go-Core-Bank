# Core Banking Microservices Platform

A modular microservice architecture for core banking operations.

## Project Structure

```
Go-Core-Bank/
├── Authentication-Service/ # Authentication microservice (user auth, JWT, password hashing)
│   ├── cmd/               # Application entry points
│   ├── internal/          # Private application code
│   │   ├── authentication/ # Authentication domain
│   │   │   ├── controller/ # HTTP controllers (Gin handlers)
│   │   │   ├── service/    # Business logic (password, JWT, validation)
│   │   │   ├── repository/ # Data access (GORM)
│   │   │   └── models/     # Domain models (user, roles)
│   │   └── external/       # External service clients (e.g., Customer Service)
│   ├── dockerfile          # Docker build file
│   ├── go.mod              # Go module definition
│   └── README.md           # Service-specific docs
├── Account-Service/        # Account microservice (accounts, balances, transactions)
│   ├── cmd/
│   │   │   ├── controller/ # HTTP controllers (Gin handlers)
│   │   │   ├── service/    # Business logic (password, JWT, validation)
│   │   │   ├── repository/ # Data access (GORM)
│   │   │   └── models/     # Domain models (user, roles)
│   │   └── external/       # External service clients (e.g., Customer Service)
│   ├── internal/
│   └── README.md
├── Customer-Service/       # Customer microservice (standalone)
│   ├── cmd/
│   │   │   ├── controller/ # HTTP controllers (Gin handlers)
│   │   │   ├── service/    # Business logic (password, JWT, validation)
│   │   │   ├── repository/ # Data access (GORM)
│   │   │   └── models/     # Domain models (user, roles)
│   ├── internal/
│   └── pkg/
│   └── README.md
├── docker-compose.yml      # Multi-service deployment
├── Makefile                # Build automation
└── README.md               # This file
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
# Build Authentication Service
make authentication-service

# Build Account Service
make account-service

# Build Customer Service
make customer-service

# Or build all services
make build
```

## Services

### Authentication Service
- **Location**: `./Authentication-Service/`
- **Port**: 8082
- **Features**:
  - User registration and login
  - Password hashing (bcrypt)
  - JWT authentication
  - Account lock/unlock
  - Role-based access (admin, support, user)
  - Inter-service validation with Customer Service
- **Documentation**: See `./Authentication-Service/README.md`

### Account Service
- **Location**: `./Account-Service/`
- **Port**: 8081
- **Features**:
  - Account creation and management
  - Balance tracking
  - Transaction history
  - Integration with Authentication and Customer Services
- **Documentation**: See `./Account-Service/README.md`

### Customer Service
- **Location**: `./Customer-Service/`
- **Port**: 8080
- **Features**:
  - Customer profile management
  - Customer data validation
  - Integration with Authentication Service
- **Documentation**: See `./Customer-Service/README.md`

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
cd Authentication-Service
# Follow the README in that directory

cd Account-Service
# Follow the README in that directory

cd Customer-Service
# Follow the README in that directory
```

## Future Services

This architecture supports adding more banking microservices:
- Transaction-Service
- Loan-Service
- Card-Service
- Notification-Service

Each service will follow the same standalone pattern with MVC architecture.