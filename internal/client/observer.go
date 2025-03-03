package client

import (
	"strings"
)

const (
	OnConnect    = "on_connect"
	OnDisconnect = "on_disconnect"
)

func (m *MpdClientImpl) startListeningIdleChannel() {
	for {
		select {
		case idleEvent, ok := <-m.idleChannel:
			if !ok {
				return
			}
			if strings.HasPrefix(idleEvent, "changed: ") {
				idleEvent = strings.TrimPrefix(idleEvent, "changed: ")
				m.Notify(idleEvent)
			}
		case <-m.context.Done():
			return
		}
	}
}
