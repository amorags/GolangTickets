package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestSignup_TableDriven demonstrates the standard table-driven testing pattern in Go.
func TestSignup_Validation(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Missing Username",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Username, email, and password are required",
		},
		{
			name: "Invalid Email",
			body: map[string]string{
				"username": "user1",
				"email":    "invalid-email",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid email format",
		},
		{
			name: "Short Password",
			body: map[string]string{
				"username": "user1",
				"email":    "test@example.com",
				"password": "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Password must be at least 6 characters",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Prepare request body
			jsonBody, _ := json.Marshal(tc.body)
			req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// 2. Create a ResponseRecorder (acts as the http.ResponseWriter)
			rr := httptest.NewRecorder()

			// 3. Call the handler directly
			// Note: We are testing the handler function directly.
			// In a real app with DB calls, you might need to mock the DB repository
			// or set up a test database connection before running this.
			handler := http.HandlerFunc(Signup)
			handler.ServeHTTP(rr, req)

			// 4. Assert Status Code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// 5. Assert Body Content
			// (Our utils.Error returns JSON, so we check if the body contains the error message)
			if !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want body containing %v",
					rr.Body.String(), tc.expectedBody)
			}
		})
	}
}
