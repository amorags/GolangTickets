package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHashPassword tests password hashing
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Short Password",
			password: "123456",
		},
		{
			name:     "Medium Password",
			password: "mySecurePassword123",
		},
		{
			name:     "Long Password with Special Chars",
			password: "MyVeryLongAndSecureP@ssw0rd!2024",
		},
		{
			name:     "Password with Spaces",
			password: "password with spaces",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := HashPassword(tc.password)

			assert.NoError(t, err, "Should hash password without error")
			assert.NotEmpty(t, hash, "Hash should not be empty")
			assert.NotEqual(t, tc.password, hash, "Hash should not equal plain password")
			assert.Greater(t, len(hash), 50, "Bcrypt hash should be at least 50 characters")
		})
	}
}

// TestHashPassword_Uniqueness ensures each hash is unique even for same password
func TestHashPassword_Uniqueness(t *testing.T) {
	password := "testPassword123"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "Two hashes of same password should be different (due to salt)")
}

// TestComparePassword tests password comparison
func TestComparePassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		compareWith    string
		expectMatch    bool
	}{
		{
			name:        "Correct Password",
			password:    "correctPassword",
			compareWith: "correctPassword",
			expectMatch: true,
		},
		{
			name:        "Wrong Password",
			password:    "correctPassword",
			compareWith: "wrongPassword",
			expectMatch: false,
		},
		{
			name:        "Empty Password",
			password:    "testPassword",
			compareWith: "",
			expectMatch: false,
		},
		{
			name:        "Case Sensitive",
			password:    "Password123",
			compareWith: "password123",
			expectMatch: false,
		},
		{
			name:        "With Special Characters",
			password:    "P@ssw0rd!#$",
			compareWith: "P@ssw0rd!#$",
			expectMatch: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Hash the original password
			hash, err := HashPassword(tc.password)
			assert.NoError(t, err)

			// Compare with the test password
			err = ComparePassword(hash, tc.compareWith)

			if tc.expectMatch {
				assert.NoError(t, err, "Should match password")
			} else {
				assert.Error(t, err, "Should not match password")
			}
		})
	}
}

// TestComparePassword_InvalidHash tests comparison with invalid hash
func TestComparePassword_InvalidHash(t *testing.T) {
	err := ComparePassword("not-a-valid-hash", "somePassword")
	assert.Error(t, err, "Should error on invalid hash format")
}

// TestHashPassword_EmptyPassword tests hashing empty password
func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("")

	// Bcrypt can hash empty strings, so this should succeed
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}