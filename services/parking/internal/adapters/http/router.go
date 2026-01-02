package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/parking-super-app/services/parking/internal/application"
)

type Router struct {
	parkingService *application.ParkingService
	router         chi.Router
}

func NewRouter(parkingService *application.ParkingService) *Router {
	r := &Router{
		parkingService: parkingService,
		router:         chi.NewRouter(),
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
	handler := NewParkingHandler(r.parkingService)

	r.router.Route("/api/v1/parking", func(router chi.Router) {
		router.Post("/sessions", handler.StartSession)
		router.Get("/sessions", handler.GetUserSessions)
		router.Get("/sessions/active", handler.GetActiveSessions)
		router.Get("/sessions/{id}", handler.GetSession)
		router.Post("/sessions/{id}/end", handler.EndSession)
		router.Delete("/sessions/{id}", handler.CancelSession)

		router.Post("/vehicles", handler.RegisterVehicle)
		router.Get("/vehicles", handler.GetUserVehicles)
	})

	r.router.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
