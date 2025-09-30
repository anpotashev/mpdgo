package mpdrwpool

import "errors"

var (
	ErrConnection     = errors.New("connection error")
	ErrSendingCommand = errors.New("error sending command")
)
