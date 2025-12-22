package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/alexs/golang_test/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. ", err)
	}

	log.Println("Connected to Database!")

	// Auto-migrate the schemas
	log.Println("Running Migrations...")
	err = DB.AutoMigrate(&models.User{}, &models.Event{}, &models.Booking{})
	if err != nil {
		log.Fatal("Failed to migrate database. ", err)
	}
	log.Println("Migrations completed!")

	// Create indexes for search and filtering
	log.Println("Creating indexes...")
	err = CreateEventIndexes()
	if err != nil {
		log.Printf("Warning: Failed to create some indexes. %v", err)
	} else {
		log.Println("Indexes created successfully!")
	}
}

// CreateEventIndexes creates database indexes for event search and filtering performance
func CreateEventIndexes() error {
	// Enable PostgreSQL trigram extension for search optimization
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm").Error; err != nil {
		return fmt.Errorf("failed to create pg_trgm extension: %w", err)
	}

	// Trigram index for full-text search across multiple fields
	if err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_events_search_trgm ON events
		USING gin ((name || ' ' || description || ' ' || venue_name || ' ' || city) gin_trgm_ops)`).Error; err != nil {
		return fmt.Errorf("failed to create search trigram index: %w", err)
	}

	// Individual field indexes for filtering
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type)").Error; err != nil {
		return fmt.Errorf("failed to create event_type index: %w", err)
	}

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_city ON events(city)").Error; err != nil {
		return fmt.Errorf("failed to create city index: %w", err)
	}

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_date ON events(date)").Error; err != nil {
		return fmt.Errorf("failed to create date index: %w", err)
	}

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_price ON events(price)").Error; err != nil {
		return fmt.Errorf("failed to create price index: %w", err)
	}

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_status ON events(status)").Error; err != nil {
		return fmt.Errorf("failed to create status index: %w", err)
	}

	// Composite indexes for common filter combinations
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_status_date ON events(status, date)").Error; err != nil {
		return fmt.Errorf("failed to create status_date composite index: %w", err)
	}

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_city_date ON events(city, date)").Error; err != nil {
		return fmt.Errorf("failed to create city_date composite index: %w", err)
	}

	return nil
}
