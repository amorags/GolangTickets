package main

import (
	"fmt"
	"net/http"

	"github.com/alexs/golang_test/internal/config"
	"github.com/alexs/golang_test/internal/repository"
	"github.com/alexs/golang_test/internal/router"
)

func main() {
	// 1. Load Configuration
	config.LoadConfig()

	// 2. Connect to Database
	repository.ConnectDB()

	// 3. Setup Router
	r := router.New()

	// 4. Start Server
	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}