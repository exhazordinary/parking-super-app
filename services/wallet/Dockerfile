# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet-service ./cmd/server

# Final stage
FROM alpine:3.19

WORKDIR /app

# Add non-root user for security
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /wallet-service .

# Copy migrations
COPY migrations ./migrations

# Use non-root user
USER appuser

EXPOSE 8081

CMD ["./wallet-service"]
