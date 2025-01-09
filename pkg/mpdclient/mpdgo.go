package mpdclient

import (
	"errors"
	"github.com/anpotashev/mpdgo/internal/mpd"
	"sync"
)

type MpdClient interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	SendCommand(command *MpdCommand) ([]string, error)
}

type MpdClientImpl struct {
	host             string
	port             uint16
	password         string
	currentMpdRWPool mpd.MpdRWPool
	sync.Mutex
	factory   mpd.MpdRWPoolFactory
	listeners mpdListeners
}

func NewMpdClient(host string, port uint16, password string) (*MpdClientImpl, error) {
	return &MpdClientImpl{
		host:     host,
		port:     port,
		password: password,
		factory:  mpd.MpdRWPoolFactoryImpl{},
	}, nil
}

func (m *MpdClientImpl) Connect() error {
	m.Lock()
	defer m.Unlock()
	if m.IsConnected() {
		return AlreadyConnected
	}
	currentMpdRWPool, err := m.factory.CreateAndConnect(m.host, m.port, m.password)
	if err != nil {
		return ConnectionError
	}
	m.currentMpdRWPool = currentMpdRWPool
	m.fireOnConnectEvent()
	return nil
}

func (m *MpdClientImpl) Disconnect() error {
	m.Lock()
	defer m.Unlock()
	if !m.IsConnected() {
		return NotConnected
	}
	go m.currentMpdRWPool.Disconnect()
	m.currentMpdRWPool = nil
	m.fireOnDisconnectEvent()
	return nil
}

func (m *MpdClientImpl) IsConnected() bool {
	return m.currentMpdRWPool != nil
}

func (m *MpdClientImpl) SendCommand(command MpdCommand) ([]string, error) {
	if !m.IsConnected() {
		return nil, NotConnected
	}
	result, err := m.currentMpdRWPool.SendCommand(command.String())
	if err != nil {
		if errors.Is(err, mpd.ConnectionError) {
			err := m.Disconnect()
			if err != nil {
				return nil, err
			}
			return nil, ConnectionError
		}
		return nil, err
	}
	return result, nil
}
