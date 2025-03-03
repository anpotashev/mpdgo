package client

import (
	"context"
	"errors"
	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/pool"
	"github.com/anpotashev/mpdgo/internal/pool/rw"
	"github.com/rs/zerolog/log"
	"sync"
)

type MpdClient interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	SendCommand(command *commands.SingleCommand) ([]string, error)
	SendBatchCommands(cmd []commands.BatchCommand) error
	observer.Subscriber[string]
}

type MpdClientImpl struct {
	host     string
	port     uint16
	password string
	pool     pool.RWPool
	sync.Mutex
	observer.Observer[string]
	idleChannel chan string
	context     context.Context
	cancel      context.CancelFunc
}

func NewMpdClient(ctx context.Context, host string, port uint16, password string) (*MpdClientImpl, error) {
	return &MpdClientImpl{
		context:  ctx,
		host:     host,
		port:     port,
		password: password,
		Observer: observer.New[string](),
	}, nil
}

func (m *MpdClientImpl) Connect() error {
	log.Debug().Msg("connecting")
	m.Lock()
	defer m.Unlock()
	if m.IsConnected() {
		return AlreadyConnected
	}
	ctx, cancel := context.WithCancel(m.context)
	m.cancel = cancel
	m.idleChannel = make(chan string)
	pool, err := pool.NewPool(ctx, m.host, m.port, m.password, m.idleChannel)
	if err != nil {
		cancel()
		return ConnectionError
	}
	m.pool = pool
	go m.startListeningIdleChannel()
	m.Notify(OnConnect)
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
	m.Notify(OnDisconnect)
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
	log.Debug().Str("command", command.String()).Msg("sending command")
	if err != nil {
		log.Debug().Msg("error!!!")
		if errors.Is(err, rw.ServerError) {
			log.Debug().Msg("disconnecting")
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
		if errors.Is(err, rw.ServerError) {
			err := m.Disconnect()
			if err != nil {
				return err
			}
			return ConnectionError
		}
		return err
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
