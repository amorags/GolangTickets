package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexs/golang_test/internal/middleware"
	"github.com/alexs/golang_test/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TestBookTicket_Validation tests booking validation without database
func TestBookTicket_Validation(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]interface{}
		withAuth       bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Missing Event ID",
			body: map[string]interface{}{
				"quantity": 2,
			},
			withAuth:       true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Valid event_id and quantity are required",
		},
		{
			name: "Zero Quantity",
			body: map[string]interface{}{
				"event_id": 1,
				"quantity": 0,
			},
			withAuth:       true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Valid event_id and quantity are required",
		},
		{
			name: "Negative Quantity",
			body: map[string]interface{}{
				"event_id": 1,
				"quantity": -5,
			},
			withAuth:       true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Valid event_id and quantity are required",
		},
		{
			name: "No User Context",
			body: map[string]interface{}{
				"event_id": 1,
				"quantity": 2,
			},
			withAuth:       false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "User not found in context",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.body)
			req, err := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Add user claims to context if needed
			if tc.withAuth {
				claims := &utils.Claims{
					UserID:   1,
					Username: "testuser",
					Email:    "test@example.com",
				}
				ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(BookTicket)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "Response body mismatch")
		})
	}
}

// TestBookTicket_InvalidJSON tests booking with malformed JSON
func TestBookTicket_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/bookings", bytes.NewBufferString("{not valid json"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Add user claims to context
	claims := &utils.Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(BookTicket)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request body")
}

// TestGetMyBookings_NoAuth tests getting bookings without auth
func TestGetMyBookings_NoAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "/bookings", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetMyBookings)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "User not found in context")
}

// TestCancelBooking_InvalidID tests canceling a booking with invalid ID
func TestCancelBooking_InvalidID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		withAuth       bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Non-numeric ID",
			id:             "abc",
			withAuth:       true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid booking ID",
		},
		{
			name:           "Negative ID",
			id:             "-1",
			withAuth:       true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid booking ID",
		},
		{
			name:           "No Auth",
			id:             "1",
			withAuth:       false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "User not found in context",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/bookings/"+tc.id, nil)
			assert.NoError(t, err)

			// Add user claims to context if needed
			if tc.withAuth {
				claims := &utils.Claims{
					UserID:   1,
					Username: "testuser",
					Email:    "test@example.com",
				}
				ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
				req = req.WithContext(ctx)
			}

			// Simulate chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CancelBooking)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "Response body mismatch")
		})
	}
}
