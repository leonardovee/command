package booking

import (
	"fmt"
	"log/slog"
	"time"

	"leonardovee.dev/command/internal/command"
)

type BookingCommand struct {
	AcommodationID string
	UserID         string
	StartAt        time.Time
	EndAt          time.Time
}

func (c *BookingCommand) GetId() string {
	return "booking-command"
}

func (c *BookingCommand) GetName() command.CommandType {
	return "booking"
}

func NewBookingCommand(acommodationID, userID string, startAt, endAt time.Time) *BookingCommand {
	return &BookingCommand{
		AcommodationID: acommodationID,
		UserID:         userID,
		StartAt:        startAt,
		EndAt:          endAt,
	}
}

type BookingCommandHandler struct {
	logger *slog.Logger
}

func NewBookingCommandHandler(logger *slog.Logger) *BookingCommandHandler {
	return &BookingCommandHandler{
		logger: logger,
	}
}

func (h *BookingCommandHandler) Handle(command command.Command) {
	h.logger.Info(fmt.Sprintf("Handling command %s", command.GetName()))
}

func (h *BookingCommandHandler) GetName() command.CommandType {
	return "booking"
}
