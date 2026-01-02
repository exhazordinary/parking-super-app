// Package main is the entry point for the auth service.
//
// MICROSERVICES PATTERN: Service Entry Point
// ==========================================
// The main package is responsible for:
// 1. Loading configuration
// 2. Creating dependencies (dependency injection)
// 3. Wiring everything together
// 4. Starting the HTTP and gRPC servers
// 5. Graceful shutdown handling
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/pkg/grpc/interceptors"
	"github.com/parking-super-app/pkg/kafka"
	"github.com/parking-super-app/pkg/middleware"
	"github.com/parking-super-app/pkg/telemetry"
	"github.com/parking-super-app/services/auth/config"
	"github.com/parking-super-app/services/auth/internal/adapters/external"
	httpAdapter "github.com/parking-super-app/services/auth/internal/adapters/http"
	"github.com/parking-super-app/services/auth/internal/adapters/repository/postgres"
	"github.com/parking-super-app/services/auth/internal/application"
	"github.com/parking-super-app/services/auth/internal/domain"
	"github.com/parking-super-app/services/auth/internal/ports"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting auth service on port %s", cfg.Server.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize OpenTelemetry tracing
	var tracerShutdown func(context.Context) error
	if cfg.OTEL.Enabled {
		shutdown, err := telemetry.InitTracer(ctx, telemetry.Config{
			ServiceName:  cfg.OTEL.ServiceName,
			OTLPEndpoint: cfg.OTEL.Endpoint,
			Insecure:     cfg.OTEL.Insecure,
			Environment:  "development",
		})
		if err != nil {
			log.Printf("warning: failed to initialize tracer: %v", err)
		} else {
			tracerShutdown = shutdown
			log.Println("OpenTelemetry tracing initialized")
		}
	}

	// Set up database connection
	dbPool, err := pgxpool.New(ctx, cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Create dependencies
	userRepo := postgres.NewUserRepository(dbPool)
	tokenRepo := postgres.NewRefreshTokenRepository(dbPool)
	otpRepo := NewInMemoryOTPRepository()

	passwordHasher := external.NewBcryptPasswordHasher(12)
	tokenService := external.NewJWTTokenService(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenTTL,
	)
	otpGenerator := external.NewSecureOTPGenerator(6)

	var smsService ports.SMSService
	switch cfg.SMS.Provider {
	case "twilio":
		smsService = external.NewTwilioSMSService(
			cfg.SMS.AccountSID,
			cfg.SMS.AuthToken,
			cfg.SMS.FromPhone,
		)
	default:
		smsService = external.NewConsoleSMSService()
	}

	// Initialize event publisher (Kafka or Noop)
	var eventPublisher ports.EventPublisher
	var kafkaPublisher *kafka.Publisher
	if cfg.Kafka.Enabled {
		kafkaPublisher = kafka.NewPublisher(kafka.DefaultPublisherConfig(cfg.Kafka.Brokers, cfg.Kafka.Topic))
		eventPublisher = &kafkaEventAdapter{publisher: kafkaPublisher}
		log.Println("Kafka event publisher initialized")
	} else {
		eventPublisher = NewNoOpEventPublisher()
	}

	logger := NewSimpleLogger()

	// Create application service
	authService := application.NewAuthService(
		userRepo,
		tokenRepo,
		otpRepo,
		passwordHasher,
		tokenService,
		smsService,
		otpGenerator,
		eventPublisher,
		logger,
	)

	// Create HTTP router with tracing middleware
	router := httpAdapter.NewRouter(authService, tokenService)
	if cfg.OTEL.Enabled {
		router.Use(middleware.Tracing(cfg.OTEL.ServiceName))
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Create gRPC server
	grpcServer := interceptors.NewServerWithDefaults()
	// Register gRPC services when proto is generated
	// authv1.RegisterAuthServiceServer(grpcServer, authGRPCServer)

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Auth gRPC server listening on port %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Printf("Auth HTTP server listening on port %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	// Close Kafka publisher
	if kafkaPublisher != nil {
		if err := kafkaPublisher.Close(); err != nil {
			log.Printf("failed to close Kafka publisher: %v", err)
		}
	}

	// Shutdown tracer
	if tracerShutdown != nil {
		if err := tracerShutdown(shutdownCtx); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}

	log.Println("Server exited")
}

// ================================================
// TEMPORARY IMPLEMENTATIONS
// ================================================

// InMemoryOTPRepository is a simple in-memory OTP store.
type InMemoryOTPRepository struct {
	otps map[string]*domain.OTP
	mu   sync.RWMutex
}

func NewInMemoryOTPRepository() *InMemoryOTPRepository {
	return &InMemoryOTPRepository{
		otps: make(map[string]*domain.OTP),
	}
}

func (r *InMemoryOTPRepository) Create(ctx context.Context, otp *domain.OTP) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.otps[otp.Phone] = otp
	return nil
}

func (r *InMemoryOTPRepository) GetLatestByPhone(ctx context.Context, phone string) (*domain.OTP, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	otp, ok := r.otps[phone]
	if !ok {
		return nil, domain.ErrTokenNotFound
	}
	return otp, nil
}

func (r *InMemoryOTPRepository) Update(ctx context.Context, otp *domain.OTP) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.otps[otp.Phone] = otp
	return nil
}

func (r *InMemoryOTPRepository) DeleteByPhone(ctx context.Context, phone string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.otps, phone)
	return nil
}

func (r *InMemoryOTPRepository) DeleteExpired(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for k, v := range r.otps {
		if now.After(v.ExpiresAt) {
			delete(r.otps, k)
		}
	}
	return nil
}

// NoOpEventPublisher is a no-op event publisher for development.
type NoOpEventPublisher struct{}

func NewNoOpEventPublisher() *NoOpEventPublisher {
	return &NoOpEventPublisher{}
}

func (p *NoOpEventPublisher) Publish(ctx context.Context, event ports.Event) error {
	log.Printf("[EVENT] %s: %v", event.Type, event.Payload)
	return nil
}

// SimpleLogger is a simple logger for development.
type SimpleLogger struct{}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

func (l *SimpleLogger) Debug(msg string, fields ...ports.Field) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (l *SimpleLogger) Info(msg string, fields ...ports.Field) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *SimpleLogger) Warn(msg string, fields ...ports.Field) {
	log.Printf("[WARN] %s %v", msg, fields)
}

func (l *SimpleLogger) Error(msg string, fields ...ports.Field) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *SimpleLogger) WithFields(fields ...ports.Field) ports.Logger {
	return l
}

// kafkaEventAdapter adapts kafka.Publisher to ports.EventPublisher
type kafkaEventAdapter struct {
	publisher *kafka.Publisher
}

func (a *kafkaEventAdapter) Publish(ctx context.Context, event ports.Event) error {
	return a.publisher.Publish(ctx, kafka.Event{
		Type:    event.Type,
		Payload: event.Payload,
	})
}
