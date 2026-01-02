FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /notification-service .
COPY migrations ./migrations
USER appuser
EXPOSE 8084
CMD ["./notification-service"]
