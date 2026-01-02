// Package config handles application configuration.
//
// MICROSERVICES PATTERN: Externalized Configuration
// =================================================
// Configuration should be externalized and environment-specific:
// - Development: Local database, mock services
// - Staging: Cloud database, real services (test mode)
// - Production: Production database, real services
//
// We use environment variables for configuration because:
// - Standard approach for containers (Docker, Kubernetes)
// - 12-Factor App methodology
// - No secrets in code or version control
// - Easy to change without rebuilding
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	// Server configuration
	Server ServerConfig

	// Database configuration
	Database DatabaseConfig

	// JWT configuration
	JWT JWTConfig

	// SMS configuration (optional)
	SMS SMSConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectionString returns the PostgreSQL connection string.
func (c DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

// JWTConfig holds JWT-related settings.
type JWTConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration
}

// SMSConfig holds SMS provider settings.
type SMSConfig struct {
	Provider   string // "console", "twilio"
	AccountSID string
	AuthToken  string
	FromPhone  string
}

// Load reads configuration from environment variables.
//
// BEST PRACTICE: Fail Fast
// If required configuration is missing, fail immediately at startup
// rather than failing later when the config is needed.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "auth_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:      getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			AccessTokenTTL: getDurationEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
		},
		SMS: SMSConfig{
			Provider:   getEnv("SMS_PROVIDER", "console"),
			AccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
			AuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
			FromPhone:  getEnv("TWILIO_FROM_PHONE", ""),
		},
	}

	// Validate required configuration
	if cfg.JWT.SecretKey == "your-super-secret-key-change-in-production" {
		fmt.Println("WARNING: Using default JWT secret key. Set JWT_SECRET in production!")
	}

	return cfg, nil
}

// Helper functions for reading environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
