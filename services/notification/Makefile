.PHONY: build run test clean migrate-up migrate-down docker-build

build:
	go build -o bin/notification-service ./cmd/server

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
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/notification_db?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/notification_db?sslmode=disable" down

docker-build:
	docker build -t notification-service:latest .

fmt:
	go fmt ./...

tidy:
	go mod tidy
