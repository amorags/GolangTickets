package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendJSON sends a JSON response with status code
func SendJSON(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Success sends a successful JSON response
func Success(w http.ResponseWriter, message string, data interface{}) {
	SendJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error JSON response
func Error(w http.ResponseWriter, statusCode int, message string) {
	SendJSON(w, statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// ErrorResponse sends an error JSON response with custom status code
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	Error(w, statusCode, message)
}

// SuccessResponse sends a successful JSON response with custom status code
func SuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	SendJSON(w, statusCode, Response{
		Success: true,
		Data:    data,
	})
}
