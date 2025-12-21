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

	// Serve static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Serve HTML pages
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})
	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/home.html")
	})
	r.Get("/event", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/event.html")
	})
	r.Get("/profile-page", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/profile.html")
	})

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
