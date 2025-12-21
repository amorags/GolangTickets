package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexs/golang_test/internal/config"
	"github.com/alexs/golang_test/internal/handlers"
	"github.com/alexs/golang_test/internal/repository"
	"github.com/alexs/golang_test/internal/router"
	"github.com/alexs/golang_test/internal/websocket"
)

func main() {
	// Health check mode for Docker
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		client := http.Client{
			Timeout: 2 * time.Second,
		}
		resp, err := client.Get("http://localhost:8080/health")
		if err != nil {
			os.Exit(1)
		}
		if resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	// 1. Load Configuration
	config.LoadConfig()

	// 2. Connect to Database
	repository.ConnectDB()

	// 3. Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// 4. Set WebSocket hub for handlers
	handlers.SetWebSocketHub(hub)

	// 5. Setup Router with WebSocket hub
	r := router.New(hub)

	// 6. Start Server
	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}