package mpdclient

import "fmt"

type MpdCommand interface {
	fmt.Stringer
}
