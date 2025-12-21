package websocket

import "time"

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeAvailabilityUpdate MessageType = "availability_update"
	MessageTypeConnectionAck      MessageType = "connection_ack"
	MessageTypeError              MessageType = "error"
	MessageTypeSubscribe          MessageType = "subscribe"
	MessageTypeUnsubscribe        MessageType = "unsubscribe"
	MessageTypePing               MessageType = "ping"
	MessageTypePong               MessageType = "pong"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType `json:"type"`
	EventID   *uint       `json:"event_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

// AvailabilityUpdate represents ticket availability data
type AvailabilityUpdate struct {
	EventID          uint   `json:"event_id"`
	AvailableTickets int    `json:"available_tickets"`
	Capacity         int    `json:"capacity"`
	LastUpdated      string `json:"last_updated"`
}

// ConnectionAck represents a connection acknowledgment
type ConnectionAck struct {
	ClientID string `json:"client_id"`
	Message  string `json:"message"`
}

// ErrorMessage represents an error message
type ErrorMessage struct {
	Message string `json:"message"`
}

// ClientMessage represents a message from client to server
type ClientMessage struct {
	Type    MessageType `json:"type"`
	EventID *uint       `json:"event_id,omitempty"`
}
