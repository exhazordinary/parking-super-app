package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/wallet/config"
	"github.com/parking-super-app/services/wallet/internal/adapters/external"
	httpAdapter "github.com/parking-super-app/services/wallet/internal/adapters/http"
	"github.com/parking-super-app/services/wallet/internal/adapters/repository/postgres"
	"github.com/parking-super-app/services/wallet/internal/application"
)

func main() {
	// Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger
	logger := external.NewStdLogger()
	logger.Info("starting wallet service", )

	// Connect to PostgreSQL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	logger.Info("connected to database")

	// Initialize repositories (adapters)
	walletRepo := postgres.NewWalletRepository(pool)
	txRepo := postgres.NewTransactionRepository(pool)

	// Initialize external services
	eventPublisher := external.NewNoopEventPublisher()
	paymentGateway := external.NewMockPaymentGateway()

	// Initialize application service (use cases)
	walletService := application.NewWalletService(
		walletRepo,
		txRepo,
		nil, // Unit of Work - not implemented yet
		paymentGateway,
		eventPublisher,
		logger,
	)

	// Initialize HTTP router
	router := httpAdapter.NewRouter(walletService)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("wallet service listening", external.NewStdLogger().WithFields())
		log.Printf("Wallet service listening on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	logger.Info("server stopped gracefully")
}
