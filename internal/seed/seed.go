package seed

import (
	"fmt"
	"log"

	"github.com/alexs/golang_test/internal/models"
	"github.com/alexs/golang_test/internal/repository"
	"gorm.io/gorm"
)

const (
	totalUsers     = 175
	organizerCount = 25
	regularUsers   = totalUsers - organizerCount
	totalEvents    = 200
	sentinelEmail  = "seed-marker@ticketdb.system"
)

// IsDatabaseSeeded checks if the database has already been seeded
func IsDatabaseSeeded() (bool, error) {
	var count int64
	err := repository.DB.Model(&models.User{}).Where("email = ?", sentinelEmail).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check seed status: %w", err)
	}
	return count > 0, nil
}

// ClearSeedData removes all seeded data from the database
func ClearSeedData() error {
	log.Println("Clearing existing seed data...")

	return repository.DB.Transaction(func(tx *gorm.DB) error {
		// Order matters: delete in reverse FK dependency order
		// bookings -> events -> users

		if err := tx.Exec("DELETE FROM bookings").Error; err != nil {
			return fmt.Errorf("failed to delete bookings: %w", err)
		}

		if err := tx.Exec("DELETE FROM events").Error; err != nil {
			return fmt.Errorf("failed to delete events: %w", err)
		}

		if err := tx.Exec("DELETE FROM users").Error; err != nil {
			return fmt.Errorf("failed to delete users: %w", err)
		}

		// Reset sequences for clean IDs
		sequences := []string{"bookings_id_seq", "events_id_seq", "users_id_seq"}
		for _, seq := range sequences {
			if err := tx.Exec(fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1", seq)).Error; err != nil {
				// Log but don't fail - sequence reset is nice-to-have
				log.Printf("Warning: Failed to reset sequence %s: %v", seq, err)
			}
		}

		log.Println("Successfully cleared all seed data")
		return nil
	})
}

// SeedUsers creates organizer and regular users
func SeedUsers() ([]uint, error) {
	log.Printf("Creating %d users (%d organizers, %d regular)...", totalUsers, organizerCount, regularUsers)

	var userIDs []uint

	for i := 1; i <= totalUsers; i++ {
		isOrganizer := i <= organizerCount

		user, err := generateUser(i, isOrganizer)
		if err != nil {
			return nil, fmt.Errorf("failed to generate user %d: %w", i, err)
		}

		// Use DB.Create directly since we already hashed the password in the generator
		if err := repository.DB.Create(user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user %d: %w", i, err)
		}

		userIDs = append(userIDs, user.ID)

		// Log progress every 50 users
		if i%50 == 0 {
			log.Printf("  Created %d/%d users...", i, totalUsers)
		}
	}

	log.Printf("Successfully created %d users", totalUsers)
	return userIDs, nil
}

// SeedEvents creates events with varied data
func SeedEvents(organizerIDs []uint) ([]uint, error) {
	log.Printf("Creating %d events...", totalEvents)

	var eventIDs []uint
	eventIndex := 0

	// Create events based on distribution
	for eventType, count := range eventTypeDistribution {
		for i := 0; i < count; i++ {
			eventIndex++

			// Rotate through organizers
			organizerID := organizerIDs[eventIndex%len(organizerIDs)]

			event := generateEvent(organizerID, eventType)

			if err := repository.CreateEvent(event); err != nil {
				return nil, fmt.Errorf("failed to create event %d (%s): %w", eventIndex, eventType, err)
			}

			eventIDs = append(eventIDs, event.ID)

			// Log progress every 50 events
			if eventIndex%50 == 0 {
				log.Printf("  Created %d/%d events...", eventIndex, totalEvents)
			}
		}
	}

	log.Printf("Successfully created %d events", totalEvents)
	return eventIDs, nil
}

// SeedBookings creates bookings based on distribution strategy
func SeedBookings(userIDs []uint, eventIDs []uint) error {
	log.Println("Creating bookings with realistic distribution...")

	// Get booking distribution plan
	distributions := generateBookingDistributions(eventIDs)

	// Filter to only regular users (not organizers)
	regularUserIDs := userIDs[organizerCount:] // Skip first 25 organizer users

	totalBookings := 0
	eventsWithBookings := len(distributions)

	for _, dist := range distributions {
		// Get the event to know its capacity
		var event models.Event
		if err := repository.DB.First(&event, dist.EventID).Error; err != nil {
			return fmt.Errorf("failed to get event %d: %w", dist.EventID, err)
		}

		// Calculate target number of tickets to book
		targetTickets := (event.Capacity * dist.TargetPercent) / 100
		bookedSoFar := 0

		// Create bookings until we hit target
		for bookedSoFar < targetTickets {
			// Random user from regular users
			userID := regularUserIDs[totalBookings%len(regularUserIDs)]

			booking := generateBooking(userID, dist.EventID, event.Capacity, bookedSoFar)
			if booking == nil {
				// No more capacity
				break
			}

			// CreateBooking handles validation and pricing
			if err := repository.CreateBooking(booking); err != nil {
				// Log but continue - might be capacity exceeded
				log.Printf("  Warning: Failed to create booking for event %d: %v", dist.EventID, err)
				break
			}

			bookedSoFar += booking.Quantity
			totalBookings++
		}
	}

	log.Printf("Successfully created %d bookings across %d events", totalBookings, eventsWithBookings)
	return nil
}

// CreateSentinelUser creates the sentinel marker user
func CreateSentinelUser() error {
	log.Println("Creating sentinel marker user...")

	user, err := generateSentinelUser()
	if err != nil {
		return fmt.Errorf("failed to generate sentinel user: %w", err)
	}

	if err := repository.DB.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create sentinel user: %w", err)
	}

	log.Println("Sentinel user created successfully")
	return nil
}

// Run executes the full database seeding process
func Run(force bool) error {
	log.Println("========================================")
	log.Println("Database Seeding Started")
	log.Println("========================================")

	// Check if already seeded
	seeded, err := IsDatabaseSeeded()
	if err != nil {
		return fmt.Errorf("failed to check seed status: %w", err)
	}

	if seeded && !force {
		log.Println("Database is already seeded. Skipping. Set FORCE_RESEED=true to reseed.")
		return nil
	}

	if seeded && force {
		log.Println("FORCE_RESEED enabled. Clearing existing data...")
		if err := ClearSeedData(); err != nil {
			return fmt.Errorf("failed to clear seed data: %w", err)
		}
	}

	// Step 1: Create users
	userIDs, err := SeedUsers()
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Step 2: Create events (use first organizerCount users as organizers)
	organizerIDs := userIDs[:organizerCount]
	eventIDs, err := SeedEvents(organizerIDs)
	if err != nil {
		return fmt.Errorf("failed to seed events: %w", err)
	}

	// Step 3: Create bookings
	if err := SeedBookings(userIDs, eventIDs); err != nil {
		return fmt.Errorf("failed to seed bookings: %w", err)
	}

	// Step 4: Create sentinel user
	if err := CreateSentinelUser(); err != nil {
		return fmt.Errorf("failed to create sentinel user: %w", err)
	}

	log.Println("========================================")
	log.Printf("Database Seeding Completed Successfully!")
	log.Printf("  Users: %d", totalUsers)
	log.Printf("  Events: %d", totalEvents)
	log.Println("  Bookings: Created with realistic distribution")
	log.Println("========================================")

	return nil
}
