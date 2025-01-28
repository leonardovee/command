package command

//go:generate mockgen -source=command.go -destination=command_mock.go -package=command

import "time"

type CommandType string

type Command interface {
	GetId() string
	GetName() CommandType
}

type CommandHandler interface {
	GetName() CommandType
	Handle(command Command) error
}

type Dispatcheable interface {
	Dispatch(command Command)
}

type Dispatcher struct {
	commands []Command
	handlers map[CommandType]CommandHandler
}

func NewCommandDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{
		commands: []Command{},
		handlers: map[CommandType]CommandHandler{},
	}
	go dispatcher.processCommands()
	return dispatcher
}

func RegisterHandler(dispatcher *Dispatcher, handler CommandHandler) {
	if _, ok := dispatcher.handlers[handler.GetName()]; !ok {
		dispatcher.handlers[handler.GetName()] = handler
	}
}

func (d *Dispatcher) Dispatch(command Command) {
	d.commands = append(d.commands, command)
}

func (d *Dispatcher) processCommands() {
	// TODO: every command must run o a separeted goroutine and
	// after the execution we gotta validate if the command was
	// successful or not, if yes we remove the command from the
	// list and publish a event
	for {
		if len(d.commands) > 0 {
			command := d.commands[0]

			if command == nil {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			handler := d.handlers[command.GetName()]

			// TODO: deal with errors
			handler.Handle(command)

			d.commands = d.commands[1:]
		}
	}
}
