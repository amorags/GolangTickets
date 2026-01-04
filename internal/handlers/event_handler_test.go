package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TestCreateEvent_Validation tests event creation validation without database
func TestCreateEvent_Validation(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Missing Name",
			body: map[string]interface{}{
				"date":     "2025-12-25T18:00:00Z",
				"capacity": 100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Name, date, and capacity are required",
		},
		{
			name: "Missing Date",
			body: map[string]interface{}{
				"name":     "Test Event",
				"capacity": 100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Name, date, and capacity are required",
		},
		{
			name: "Zero Capacity",
			body: map[string]interface{}{
				"name":     "Test Event",
				"date":     "2025-12-25T18:00:00Z",
				"capacity": 0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Name, date, and capacity are required",
		},
		{
			name: "Negative Capacity",
			body: map[string]interface{}{
				"name":     "Test Event",
				"date":     "2025-12-25T18:00:00Z",
				"capacity": -10,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Name, date, and capacity are required",
		},
		{
			name: "Invalid Date Format",
			body: map[string]interface{}{
				"name":     "Test Event",
				"date":     "2025-12-25",
				"capacity": 100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid date format",
		},
		{
			name: "Invalid Date Format 2",
			body: map[string]interface{}{
				"name":     "Test Event",
				"date":     "not a date",
				"capacity": 100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid date format",
		},
		{
			name: "Empty Name",
			body: map[string]interface{}{
				"name":     "",
				"date":     "2025-12-25T18:00:00Z",
				"capacity": 100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Name, date, and capacity are required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.body)
			req, err := http.NewRequest("POST", "/events", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateEvent)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "Response body mismatch")
		})
	}
}

// TestCreateEvent_InvalidJSON tests event creation with malformed JSON
func TestCreateEvent_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/events", bytes.NewBufferString("{invalid json}"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateEvent)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request body")
}

// TestParseDate tests the parseDate helper function
func TestParseDate(t *testing.T) {
	tests := []struct {
		name      string
		dateStr   string
		expectErr bool
	}{
		{
			name:      "Valid RFC3339 Date",
			dateStr:   "2025-12-25T18:00:00Z",
			expectErr: false,
		},
		{
			name:      "Valid RFC3339 with Timezone",
			dateStr:   "2025-12-25T18:00:00+05:00",
			expectErr: false,
		},
		{
			name:      "Invalid Format - Date Only",
			dateStr:   "2025-12-25",
			expectErr: true,
		},
		{
			name:      "Invalid Format - Random String",
			dateStr:   "not a date",
			expectErr: true,
		},
		{
			name:      "Invalid Format - Empty",
			dateStr:   "",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseDate(tc.dateStr)
			if tc.expectErr {
				assert.Error(t, err, "Expected error for invalid date")
			} else {
				assert.NoError(t, err, "Expected no error for valid date")
			}
		})
	}
}

// TestGetEvent_InvalidID tests getting an event with invalid ID
func TestGetEvent_InvalidID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Non-numeric ID",
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid event ID",
		},
		{
			name:           "Negative ID",
			id:             "-1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid event ID",
		},
		{
			name:           "Empty ID",
			id:             "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid event ID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/events/"+tc.id, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			// Simulate chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler := http.HandlerFunc(GetEvent)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "Response body mismatch")
		})
	}
}

// TestDeleteEvent_InvalidID tests deleting an event with invalid ID
func TestDeleteEvent_InvalidID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Non-numeric ID",
			id:             "xyz",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid event ID",
		},
		{
			name:           "Negative ID",
			id:             "-5",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid event ID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/events/"+tc.id, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			// Simulate chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler := http.HandlerFunc(DeleteEvent)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "Response body mismatch")
		})
	}
}
