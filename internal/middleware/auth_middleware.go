package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexs/golang_test/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

// RequireAuth validates JWT token and adds user info to context
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Check Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(w, http.StatusUnauthorized, "Invalid authorization format. Use: Bearer <token>")
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves user claims from request context
func GetUserFromContext(r *http.Request) *utils.Claims {
	claims, ok := r.Context().Value(UserContextKey).(*utils.Claims)
	if !ok {
		return nil
	}
	return claims
}
