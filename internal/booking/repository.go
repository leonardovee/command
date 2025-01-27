package booking

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type reservation struct {
	ID              string    `json:"id" db:"id"`
	AccommodationID string    `json:"accommodation_id" db:"accommodation_id"`
	UserID          string    `json:"user_id" db:"user_id"`
	StartAt         time.Time `json:"start_at" db:"start_at"`
	EndAt           time.Time `json:"end_at" db:"end_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) NewBooking(ctx context.Context, booking BookingCommand) (*reservation, error) {
	var reservation reservation

	err := r.database.QueryRow(ctx,
		`INSERT INTO reservations (
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
		&reservation.ID,
		&reservation.AccommodationID,
		&reservation.UserID,
		&reservation.StartAt,
		&reservation.EndAt,
		&reservation.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}
