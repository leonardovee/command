package infrastructure

import "context"

type CancellationRepository interface {
	NewCancellation(ctx context.Context, booking NewCancellationInput) error
}

type NewCancellationInput struct {
	BookingID string
	UserID    string
}
