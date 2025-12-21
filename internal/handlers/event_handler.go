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
	events, err := repository.GetAllEvents()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch events")
		return
	}

	// Add available tickets count to each event
	eventResponses := make([]EventResponse, len(events))
	for i, event := range events {
		available, _ := event.AvailableTickets(repository.DB)
		eventResponses[i] = EventResponse{
			Event:            event,
			AvailableTickets: available,
		}
	}

	utils.SuccessResponse(w, http.StatusOK, eventResponses)
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

	if err := repository.DeleteEvent(uint(id)); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Event deleted successfully"})
}
