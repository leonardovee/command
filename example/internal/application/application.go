package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jackc/pgx/v5/pgxpool"
	"leonardovee.dev/command/internal/booking"
	"leonardovee.dev/command/internal/cancellation"
	"leonardovee.dev/command/internal/command"
	"leonardovee.dev/command/internal/infrastructure"
)

type application struct {
	config config
	logger *slog.Logger
	server *http.Server
}

func NewApplication(logger *slog.Logger) *application {
	return &application{
		config: loadConfig(),
		logger: logger,
	}
}

func (app *application) setupDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, app.config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func (app *application) setupAWS(ctx context.Context) (*sns.Client, error) {
	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := sns.NewFromConfig(awsConfig, func(o *sns.Options) {
		o.BaseEndpoint = aws.String(app.config.AWSEndpoint)
	})

	return client, nil
}

func (app *application) setupRoutes(bookingHandler *booking.Handler, cancellationHandler *cancellation.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/bookings", bookingHandler.NewBooking)
	mux.HandleFunc("POST /api/v1/cancellations", cancellationHandler.NewCancellation)

	return mux
}

func (app *application) Run(ctx context.Context) error {
	db, err := app.setupDatabase(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	defer db.Close()

	snsClient, err := app.setupAWS(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup AWS SNS: %w", err)
	}

	repository := infrastructure.NewRepository(db)
	snsService := infrastructure.NewSNS(app.logger, snsClient)

	dispatcher := command.NewDispatcher(app.logger, []command.CallbackFn{
		snsService.Publish,
	})

	command.RegisterHandler(dispatcher, booking.NewBookingCommandHandler(app.logger, repository))
	command.RegisterHandler(dispatcher, cancellation.NewCancellationCommandHandler(app.logger, repository))

	bookingHandler := booking.NewHandler(app.logger, dispatcher)
	cancellationHandler := cancellation.NewHandler(app.logger, dispatcher)

	app.server = &http.Server{
		Addr:         ":" + app.config.Port,
		Handler:      app.setupRoutes(bookingHandler, cancellationHandler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		app.logger.Info("starting server", "port", app.config.Port)
		serverErrors <- app.server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		return app.shutdown()
	}
}

func (app *application) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.config.ShutdownTimeout)
	defer cancel()

	app.logger.Info("starting graceful shutdown")
	if err := app.server.Shutdown(ctx); err != nil {
		app.server.Close()
		return fmt.Errorf("could not stop server gracefully: %w", err)
	}

	return nil
}
