# Multi-stage Dockerfile for Go microservice
#
# PATTERN: Multi-stage Build
# ==========================
# Stage 1: Build the binary with all dependencies
# Stage 2: Copy only the binary to a minimal image
#
# Benefits:
# - Small final image (~10MB vs ~1GB)
# - No build tools in production
# - Reduced attack surface
# - Faster deployments

# ================================================
# Stage 1: Build
# ================================================
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
# This layer is cached if go.mod/go.sum don't change
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=0: Disable CGO for static binary
# -ldflags="-w -s": Strip debug info for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /auth-service \
    ./cmd/server

# ================================================
# Stage 2: Runtime
# ================================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /auth-service .

# Copy migrations (for running migrations on startup)
COPY --from=builder /app/migrations ./migrations

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the service port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./auth-service"]
