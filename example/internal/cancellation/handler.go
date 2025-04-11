package cancellation

import (
	"encoding/json"
	"log/slog"
	"net/http"

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

func (h *Handler) NewCancellation(w http.ResponseWriter, r *http.Request) {
	type request struct {
		BookingID string `json:"booking_id" validate:"required"`
		UserID    string `json:"user_id" validate:"required"`
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

	command := NewCancellationCommand(req.BookingID, req.UserID)
	h.dispatcheable.Dispatch(command)

	w.WriteHeader(http.StatusAccepted)
}
