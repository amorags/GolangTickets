package main

import (
	"fmt"
	"net/http"

	"github.com/alexs/golang_test/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Connect to Database
	repository.ConnectDB()

	// 2. Setup Router
	r := chi.NewRouter()

	// 3. Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 4. Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticket API is running!"))
	})

	// 5. Start Server
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
