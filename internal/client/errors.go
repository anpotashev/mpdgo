package client

import (
	"errors"
)

var (
	ValidationError  = errors.New("validation error")
	AlreadyConnected = errors.New("already connected")
	NotConnected     = errors.New("not connected")
	ConnectionError  = errors.New("connection error")
)
