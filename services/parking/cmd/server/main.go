package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/pkg/grpc/interceptors"
	"github.com/parking-super-app/pkg/kafka"
	"github.com/parking-super-app/pkg/middleware"
	"github.com/parking-super-app/pkg/telemetry"
	"github.com/parking-super-app/services/parking/config"
	"github.com/parking-super-app/services/parking/internal/adapters/external"
	grpcClients "github.com/parking-super-app/services/parking/internal/adapters/grpc"
	httpAdapter "github.com/parking-super-app/services/parking/internal/adapters/http"
	"github.com/parking-super-app/services/parking/internal/adapters/repository/postgres"
	"github.com/parking-super-app/services/parking/internal/application"
	"github.com/parking-super-app/services/parking/internal/ports"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := external.NewStdLogger()
	logger.Info("starting parking service")

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
			logger.Info("OpenTelemetry tracing initialized")
		}
	}

	// Connect to PostgreSQL
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

	// Initialize gRPC clients for dependent services or fallback to mock
	var providerClient ports.ProviderClient
	var walletClient ports.WalletClient
	var providerGRPCClient *grpcClients.ProviderGRPCClient
	var walletGRPCClient *grpcClients.WalletGRPCClient

	if cfg.Services.ProviderGRPC != "" && cfg.Services.WalletGRPC != "" {
		// Try to connect via gRPC
		providerGRPCClient, err = grpcClients.NewProviderGRPCClient(cfg.Services.ProviderGRPC)
		if err != nil {
			log.Printf("warning: failed to connect to provider service, using mock: %v", err)
			providerClient = external.NewMockProviderClient()
		} else {
			providerClient = providerGRPCClient
			logger.Info("connected to provider service via gRPC")
		}

		walletGRPCClient, err = grpcClients.NewWalletGRPCClient(cfg.Services.WalletGRPC)
		if err != nil {
			log.Printf("warning: failed to connect to wallet service, using mock: %v", err)
			walletClient = external.NewMockWalletClient()
		} else {
			walletClient = walletGRPCClient
			logger.Info("connected to wallet service via gRPC")
		}
	} else {
		// Use mock clients for development
		providerClient = external.NewMockProviderClient()
		walletClient = external.NewMockWalletClient()
		logger.Info("using mock clients for provider and wallet services")
	}

	// Initialize event publisher (Kafka or Noop)
	var eventPublisher ports.EventPublisher
	var kafkaPublisher *kafka.Publisher
	if cfg.Kafka.Enabled {
		kafkaPublisher = kafka.NewPublisher(kafka.DefaultPublisherConfig(cfg.Kafka.Brokers, cfg.Kafka.Topic))
		eventPublisher = &kafkaEventAdapter{publisher: kafkaPublisher}
		logger.Info("Kafka event publisher initialized")
	} else {
		eventPublisher = external.NewNoopEventPublisher()
	}

	// Initialize application service
	parkingService := application.NewParkingService(
		sessionRepo,
		vehicleRepo,
		providerClient,
		walletClient,
		eventPublisher,
		logger,
	)

	// Initialize HTTP router with tracing middleware
	router := httpAdapter.NewRouter(parkingService)
	if cfg.OTEL.Enabled {
		router.Use(middleware.Tracing(cfg.OTEL.ServiceName))
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Create gRPC server (for future use when parking exposes gRPC)
	grpcServer := interceptors.NewServerWithDefaults()

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Parking gRPC server listening on port %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Printf("Parking HTTP server listening on port %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down servers")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	// Close gRPC clients
	if providerGRPCClient != nil {
		providerGRPCClient.Close()
	}
	if walletGRPCClient != nil {
		walletGRPCClient.Close()
	}

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

	logger.Info("server stopped gracefully")
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
