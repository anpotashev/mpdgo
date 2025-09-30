package mpdclient

import (
	"fmt"
)

var (
	ErrNotConnected     = fmt.Errorf("not connected")
	ErrAlreadyConnected = fmt.Errorf("already connected")
	ErrOnConnection     = fmt.Errorf("connection error")
	ErrSendCommand      = fmt.Errorf("command send error")
)
