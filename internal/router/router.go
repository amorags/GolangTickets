package router

import (
	"net/http"

	"github.com/alexs/golang_test/internal/handlers"
	"github.com/alexs/golang_test/internal/middleware"
	"github.com/alexs/golang_test/internal/websocket"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// New returns the configured router with all routes and middleware
func New(hub *websocket.Hub) http.Handler {
	r := chi.NewRouter()

	// Initialize handlers with the WebSocket hub
	handlers.SetWebSocketHub(hub)

	// Global Middlewares
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// CORS middleware for Nuxt frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Note: Frontend is now a separate Nuxt app running on port 3000
	// No need to serve static files from Go backend anymore

	// Public Routes (no auth required)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticket API is running!"))
	})

	// Auth routes
	r.Post("/auth/signup", handlers.Signup)
	r.Post("/auth/login", handlers.Login)

	// Public event routes (no auth required for browsing)
	r.Get("/events", handlers.GetEvents)
	r.Get("/events/{id}", handlers.GetEvent)

	// WebSocket endpoint (auth via token query parameter)
	r.Get("/ws", websocket.HandleWebSocket(hub))

	// Protected Routes (auth required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		// User profile
		r.Get("/profile", handlers.GetProfile)

		// Event management (protected - could add admin check later)
		r.Post("/events", handlers.CreateEvent)
		r.Delete("/events/{id}", handlers.DeleteEvent)

		// Booking routes
		r.Post("/bookings", handlers.BookTicket)
		r.Get("/bookings", handlers.GetMyBookings)
		r.Delete("/bookings/{id}", handlers.CancelBooking)
	})

	return r
}
