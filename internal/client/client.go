package client

import (
	"context"
	"errors"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/mpdconnect"
	"sync"
)

type MpdClient interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	SendCommand(command *commands.SingleCommand) ([]string, error)
	SendBatchCommands(cmd []commands.BatchCommand) error
}

type MpdClientImpl struct {
	host     string
	port     uint16
	password string
	pool     mpdconnect.RWPool
	sync.Mutex
	listeners mpdListeners
	context   context.Context
	cancel    context.CancelFunc
}

func NewMpdClient(ctx context.Context, host string, port uint16, password string) (*MpdClientImpl, error) {
	return &MpdClientImpl{
		context:  ctx,
		host:     host,
		port:     port,
		password: password,
		pool:     nil,
	}, nil
}

func (m *MpdClientImpl) Connect() error {
	m.Lock()
	defer m.Unlock()
	if m.IsConnected() {
		return AlreadyConnected
	}
	ctx, cancel := context.WithCancel(m.context)
	m.cancel = cancel
	pool, err := mpdconnect.NewPool(ctx, m.host, m.port, m.password)
	if err != nil {
		cancel()
		return ConnectionError
	}
	m.pool = pool
	m.fireOnConnectEvent()
	return nil
}

func (m *MpdClientImpl) Disconnect() error {
	m.Lock()
	defer m.Unlock()
	if !m.IsConnected() {
		return NotConnected
	}
	go m.cancel()
	m.pool = nil
	m.fireOnDisconnectEvent()
	return nil
}

func (m *MpdClientImpl) IsConnected() bool {
	return m.pool != nil
}

func (m *MpdClientImpl) SendCommand(command *commands.SingleCommand) ([]string, error) {
	if !m.IsConnected() {
		return nil, NotConnected
	}
	result, err := m.pool.SendCommand(command)
	if err != nil {
		if errors.Is(err, mpdconnect.ConnectionError) {
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

func (m *MpdClientImpl) sendBatchCommand(command *commands.BatchCommand) error {
	err := m.pool.SendBatchCommand(command)
	if err != nil {
		if errors.Is(err, mpdconnect.ConnectionError) {
			err := m.Disconnect()
			if err != nil {
				return err
			}
			return ConnectionError
		}
		return nil
	}
	return nil
}

func (m *MpdClientImpl) SendBatchCommands(cmds []commands.BatchCommand) error {
	for _, batchCmd := range cmds {
		err := m.sendBatchCommand(&batchCmd)
		if err != nil {
			return err
		}
	}
	return nil
}
