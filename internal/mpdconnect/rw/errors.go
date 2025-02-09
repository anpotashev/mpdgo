package rw

import (
	"errors"
	"fmt"
)

type CommandError struct {
	command      string
	errorMessage string
}

func newCommandError(command, errorMessage string) *CommandError {
	return &CommandError{
		command:      command,
		errorMessage: errorMessage,
	}
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("%s: %s", e.command, e.errorMessage)
}

func (e *CommandError) Is(err error) bool {
	var commandError *CommandError
	ok := errors.As(err, &commandError)
	return ok
}
