// Package http provides HTTP handlers for the auth service.
//
// MICROSERVICES PATTERN: Primary Adapter
// =====================================
// This is a PRIMARY ADAPTER - it receives requests from the outside world
// and translates them into calls to our application layer.
//
// The handler doesn't contain business logic - it only:
// 1. Parses HTTP requests into DTOs
// 2. Calls the application layer
// 3. Translates responses back to HTTP
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/parking-super-app/services/auth/internal/application"
	"github.com/parking-super-app/services/auth/internal/ports"
)

// Router holds the HTTP router and dependencies.
type Router struct {
	authService  *application.AuthService
	tokenService ports.TokenService
	router       chi.Router
}

// NewRouter creates a new HTTP router with all routes configured.
//
// PATTERN: Chi Router
// ===================
// Chi is a lightweight, idiomatic router for Go. It's:
// - Fast (uses radix tree)
// - Compatible with net/http
// - Has great middleware support
// - Easy to test
func NewRouter(authService *application.AuthService, tokenService ports.TokenService) *Router {
	r := &Router{
		authService:  authService,
		tokenService: tokenService,
		router:       chi.NewRouter(),
	}

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

// setupMiddleware configures global middleware.
//
// MICROSERVICES PATTERN: Cross-Cutting Concerns
// =============================================
// Middleware handles concerns that apply to all requests:
// - Logging
// - Request ID generation
// - Panic recovery
// - CORS
// - Compression
// - Rate limiting (often done at API Gateway level)
func (r *Router) setupMiddleware() {
	// RequestID adds a unique ID to each request for tracing
	r.router.Use(middleware.RequestID)

	// RealIP extracts the real client IP from X-Forwarded-For
	r.router.Use(middleware.RealIP)

	// Logger logs the start and end of each request
	r.router.Use(middleware.Logger)

	// Recoverer catches panics and returns 500 instead of crashing
	r.router.Use(middleware.Recoverer)

	// Content-Type enforcement
	r.router.Use(middleware.AllowContentType("application/json"))

	// Set response content type
	r.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	})
}

// setupRoutes configures all HTTP routes.
//
// REST API DESIGN:
// ================
// - Use nouns for resources (users, tokens)
// - Use HTTP methods to indicate action (POST = create, GET = read)
// - Use proper status codes (201 Created, 401 Unauthorized, etc.)
// - Version your API (/api/v1/)
func (r *Router) setupRoutes() {
	handler := NewAuthHandler(r.authService)
	handler.SetTokenService(r.tokenService)

	r.router.Route("/api/v1/auth", func(router chi.Router) {
		// Public routes (no authentication required)
		router.Post("/register", handler.Register)
		router.Post("/login", handler.Login)
		router.Post("/refresh", handler.RefreshToken)
		router.Post("/otp/request", handler.RequestOTP)
		router.Post("/otp/verify", handler.VerifyOTP)

		// Protected routes (require valid access token)
		router.Group(func(protected chi.Router) {
			protected.Use(handler.AuthMiddleware)

			protected.Get("/me", handler.GetProfile)
			protected.Post("/logout", handler.Logout)
			protected.Post("/logout/all", handler.LogoutAllDevices)
		})
	})

	// Health check endpoint (for Kubernetes probes)
	r.router.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Ready check endpoint (for Kubernetes probes)
	r.router.Get("/ready", func(w http.ResponseWriter, req *http.Request) {
		// In production, check database connection, etc.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})
}

// ServeHTTP implements http.Handler interface.
// This allows our Router to be used with standard http.Server.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
