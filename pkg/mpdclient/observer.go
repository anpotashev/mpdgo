package mpdclient

import (
	"sync"
)

type MpdListener interface {
	OnEvent()
}

type mpdListeners struct {
	onConnectListeners    []MpdListener
	onDisconnectListeners []MpdListener
	sync.Mutex
}

type MpdClientObserver interface {
	AddOnConnectListener(listener MpdListener)
	RemoveOnConnectListener(listener MpdListener)
	AddOnDisconnectListener(listener MpdListener)
	RemoveOnDisconnectListener(listener MpdListener)
}

func (m *MpdClientImpl) AddOnConnectListener(listener MpdListener) {
	m.listeners.Lock()
	defer m.listeners.Unlock()
	m.listeners.onConnectListeners = append(m.listeners.onConnectListeners, listener)
}

func (m *MpdClientImpl) RemoveOnConnectListener(listener MpdListener) {
	m.listeners.Lock()
	defer m.listeners.Unlock()
	for i, l := range m.listeners.onConnectListeners {
		if l == listener {
			m.listeners.onConnectListeners = append(m.listeners.onConnectListeners[:i], m.listeners.onConnectListeners[i+1:]...)
			break
		}
	}
}

func (m *MpdClientImpl) AddOnDisconnectListener(listener MpdListener) {
	m.listeners.Lock()
	defer m.listeners.Unlock()
	m.listeners.onDisconnectListeners = append(m.listeners.onDisconnectListeners, listener)
}

func (m *MpdClientImpl) RemoveOnDisconnectListener(listener MpdListener) {
	m.listeners.Lock()
	defer m.listeners.Unlock()
	for i, l := range m.listeners.onDisconnectListeners {
		if l == listener {
			m.listeners.onDisconnectListeners = append(m.listeners.onDisconnectListeners[:i], m.listeners.onDisconnectListeners[i+1:]...)
			break
		}
	}
}

func (m *MpdClientImpl) fireOnConnectEvent() {
	m.listeners.Lock()
	listeners := make([]MpdListener, len(m.listeners.onConnectListeners))
	copy(listeners, m.listeners.onConnectListeners)
	m.listeners.Unlock()
	for _, l := range m.listeners.onConnectListeners {
		go (l).OnEvent()
	}
}

func (m *MpdClientImpl) fireOnDisconnectEvent() {
	m.listeners.Lock()
	listeners := make([]MpdListener, len(m.listeners.onDisconnectListeners))
	copy(listeners, m.listeners.onDisconnectListeners)
	m.listeners.Unlock()
	for _, l := range listeners {
		go l.OnEvent()
	}
}
