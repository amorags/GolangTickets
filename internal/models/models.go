package models

import (
	"time"

	"gorm.io/gorm"
)

// Event represents a concert, game, or show.
type Event struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Date        time.Time `json:"date"`
	Tickets     []Ticket  `json:"tickets"` // One-to-Many relationship
}

// Ticket represents a single seat for an event.
type Ticket struct {
	gorm.Model
	EventID uint    `json:"event_id"`
	Price   float64 `json:"price"`
	Seat    string  `json:"seat"`
	IsSold  bool    `json:"is_sold" gorm:"default:false"`
}

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
}