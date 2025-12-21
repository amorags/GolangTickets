package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client represents a WebSocket client connection
type Client struct {
	// Unique identifier for this client
	ID string

	// Reference to the hub
	hub *Hub

	// The WebSocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Map of event IDs this client is subscribed to
	eventIDs map[uint]bool

	// User ID from JWT authentication
	userID uint

	// Mutex for thread-safe access
	mutex sync.RWMutex
}

// NewClient creates a new Client instance
func NewClient(id string, hub *Hub, conn *websocket.Conn, userID uint) *Client {
	return &Client{
		ID:       id,
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		eventIDs: make(map[uint]bool),
		userID:   userID,
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for client %s: %v", c.ID, err)
			}
			break
		}

		// Parse the message
		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			log.Printf("Error parsing message from client %s: %v", c.ID, err)
			c.sendError("Invalid message format")
			continue
		}

		// Handle the message based on type
		c.handleMessage(&clientMsg)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages from the client
func (c *Client) handleMessage(msg *ClientMessage) {
	switch msg.Type {
	case MessageTypeSubscribe:
		if msg.EventID != nil {
			c.hub.SubscribeToEvent(c, *msg.EventID)
		} else {
			c.sendError("Event ID required for subscription")
		}

	case MessageTypeUnsubscribe:
		if msg.EventID != nil {
			c.hub.UnsubscribeFromEvent(c, *msg.EventID)
		} else {
			c.sendError("Event ID required for unsubscription")
		}

	case MessageTypePing:
		// Respond with pong
		pongMsg := &Message{
			Type:      MessageTypePong,
			Timestamp: time.Now(),
		}
		data, err := json.Marshal(pongMsg)
		if err == nil {
			select {
			case c.send <- data:
			default:
				log.Printf("Failed to send pong to client %s", c.ID)
			}
		}

	default:
		log.Printf("Unknown message type from client %s: %s", c.ID, msg.Type)
		c.sendError("Unknown message type")
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(errorMsg string) {
	msg := &Message{
		Type:      MessageTypeError,
		Timestamp: time.Now(),
		Data: ErrorMessage{
			Message: errorMsg,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling error message: %v", err)
		return
	}

	select {
	case c.send <- data:
	default:
		log.Printf("Failed to send error message to client %s (buffer full)", c.ID)
	}
}
