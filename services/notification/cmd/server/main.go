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
	"github.com/parking-super-app/services/notification/config"
	"github.com/parking-super-app/services/notification/internal/adapters/external"
	httpAdapter "github.com/parking-super-app/services/notification/internal/adapters/http"
	"github.com/parking-super-app/services/notification/internal/adapters/repository/postgres"
	"github.com/parking-super-app/services/notification/internal/application"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := external.NewStdLogger()
	logger.Info("starting notification service")

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

	// Connect to database
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
	notificationRepo := postgres.NewNotificationRepository(pool)
	preferenceRepo := postgres.NewPreferenceRepository(pool)

	// Initialize providers
	pushProvider := external.NewMockPushProvider()
	smsProvider := external.NewMockSMSProvider()
	emailProvider := external.NewMockEmailProvider()

	// Initialize application service
	notificationService := application.NewNotificationService(
		notificationRepo,
		nil, // template repo
		preferenceRepo,
		pushProvider,
		smsProvider,
		emailProvider,
		logger,
	)

	// Initialize Kafka consumer for event-driven notifications
	var kafkaConsumer *kafka.Consumer
	if cfg.Kafka.Enabled && len(cfg.Kafka.Topics) > 0 {
		// Create consumer for first topic (would need multiple consumers for multiple topics)
		kafkaConsumer = kafka.NewConsumer(kafka.DefaultConsumerConfig(
			cfg.Kafka.Brokers,
			cfg.Kafka.Topics[0],
			cfg.Kafka.ConsumerGroup,
		))

		// Register event handlers
		kafkaConsumer.RegisterHandler("parking.session.started", func(ctx context.Context, event kafka.Event) error {
			logger.Info("received parking session started event")
			// Handle event - send notification to user
			return nil
		})

		kafkaConsumer.RegisterHandler("parking.session.ended", func(ctx context.Context, event kafka.Event) error {
			logger.Info("received parking session ended event")
			// Handle event - send notification to user
			return nil
		})

		kafkaConsumer.RegisterHandler("wallet.payment.completed", func(ctx context.Context, event kafka.Event) error {
			logger.Info("received payment completed event")
			// Handle event - send notification to user
			return nil
		})

		// Start consumer in background
		go func() {
			logger.Info("starting Kafka consumer")
			if err := kafkaConsumer.Start(ctx); err != nil {
				log.Printf("Kafka consumer error: %v", err)
			}
		}()
	}

	// Initialize HTTP router with tracing middleware
	router := httpAdapter.NewRouter(notificationService)
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
	// Register gRPC services when proto is generated
	// notificationv1.RegisterNotificationServiceServer(grpcServer, notificationGRPCServer)

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Notification gRPC server listening on port %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Printf("Notification HTTP server listening on port %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
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

	// Close Kafka consumer
	if kafkaConsumer != nil {
		if err := kafkaConsumer.Close(); err != nil {
			log.Printf("failed to close Kafka consumer: %v", err)
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
