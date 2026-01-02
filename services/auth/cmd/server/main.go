// Package main is the entry point for the auth service.
//
// MICROSERVICES PATTERN: Service Entry Point
// ==========================================
// The main package is responsible for:
// 1. Loading configuration
// 2. Creating dependencies (dependency injection)
// 3. Wiring everything together
// 4. Starting the HTTP server
// 5. Graceful shutdown handling
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/auth/config"
	"github.com/parking-super-app/services/auth/internal/adapters/external"
	httpAdapter "github.com/parking-super-app/services/auth/internal/adapters/http"
	"github.com/parking-super-app/services/auth/internal/adapters/repository/postgres"
	"github.com/parking-super-app/services/auth/internal/application"
	"github.com/parking-super-app/services/auth/internal/domain"
	"github.com/parking-super-app/services/auth/internal/ports"
)

func main() {
	// ================================================
	// 1. LOAD CONFIGURATION
	// ================================================
	// Configuration is loaded from environment variables
	// See config/config.go for available options
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting auth service on port %s", cfg.Server.Port)

	// ================================================
	// 2. SET UP DATABASE CONNECTION
	// ================================================
	// We use pgxpool for connection pooling
	// This manages a pool of connections efficiently
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Verify database connection
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// ================================================
	// 3. CREATE DEPENDENCIES (Dependency Injection)
	// ================================================
	// This is where we wire everything together.
	// Each component depends on interfaces, and we provide
	// concrete implementations here.
	//
	// PATTERN: Composition Root
	// All dependency injection happens in one place (main)
	// This makes it easy to see the application structure
	// and swap implementations for testing.

	// Repositories (database access)
	userRepo := postgres.NewUserRepository(dbPool)
	tokenRepo := postgres.NewRefreshTokenRepository(dbPool)
	otpRepo := NewInMemoryOTPRepository() // Simple in-memory for now

	// External services
	passwordHasher := external.NewBcryptPasswordHasher(12) // cost = 12
	tokenService := external.NewJWTTokenService(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenTTL,
	)
	otpGenerator := external.NewSecureOTPGenerator(6)

	// SMS service (based on configuration)
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

	// Event publisher (using a no-op for now)
	eventPublisher := NewNoOpEventPublisher()

	// Logger
	logger := NewSimpleLogger()

	// ================================================
	// 4. CREATE APPLICATION SERVICE
	// ================================================
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

	// ================================================
	// 5. CREATE HTTP ROUTER
	// ================================================
	router := httpAdapter.NewRouter(authService, tokenService)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// ================================================
	// 6. START SERVER WITH GRACEFUL SHUTDOWN
	// ================================================
	// PATTERN: Graceful Shutdown
	// When receiving SIGTERM (Kubernetes) or SIGINT (Ctrl+C):
	// 1. Stop accepting new connections
	// 2. Wait for existing requests to complete
	// 3. Close database connections
	// 4. Exit cleanly

	// Channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Auth service listening on :%s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// ================================================
// TEMPORARY IMPLEMENTATIONS
// ================================================
// These are simple implementations for development.
// Replace with proper implementations for production.

// InMemoryOTPRepository is a simple in-memory OTP store.
// TODO: Replace with Redis or PostgreSQL in production.
type InMemoryOTPRepository struct {
	otps map[string]*domain.OTP
	mu   sync.RWMutex
}

func NewInMemoryOTPRepository() *InMemoryOTPRepository {
	return &InMemoryOTPRepository{
		otps: make(map[string]*domain.OTP),
	}
}

// Implement ports.OTPRepository interface

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
