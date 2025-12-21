package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSignup_Validation tests the signup endpoint validation logic without database
func TestSignup_Validation(t *testing.T) {
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
			name: "Missing Email",
			body: map[string]string{
				"username": "user1",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Username, email, and password are required",
		},
		{
			name: "Missing Password",
			body: map[string]string{
				"username": "user1",
				"email":    "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Username, email, and password are required",
		},
		{
			name: "Invalid Email Format",
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
		{
			name: "Empty Username",
			body: map[string]string{
				"username": "",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Username, email, and password are required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.body)
			req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Signup)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "Response body mismatch")
		})
	}
}

// TestSignup_InvalidJSON tests signup with malformed JSON
func TestSignup_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/signup", bytes.NewBufferString("{invalid json"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Signup)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request body")
}

// TestLogin_Validation tests the login endpoint validation logic
func TestLogin_Validation(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Missing Email",
			body: map[string]string{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Email and password are required",
		},
		{
			name: "Missing Password",
			body: map[string]string{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Email and password are required",
		},
		{
			name: "Empty Email",
			body: map[string]string{
				"email":    "",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Email and password are required",
		},
		{
			name: "Empty Password",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Email and password are required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.body)
			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Login)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "Response body mismatch")
		})
	}
}

// TestLogin_InvalidJSON tests login with malformed JSON
func TestLogin_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString("not json at all"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request body")
}
