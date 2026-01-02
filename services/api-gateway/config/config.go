package config

import (
	"os"
	"strconv"
)

// Config holds API Gateway configuration
type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	Auth     AuthConfig
	OTEL     OTELConfig
}

type ServerConfig struct {
	Port string
}

type ServicesConfig struct {
	// HTTP URLs for proxying REST requests
	AuthURL         string
	WalletURL       string
	ProviderURL     string
	ParkingURL      string
	NotificationURL string

	// gRPC addresses for internal communication
	AuthGRPC         string
	WalletGRPC       string
	ProviderGRPC     string
	ParkingGRPC      string
	NotificationGRPC string
}

type AuthConfig struct {
	JWTSecret string
}

type OTELConfig struct {
	Enabled     bool
	Endpoint    string
	ServiceName string
	Insecure    bool
}

func Load() (*Config, error) {
	otelEnabled, _ := strconv.ParseBool(getEnv("OTEL_ENABLED", "false"))
	otelInsecure, _ := strconv.ParseBool(getEnv("OTEL_INSECURE", "true"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Services: ServicesConfig{
			// HTTP URLs
			AuthURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
			WalletURL:       getEnv("WALLET_SERVICE_URL", "http://localhost:8082"),
			ProviderURL:     getEnv("PROVIDER_SERVICE_URL", "http://localhost:8083"),
			ParkingURL:      getEnv("PARKING_SERVICE_URL", "http://localhost:8084"),
			NotificationURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8085"),
			// gRPC addresses
			AuthGRPC:         getEnv("AUTH_SERVICE_GRPC", "localhost:9081"),
			WalletGRPC:       getEnv("WALLET_SERVICE_GRPC", "localhost:9082"),
			ProviderGRPC:     getEnv("PROVIDER_SERVICE_GRPC", "localhost:9083"),
			ParkingGRPC:      getEnv("PARKING_SERVICE_GRPC", "localhost:9084"),
			NotificationGRPC: getEnv("NOTIFICATION_SERVICE_GRPC", "localhost:9085"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		},
		OTEL: OTELConfig{
			Enabled:     otelEnabled,
			Endpoint:    getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			ServiceName: getEnv("OTEL_SERVICE_NAME", "api-gateway"),
			Insecure:    otelInsecure,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
