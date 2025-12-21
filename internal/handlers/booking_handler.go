package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexs/golang_test/internal/middleware"
	"github.com/alexs/golang_test/internal/models"
	"github.com/alexs/golang_test/internal/repository"
	"github.com/alexs/golang_test/internal/utils"
	"github.com/alexs/golang_test/internal/websocket"
	"github.com/go-chi/chi/v5"
)

// WebSocket hub for broadcasting updates
var wsHub *websocket.Hub

// SetWebSocketHub sets the WebSocket hub for broadcasting
func SetWebSocketHub(hub *websocket.Hub) {
	wsHub = hub
}

// broadcastUpdate is a helper to send availability updates
func broadcastUpdate(eventID uint) {
	if wsHub != nil {
		event, err := repository.GetEventByID(eventID)
		if err == nil {
			available, _ := event.AvailableTickets(repository.DB)
			wsHub.BroadcastAvailabilityUpdate(eventID, available, event.Capacity)
		}
	}
}

type BookTicketRequest struct {
	EventID  uint `json:"event_id"`
	Quantity int  `json:"quantity"`
}

func BookTicket(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	var req BookTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.EventID == 0 || req.Quantity <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, "Valid event_id and quantity are required")
		return
	}

	// Create booking
	booking := models.Booking{
		UserID:   claims.UserID,
		EventID:  req.EventID,
		Quantity: req.Quantity,
	}

	if err := repository.CreateBooking(&booking); err != nil {
		if err.Error() == "not enough tickets available" {
			utils.ErrorResponse(w, http.StatusBadRequest, "Not enough tickets available")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create booking: "+err.Error())
		return
	}

	// Fetch complete booking with event details
	completeBooking, err := repository.GetBookingByID(booking.ID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Booking created but failed to fetch details")
		return
	}

	// Broadcast availability update via WebSocket
	broadcastUpdate(req.EventID)

	utils.SuccessResponse(w, http.StatusCreated, completeBooking)
}

func GetMyBookings(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	bookings, err := repository.GetUserBookings(claims.UserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, bookings)
}

func CancelBooking(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	// Get booking before cancellation to get event ID
	booking, err := repository.GetBookingByID(uint(id))
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Booking not found")
		return
	}

	if err := repository.CancelBooking(uint(id), claims.UserID); err != nil {
		if err.Error() == "unauthorized to cancel this booking" {
			utils.ErrorResponse(w, http.StatusForbidden, "You are not authorized to cancel this booking")
			return
		}
		if err.Error() == "booking is already cancelled" {
			utils.ErrorResponse(w, http.StatusBadRequest, "Booking is already cancelled")
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to cancel booking")
		return
	}

	// Broadcast availability update via WebSocket
	broadcastUpdate(booking.EventID)

	utils.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Booking cancelled successfully"})
}
