package mpdrwpool

import (
	"context"
	"errors"
	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/commands"
	log "github.com/anpotashev/mpdgo/internal/logger"
	"github.com/anpotashev/mpdgo/internal/mpdrw"
	"time"
)

type Impl struct {
	cancel context.CancelFunc
	idleRW mpdrw.MpdRW
	rws    chan mpdrw.MpdRW
	observer.Observer[[]string]
	pingInterval time.Duration
	ctx          context.Context
}

type mpdRWFactory func() (mpdrw.MpdRW, error)

// NewMpdRWPool creates a new MPD RW pool.
//
// The requestContext is used for logging.
// poolSize specifies the number of active connections; the total includes
// one additional connection for idle listening, so the actual count is poolSize + 1.
// host and port specify the MPD server address.
// password is used for authentication.
// readTimeout defines the maximum time to read a line from the MPD server response.
// onDisconnect is a callback invoked when the connection is disconnected.
//
// Can return the following errors:
// - ConnectionError
func NewMpdRWPool(requestContext, ctx context.Context,
	poolSize uint8,
	host string,
	port uint16,
	password string,
	readTimeout, pingInterval time.Duration,
	onDisconnect func(),
) (*Impl, error) {
	dialer := mpdrw.NewDialer(host, port)
	var mpdRWFactoryFunction mpdRWFactory = func() (mpdrw.MpdRW, error) {
		return dialer.NewMpdRW(requestContext, ctx, password, readTimeout)
	}
	return newMpdRWPool(mpdRWFactoryFunction, requestContext, ctx, poolSize, pingInterval, onDisconnect)
}

func newMpdRWPool(
	mpdRWFactoryFunction mpdRWFactory,
	requestContext, ctx context.Context,
	poolSize uint8,
	pingInterval time.Duration,
	onDisconnect func(),
) (*Impl, error) {
	ctx, cancel := context.WithCancel(ctx)
	idleRW, err := mpdRWFactoryFunction()
	if err != nil {
		log.ErrorContext(requestContext, "Error creating idleRW", "err", err)
		cancel()
		return nil, errors.Join(ConnectionError, err)
	}
	rws := make(chan mpdrw.MpdRW, poolSize)
	for range poolSize {
		rw, err := mpdRWFactoryFunction()
		if err != nil {
			cancel()
			return nil, errors.Join(ConnectionError, err)
		}
		rws <- rw
	}
	result := &Impl{
		cancel:       cancel,
		idleRW:       idleRW,
		rws:          rws,
		Observer:     observer.New[[]string](),
		pingInterval: pingInterval,
		ctx:          ctx,
	}
	go func() {
		<-ctx.Done()
		onDisconnect()
	}()
	go result.startIdleWatching()
	go result.startPeriodicPing()
	return result, nil
}

func (p *Impl) SendSingleCommand(requestContext context.Context, command commands.SingleCommand) ([]string, error) {
	rw := <-p.rws
	defer func() {
		p.rws <- rw
	}()
	result, err := rw.SendSingleCommand(requestContext, command)
	if err != nil {
		if errors.Is(err, mpdrw.IOError) {
			p.cancel()
		}
		return nil, errors.Join(SendCommandError, err)
	}
	return result, nil
}

func (p *Impl) SendBatchCommand(requestContext context.Context, command commands.BatchCommand) error {
	rw := <-p.rws
	defer func() {
		p.rws <- rw
	}()
	err := rw.SendBatchCommand(requestContext, command)
	if err != nil {
		if errors.Is(err, mpdrw.IOError) {
			p.cancel()
		}
		return errors.Join(SendCommandError, err)
	}
	return nil
}

func (p *Impl) startIdleWatching() {
	for {
		result, err := p.idleRW.SendIdleCommand()
		if err != nil {
			p.cancel()
			return
		}
		p.Notify(result)
	}
}

func (p *Impl) startPeriodicPing() {
	tick := time.Tick(p.pingInterval)
	for {
		select {
		case <-tick:
			for i := 0; i < cap(p.rws); i++ {
				p.SendSingleCommand(nil, commands.NewSingleCommand(commands.PING))
			}
		case <-p.ctx.Done():
			return
		}
	}
}
