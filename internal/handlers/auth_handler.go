package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alexs/golang_test/internal/middleware"
	"github.com/alexs/golang_test/internal/repository"
	"github.com/alexs/golang_test/internal/utils"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup handles user registration
func Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.Error(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// Validate email format (basic check)
	if !strings.Contains(req.Email, "@") {
		utils.Error(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Validate password length
	if len(req.Password) < 6 {
		utils.Error(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	// Create user in database
	user, err := repository.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		// Check for duplicate username/email
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			utils.Error(w, http.StatusConflict, "Username or email already exists")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return success response with token
	utils.Success(w, "User created successfully", map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		utils.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Find user by email
	user, err := repository.FindUserByEmail(req.Email)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Compare passwords
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		utils.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return success response with token
	utils.Success(w, "Login successful", map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetProfile returns the authenticated user's profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (populated by middleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Optionally fetch full user from database
	user, err := repository.FindUserByID(claims.UserID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "User not found")
		return
	}

	// Return user profile
	utils.Success(w, "Profile retrieved", map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
