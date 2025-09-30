package mpdrwpool

import "errors"

var (
	ConnectionError  = errors.New("connection error")
	SendCommandError = errors.New("error sending command")
)
