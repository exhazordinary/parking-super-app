# Parking Super App - Backend Microservices

A cloud-agnostic backend microservices system for a centralized parking application in Malaysia. Integrates multiple parking providers through a unified API.

## Architecture

```
                              Mobile App
                                  │
                                  ▼
                          ┌──────────────┐
                          │ API Gateway  │
                          │  REST :8080  │
                          └──────────────┘
                                  │
        ┌─────────┬───────┬───────┼───────┬───────────┐
        │         │       │       │       │           │
        ▼         ▼       ▼       ▼       ▼           ▼
   ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────────┐
   │  Auth  │ │ Wallet │ │Provider│ │Parking │ │Notification│
   │:8081   │ │:8082   │ │:8083   │ │:8084   │ │:8085       │
   │:9081   │ │:9082   │ │:9083   │ │:9084   │ │:9085       │
   └────────┘ └────────┘ └────────┘ └────────┘ └────────────┘
        │         │       │       │       │           │
        │         └───────┼───────┼───────┼───────────┤
        │                 │       │       │           │
        ▼                 ▼       ▼       ▼           ▼
   ┌─────────────────────────────────────────────────────┐
   │                    PostgreSQL                        │
   │              (database per service)                  │
   └─────────────────────────────────────────────────────┘
        │                 │       │       │           │
        └─────────────────┼───────┼───────┼───────────┘
                          │       │       │
                          ▼       ▼       ▼
                    ┌─────────────────────────┐
                    │         Kafka           │
                    │    (Event Streaming)    │
                    └─────────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────┐
                    │         Jaeger          │
                    │  (Distributed Tracing)  │
                    └─────────────────────────┘
```

## Services

| Service | REST | gRPC | Description |
|---------|------|------|-------------|
| API Gateway | 8080 | - | Request routing, JWT validation, rate limiting |
| Auth | 8081 | 9081 | User registration, login, JWT tokens |
| Wallet | 8082 | 9082 | Digital wallet, transactions, payments |
| Provider | 8083 | 9083 | Parking provider management |
| Parking | 8084 | 9084 | Parking session management |
| Notification | 8085 | 9085 | Push, SMS, email notifications |

## Tech Stack

- **Language**: Go 1.24+
- **Architecture**: Hexagonal (Ports & Adapters)
- **HTTP Router**: Chi
- **Database**: PostgreSQL with pgx driver
- **Auth**: JWT (golang-jwt/jwt)
- **Inter-Service Communication**: gRPC
- **Event Streaming**: Apache Kafka
- **Distributed Tracing**: OpenTelemetry + Jaeger
- **Containerization**: Docker + Docker Compose

## Project Structure

```
parking-super-app/
├── go.work                    # Go workspace
├── pkg/                       # Shared packages
│   ├── kafka/                 # Kafka publisher/consumer
│   ├── telemetry/             # OpenTelemetry setup
│   ├── grpc/interceptors/     # gRPC middleware
│   ├── middleware/            # HTTP middleware
│   └── proto/                 # Protocol Buffer definitions
├── services/
│   ├── api-gateway/           # API Gateway
│   ├── auth/                  # Auth Service
│   ├── wallet/                # Wallet Service
│   ├── provider/              # Provider Service
│   ├── parking/               # Parking Service
│   └── notification/          # Notification Service
└── deployments/
    └── docker/                # Docker Compose configs
```

### Service Structure (Hexagonal Architecture)

```
service/
├── cmd/server/          # Entry point
├── config/              # Configuration
├── internal/
│   ├── domain/          # Business entities
│   ├── ports/           # Interfaces (repository, services)
│   ├── application/     # Use cases / business logic
│   └── adapters/
│       ├── http/        # REST handlers
│       ├── grpc/        # gRPC server/clients
│       ├── repository/  # Database implementations
│       └── external/    # External service clients
├── migrations/          # Database migrations
├── Dockerfile
└── Makefile
```

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL (if running locally)

### Run Full Stack with Docker Compose

```bash
cd deployments/docker
docker-compose up -d
```

This starts:
- PostgreSQL (port 5432)
- Zookeeper + Kafka (port 9092)
- Jaeger UI (port 16686)
- All 6 microservices

### Access Services

| Service | URL |
|---------|-----|
| API Gateway | http://localhost:8080 |
| Jaeger UI | http://localhost:16686 |
| Kafka | localhost:9092 |

### Run a Single Service Locally

```bash
cd services/auth
cp .env.example .env   # Configure environment
make run
```

### Run Tests

```bash
cd services/auth
make test
```

## API Endpoints

### Auth Service

```
POST /api/v1/auth/register     Register new user
POST /api/v1/auth/login        Login
POST /api/v1/auth/refresh      Refresh access token
POST /api/v1/auth/logout       Logout
GET  /api/v1/auth/me           Get current user
```

### Wallet Service

```
GET  /api/v1/wallet            Get wallet balance
POST /api/v1/wallet/topup      Top-up wallet
POST /api/v1/wallet/pay        Make payment
GET  /api/v1/wallet/txns       Transaction history
```

### Provider Service

```
GET  /api/v1/providers         List providers
GET  /api/v1/providers/:id     Get provider details
POST /api/v1/providers         Register provider (admin)
```

### Parking Service

```
POST /api/v1/parking/sessions          Start session
GET  /api/v1/parking/sessions          List user sessions
GET  /api/v1/parking/sessions/:id      Get session details
PUT  /api/v1/parking/sessions/:id      Update session
DELETE /api/v1/parking/sessions/:id    End session
```

### Notification Service

```
POST /api/v1/notifications             Send notification
GET  /api/v1/preferences               Get user preferences
PUT  /api/v1/preferences               Update preferences
```

## Configuration

### Environment Variables

Each service requires its own `.env` file. See `.env.example` in each service directory.

```bash
# Server
SERVER_PORT=8081
GRPC_PORT=9081

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=auth_db

# JWT
JWT_SECRET=your-secret-key

# Kafka (optional)
KAFKA_ENABLED=true
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=auth.events

# OpenTelemetry (optional)
OTEL_ENABLED=true
OTEL_SERVICE_NAME=auth-service
OTEL_ENDPOINT=localhost:4317
```

## Database Migrations

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
cd services/auth
migrate -path ./migrations -database "postgres://user:pass@localhost/auth_db?sslmode=disable" up
```

## Observability

### Distributed Tracing

All services are instrumented with OpenTelemetry. Traces are exported to Jaeger.

1. Start the stack: `docker-compose up -d`
2. Open Jaeger UI: http://localhost:16686
3. Select a service from the dropdown
4. View traces across service boundaries

### Kafka Events

Services publish domain events to Kafka topics:

| Topic | Publisher | Events |
|-------|-----------|--------|
| `auth.events` | Auth | user.registered, user.logged_in |
| `wallet.events` | Wallet | payment.completed, topup.completed |
| `parking.events` | Parking | session.started, session.ended |
| `provider.events` | Provider | provider.registered |

## Development

### Add New Service to Workspace

```bash
go work use ./services/new-service
```

### Update Dependencies

```bash
cd services/auth
go mod tidy
```

### Build All Services

```bash
cd deployments/docker
docker-compose build
```

## License

MIT
