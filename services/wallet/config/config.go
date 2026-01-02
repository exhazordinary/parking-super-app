package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the wallet service.
// Configuration is loaded from environment variables following 12-factor app principles.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
	GRPC     GRPCConfig
	OTEL     OTELConfig
}

type ServerConfig struct {
	Port string
}

type GRPCConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	Enabled bool
}

type OTELConfig struct {
	Enabled     bool
	Endpoint    string
	ServiceName string
	Insecure    bool
}

func (d DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

func Load() (*Config, error) {
	kafkaEnabled, _ := strconv.ParseBool(getEnv("KAFKA_ENABLED", "false"))
	otelEnabled, _ := strconv.ParseBool(getEnv("OTEL_ENABLED", "false"))
	otelInsecure, _ := strconv.ParseBool(getEnv("OTEL_INSECURE", "true"))

	// Parse Kafka brokers (comma-separated)
	brokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "9000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "wallet_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers: brokers,
			Topic:   getEnv("KAFKA_TOPIC", "wallet.events"),
			Enabled: kafkaEnabled,
		},
		OTEL: OTELConfig{
			Enabled:     otelEnabled,
			Endpoint:    getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			ServiceName: getEnv("OTEL_SERVICE_NAME", "wallet-service"),
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
