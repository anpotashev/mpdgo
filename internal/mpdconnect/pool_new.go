package mpdconnect

import (
	"context"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/mpdconnect/rw"
)

type RWPool interface {
	SendCommand(command *commands.SingleCommand) ([]string, error)
	SendBatchCommand(command *commands.BatchCommand) error
	//Connect() error
	//Disconnect() error
}

type RWPoolImpl struct {
	rw     chan rw.MpdRW
	ctx    context.Context
	cancel context.CancelFunc
}

const poolSize = 3

func NewPool(ctx context.Context, host string, port uint16, password string) (RWPool, error) {
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
	return &RWPoolImpl{
		rw:     rwChan,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (pool *RWPoolImpl) SendCommand(command *commands.SingleCommand) ([]string, error) {
	select {
	case mpdRW := <-pool.rw:
		result, err := mpdRW.SendCommand(command)
		// todo обработка ошибок. Надо ли закрывать соединение
		pool.rw <- mpdRW
		return result, err
	case <-pool.ctx.Done():
		// todo возврат ошибки
		return nil, nil
	}
}

func (pool *RWPoolImpl) SendBatchCommand(command *commands.BatchCommand) error {
	select {
	case mpdRW := <-pool.rw:
		err := mpdRW.SendBatchCommand(command)
		// todo обработка ошибок. Надо ли закрывать соединение
		pool.rw <- mpdRW
		return err
	case <-pool.ctx.Done():
		// todo возврат ошибки
		return nil
	}
}

//func (pool *RWPoolImpl) Connect() error {
//	go pool.cancel()
//	return nil
//}
//
//func (pool *RWPoolImpl) Disconnect() error {
//	go pool.cancel()
//	return nil
//}
