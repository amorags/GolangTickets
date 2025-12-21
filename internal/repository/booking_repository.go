package repository

import (
	"errors"

	"github.com/alexs/golang_test/internal/models"
	"gorm.io/gorm"
)

func CreateBooking(booking *models.Booking) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Get the event to check availability
		var event models.Event
		if err := tx.First(&event, booking.EventID).Error; err != nil {
			return err
		}

		// Calculate available tickets
		available, err := event.AvailableTickets(tx)
		if err != nil {
			return err
		}

		// Check if enough tickets are available
		if available < booking.Quantity {
			return errors.New("not enough tickets available")
		}

		// Set pricing information
		booking.PricePerTicket = event.Price
		booking.TotalPrice = event.Price * float64(booking.Quantity)
		booking.Status = "confirmed"

		// Create the booking
		return tx.Create(booking).Error
	})
}

func GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := DB.Preload("Event").Preload("User").First(&booking, id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func GetUserBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	err := DB.Preload("Event").Where("user_id = ?", userID).Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}

func CancelBooking(bookingID uint, userID uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.First(&booking, bookingID).Error; err != nil {
			return err
		}

		// Verify the booking belongs to the user
		if booking.UserID != userID {
			return errors.New("unauthorized to cancel this booking")
		}

		// Check if already cancelled
		if booking.Status == "cancelled" {
			return errors.New("booking is already cancelled")
		}

		// Update status to cancelled
		booking.Status = "cancelled"
		return tx.Save(&booking).Error
	})
}
