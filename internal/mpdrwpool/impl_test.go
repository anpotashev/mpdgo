package mpdrwpool

import (
	"context"
	"fmt"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/mpdrw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync/atomic"
	"testing"
	"time"
)

var defaultConnectParams = struct {
	poolSize       uint8
	host           string
	port           uint16
	password       string
	readTimeout    time.Duration
	pingInterval   time.Duration
	ctx            context.Context
	requestContext context.Context
}{
	poolSize:       3,
	host:           "localhost",
	port:           6600,
	password:       "12345",
	readTimeout:    time.Millisecond * 100,
	pingInterval:   time.Second,
	ctx:            context.Background(),
	requestContext: context.Background(),
}

func TestNewMpdRWPool(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		mockRw := &mockMpdRW{}
		var mpdRWFactoryFunc mpdRWFactory = func() (mpdrw.MpdRW, error) {
			return mockRw, nil
		}
		idleChan := make(chan struct{})
		mockRw.On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		pool, err := newMpdRWPool(mpdRWFactoryFunc, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		select {
		case <-onDisconnectCalled:
			t.Error("onDisconnect was called")
		case <-time.NewTimer(time.Microsecond * 100).C:
		}
		pool.cancel()
		select {
		case <-time.NewTimer(time.Microsecond * 100).C:
			t.Error("onDisconnect was not called")
		case <-onDisconnectCalled:
		}
	})
	t.Run("error when creating rw (idleRW)", func(t *testing.T) {
		var mpdRWFactoryFunc mpdRWFactory = func() (mpdrw.MpdRW, error) {
			return nil, fmt.Errorf("error")
		}
		onDisconnect := func() {}
		pool, err := newMpdRWPool(mpdRWFactoryFunc, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ConnectionError)
		assert.Nil(t, pool)
	})
	t.Run("error when creating rw (non idleRW)", func(t *testing.T) {
		first := true
		var mpdRWFactoryFunc mpdRWFactory = func() (mpdrw.MpdRW, error) {
			if first {
				first = false
				return &mockMpdRW{}, nil
			}
			return nil, fmt.Errorf("error")
		}
		onDisconnect := func() {}
		pool, err := newMpdRWPool(mpdRWFactoryFunc, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ConnectionError)
		assert.Nil(t, pool)
	})
}

func TestImpl_VerifyPingGoroutine(t *testing.T) {
	t.Run("check ping command (no error)", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Mocking responses for PING command
		pingCmd := commands.NewSingleCommand(commands.PING)
		var pingCount = atomic.Uint32{}
		for _, rw := range rws[1:] {
			rw.On("SendSingleCommand", nil, pingCmd).
				Run(func(args mock.Arguments) {
					pingCount.Add(1)
				}).
				Return([]string{}, nil)
		}
		// Sleeping for longer than the ping interval
		time.Sleep(defaultConnectParams.pingInterval + (defaultConnectParams.pingInterval / 10))
		// Verifying than PING was called poolSize times
		assert.Equal(t, defaultConnectParams.poolSize, uint8(pingCount.Load()))
		// Verifying that onDisconnect was not called.
		select {
		case <-onDisconnectCalled:
			t.Error("onDisconnect was called")
		case <-time.NewTimer(time.Microsecond * 100).C:
		}

		pool.cancel()
	})
	t.Run("check ping command (IO error)", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Mocking responses for PING command
		pingCmd := commands.NewSingleCommand(commands.PING)
		for _, rw := range rws[1:] {
			rw.On("SendSingleCommand", nil, pingCmd).
				Return(nil, mpdrw.IOError)
		}
		// Sleeping for longer than the ping interval
		time.Sleep(defaultConnectParams.pingInterval + (defaultConnectParams.pingInterval / 10))

		// Verifying that onDisconnect was called.
		select {
		case <-onDisconnectCalled:
		case <-time.NewTimer(time.Microsecond * 100).C:
			t.Error("onDisconnect was not called")
		}
	})
}

func TestImpl_SendSingleCommand(t *testing.T) {
	t.Run("sending command (no error)", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Mocking responses for LISTALL command
		cmd := commands.NewSingleCommand(commands.LISTALL)
		response := []string{"aaaa", "bbbb"}
		for _, rw := range rws[1:] {
			rw.On("SendSingleCommand", defaultConnectParams.requestContext, cmd).
				Return(response, nil)
		}
		actualResponse, err := pool.SendSingleCommand(defaultConnectParams.requestContext, cmd)
		assert.Nil(t, err)
		assert.NotNil(t, actualResponse, response)

		pool.cancel()
	})
	t.Run("sending command (ACK error)", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Mocking responses for LISTALL command
		cmd := commands.NewSingleCommand(commands.LISTALL)
		for _, rw := range rws[1:] {
			rw.On("SendSingleCommand", defaultConnectParams.requestContext, cmd).
				Return(nil, mpdrw.ACKError)
		}
		actualResponse, err := pool.SendSingleCommand(defaultConnectParams.requestContext, cmd)
		assert.Error(t, err)
		assert.Nil(t, actualResponse)
		// Verifying that onDisconnect was not called.
		select {
		case <-onDisconnectCalled:
			t.Error("onDisconnect was called")
		case <-time.NewTimer(time.Microsecond * 100).C:
		}

		pool.cancel()
	})
	t.Run("sending command (IO error)", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Mocking responses for LISTALL command
		cmd := commands.NewSingleCommand(commands.LISTALL)
		for _, rw := range rws[1:] {
			rw.On("SendSingleCommand", defaultConnectParams.requestContext, cmd).
				Return(nil, mpdrw.IOError)
		}
		actualResponse, err := pool.SendSingleCommand(defaultConnectParams.requestContext, cmd)
		assert.Error(t, err)
		assert.Nil(t, actualResponse)
		// Verifying that onDisconnect was called.
		select {
		case <-onDisconnectCalled:
		case <-time.NewTimer(time.Microsecond * 100).C:
			t.Error("onDisconnect was not called")
		}
	})
}

