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

		command := NewMockCommand(ctrl)
		command.EXPECT().GetName().AnyTimes().Return(CommandType("test"))

		commandHandler := NewMockCommandHandler(ctrl)
		commandHandler.EXPECT().GetName().AnyTimes().Return(CommandType("test"))
		commandHandler.EXPECT().Handle(command).Return(nil)

		RegisterHandler(dispatcher, commandHandler)

		dispatcher.Dispatch(command)
	})
}
