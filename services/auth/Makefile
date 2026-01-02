# Makefile for Auth Service
#
# MICROSERVICES PATTERN: Build Automation
# =======================================
# Makefiles provide consistent, documented commands for:
# - Building the service
# - Running tests
# - Managing database migrations
# - Starting development environment

.PHONY: all build run test clean docker migrate help

# Default target
all: build

# Variables
SERVICE_NAME := auth-service
MAIN_PATH := ./cmd/server
DOCKER_IMAGE := parking-super-app/$(SERVICE_NAME)
DB_URL := postgres://postgres:postgres@localhost:5433/auth_db?sslmode=disable

# ================================================
# BUILD COMMANDS
# ================================================

## build: Build the service binary
build:
	@echo "Building $(SERVICE_NAME)..."
	go build -o bin/$(SERVICE_NAME) $(MAIN_PATH)

## run: Run the service locally
run:
	@echo "Running $(SERVICE_NAME)..."
	go run $(MAIN_PATH)/main.go

## run-watch: Run with live reload (requires air)
run-watch:
	@echo "Running $(SERVICE_NAME) with live reload..."
	air -c .air.toml

# ================================================
# TEST COMMANDS
# ================================================

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## test-unit: Run only unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v -short ./...

# ================================================
# DOCKER COMMANDS
# ================================================

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .

## docker-run: Run in Docker
docker-run:
	@echo "Running in Docker..."
	docker run -p 8080:8080 \
		-e DB_HOST=host.docker.internal \
		-e DB_PORT=5432 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=postgres \
		-e DB_NAME=auth_db \
		$(DOCKER_IMAGE):latest

## docker-up: Start all dependencies
docker-up:
	@echo "Starting dependencies..."
	docker-compose -f ../../deployments/docker/docker-compose.yml up -d

## docker-down: Stop all dependencies
docker-down:
	@echo "Stopping dependencies..."
	docker-compose -f ../../deployments/docker/docker-compose.yml down

## docker-logs: View logs
docker-logs:
	docker-compose -f ../../deployments/docker/docker-compose.yml logs -f

# ================================================
# DATABASE COMMANDS
# ================================================

## migrate-up: Run all migrations
migrate-up:
	@echo "Running migrations..."
	migrate -path migrations -database "$(DB_URL)" up

## migrate-down: Rollback last migration
migrate-down:
	@echo "Rolling back migration..."
	migrate -path migrations -database "$(DB_URL)" down 1

## migrate-create: Create new migration (usage: make migrate-create name=add_column)
migrate-create:
	@echo "Creating migration..."
	migrate create -ext sql -dir migrations -seq $(name)

## migrate-force: Force migration version (usage: make migrate-force version=1)
migrate-force:
	@echo "Forcing migration version..."
	migrate -path migrations -database "$(DB_URL)" force $(version)

# ================================================
# DEVELOPMENT COMMANDS
# ================================================

## deps: Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

## generate: Run code generation
generate:
	@echo "Running code generation..."
	go generate ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# ================================================
# HELP
# ================================================

## help: Show this help message
help:
	@echo "Available commands:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
