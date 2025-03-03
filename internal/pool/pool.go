package pool

import (
	"context"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/pool/rw"
	"github.com/rs/zerolog/log"
)

type RWPool interface {
	SendCommand(command *commands.SingleCommand) ([]string, error)
	SendBatchCommand(command *commands.BatchCommand) error
}

type Impl struct {
	rw     chan rw.MpdRW
	idleRw rw.MpdRW
	ctx    context.Context
	cancel context.CancelFunc
}

const poolSize = 3

func NewPool(ctx context.Context, host string, port uint16, password string, channel chan string) (RWPool, error) {
	ctx, cancel := context.WithCancel(ctx)
	rwChan := make(chan rw.MpdRW, poolSize)
	for i := 0; i < poolSize; i++ {
		newRW, err := rw.NewMpdRW(ctx, host, port, password)
		if err != nil {
			cancel()
			return nil, err
		}
		rwChan <- newRW
	}
	idleRw, err := rw.NewMpdRW(ctx, host, port, password)
	if err != nil {
		cancel()
		return nil, err
	}
	result := &Impl{
		rw:     rwChan,
		ctx:    ctx,
		cancel: cancel,
		idleRw: idleRw,
	}
	go func() {
		idleChan := make(chan []string)
		errChan := make(chan error)
		defer close(idleChan)
		defer close(errChan)
		for {
			done := false
			go func() {
				results, err := idleRw.SendCommand(commands.NewSingleCommand(commands.IDLE))
				if err != nil {
					log.Err(err).Msg("error on idle")
					if !done {
						errChan <- err
					}
					return
				}
				// todo обработка ошибки
				idleChan <- results
			}()
			select {
			case lines := <-idleChan:
				for _, line := range lines {
					channel <- line
				}
			case <-ctx.Done():
				done = true
				return
			case <-errChan:
				ctx.Done()
				return
			}
		}
	}()
	return result, nil
}

func (pool *Impl) SendCommand(command *commands.SingleCommand) ([]string, error) {
	select {
	case mpdRW := <-pool.rw:
		result, err := mpdRW.SendCommand(command)
		pool.rw <- mpdRW
		return result, err
	case <-pool.ctx.Done():
		// todo возврат ошибки?
		return nil, nil
	}
}

func (pool *Impl) SendBatchCommand(command *commands.BatchCommand) error {
	select {
	case mpdRW := <-pool.rw:
		err := mpdRW.SendBatchCommand(command)
		pool.rw <- mpdRW
		return err
	case <-pool.ctx.Done():
		// todo возврат ошибки?
		return nil
	}
}
