package mpdrwpool

import (
	"context"
	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/commands"
)

type MpdRWPool interface {
	// SendSingleCommand sends a command to the MPD server
	//
	// The requestContext is used for logging.
	// returns a slice of strings containing the raw response from the MPD server
	//
	// Can return the following errors:
	// - SendCommandError
	SendSingleCommand(requestContext context.Context, command commands.SingleCommand) ([]string, error)

	// SendBatchCommand sends a batch command to the MPD server
	//
	// The requestContext is used for logging.
	//
	// Can return the following errors:
	// - SendCommandError
	SendBatchCommand(requestContext context.Context, command commands.BatchCommand) error
	observer.Observer[[]string]
}
