package mpdclient

import (
	"context"
	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/commands"
)

type MpdClient interface {
	// Connect connects to the MPD server.
	//
	// The requestContext is used for logging.
	//
	// Can return the following errors:
	// - ConnectionError
	// - AlreadyConnected
	Connect(requestContext context.Context) error
	// Disconnect disconnects from the MPD server.
	//
	// The requestContext is used for logging.
	//
	// Can return the following errors:
	// - NotConnected
	Disconnect(requestContext context.Context) error
	// IsConnected reports whether the client is currently connected.
	//
	// The requestContext is used for logging.
	IsConnected(requestContext context.Context) bool
	// SendSingleCommand sends a command to the MPD server
	//
	// The requestContext is used for logging.
	// returns a slice of strings containg the raw response from the MPD server
	//
	// Can return the following errors:
	// - NotConnected
	// - CommandSendError
	SendSingleCommand(requestContext context.Context, command commands.SingleCommand) ([]string, error)
	// SendBatchCommand sends a batch command to the MPD server
	//
	// The requestContext is used for logging.
	//
	// Can return the following errors:
	// - NotConnected
	// - CommandSendError
	SendBatchCommand(requestContext context.Context, cmd []commands.SingleCommand) error
	observer.Observer[string]
}

const (
	OnConnect    = "on_connect"
	OnDisconnect = "on_disconnect"
)
