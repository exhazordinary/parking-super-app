package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/parking-super-app/services/provider/internal/application"
)

type Router struct {
	providerService *application.ProviderService
	router          chi.Router
}

func NewRouter(providerService *application.ProviderService) *Router {
	r := &Router{
		providerService: providerService,
		router:          chi.NewRouter(),
	}

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

func (r *Router) setupMiddleware() {
	r.router.Use(middleware.RequestID)
	r.router.Use(middleware.RealIP)
	r.router.Use(middleware.Logger)
	r.router.Use(middleware.Recoverer)
	r.router.Use(middleware.AllowContentType("application/json"))

	r.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	})
}

func (r *Router) setupRoutes() {
	handler := NewProviderHandler(r.providerService)

	r.router.Route("/api/v1/providers", func(router chi.Router) {
		router.Post("/", handler.RegisterProvider)
		router.Get("/", handler.ListProviders)
		router.Get("/code/{code}", handler.GetProviderByCode)
		router.Get("/{id}", handler.GetProvider)
		router.Post("/{id}/activate", handler.ActivateProvider)
		router.Post("/{id}/deactivate", handler.DeactivateProvider)
		router.Post("/{id}/credentials", handler.GenerateCredentials)
		router.Post("/{id}/locations", handler.AddLocation)
		router.Get("/{id}/locations", handler.GetProviderLocations)
	})

	r.router.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
