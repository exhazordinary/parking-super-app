.PHONY: build run test clean docker-build

build:
	go build -o bin/api-gateway ./cmd/server

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

docker-build:
	docker build -t api-gateway:latest .

fmt:
	go fmt ./...

tidy:
	go mod tidy
