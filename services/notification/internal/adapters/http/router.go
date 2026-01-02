package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/parking-super-app/services/notification/internal/application"
)

type Router struct {
	service *application.NotificationService
	router  chi.Router
}

func NewRouter(service *application.NotificationService) *Router {
	r := &Router{
		service: service,
		router:  chi.NewRouter(),
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
	handler := NewNotificationHandler(r.service)

	r.router.Route("/api/v1/notifications", func(router chi.Router) {
		router.Post("/", handler.SendNotification)
		router.Post("/template", handler.SendFromTemplate)
		router.Get("/", handler.GetUserNotifications)
		router.Get("/{id}", handler.GetNotification)
	})

	r.router.Route("/api/v1/preferences", func(router chi.Router) {
		router.Get("/", handler.GetPreferences)
		router.Put("/", handler.UpdatePreferences)
	})

	r.router.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
