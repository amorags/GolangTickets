package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alexs/golang_test/internal/middleware"
	"github.com/alexs/golang_test/internal/models"
	"github.com/alexs/golang_test/internal/repository"
	"github.com/alexs/golang_test/internal/utils"
	"github.com/go-chi/chi/v5"
)

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateStr)
}

// parseEventFilters extracts filter parameters from the HTTP request query string
func parseEventFilters(r *http.Request) repository.EventFilters {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	return repository.EventFilters{
		Search:    query.Get("search"),
		EventType: query.Get("type"),
		City:      query.Get("city"),
		Status:    query.Get("status"),
		Page:      page,
		Limit:     limit,
		Sort:      query.Get("sort"),
		Order:     query.Get("order"),
	}
}

type CreateEventRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	EventType   string  `json:"event_type"`
	VenueName   string  `json:"venue_name"`
	City        string  `json:"city"`
	Address     string  `json:"address"`
	Date        string  `json:"date"` // RFC3339 format
	Price       float64 `json:"price"`
	Capacity    int     `json:"capacity"`
	ImageURL    string  `json:"image_url"`
}

type EventResponse struct {
	models.Event
	AvailableTickets int `json:"available_tickets"`
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" || req.Date == "" || req.Capacity <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, "Name, date, and capacity are required")
		return
	}

	// Parse date
	date, err := parseDate(req.Date)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid date format. Use RFC3339 format")
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	event := models.Event{
		Name:        req.Name,
		Description: req.Description,
		EventType:   req.EventType,
		Status:      "published",
		OrganizerID: user.UserID,
		VenueName:   req.VenueName,
		City:        req.City,
		Address:     req.Address,
		Date:        date,
		Price:       req.Price,
		Capacity:    req.Capacity,
		ImageURL:    req.ImageURL,
	}

	if err := repository.CreateEvent(&event); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, event)
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filters := parseEventFilters(r)

	// Parse date filters
	if dateFromStr := r.URL.Query().Get("date_from"); dateFromStr != "" {
		dateFrom, err := parseDate(dateFromStr)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid date_from format. Use RFC3339 format (e.g., 2025-07-15T00:00:00Z)")
			return
		}
		filters.DateFrom = &dateFrom
	}

	if dateToStr := r.URL.Query().Get("date_to"); dateToStr != "" {
		dateTo, err := parseDate(dateToStr)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid date_to format. Use RFC3339 format (e.g., 2025-07-15T00:00:00Z)")
			return
		}
		filters.DateTo = &dateTo
	}

	// Parse price filters
	if priceMinStr := r.URL.Query().Get("price_min"); priceMinStr != "" {
		priceMin, err := strconv.ParseFloat(priceMinStr, 64)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid price_min value. Must be a number.")
			return
		}
		if priceMin < 0 {
			priceMin = 0
		}
		filters.PriceMin = &priceMin
	}

	if priceMaxStr := r.URL.Query().Get("price_max"); priceMaxStr != "" {
		priceMax, err := strconv.ParseFloat(priceMaxStr, 64)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid price_max value. Must be a number.")
			return
		}
		if priceMax < 0 {
			priceMax = 0
		}
		filters.PriceMax = &priceMax
	}

	// Get events with filters
	result, err := repository.GetEventsWithFilters(filters)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch events")
		return
	}

	// Transform to EventResponse with available tickets
	eventResponses := make([]EventResponse, len(result.Events))
	for i, event := range result.Events {
		eventResponses[i] = EventResponse{
			Event:            event.Event,
			AvailableTickets: event.Capacity - int(event.TotalBooked),
		}
	}

	// Return paginated response with metadata
	response := map[string]interface{}{
		"events":       eventResponses,
		"total":        result.Total,
		"page":         result.Page,
		"limit":        result.Limit,
		"total_pages":  result.TotalPages,
		"has_next":     result.HasNext,
		"has_previous": result.HasPrevious,
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	event, err := repository.GetEventByID(uint(id))
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Event not found")
		return
	}

	// Add available tickets count
	available, _ := event.AvailableTickets(repository.DB)
	eventResponse := EventResponse{
		Event:            *event,
		AvailableTickets: available,
	}

	utils.SuccessResponse(w, http.StatusOK, eventResponse)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	event, err := repository.GetEventByID(uint(id))
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Event not found")
		return
	}

	if event.OrganizerID != user.UserID {
		utils.ErrorResponse(w, http.StatusForbidden, "You are not authorized to delete this event")
		return
	}

	if err := repository.DeleteEvent(uint(id)); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Event deleted successfully"})
}
