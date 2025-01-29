package command

import (
	"context"
	"log/slog"
	"testing"

	gomock "go.uber.org/mock/gomock"
)

func TestCommand_processCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := slog.New(slog.NewJSONHandler(nil, nil))

	defer ctrl.Finish()

	t.Run("should process commands", func(t *testing.T) {
		done := make(chan struct{}, 1)
		dispatcher := NewCommandDispatcher(logger, []CallbackFn{
			func(context.Context, Command) {
				defer close(done)
			},
		})

		command := NewMockCommand(ctrl)
		command.EXPECT().GetName().AnyTimes().Return(CommandType("test"))

		commandHandler := NewMockCommandHandler(ctrl)
		commandHandler.EXPECT().GetName().AnyTimes().Return(CommandType("test"))
		commandHandler.EXPECT().Handle(command).Return(nil)

		RegisterHandler(dispatcher, commandHandler)

		dispatcher.Dispatch(command)
		<-done
	})
}
