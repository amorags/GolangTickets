package router

import (
	"net/http"

	"github.com/alexs/golang_test/internal/handlers"
	"github.com/alexs/golang_test/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// New returns the configured router with all routes and middleware
func New() http.Handler {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Public Routes (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticket API is running!"))
	})

	// Auth routes
	r.Post("/auth/signup", handlers.Signup)
	r.Post("/auth/login", handlers.Login)

	// Protected Routes (auth required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		// Example protected endpoint
		r.Get("/profile", handlers.GetProfile)

		// Future: Add ticket/event routes here
		// r.Get("/events", handlers.GetEvents)
		// r.Post("/tickets", handlers.CreateTicket)
	})

	return r
}
