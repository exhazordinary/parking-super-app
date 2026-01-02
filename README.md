# Parking Super App - Backend Microservices

A cloud-agnostic backend microservices system for a centralized parking application in Malaysia. Integrates multiple parking providers through a unified API.

## Architecture

```
                         Mobile App
                             |
                             v
                     +--------------+
                     | API Gateway  |
                     |   (8080)     |
                     +--------------+
                             |
     +--------+--------+-----+-----+--------+
     |        |        |          |         |
     v        v        v          v         v
+-------+ +--------+ +--------+ +-------+ +--------------+
| Auth  | | Wallet | |Provider| |Parking| | Notification |
| 8080  | | 8081   | | 8082   | | 8083  | |    8084      |
+-------+ +--------+ +--------+ +-------+ +--------------+
     |        |        |          |         |
     v        v        v          v         v
   PostgreSQL (separate database per service)
```

## Services

| Service | Port | Description |
|---------|------|-------------|
| API Gateway | 8080 | Request routing, JWT validation, rate limiting |
| Auth | 8080 | User registration, login, JWT tokens |
| Wallet | 8081 | Digital wallet, transactions, payments |
| Provider | 8082 | Parking provider management |
| Parking | 8083 | Parking session management |
| Notification | 8084 | Push, SMS, email notifications |

## Tech Stack

- **Language**: Go 1.22+
- **Architecture**: Hexagonal (Ports & Adapters)
- **HTTP Router**: Chi
- **Database**: PostgreSQL with pgx driver
- **Auth**: JWT (golang-jwt/jwt)
- **Decimal**: shopspring/decimal

## Project Structure

```
parking-super-app/
├── go.work              # Go workspace
├── services/
│   ├── api-gateway/     # API Gateway
│   ├── auth/            # Auth Service
│   ├── wallet/          # Wallet Service
│   ├── provider/        # Provider Service
│   ├── parking/         # Parking Service
│   └── notification/    # Notification Service
└── deployments/         # Docker, K8s configs
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
│       ├── repository/  # Database implementations
│       └── external/    # External service clients
├── migrations/          # Database migrations
├── Dockerfile
└── Makefile
```

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL
- Docker (optional)

### Run a Service

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

### Build Docker Image

```bash
cd services/auth
make docker-build
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

## Environment Variables

Each service requires its own `.env` file. See `.env.example` in each service directory.

Common variables:
```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=service_db
JWT_SECRET=your-secret-key
```

## Database Migrations

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
cd services/auth
migrate -path ./migrations -database "postgres://user:pass@localhost/auth_db?sslmode=disable" up
```

## Development

### Add New Service to Workspace

```bash
go work use ./services/new-service
```

### Format Code

```bash
cd services/auth
make fmt
```

### Tidy Dependencies

```bash
cd services/auth
make tidy
```

## License

MIT
