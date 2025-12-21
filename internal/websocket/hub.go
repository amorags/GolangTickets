package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// eventClients maps eventID to a map of clientID to Client
	// This allows efficient broadcasting to event-specific subscribers
	eventClients map[uint]map[string]*Client

	// clients is a map of all connected clients
	clients map[*Client]bool

	// broadcast channel for sending messages to clients
	broadcast chan *Message

	// register channel for registering new clients
	register chan *Client

	// unregister channel for unregistering clients
	unregister chan *Client

	// mutex for thread-safe access to maps
	mutex sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		eventClients: make(map[uint]map[string]*Client),
		clients:      make(map[*Client]bool),
		broadcast:    make(chan *Message, 256),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	log.Println("WebSocket Hub started")
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true
	log.Printf("Client %s connected. Total clients: %d", client.ID, len(h.clients))

	// Send connection acknowledgment
	ack := &Message{
		Type:      MessageTypeConnectionAck,
		Timestamp: time.Now(),
		Data: ConnectionAck{
			ClientID: client.ID,
			Message:  "Connected successfully",
		},
	}

	data, err := json.Marshal(ack)
	if err == nil {
		select {
		case client.send <- data:
		default:
			log.Printf("Failed to send connection ack to client %s", client.ID)
		}
	}
}

// unregisterClient removes a client from the hub
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		// Remove from global clients map
		delete(h.clients, client)

		// Remove from all event subscriptions
		for eventID := range client.eventIDs {
			if eventClients, exists := h.eventClients[eventID]; exists {
				delete(eventClients, client.ID)

				// Clean up empty event maps
				if len(eventClients) == 0 {
					delete(h.eventClients, eventID)
				}
			}
		}

		// Close the client's send channel
		close(client.send)

		log.Printf("Client %s disconnected. Total clients: %d", client.ID, len(h.clients))
	}
}

// broadcastMessage sends a message to all relevant subscribers
func (h *Hub) broadcastMessage(message *Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	// If message has an event ID, broadcast only to subscribers of that event
	if message.EventID != nil {
		eventID := *message.EventID
		if eventClients, exists := h.eventClients[eventID]; exists {
			for _, client := range eventClients {
				select {
				case client.send <- data:
				default:
					log.Printf("Failed to send message to client %s (buffer full)", client.ID)
				}
			}
			log.Printf("Broadcast message to %d clients subscribed to event %d", len(eventClients), eventID)
		}
	} else {
		// Broadcast to all clients
		for client := range h.clients {
			select {
			case client.send <- data:
			default:
				log.Printf("Failed to send message to client %s (buffer full)", client.ID)
			}
		}
		log.Printf("Broadcast message to all %d clients", len(h.clients))
	}
}

// SubscribeToEvent adds a client to the subscribers list for a specific event
func (h *Hub) SubscribeToEvent(client *Client, eventID uint) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Initialize the event's client map if it doesn't exist
	if _, exists := h.eventClients[eventID]; !exists {
		h.eventClients[eventID] = make(map[string]*Client)
	}

	// Add client to event subscribers
	h.eventClients[eventID][client.ID] = client

	// Add event to client's subscriptions
	client.mutex.Lock()
	client.eventIDs[eventID] = true
	client.mutex.Unlock()

	log.Printf("Client %s subscribed to event %d", client.ID, eventID)
}

// UnsubscribeFromEvent removes a client from the subscribers list for a specific event
func (h *Hub) UnsubscribeFromEvent(client *Client, eventID uint) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if eventClients, exists := h.eventClients[eventID]; exists {
		delete(eventClients, client.ID)

		// Clean up empty event maps
		if len(eventClients) == 0 {
			delete(h.eventClients, eventID)
		}
	}

	// Remove event from client's subscriptions
	client.mutex.Lock()
	delete(client.eventIDs, eventID)
	client.mutex.Unlock()

	log.Printf("Client %s unsubscribed from event %d", client.ID, eventID)
}

// BroadcastAvailabilityUpdate broadcasts a ticket availability update to all subscribers
func (h *Hub) BroadcastAvailabilityUpdate(eventID uint, availableTickets int, capacity int) {
	message := &Message{
		Type:      MessageTypeAvailabilityUpdate,
		EventID:   &eventID,
		Timestamp: time.Now(),
		Data: AvailabilityUpdate{
			EventID:          eventID,
			AvailableTickets: availableTickets,
			Capacity:         capacity,
			LastUpdated:      time.Now().Format(time.RFC3339),
		},
	}

	h.broadcast <- message
}
