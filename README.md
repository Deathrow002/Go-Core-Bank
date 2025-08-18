# Core Banking Microservices Platform

A modular microservice architecture for core banking operations.

---

## Project Structure

```
Go-Core-Bank/
├── Authentication-Service/     # Authentication microservice (user auth, JWT, password hashing)
│   ├── cmd/                    # Application entry points
│   ├── internal/               # Private application code
│   │   ├── authentication/     # Authentication domain
│   │   │   ├── controller/     # HTTP controllers (Gin handlers)
│   │   │   ├── service/        # Business logic (password, JWT, validation)
│   │   │   ├── repository/     # Data access (GORM)
│   │   │   └── models/         # Domain models (user, roles)
│   │   └── external/           # External service clients (e.g., Customer Service)
│   ├── dockerfile              # Docker build file
│   ├── go.mod                  # Go module definition
│   └── README.md               # Service-specific docs
├── Account-Service/            # Account microservice (accounts, balances, transactions)
│   ├── cmd/
│   ├── internal/
│   │   ├── account/
│   │   │   ├── controllers/    # HTTP controllers (Gin handlers)
│   │   │   ├── service/        # Business logic
│   │   │   ├── repository/     # Data access (GORM)
│   │   │   └── models/         # Domain models
│   │   └── external/           # External service clients (e.g., Customer Service)
│   │       └── consumer/       # Kafka consumer for balance updates
│   ├── dockerfile
│   ├── go.mod
│   └── README.md
├── Customer-Service/           # Customer microservice (standalone)
│   ├── cmd/
│   ├── internal/
│   │   ├── customer/
│   │   │   ├── controllers/    # HTTP controllers (Gin handlers)
│   │   │   ├── service/        # Business logic
│   │   │   ├── repository/     # Data access (GORM)
│   │   │   └── models/         # Domain models
│   ├── pkg/
│   ├── dockerfile
│   ├── go.mod
│   └── README.md
├── Transaction-Service/        # Transaction microservice (transfers, deposits, withdrawals)
│   ├── cmd/
│   ├── internal/
│   │   ├── transaction/
│   │   │   ├── controller/     # HTTP controllers (Gin handlers)
│   │   │   ├── service/        # Business logic
│   │   │   ├── repository/     # Data access (GORM)
│   │   │   └── models/         # Domain models
│   │   └── external/
│   │       └── producer/       # Kafka producer for balance updates
│   ├── dockerfile
│   ├── go.mod
│   └── README.md
├── PostmanCollection/          # API documentation and test collections
│   └── Core-Bank_Collections   # Postman collection for all services
├── K8s/                        # Kubernetes manifests
├── prometheus/
│   └── prometheus.yml          # Prometheus configuration
├── docker-compose.yml          # Multi-service deployment with Kafka
├── Makefile                    # Build automation
└── README.md                   # This file
```

---

## Quick Start


### Using Docker (Recommended)
```bash
# Run all services (including monitoring)
make docker-run

# Stop all services
make docker-stop
```

### Monitoring with Prometheus & Grafana

Prometheus and Grafana are included in `docker-compose.yml` for monitoring and visualization.

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000      (login: admin / password: admin)

Prometheus scrapes metrics from each microservice (if `/metrics` endpoint is exposed). You can add custom metrics to your Go services using Prometheus client libraries.

Grafana can be used to visualize metrics and create dashboards. Add Prometheus as a data source in Grafana and import dashboards as needed.

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

### Transaction Service
- **Location**: `./Transaction-Service/`
- **Port**: 8083
- **Features**:
  - Funds transfers between accounts
  - Deposits and withdrawals
  - Transaction history
  - Integration with Account and Customer Services
- **Environment**:
  - `KAFKA_BROKER_URL` should be set to `kafka:9092` when running in Docker Compose.
- **Documentation**: See `./Transaction-Service/README.md`

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

**Monitoring** is provided by Prometheus (metrics collection) and Grafana (visualization). See `prometheus.yml` for scrape targets and configure your services to expose metrics endpoints for best results.

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

cd Transaction-Service
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
- Notification-Service

Each service will follow the same standalone pattern with MVC architecture.