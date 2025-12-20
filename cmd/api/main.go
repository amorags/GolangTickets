package main

import (
	"fmt"
	"net/http"

	"github.com/alexs/golang_test/internal/config"
	"github.com/alexs/golang_test/internal/handlers"
	"github.com/alexs/golang_test/internal/middleware"
	"github.com/alexs/golang_test/internal/repository"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Load Configuration
	config.LoadConfig()

	// 2. Connect to Database
	repository.ConnectDB()

	// 3. Setup Router
	r := chi.NewRouter()

	// 4. Global Middlewares
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// 5. Public Routes (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticket API is running!"))
	})

	// Auth routes
	r.Post("/auth/signup", handlers.Signup)
	r.Post("/auth/login", handlers.Login)

	// 6. Protected Routes (auth required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		// Example protected endpoint
		r.Get("/profile", handlers.GetProfile)

		// Future: Add ticket/event routes here
		// r.Get("/events", handlers.GetEvents)
		// r.Post("/tickets", handlers.CreateTicket)
	})

	// 7. Start Server
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
