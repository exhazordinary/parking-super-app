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
	"github.com/parking-super-app/services/wallet/config"
	"github.com/parking-super-app/services/wallet/internal/adapters/external"
	grpcAdapter "github.com/parking-super-app/services/wallet/internal/adapters/grpc"
	httpAdapter "github.com/parking-super-app/services/wallet/internal/adapters/http"
	"github.com/parking-super-app/services/wallet/internal/adapters/repository/postgres"
	"github.com/parking-super-app/services/wallet/internal/application"
	"github.com/parking-super-app/services/wallet/internal/ports"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger
	logger := external.NewStdLogger()
	logger.Info("starting wallet service")

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

	// Verify database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	logger.Info("connected to database")

	// Initialize repositories (adapters)
	walletRepo := postgres.NewWalletRepository(pool)
	txRepo := postgres.NewTransactionRepository(pool)

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

	// Initialize external services
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

	// Initialize HTTP router with tracing middleware
	router := httpAdapter.NewRouter(walletService)
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

	// Create gRPC server
	grpcServer := interceptors.NewServerWithDefaults()
	walletGRPCServer := grpcAdapter.NewWalletServiceServer(walletService)
	_ = walletGRPCServer // Register when proto is generated
	// walletv1.RegisterWalletServiceServer(grpcServer, walletGRPCServer)

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Wallet gRPC server listening on port %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Printf("Wallet HTTP server listening on port %s", cfg.Server.Port)
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
