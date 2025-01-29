package command

import "log/slog"

//go:generate mockgen -source=command.go -destination=command_mock.go -package=command

type CommandType string
type CallbackFn func(Command)

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
	logger    *slog.Logger
	commands  chan Command
	handlers  map[CommandType]CommandHandler
	callbacks []CallbackFn
}

func NewCommandDispatcher(logger *slog.Logger, callbacks []CallbackFn) *Dispatcher {
	dispatcher := &Dispatcher{
		logger:    logger,
		commands:  make(chan Command, 100),
		handlers:  map[CommandType]CommandHandler{},
		callbacks: callbacks,
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
	d.commands <- command
}

func (d *Dispatcher) processCommands() {
	for command := range d.commands {
		handler, ok := d.handlers[command.GetName()]
		if !ok {
			continue
		}
		go func() {
			err := handler.Handle(command)
			if err != nil {
				d.logger.Error("error handling command", "error", err)
				return
			}
			for _, callback := range d.callbacks {
				go callback(command)
			}
		}()
	}
}
