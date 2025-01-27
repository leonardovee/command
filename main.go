package main

import (
	"log/slog"
	"net/http"
	"os"

	"leonardovee.dev/command/internal/booking"
	"leonardovee.dev/command/internal/command"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dispatcher := command.NewCommandDispatcher()
	command.RegisterHandler(dispatcher, booking.NewBookingCommandHandler(logger))

	bookingHandler := booking.NewHandler(logger, dispatcher)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /bookings", bookingHandler.NewBooking)
	http.ListenAndServe(":8080", mux)
}
