package infrastructure

import (
	"context"
	"time"
)

type BookingRepository interface {
	NewBooking(ctx context.Context, booking NewBookingInput) (*Booking, error)
}

type NewBookingInput struct {
	AccommodationID string
	UserID          string
	StartAt         time.Time
	EndAt           time.Time
}

type Booking struct {
	ID              string    `json:"id"`
	AccommodationID string    `json:"accommodation_id"`
	UserID          string    `json:"user_id"`
	StartAt         time.Time `json:"start_at"`
	EndAt           time.Time `json:"end_at"`
	CreatedAt       time.Time `json:"created_at"`
}
