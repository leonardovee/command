package command

import (
	"testing"

	gomock "go.uber.org/mock/gomock"
)

func TestCommand_processCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("should process commands", func(t *testing.T) {
		dispatcher := NewCommandDispatcher()
		done := make(chan struct{}, 1)

		command := NewMockCommand(ctrl)
		command.EXPECT().GetName().AnyTimes().Return(CommandType("test"))

		commandHandler := NewMockCommandHandler(ctrl)
		commandHandler.EXPECT().GetName().AnyTimes().Return(CommandType("test"))
		commandHandler.EXPECT().Handle(command).DoAndReturn(func(cmd Command) error {
			defer close(done)
			return nil
		})

		RegisterHandler(dispatcher, commandHandler)

		dispatcher.Dispatch(command)
		<-done
	})
}
