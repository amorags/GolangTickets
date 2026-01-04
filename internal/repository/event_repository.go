package repository

import (
	"math"
	"time"

	"github.com/alexs/golang_test/internal/models"
	"gorm.io/gorm"
)

type EventWithStats struct {
	models.Event
	TotalBooked int64 `gorm:"column:total_booked"`
}

// EventFilters contains all possible filtering, pagination, and sorting options for events
type EventFilters struct {
	Search    string
	EventType string
	City      string
	DateFrom  *time.Time
	DateTo    *time.Time
	PriceMin  *float64
	PriceMax  *float64
	Status    string
	Page      int
	Limit     int
	Sort      string
	Order     string
}

// PaginatedEventsResponse contains paginated event results with metadata
type PaginatedEventsResponse struct {
	Events      []EventWithStats `json:"events"`
	Total       int64            `json:"total"`
	Page        int              `json:"page"`
	Limit       int              `json:"limit"`
	TotalPages  int              `json:"total_pages"`
	HasNext     bool             `json:"has_next"`
	HasPrevious bool             `json:"has_previous"`
}

func CreateEvent(event *models.Event) error {
	return DB.Create(event).Error
}

func GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	err := DB.Order("date ASC").Find(&events).Error
	return events, err
}

func GetAllEventsWithStats() ([]EventWithStats, error) {
	var results []EventWithStats
	err := DB.Model(&models.Event{}).
		Select("events.*, COALESCE(SUM(bookings.quantity), 0) as total_booked").
		Joins("LEFT JOIN bookings ON bookings.event_id = events.id AND bookings.status = ?", "confirmed").
		Group("events.id").
		Order("events.date ASC").
		Scan(&results).Error
	return results, err
}

func GetEventByID(id uint) (*models.Event, error) {
	var event models.Event
	err := DB.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UpdateEvent(event *models.Event) error {
	return DB.Save(event).Error
}

func DeleteEvent(id uint) error {
	return DB.Delete(&models.Event{}, id).Error
}

// GetEventsWithFilters retrieves events with search, filtering, pagination, and sorting
func GetEventsWithFilters(filters EventFilters) (*PaginatedEventsResponse, error) {
	// Apply defaults
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 20
	}
	if filters.Status == "" {
		filters.Status = "published"
	}

	// Start with base query joining bookings for stats
	query := DB.Model(&models.Event{}).
		Select("events.*, COALESCE(SUM(bookings.quantity), 0) as total_booked").
		Joins("LEFT JOIN bookings ON bookings.event_id = events.id AND bookings.status = ?", "confirmed").
		Group("events.id")

	// Apply search filter (ILIKE for case-insensitive pattern matching)
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where(
			"events.name ILIKE ? OR events.description ILIKE ? OR events.venue_name ILIKE ? OR events.city ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Apply individual filters conditionally
	if filters.EventType != "" {
		query = query.Where("events.event_type = ?", filters.EventType)
	}

	if filters.City != "" {
		query = query.Where("LOWER(events.city) = LOWER(?)", filters.City)
	}

	if filters.DateFrom != nil {
		query = query.Where("events.date >= ?", filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("events.date <= ?", filters.DateTo)
	}

	if filters.PriceMin != nil {
		query = query.Where("events.price >= ?", filters.PriceMin)
	}

	if filters.PriceMax != nil {
		query = query.Where("events.price <= ?", filters.PriceMax)
	}

	// Status filter
	query = query.Where("events.status = ?", filters.Status)

	// Count total results BEFORE pagination
	var total int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting with validation
	sortField := "events.date"
	sortOrder := "ASC"

	validSortFields := map[string]bool{
		"date": true, "price": true, "created_at": true, "name": true,
	}

	if validSortFields[filters.Sort] {
		sortField = "events." + filters.Sort
	}

	if filters.Order == "desc" {
		sortOrder = "DESC"
	}

	query = query.Order(sortField + " " + sortOrder)

	// Apply pagination
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)

	// Execute query
	var results []EventWithStats
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	totalPages := int(math.Ceil(float64(total) / float64(filters.Limit)))

	return &PaginatedEventsResponse{
		Events:      results,
		Total:       total,
		Page:        filters.Page,
		Limit:       filters.Limit,
		TotalPages:  totalPages,
		HasNext:     filters.Page < totalPages,
		HasPrevious: filters.Page > 1,
	}, nil
}
