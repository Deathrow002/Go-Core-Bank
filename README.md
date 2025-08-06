# Core Banking Microservices Platform

A modular microservice architecture for core banking operations.

---

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
│   ├── internal/
│   │   ├── account/
│   │   │   ├── controllers/ # HTTP controllers (Gin handlers)
│   │   │   ├── service/     # Business logic
│   │   │   ├── repository/  # Data access (GORM)
│   │   │   └── models/      # Domain models
│   │   └── external/        # External service clients (e.g., Customer Service)
│   ├── dockerfile
│   ├── go.mod
│   └── README.md
├── Customer-Service/       # Customer microservice (standalone)
│   ├── cmd/
│   ├── internal/
│   │   ├── customer/
│   │   │   ├── controllers/ # HTTP controllers (Gin handlers)
│   │   │   ├── service/     # Business logic
│   │   │   ├── repository/  # Data access (GORM)
│   │   │   └── models/      # Domain models
│   ├── pkg/
│   ├── dockerfile
│   ├── go.mod
│   └── README.md
├── docker-compose.yml      # Multi-service deployment
├── Makefile                # Build automation
└── README.md               # This file
```

---

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

---

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
- **Environment**:
  - `CUSTOMER_SERVICE_URL` should be set to `http://customer-service:8080` when running in Docker Compose.
- **Documentation**: See `./Authentication-Service/README.md`

### Account Service
- **Location**: `./Account-Service/`
- **Port**: 8081
- **Features**:
  - Account creation and management
  - Balance tracking
  - Transaction history
  - Integration with Authentication and Customer Services
- **Environment**:
  - `CUSTOMER_SERVICE_URL` should be set to `http://customer-service:8080` when running in Docker Compose.
- **Documentation**: See `./Account-Service/README.md`

### Customer Service
- **Location**: `./Customer-Service/`
- **Port**: 8080
- **Features**:
  - Customer profile management
  - Customer data validation
  - Integration with Authentication Service
- **Documentation**: See `./Customer-Service/README.md`

---

## Architecture

Each microservice is completely standalone with its own:
- Go module (`go.mod`)
- Database schema
- Docker configuration
- Documentation
- Build system
- MVC architecture (Controllers, Services, Repositories)

**Inter-service communication** is done via HTTP using Docker Compose service names as hostnames (e.g., `http://customer-service:8080`).

---

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

---

## Troubleshooting

- **Service-to-service communication in Docker Compose:**  
  Use service names (e.g., `customer-service:8080`) instead of `localhost` for URLs between services.
- **Connection refused errors:**  
  Ensure the target service is running and listening on `0.0.0.0`, and the correct environment variable is set.
- **JWT errors:**  
  Ensure all services use the same `JWT_SECRET` and token format.

---

## Future Services

This architecture supports adding more banking microservices:
- Transaction-Service
- Distributed Message Queue
- Loan-Service
- Card-Service
- Notification-Service

Each service will follow the same standalone pattern with MVC architecture.