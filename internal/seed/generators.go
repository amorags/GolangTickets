package seed

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alexs/golang_test/internal/models"
	"github.com/alexs/golang_test/internal/utils"
	"github.com/brianvoe/gofakeit/v6"
)

// generateUser creates a new user with unique email and username
func generateUser(index int, isOrganizer bool) (*models.User, error) {
	// Use index to ensure uniqueness
	username := fmt.Sprintf("%s%d", gofakeit.Username(), index)
	email := fmt.Sprintf("user%d-%s@%s", index, gofakeit.Username(), gofakeit.DomainName())

	// Hash password - all test users have password "password123"
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	return user, nil
}

// generateEvent creates a new event with realistic data
func generateEvent(organizerID uint, eventType string) *models.Event {
	city := randomCity()
	eventName := generateEventName(eventType)
	description := generateEventDescription(eventType, eventName)
	venueName := randomVenueName(eventType)
	address := gofakeit.Street()
	date := generateEventDate()
	price := randomPrice(eventType)
	capacity := randomCapacity()
	imageURL := randomImageForType(eventType)

	event := &models.Event{
		Name:        eventName,
		Description: description,
		EventType:   eventType,
		Status:      "published",
		OrganizerID: organizerID,
		VenueName:   venueName,
		City:        city,
		Address:     address,
		Date:        date,
		Price:       price,
		Capacity:    capacity,
		ImageURL:    imageURL,
	}

	return event
}

// generateBooking creates a booking with capacity awareness
func generateBooking(userID uint, eventID uint, eventCapacity int, currentBookedCount int) *models.Booking {
	// Calculate remaining capacity
	remainingCapacity := eventCapacity - currentBookedCount

	// Don't create booking if no capacity left
	if remainingCapacity <= 0 {
		return nil
	}

	// Generate quantity (1-10 tickets, but not exceeding remaining capacity)
	maxQuantity := 10
	if remainingCapacity < maxQuantity {
		maxQuantity = remainingCapacity
	}

	quantity := 1 + rand.Intn(maxQuantity)

	// Price will be set by the repository's CreateBooking function
	// but we need to pass the event's price
	booking := &models.Booking{
		UserID:   userID,
		EventID:  eventID,
		Quantity: quantity,
		Status:   "confirmed",
	}

	return booking
}

// BookingDistribution determines which events should have bookings
type BookingDistribution struct {
	EventID       uint
	TargetPercent int // 0-90% of capacity
}

// generateBookingDistributions creates a distribution plan for bookings
func generateBookingDistributions(eventIDs []uint) []BookingDistribution {
	distributions := []BookingDistribution{}

	for _, eventID := range eventIDs {
		roll := rand.Intn(100)

		var targetPercent int
		if roll < 40 {
			// 40% of events: No bookings (skip this event)
			continue
		} else if roll < 70 {
			// 30% of events: Light bookings (10-30% capacity)
			targetPercent = 10 + rand.Intn(21)
		} else if roll < 90 {
			// 20% of events: Medium bookings (30-60% capacity)
			targetPercent = 30 + rand.Intn(31)
		} else {
			// 10% of events: Heavy bookings (60-90% capacity)
			targetPercent = 60 + rand.Intn(31)
		}

		distributions = append(distributions, BookingDistribution{
			EventID:       eventID,
			TargetPercent: targetPercent,
		})
	}

	return distributions
}

// generateSentinelUser creates the sentinel marker user for idempotency checks
func generateSentinelUser() (*models.User, error) {
	hashedPassword, err := utils.HashPassword("sentinel-not-for-login")
	if err != nil {
		return nil, fmt.Errorf("failed to hash sentinel password: %w", err)
	}

	user := &models.User{
		Username: "seed-marker",
		Email:    "seed-marker@ticketdb.system",
		Password: hashedPassword,
	}

	return user, nil
}

// init initializes the random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}
