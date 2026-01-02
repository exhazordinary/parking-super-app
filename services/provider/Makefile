.PHONY: build run test clean migrate-up migrate-down docker-build

build:
	go build -o bin/provider-service ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/provider_db?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/provider_db?sslmode=disable" down

docker-build:
	docker build -t provider-service:latest .

fmt:
	go fmt ./...

tidy:
	go mod tidy
