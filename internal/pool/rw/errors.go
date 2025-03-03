package rw

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"regexp"
)

var (
	// ServerError Фатальная ошибка, ведущая к дисконнекту (разрыв соединения, таймаут, неожиданный ответ и т.п.)
	ServerError = errors.New("server error")
)

func wrapIntoServerError(err error) error {
	return fmt.Errorf("%w: %v", ServerError, err)
}

// CommandError - ошибка, полученная от mpd-server-а после отправки команды (ответ ACK ...)
type CommandError struct {
	command      string
	errorMessage string
}

func parseACKAnswer(answer string) *CommandError {
	re := regexp.MustCompile(`ACK .*\{(.*)\} (.+)`)
	matches := re.FindStringSubmatch(answer)
	if len(matches) == 3 {
		command := matches[1]
		errorMessage := matches[2]
		return newCommandError(command, errorMessage)
	}
	log.Error().Str("answer", answer).Msg("unexpected answer format")
	return newCommandError("unexpected answer", answer)
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
