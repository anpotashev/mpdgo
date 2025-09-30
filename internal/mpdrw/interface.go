package mpdrw

import (
	"context"
	"github.com/anpotashev/mpdgo/internal/commands"
)

type MpdRW interface {
	// SendIdleCommand sends an IDLE command to the MPD server
	//
	// Returns a slice of strings containing the raw response from the MPD server
	// This function should be called in a goroutine. It returns the response
	// after receiving an IDLE event from MPD server.
	// Can return the following errors:
	// - ErrIO: returned if connection is lost
	// - ErrACK: theoretically possible if the client receives an ACK response for the IDLE command; practically, this should not happen.
	SendIdleCommand() ([]string, error)

	// SendSingleCommand sends a command to the MPD server
	//
	// The requestContext is used for logging.
	// returns a slice of strings containing the raw response from the MPD server
	// Can return the following errors:
	// - ErrIO: returned if connection is lost.
	// - ErrACK: returned if an ACK response is received from the MPD server.
	SendSingleCommand(requestContext context.Context, command commands.SingleCommand) ([]string, error)

	// SendBatchCommand sends a batch command to the MPD server
	//
	// The requestContext
	// Can return the following errors:
	// - ErrIO: returned if connection is lost.
	// - ErrACK: returned if an ACK response is received from the MPD server.
	SendBatchCommand(requestContext context.Context, command commands.BatchCommand) error
}
