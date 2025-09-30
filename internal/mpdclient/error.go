package mpdclient

import (
	"fmt"
)

var (
	NotConnected     = fmt.Errorf("not connected")
	AlreadyConnected = fmt.Errorf("already connected")
	ConnectionError  = fmt.Errorf("connection error")
	CommandSendError = fmt.Errorf("command send error")
)
