package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/alexs/golang_test/internal/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from the frontend
		// In production, you should check the origin more strictly
		return true
	},
}

// HandleWebSocket handles WebSocket upgrade requests
func HandleWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT token from query parameter
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		// Validate JWT token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			log.Printf("WebSocket authentication failed: %v", err)
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		// Create a new client
		clientID := uuid.New().String()
		client := NewClient(clientID, hub, conn, claims.UserID)

		// Register the client with the hub
		hub.register <- client

		// Start the client's read and write pumps in separate goroutines
		go client.writePump()
		go client.readPump()

		log.Printf("WebSocket connection established for user %d (client %s)", claims.UserID, clientID)
	}
}
