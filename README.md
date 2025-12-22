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


## Kubernetes Deployment

Use helper scripts under `Deployment/`.

### Deploy (`Deployment/deploy.sh`)

**Flags:**
- `-n <namespace>`: Target namespace (default: `core-bank`)
- `-y`: Auto-confirm (skip interactive prompt)
- `-w`: Wait for deployments to be ready
- `-c`: Check images exist before deploy
- `-b`: Auto-build missing images using `dockerbuild.sh`
- `-L <target>`: Load built images to local cluster (`docker-desktop`, `minikube`, or `kind`)
- `-I`: Install nginx ingress controller (ingress-nginx namespace)
- `-G`: Port-forward ingress controller to `localhost:8080` after deploy
- `-S`: Port-forward services to `localhost:18080-18083` after deploy
- `-h`: Show help

**Running from PowerShell (Windows):**
```powershell
# Use Git Bash to run the script
& "C:\Program Files\Git\bin\bash.exe" Deployment/deploy.sh -n core-bank -c -b -L docker-desktop -I -w -G -S -y
```

**Running from Git Bash or Linux/macOS:**
```bash
# Full deployment with auto-build, ingress, and port-forwards
bash Deployment/deploy.sh -n core-bank -c -b -L docker-desktop -I -w -G -S -y

# Minimal deployment (no ingress, no port-forwards)
bash Deployment/deploy.sh -n core-bank -y

# Deploy with auto-build for missing images
bash Deployment/deploy.sh -n core-bank -c -b -L docker-desktop -w -y
```

**Access Methods:**

After deploying with `-S` (service port-forwards):
- Account Service: `http://localhost:18080/health`
- Authentication Service: `http://localhost:18081/health`
- Customer Service: `http://localhost:18082/health`
- Transaction Service: `http://localhost:18083/health`

After deploying with `-G` (ingress port-forward), use Host headers:
```bash
curl -H 'Host: account-service.local' http://localhost:8080/health
curl -H 'Host: authentication-service.local' http://localhost:8080/health
curl -H 'Host: customer-service.local' http://localhost:8080/health
curl -H 'Host: transaction-service.local' http://localhost:8080/health
```

### Build Images (`Deployment/dockerbuild.sh`)

**Flags:**
- `-t <tag>`: Image tag (defaults to git SHA or `latest`)
- `-r <registry>`: Registry prefix (e.g., `ghcr.io/yourorg/core-bank`)
- `-p`: Push images after build
- `-s <svc1,svc2>`: Build subset of services (comma-separated)
- `-n`: Disable build cache (adds `--no-cache`)
- `-j`: Parallel build (background jobs)
- `-L <target>`: Load built images to local cluster (`docker-desktop`, `minikube`, or `kind`)
- `-h`: Show help

**Examples:**
```bash
# Build and push all services with version tag
bash Deployment/dockerbuild.sh -r ghcr.io/yourorg/core-bank -t v1.0.0 -p

# Build subset for local development
bash Deployment/dockerbuild.sh -s account-service,transaction-service -t dev

# Build all and load to Docker Desktop (for K8s)
bash Deployment/dockerbuild.sh -t latest -L docker-desktop

# Build for minikube with parallel jobs
bash Deployment/dockerbuild.sh -t latest -L minikube -j
```

### Cleanup (`Deployment/clearup.sh`)
### Cleanup (`Deployment/clearup.sh`)

**Flags:**
- `-n <namespace>`: Target namespace (default: `core-bank`)
- `-F`: Delete namespace after resources are removed
- `-P`: Delete persistent storage (PVCs) and related PVs
- `-d`: Dry run (show actions without executing)
- `-y`: Auto-confirm (skip interactive prompt)
- `-w`: Wait for workload termination (and optionally namespace deletion)
- `-i`: Remove local Docker images for first-party services
- `-I`: Remove local Docker images including third-party dependencies

**Examples:**
```bash
# Remove workloads and manifests only (keep namespace & storage)
bash Deployment/clearup.sh -n core-bank -y -w

# Full cleanup: workloads + namespace + PVC/PV
bash Deployment/clearup.sh -n core-bank -F -P -y -w

# Dry-run to preview what will be deleted
bash Deployment/clearup.sh -n core-bank -d

# Also remove locally built service images
bash Deployment/clearup.sh -n core-bank -i -y

# Also remove third-party images (e.g., postgres, kafka)
bash Deployment/clearup.sh -n core-bank -I -y
```

**Windows (PowerShell) via Git Bash:**
```powershell
& "C:\Program Files\Git\bin\bash.exe" Deployment/clearup.sh -n core-bank -F -P -y -w
```

**Safety Notes:**
- Deleting PVs (`-P`) is irreversible for local clusters; data will be lost.
- Use `-d` first if you are unsure; it prints intended actions.
- Image removal (`-i`/`-I`) affects only local Docker images, not remote registries.

### Service Discovery on K8s

**Inside the cluster (pod-to-pod):**
- `account-service:8080`
- `authentication-service:8081`
- `customer-service:8082`
- `transaction-service:8083`

**Environment variables for inter-service communication:**
- `ACCOUNT_SERVICE_URL=http://account-service:8080`
- `AUTHENTICATION_SERVICE_URL=http://authentication-service:8081`
- `CUSTOMER_SERVICE_URL=http://customer-service:8082`
- `TRANSACTION_SERVICE_URL=http://transaction-service:8083`

**From your local machine (with port-forwards via `-S`):**
- Account: `http://localhost:18080`
- Authentication: `http://localhost:18081`
- Customer: `http://localhost:18082`
- Transaction: `http://localhost:18083`

### Network Policies (recommended)
Add baseline deny-all, allow intra-namespace app ports, and allow egress to data services (Postgres 5432, Redis 6379, Kafka 9092). Place under `K8s/network-policy.yaml` and include in deploy manifests.


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
  - `CUSTOMER_SERVICE_URL` should be set to `http://account-service:8080` when running in Docker Compose.
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