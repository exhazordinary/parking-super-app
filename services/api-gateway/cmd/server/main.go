package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/parking-super-app/pkg/middleware"
	"github.com/parking-super-app/pkg/telemetry"
	"github.com/parking-super-app/services/api-gateway/config"
	"github.com/parking-super-app/services/api-gateway/internal/health"
	gatewaymw "github.com/parking-super-app/services/api-gateway/internal/middleware"
	"github.com/parking-super-app/services/api-gateway/internal/proxy"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("Starting API Gateway...")

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

	// Initialize components
	authMw := gatewaymw.NewAuthMiddleware(cfg.Auth.JWTSecret)
	rateLimiter := gatewaymw.NewRateLimiter(100, time.Minute)
	serviceProxy := proxy.NewServiceProxy()

	// Initialize health checker
	healthChecker := health.NewServiceHealth(map[string]string{
		"auth":         cfg.Services.AuthURL,
		"wallet":       cfg.Services.WalletURL,
		"provider":     cfg.Services.ProviderURL,
		"parking":      cfg.Services.ParkingURL,
		"notification": cfg.Services.NotificationURL,
	})

	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(gatewaymw.CORS)
	r.Use(rateLimiter.Limit)

	// Add tracing middleware
	if cfg.OTEL.Enabled {
		r.Use(middleware.Tracing(cfg.OTEL.ServiceName))
	}

	// Health endpoint
	r.Get("/health", healthChecker.Handler())

	// Auth routes (public)
	r.Route("/api/v1/auth", func(router chi.Router) {
		router.HandleFunc("/*", serviceProxy.Forward(cfg.Services.AuthURL))
	})

	// Protected routes
	r.Group(func(router chi.Router) {
		router.Use(authMw.Authenticate)

		// Wallet routes
		router.Route("/api/v1/wallet", func(r chi.Router) {
			r.HandleFunc("/*", serviceProxy.Forward(cfg.Services.WalletURL))
		})

		// Parking routes
		router.Route("/api/v1/parking", func(r chi.Router) {
			r.HandleFunc("/*", serviceProxy.Forward(cfg.Services.ParkingURL))
		})

		// Notification routes
		router.Route("/api/v1/notifications", func(r chi.Router) {
			r.HandleFunc("/*", serviceProxy.Forward(cfg.Services.NotificationURL))
		})

		router.Route("/api/v1/preferences", func(r chi.Router) {
			r.HandleFunc("/*", serviceProxy.Forward(cfg.Services.NotificationURL))
		})
	})

	// Provider routes (partially public)
	r.Route("/api/v1/providers", func(router chi.Router) {
		// Public: list providers
		router.With(authMw.OptionalAuth).Get("/", serviceProxy.Forward(cfg.Services.ProviderURL))
		router.With(authMw.OptionalAuth).Get("/{id}", serviceProxy.Forward(cfg.Services.ProviderURL))
		router.With(authMw.OptionalAuth).Get("/code/{code}", serviceProxy.Forward(cfg.Services.ProviderURL))

		// Protected: admin operations
		router.Group(func(r chi.Router) {
			r.Use(authMw.Authenticate)
			r.Post("/", serviceProxy.Forward(cfg.Services.ProviderURL))
			r.Post("/{id}/*", serviceProxy.Forward(cfg.Services.ProviderURL))
		})
	})

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("API Gateway listening on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	// Shutdown tracer
	if tracerShutdown != nil {
		if err := tracerShutdown(shutdownCtx); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}

	log.Println("API Gateway stopped")
}
