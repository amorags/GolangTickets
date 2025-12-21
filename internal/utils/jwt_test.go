package utils

import (
	"testing"
	"time"

	"github.com/alexs/golang_test/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestGenerateJWT tests JWT token generation
func TestGenerateJWT(t *testing.T) {
	// Ensure config is loaded for tests
	if config.AppConfig == nil {
		config.AppConfig = &config.Config{
			JWTSecret:     "test-secret-key-for-testing",
			JWTExpiration: 24 * time.Hour,
		}
	}

	tests := []struct {
		name     string
		userID   uint
		username string
		email    string
	}{
		{
			name:     "Valid User Data",
			userID:   1,
			username: "testuser",
			email:    "test@example.com",
		},
		{
			name:     "User with ID Zero",
			userID:   0,
			username: "zerouser",
			email:    "zero@example.com",
		},
		{
			name:     "User with Special Characters in Username",
			userID:   123,
			username: "user_name-123",
			email:    "special@example.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GenerateJWT(tc.userID, tc.username, tc.email)

			assert.NoError(t, err, "Should generate token without error")
			assert.NotEmpty(t, token, "Token should not be empty")
			assert.Greater(t, len(token), 20, "JWT should be reasonably long")
		})
	}
}

// TestValidateJWT tests JWT token validation
func TestValidateJWT(t *testing.T) {
	// Ensure config is loaded for tests
	if config.AppConfig == nil {
		config.AppConfig = &config.Config{
			JWTSecret:     "test-secret-key-for-testing",
			JWTExpiration: 24 * time.Hour,
		}
	}

	tests := []struct {
		name          string
		userID        uint
		username      string
		email         string
		expectValid   bool
	}{
		{
			name:        "Valid Token",
			userID:      1,
			username:    "testuser",
			email:       "test@example.com",
			expectValid: true,
		},
		{
			name:        "Valid Token with Different Data",
			userID:      999,
			username:    "anotheruser",
			email:       "another@example.com",
			expectValid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Generate a token
			token, err := GenerateJWT(tc.userID, tc.username, tc.email)
			assert.NoError(t, err)

			// Validate the token
			claims, err := ValidateJWT(token)

			if tc.expectValid {
				assert.NoError(t, err, "Should validate token without error")
				assert.NotNil(t, claims, "Claims should not be nil")
				assert.Equal(t, tc.userID, claims.UserID, "UserID should match")
				assert.Equal(t, tc.username, claims.Username, "Username should match")
				assert.Equal(t, tc.email, claims.Email, "Email should match")
			} else {
				assert.Error(t, err, "Should error on invalid token")
			}
		})
	}
}

// TestValidateJWT_InvalidToken tests validation of invalid tokens
func TestValidateJWT_InvalidToken(t *testing.T) {
	// Ensure config is loaded for tests
	if config.AppConfig == nil {
		config.AppConfig = &config.Config{
			JWTSecret:     "test-secret-key-for-testing",
			JWTExpiration: 24 * time.Hour,
		}
	}

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty Token",
			token: "",
		},
		{
			name:  "Random String",
			token: "not-a-valid-token",
		},
		{
			name:  "Malformed JWT",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
		},
		{
			name:  "Token with Wrong Format",
			token: "Bearer sometoken",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := ValidateJWT(tc.token)

			assert.Error(t, err, "Should error on invalid token")
			assert.Nil(t, claims, "Claims should be nil for invalid token")
		})
	}
}

// TestValidateJWT_WrongSecret tests validation with wrong secret
func TestValidateJWT_WrongSecret(t *testing.T) {
	// Ensure config is initialized
	if config.AppConfig == nil {
		config.AppConfig = &config.Config{}
	}

	originalSecret := config.AppConfig.JWTSecret
	defer func() { config.AppConfig.JWTSecret = originalSecret }()

	// Generate token with one secret
	config.AppConfig.JWTSecret = "secret1"
	config.AppConfig.JWTExpiration = 24 * time.Hour
	token, err := GenerateJWT(1, "testuser", "test@example.com")
	assert.NoError(t, err)

	// Try to validate with different secret
	config.AppConfig.JWTSecret = "secret2"
	claims, err := ValidateJWT(token)

	assert.Error(t, err, "Should error when validating with wrong secret")
	assert.Nil(t, claims, "Claims should be nil when secret doesn't match")
}

// TestJWT_RoundTrip tests full generate and validate cycle
func TestJWT_RoundTrip(t *testing.T) {
	// Ensure config is loaded for tests
	if config.AppConfig == nil {
		config.AppConfig = &config.Config{
			JWTSecret:     "test-secret-key-for-testing",
			JWTExpiration: 24 * time.Hour,
		}
	}

	userID := uint(42)
	username := "roundtripuser"
	email := "roundtrip@example.com"

	// Generate token
	token, err := GenerateJWT(userID, username, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	claims, err := ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Verify all fields match
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)

	// Verify timestamps
	assert.True(t, claims.ExpiresAt.After(time.Now()), "Token should not be expired")
	assert.True(t, claims.IssuedAt.Before(time.Now().Add(time.Second)), "IssuedAt should be in the past")
}
