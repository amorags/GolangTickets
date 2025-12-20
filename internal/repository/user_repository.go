package repository

import (
	"errors"

	"github.com/alexs/golang_test/internal/models"
	"github.com/alexs/golang_test/internal/utils"
	"gorm.io/gorm"
)

// CreateUser creates a new user with hashed password
func CreateUser(username, email, password string) (*models.User, error) {
	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	result := DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// FindUserByEmail retrieves a user by email
func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// FindUserByUsername retrieves a user by username
func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// FindUserByID retrieves a user by ID
func FindUserByID(id uint) (*models.User, error) {
	var user models.User
	result := DB.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}
