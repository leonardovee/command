package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"leonardovee.dev/command/internal/booking"
	"leonardovee.dev/command/internal/command"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	database, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer database.Close()

	repository := booking.NewRepository(database)

	dispatcher := command.NewCommandDispatcher()
	command.RegisterHandler(dispatcher, booking.NewBookingCommandHandler(logger, repository))

	bookingHandler := booking.NewHandler(logger, dispatcher)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /bookings", bookingHandler.NewBooking)
	http.ListenAndServe(":8080", mux)
}
