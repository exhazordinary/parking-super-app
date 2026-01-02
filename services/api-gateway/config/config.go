package config

import (
	"os"
)

// Config holds API Gateway configuration
type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type ServicesConfig struct {
	AuthURL         string
	WalletURL       string
	ProviderURL     string
	ParkingURL      string
	NotificationURL string
}

type AuthConfig struct {
	JWTSecret string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Services: ServicesConfig{
			AuthURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8000"),
			WalletURL:       getEnv("WALLET_SERVICE_URL", "http://localhost:8081"),
			ProviderURL:     getEnv("PROVIDER_SERVICE_URL", "http://localhost:8082"),
			ParkingURL:      getEnv("PARKING_SERVICE_URL", "http://localhost:8083"),
			NotificationURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8084"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
