package booking

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"leonardovee.dev/command/internal/command"
	"leonardovee.dev/command/internal/infrastructure"
)

type BookingCommand struct {
	AccommodationID string
	UserID          string
	StartAt         time.Time
	EndAt           time.Time
}

func (c *BookingCommand) GetId() string {
	return "booking-command"
}

func (c *BookingCommand) GetName() command.CommandType {
	return "booking"
}

func NewBookingCommand(acommodationID, userID string, startAt, endAt time.Time) *BookingCommand {
	return &BookingCommand{
		AccommodationID: acommodationID,
		UserID:          userID,
		StartAt:         startAt,
		EndAt:           endAt,
	}
}

type BookingCommandHandler struct {
	logger     *slog.Logger
	repository infrastructure.BookingRepository
}

func NewBookingCommandHandler(logger *slog.Logger, repository infrastructure.BookingRepository) *BookingCommandHandler {
	return &BookingCommandHandler{
		logger:     logger,
		repository: repository,
	}
}

func (h *BookingCommandHandler) GetName() command.CommandType {
	return "booking"
}

func (h *BookingCommandHandler) Handle(command command.Command) error {
	ctx := context.Background()
	h.logger.Info("handling command", "name", command.GetName())

	bookingCmd, ok := command.(*BookingCommand)
	if !ok {
		return fmt.Errorf("invalid command type: expected BookingCommand")
	}

	booking, err := h.repository.NewBooking(ctx, infrastructure.NewBookingInput{
		AccommodationID: bookingCmd.AccommodationID,
		UserID:          bookingCmd.UserID,
		StartAt:         bookingCmd.StartAt,
		EndAt:           bookingCmd.EndAt,
	})
	if err != nil {
		h.logger.Error("failed to create booking", "error", err)
		return err
	}

	h.logger.Info("booking created", "id", booking.ID)

	return nil
}
