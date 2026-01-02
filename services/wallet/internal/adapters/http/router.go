package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/parking-super-app/services/wallet/internal/application"
)

type Router struct {
	walletService *application.WalletService
	router        chi.Router
}

func NewRouter(walletService *application.WalletService) *Router {
	r := &Router{
		walletService: walletService,
		router:        chi.NewRouter(),
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
	handler := NewWalletHandler(r.walletService)

	r.router.Route("/api/v1/wallet", func(router chi.Router) {
		router.Post("/", handler.CreateWallet)
		router.Get("/", handler.GetWallet)
		router.Post("/topup", handler.TopUp)
		router.Post("/pay", handler.Pay)
		router.Get("/transactions", handler.GetTransactions)
	})

	r.router.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
