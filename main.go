package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jackc/pgx/v5/pgxpool"
	"leonardovee.dev/command/internal/booking"
	"leonardovee.dev/command/internal/cancellation"
	"leonardovee.dev/command/internal/command"
	"leonardovee.dev/command/internal/infrastructure"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	database, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer database.Close()
	repository := infrastructure.NewRepository(database)

	ac, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Error("failed to load aws config", "error", err)
		os.Exit(1)
	}

	client := sns.NewFromConfig(ac, func(o *sns.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
	})
	sns := infrastructure.NewSNS(logger, client)

	dispatcher := command.NewCommandDispatcher(logger, []command.CallbackFn{
		sns.Publish,
	})
	command.RegisterHandler(dispatcher, booking.NewBookingCommandHandler(logger, repository))
	command.RegisterHandler(dispatcher, cancellation.NewCancellationCommandHandler(logger, repository))

	bookingHandler := booking.NewHandler(logger, dispatcher)
	cancellationHandler := cancellation.NewHandler(logger, dispatcher)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/bookings", bookingHandler.NewBooking)
	mux.HandleFunc("POST /api/v1/cancellations", cancellationHandler.NewCancellation)
	http.ListenAndServe(":8080", mux)
}
