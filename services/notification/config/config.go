package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	GRPC     GRPCConfig
	Kafka    KafkaConfig
	OTEL     OTELConfig
	Provider ProviderConfig
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
	Brokers       []string
	Topics        []string // Topics to consume from
	ConsumerGroup string
	Enabled       bool
}

type OTELConfig struct {
	Enabled     bool
	Endpoint    string
	ServiceName string
	Insecure    bool
}

// ProviderConfig holds notification provider settings
type ProviderConfig struct {
	SMS   string // "console", "twilio"
	Email string // "console", "sendgrid"
	Push  string // "console", "firebase"
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

	brokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	topics := strings.Split(getEnv("KAFKA_TOPICS", "parking.events,wallet.events,auth.events"), ",")

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
			DBName:   getEnv("DB_NAME", "notification_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers:       brokers,
			Topics:        topics,
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "notification-service"),
			Enabled:       kafkaEnabled,
		},
		OTEL: OTELConfig{
			Enabled:     otelEnabled,
			Endpoint:    getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			ServiceName: getEnv("OTEL_SERVICE_NAME", "notification-service"),
			Insecure:    otelInsecure,
		},
		Provider: ProviderConfig{
			SMS:   getEnv("SMS_PROVIDER", "console"),
			Email: getEnv("EMAIL_PROVIDER", "console"),
			Push:  getEnv("PUSH_PROVIDER", "console"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
