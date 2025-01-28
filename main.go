package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"leonardovee.dev/command/internal/booking"
	"leonardovee.dev/command/internal/cancellation"
	"leonardovee.dev/command/internal/command"
	"leonardovee.dev/command/internal/infrastructure"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	database, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer database.Close()

	repository := infrastructure.NewRepository(database)

	dispatcher := command.NewCommandDispatcher()
	command.RegisterHandler(dispatcher, booking.NewBookingCommandHandler(logger, repository))
	command.RegisterHandler(dispatcher, cancellation.NewCancellationCommandHandler(logger, repository))

	bookingHandler := booking.NewHandler(logger, dispatcher)
	cancellationHandler := cancellation.NewHandler(logger, dispatcher)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/bookings", bookingHandler.NewBooking)
	mux.HandleFunc("POST /api/v1/cancellations", cancellationHandler.NewCancellation)
	http.ListenAndServe(":8080", mux)
}
