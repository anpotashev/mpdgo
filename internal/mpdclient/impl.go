package mpdclient

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/commands"
	log "github.com/anpotashev/mpdgo/internal/logger"
	"github.com/anpotashev/mpdgo/internal/mpdrwpool"
)

type config struct {
	host                  string
	port                  uint16
	password              string
	readTimeout           time.Duration
	pingPeriod            time.Duration
	maxBatchCommandLength uint16
	poolSize              uint8
}

type Impl struct {
	mu         sync.Mutex
	pool       mpdrwpool.MpdRWPool
	ctx        context.Context
	config     config
	cancelFunc context.CancelFunc
	observer.Observer[string]
}

func NewMpdClientImpl(ctx context.Context,
	host string,
	port uint16,
	password string,
	maxBatchCommandLength uint16,
	poolSize uint8,
	readTimeout, pingPeriod time.Duration) *Impl {
	return &Impl{
		ctx: ctx,
		config: config{
			host:                  host,
			port:                  port,
			password:              password,
			maxBatchCommandLength: maxBatchCommandLength,
			poolSize:              poolSize,
			readTimeout:           readTimeout,
			pingPeriod:            pingPeriod,
		},
		cancelFunc: nil,
		Observer:   observer.New[string](),
	}
}

type newMpdRWPoolFactory func(requestContext, ctx context.Context, onDisconnect func()) (mpdrwpool.MpdRWPool, error)

func (m *Impl) newMpdRWPoolFactory(requestContext, ctx context.Context, onDisconnect func()) (mpdrwpool.MpdRWPool, error) {
	return mpdrwpool.NewMpdRWPool(
		requestContext,
		ctx,
		m.config.poolSize,
		m.config.host,
		m.config.port,
		m.config.password,
		m.config.readTimeout,
		m.config.pingPeriod,
		onDisconnect)
}

func (m *Impl) Connect(requestContext context.Context) error {
	return m.connect(requestContext, m.newMpdRWPoolFactory)
}

func (m *Impl) connect(requestContext context.Context, newMpdRWPoolFactoryFunc newMpdRWPoolFactory) error {
	log.DebugContext(requestContext, "Trying to connect")
	log.Debug("Locking resource")
	m.mu.Lock()
	defer func() {
		log.DebugContext(requestContext, "Unlocking resource")
		m.mu.Unlock()
	}()
	if m.pool != nil {
		log.DebugContext(requestContext, "Already connected")
		return ErrAlreadyConnected
	}
	log.DebugContext(requestContext, "Creating a cancel context")
	ctx, cancel := context.WithCancel(m.ctx)
	log.DebugContext(requestContext, "Creating an onDisconnect function")
	//lint:ignore SA1012 ignore
	onDisconnect := func() { m.Disconnect(nil) }
	log.DebugContext(requestContext, "Creating an mpdRWPool")
	pool, err := newMpdRWPoolFactoryFunc(requestContext, ctx, onDisconnect)
	if err != nil {
		log.ErrorContext(requestContext, "Error creating new mpd rw pool", "err", err)
		cancel()
		return errors.Join(ErrOnConnection, err)
	}
	log.DebugContext(requestContext, "pool successfully created")
	log.DebugContext(requestContext, "Subscribing to IDLE events")
	subscribe := pool.Subscribe(time.Millisecond * 100)
	go func() {
		for {
			select {
			case idleEvents := <-subscribe:
				for _, event := range idleEvents {
					if len(event) > 9 && event[:9] == "changed: " {
						m.Notify(event[9:])
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	m.pool = pool
	m.cancelFunc = cancel
	log.DebugContext(requestContext, "Sending an onConnect event")
	m.Notify(OnConnect)
	return nil
}

func (m *Impl) Disconnect(requestContext context.Context) error {
	log.DebugContext(requestContext, "Disconnecting")
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.pool == nil {
		return ErrNotConnected
	}
	m.cancelFunc()
	m.pool = nil
	m.cancelFunc = nil
	m.Notify(OnDisconnect)
	return nil
}

func (m *Impl) IsConnected(requestContext context.Context) bool {
	log.DebugContext(requestContext, "Getting connected state")
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.pool != nil
}

func (m *Impl) SendSingleCommand(requestContext context.Context, command commands.SingleCommand) ([]string, error) {
	log.DebugContext(requestContext, "Sending single command", "command", log.Truncate(command.String(), 100))
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.pool == nil {
		return nil, ErrNotConnected
	}
	response, err := m.pool.SendSingleCommand(requestContext, command)
	if err != nil {
		return nil, errors.Join(ErrSendCommand, err)
	}
	return response, nil
}

func (m *Impl) SendBatchCommand(requestContext context.Context, cmds []commands.SingleCommand) error {
	log.DebugContext(requestContext, "Sending batch commands", "commands", log.JoinAndTruncateSingleCommands(cmds, "\n", 100))
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.pool == nil {
		return ErrNotConnected
	}
	for _, batchCommand := range commands.NewBatchCommands(cmds, int(m.config.maxBatchCommandLength)) {
		err := m.pool.SendBatchCommand(requestContext, batchCommand)
		if err != nil {
			return errors.Join(err, ErrSendCommand)
		}
	}
	return nil
}
