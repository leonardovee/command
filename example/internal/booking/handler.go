package booking

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"leonardovee.dev/command/internal/command"
)

type Handler struct {
	logger        *slog.Logger
	dispatcheable command.Dispatcheable
}

func NewHandler(logger *slog.Logger, dispatcheable command.Dispatcheable) *Handler {
	return &Handler{
		logger:        logger,
		dispatcheable: dispatcheable,
	}
}

func (h *Handler) NewBooking(w http.ResponseWriter, r *http.Request) {
	type request struct {
		AcommodationID string    `json:"acommodation_id" validate:"required"`
		UserID         string    `json:"user_id" validate:"required"`
		StartAt        time.Time `json:"start_at" validate:"required"`
		EndAt          time.Time `json:"end_at" validate:"required"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	command := NewBookingCommand(req.AcommodationID, req.UserID, req.StartAt, req.EndAt)
	h.dispatcheable.Dispatch(command)

	w.WriteHeader(http.StatusAccepted)
}
