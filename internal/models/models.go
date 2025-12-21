package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string    `json:"username" gorm:"uniqueIndex;not null"`
	Email    string    `json:"email" gorm:"uniqueIndex;not null"`
	Password string    `json:"-" gorm:"not null"`
	Bookings []Booking `json:"bookings,omitempty" gorm:"foreignKey:UserID"`
}

// Event represents a concert, tour, standup show, lecture, musical, etc.
type Event struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	EventType   string    `json:"event_type"` // concert, tour, standup, lecture, musical, etc.
	VenueName   string    `json:"venue_name"`
	City        string    `json:"city"`
	Address     string    `json:"address"`
	Date        time.Time `json:"date" gorm:"not null"`
	Price       float64   `json:"price" gorm:"default:0"` // Price per ticket, 0 for free events
	Capacity    int       `json:"capacity" gorm:"not null"`
	ImageURL    string    `json:"image_url"`
	Bookings    []Booking `json:"bookings,omitempty" gorm:"foreignKey:EventID"`
}

// AvailableTickets calculates remaining tickets
func (e *Event) AvailableTickets(db *gorm.DB) (int, error) {
	var totalBooked int64
	err := db.Model(&Booking{}).
		Where("event_id = ? AND status = ?", e.ID, "confirmed").
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalBooked).Error

	if err != nil {
		return 0, err
	}

	return e.Capacity - int(totalBooked), nil
}

// Booking represents a user's ticket booking for an event
type Booking struct {
	gorm.Model
	UserID      uint    `json:"user_id" gorm:"not null"`
	EventID     uint    `json:"event_id" gorm:"not null"`
	Quantity    int     `json:"quantity" gorm:"not null"`
	PricePerTicket float64 `json:"price_per_ticket"`
	TotalPrice  float64 `json:"total_price"`
	Status      string  `json:"status" gorm:"default:'confirmed'"` // confirmed, cancelled
	User        User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Event       Event   `json:"event,omitempty" gorm:"foreignKey:EventID"`
}