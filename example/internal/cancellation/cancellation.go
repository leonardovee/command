package cancellation

import (
	"context"
	"fmt"
	"log/slog"

	"leonardovee.dev/command/internal/command"
	"leonardovee.dev/command/internal/infrastructure"
)

type CancellationCommand struct {
	BookingID string
	UserID    string
}

func (c *CancellationCommand) GetId() string {
	return "cancellation-command"
}

func (c *CancellationCommand) GetName() command.CommandType {
	return "cancellation"
}

func NewCancellationCommand(reservationID, userID string) *CancellationCommand {
	return &CancellationCommand{
		BookingID: reservationID,
		UserID:    userID,
	}
}

type CancellationCommandHandler struct {
	logger     *slog.Logger
	repository infrastructure.CancellationRepository
}

func NewCancellationCommandHandler(logger *slog.Logger, repository infrastructure.CancellationRepository) *CancellationCommandHandler {
	return &CancellationCommandHandler{
		logger:     logger,
		repository: repository,
	}
}

func (h *CancellationCommandHandler) GetName() command.CommandType {
	return "cancellation"
}

func (h *CancellationCommandHandler) Handle(command command.Command) error {
	ctx := context.Background()
	h.logger.Info("handling command", "name", command.GetName())

	cancellationCmd, ok := command.(*CancellationCommand)
	if !ok {
		return fmt.Errorf("invalid command type: expected CancellationCommand")
	}

	err := h.repository.NewCancellation(ctx, infrastructure.NewCancellationInput{
		BookingID: cancellationCmd.BookingID,
		UserID:    cancellationCmd.UserID,
	})
	if err != nil {
		return err
	}

	h.logger.Info("cancellation done", "id", cancellationCmd.BookingID)

	return nil
}
