package main

import (
	"context"
	"log/slog"
	"os"

	"leonardovee.dev/command/internal/application"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	app := application.NewApplication(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx); err != nil {
		logger.Error("application error", "error", err)
		os.Exit(1)
	}
}
