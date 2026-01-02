.PHONY: build run test clean migrate-up migrate-down docker-build

# Build binary
build:
	go build -o bin/wallet-service ./cmd/server

# Run the service
run:
	go run ./cmd/server

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run migrations up
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/wallet_db?sslmode=disable" up

# Run migrations down
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/wallet_db?sslmode=disable" down

# Build Docker image
docker-build:
	docker build -t wallet-service:latest .

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy
