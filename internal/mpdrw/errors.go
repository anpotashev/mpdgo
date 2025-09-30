package mpdrw

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrACK = fmt.Errorf("ACK error")
	ErrIO  = errors.New("IO error on sending command")
)

func parseACKAnswer(answer string) error {
	re := regexp.MustCompile(`ACK .*\{(.*)\} (.+)`)
	matches := re.FindStringSubmatch(answer)
	if len(matches) == 3 {
		command := matches[1]
		errorMessage := matches[2]
		return errors.Join(fmt.Errorf("ACK error on sending command %s: %s", command, errorMessage), ErrACK)
	}
	return errors.Join(fmt.Errorf("unexpected answer format: %s", answer), ErrACK)
}
