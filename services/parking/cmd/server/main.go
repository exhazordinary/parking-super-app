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
	"github.com/parking-super-app/services/parking/config"
	"github.com/parking-super-app/services/parking/internal/adapters/external"
	httpAdapter "github.com/parking-super-app/services/parking/internal/adapters/http"
	"github.com/parking-super-app/services/parking/internal/adapters/repository/postgres"
	"github.com/parking-super-app/services/parking/internal/application"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := external.NewStdLogger()
	logger.Info("starting parking service")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	logger.Info("connected to database")

	// Initialize repositories
	sessionRepo := postgres.NewSessionRepository(pool)
	vehicleRepo := postgres.NewVehicleRepository(pool)

	// Initialize external clients
	providerClient := external.NewMockProviderClient()
	walletClient := external.NewMockWalletClient()
	eventPublisher := external.NewNoopEventPublisher()

	// Initialize application service
	parkingService := application.NewParkingService(
		sessionRepo,
		vehicleRepo,
		providerClient,
		walletClient,
		eventPublisher,
		logger,
	)

	// Initialize HTTP router
	router := httpAdapter.NewRouter(parkingService)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Parking service listening on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

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
