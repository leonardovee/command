package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) NewBooking(ctx context.Context, booking NewBookingInput) (*Booking, error) {
	var b Booking

	err := r.database.QueryRow(ctx,
		`INSERT INTO bookings (
            accommodation_id,
            user_id,
            start_at,
            end_at
        ) VALUES ($1, $2, $3, $4)
        RETURNING id, accommodation_id, user_id, start_at, end_at, created_at`,
		booking.AccommodationID,
		booking.UserID,
		booking.StartAt,
		booking.EndAt,
	).Scan(
		&b.ID,
		&b.AccommodationID,
		&b.UserID,
		&b.StartAt,
		&b.EndAt,
		&b.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *Repository) NewCancellation(ctx context.Context, booking NewCancellationInput) error {
	_, err := r.database.Exec(ctx,
		`DELETE FROM bookings WHERE id = $1 AND user_id = $2`,
		booking.BookingID,
		booking.UserID,
	)

	return err
}