func TestImpl_SendBatchCommand(t *testing.T) {
	t.Run("sending batch command (no error)", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Mocking responses for batch command
		cmds := commands.NewBatchCommands(
			[]commands.SingleCommand{
				commands.NewSingleCommand(commands.PAUSE),
				commands.NewSingleCommand(commands.NEXT),
				commands.NewSingleCommand(commands.PLAY),
			}, 100)[0]
		for _, rw := range rws[1:] {
			rw.On("SendBatchCommand", defaultConnectParams.requestContext, cmds).
				Return(nil)
		}
		err = pool.SendBatchCommand(defaultConnectParams.requestContext, cmds)
		assert.Nil(t, err)

		pool.cancel()
	})
	t.Run("sending batch command (ACK error)", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Mocking responses for batch command
		cmds := commands.NewBatchCommands(
			[]commands.SingleCommand{
				commands.NewSingleCommand(commands.PAUSE),
				commands.NewSingleCommand(commands.NEXT),
				commands.NewSingleCommand(commands.PLAY),
			}, 100)[0]
		for _, rw := range rws[1:] {
			rw.On("SendBatchCommand", defaultConnectParams.requestContext, cmds).
				Return(mpdrw.ACKError)
		}
		err = pool.SendBatchCommand(defaultConnectParams.requestContext, cmds)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, SendCommandError)
		// Verifying that onDisconnect was not called.
		select {
		case <-onDisconnectCalled:
			t.Error("onDisconnect was called")
		case <-time.NewTimer(time.Microsecond * 100).C:
		}

		pool.cancel()
	})
	t.Run("sending batch command (IO error)", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return([]string{}, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Mocking responses for batch command
		cmds := commands.NewBatchCommands(
			[]commands.SingleCommand{
				commands.NewSingleCommand(commands.PAUSE),
				commands.NewSingleCommand(commands.NEXT),
				commands.NewSingleCommand(commands.PLAY),
			}, 100)[0]
		for _, rw := range rws[1:] {
			rw.On("SendBatchCommand", defaultConnectParams.requestContext, cmds).
				Return(mpdrw.IOError)
		}
		err = pool.SendBatchCommand(defaultConnectParams.requestContext, cmds)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, SendCommandError)
		// Verifying that onDisconnect was not called.
		select {
		case <-onDisconnectCalled:
		case <-time.NewTimer(time.Microsecond * 100).C:
			t.Error("onDisconnect was not called")
		}
	})
}

func TestImpl_IdleListener(t *testing.T) {
	t.Run("get IDLE event", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		idleResponse := []string{"aaaa", "bbbb"}
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return(idleResponse, nil)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnect := func() {}
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Subscribing for IDLE events
		idleSubscribeChannel := pool.Subscribe(time.Millisecond * 100)
		// Triggering the idleRW
		idleChan <- struct{}{}
		// Verifying receiving IDLE events
		select {
		case event := <-idleSubscribeChannel:
			assert.Equal(t, idleResponse, event)
		case <-time.NewTimer(time.Microsecond * 100).C:
			t.Error("No IDLE events were received")
		}

		pool.cancel()
	})
	t.Run("received IOError in IDLE listening", func(t *testing.T) {
		// Creating an mpdRW slice
		rws := make([]*mockMpdRW, defaultConnectParams.poolSize+1)
		for i := range rws {
			rws[i] = &mockMpdRW{}
		}
		idleChan := make(chan struct{})
		// The first element of the slice is idleRw. Mocking its behavior.
		rws[0].On("SendIdleCommand").Run(func(args mock.Arguments) {
			<-idleChan
		}).Return(nil, mpdrw.IOError)
		// Creating a mpdRWFactoryFunction
		mpdRWCounter := -1
		f := func() (mpdrw.MpdRW, error) {
			mpdRWCounter++
			return rws[mpdRWCounter], nil
		}
		onDisconnectCalled := make(chan struct{})
		onDisconnect := func() { onDisconnectCalled <- struct{}{} }
		// Creating an mpdRWPool
		pool, err := newMpdRWPool(f, defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.poolSize, defaultConnectParams.pingInterval, onDisconnect)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
		// Triggering the idleRW
		idleChan <- struct{}{}
		// Verifying that onDisconnect was called.
		select {
		case <-onDisconnectCalled:
		case <-time.NewTimer(time.Microsecond * 100).C:
			t.Error("onDisconnect was not called")
		}
	})
}
