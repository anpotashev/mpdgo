package mpdclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/anpotashev/mpdgo/internal/mpdrwpool"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// implementation newMpdRWPoolFactory
func newMockMpdRWPoolFactory(requestContext, ctx context.Context, onDisconnect func()) (mpdrwpool.MpdRWPool, error) {
	return &mockMpdRWPool{Observer: observer.New[[]string]()}, nil
}

var defaultClientParams = struct {
	ctx                   context.Context
	host                  string
	port                  uint16
	password              string
	maxBatchCommandLength uint16
	poolSize              uint8
	readTimeout           time.Duration
	pingTimeout           time.Duration
}{
	ctx:                   context.Background(),
	host:                  "localhost",
	port:                  6600,
	password:              "1234",
	maxBatchCommandLength: 100,
	poolSize:              3,
	readTimeout:           time.Millisecond * 100,
	pingTimeout:           time.Second * 10,
}

func createClientWithDefaultValues() *Impl {
	return NewMpdClientImpl(
		defaultClientParams.ctx,
		defaultClientParams.host,
		defaultClientParams.port,
		defaultClientParams.password,
		defaultClientParams.maxBatchCommandLength,
		defaultClientParams.poolSize,
		defaultClientParams.readTimeout,
		defaultClientParams.pingTimeout,
	)
}

func TestImpl_Connect(t *testing.T) {
	t.Run("establishes connection once, returns error on repeated connect", func(t *testing.T) {
		client := createClientWithDefaultValues()
		connectTestClient(t, client)
		// checking connect call, when already connected
		err := client.connect(context.Background(), newMockMpdRWPoolFactory)
		assert.Error(t, err)
		assert.ErrorIs(t, err, AlreadyConnected)
		client.cancelFunc()
	})
	t.Run("error on creating mpdRWPool", func(t *testing.T) {
		pool := createClientWithDefaultValues()
		requestContext := context.Background()
		f := func(requestContext, ctx context.Context, onDisconnect func()) (mpdrwpool.MpdRWPool, error) {
			return nil, errors.New("some error")
		}
		err := pool.connect(requestContext, f)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ConnectionError)
		assert.Nil(t, pool.pool)
	})
}

func TestImpl_Disconnect(t *testing.T) {
	t.Run("Successful disconnect once, returns error on repeated disconnect", func(t *testing.T) {
		client := createClientWithDefaultValues()
		connectTestClient(t, client)
		err := client.Disconnect(context.Background())
		assert.NoError(t, err)
		err = client.Disconnect(context.Background())
		assert.Error(t, err)
		assert.ErrorIs(t, err, NotConnected)
	})
}

func TestImpl_IsConnected(t *testing.T) {
	t.Run("check connection state", func(t *testing.T) {
		client := createClientWithDefaultValues()
		state := client.IsConnected(context.Background())
		assert.False(t, state)
		connectTestClient(t, client)
		state = client.IsConnected(context.Background())
		assert.True(t, state)
		err := client.Disconnect(context.Background())
		assert.NoError(t, err)
		state = client.IsConnected(context.Background())
		assert.False(t, state)
	})
}

func TestImpl_Events(t *testing.T) {
	t.Run("getting events on_connect, on_disconnect, idle-events", func(t *testing.T) {
		client := createClientWithDefaultValues()
		subscribeChan := client.Subscribe(time.Millisecond * 100)
		connectTestClient(t, client)
		select {
		case event := <-subscribeChan:
			assert.Equal(t, OnConnect, event)
		case <-time.After(time.Millisecond * 100):
			t.Errorf("%s event was not received", OnConnect)
		}
		idleEvents := []string{"aaa", "bbb"}
		client.pool.Notify(idleEvents)
		var receivedIdleEvent []string
		for range idleEvents {
			select {
			case event := <-subscribeChan:
				receivedIdleEvent = append(receivedIdleEvent, event)
			case <-time.After(time.Millisecond * 100):
				t.Error("expected IDLE event was not received")
			}
		}
		assert.ElementsMatch(t, idleEvents, receivedIdleEvent)
		client.Disconnect(context.Background())
		select {
		case event := <-subscribeChan:
			assert.Equal(t, OnDisconnect, event)
		case <-time.After(time.Millisecond * 100):
			t.Errorf("%s event was not received", OnDisconnect)
		}
	})
	t.Run("getting events on_disconnect after cancel call", func(t *testing.T) {
		// todo client disconnects on callback onDisconnect
	})
}

func TestImpl_SendSingleCommand(t *testing.T) {
	t.Run("send single command. No error", func(t *testing.T) {
		client := createClientWithDefaultValues()
		connectTestClient(t, client)
		pool := client.pool.(*mockMpdRWPool)
		cmd := commands.NewSingleCommand(commands.PLAY)
		response := []string{"aaa", "bbb"}
		pool.On("SendSingleCommand").Return(response, nil)
		actual, err := client.SendSingleCommand(context.Background(), cmd)
		assert.NoError(t, err)
		assert.Equal(t, response, actual)
		client.cancelFunc()
	})
	t.Run("send single command. Error", func(t *testing.T) {
		client := createClientWithDefaultValues()
		connectTestClient(t, client)
		pool := client.pool.(*mockMpdRWPool)
		cmd := commands.NewSingleCommand(commands.PLAY)
		pool.On("SendSingleCommand").Return(nil, fmt.Errorf("error"))
		actual, err := client.SendSingleCommand(context.Background(), cmd)
		assert.Error(t, err)
		assert.ErrorIs(t, err, CommandSendError)
		assert.Nil(t, actual)
	})
	t.Run("send single command when not connected", func(t *testing.T) {
		client := createClientWithDefaultValues()
		cmd := commands.NewSingleCommand(commands.PLAY)
		actual, err := client.SendSingleCommand(context.Background(), cmd)
		assert.Error(t, err)
		assert.ErrorIs(t, err, NotConnected)
		assert.Nil(t, actual)
	})
}

func TestImpl_SendBatchCommand(t *testing.T) {
	t.Run("send single command. No error", func(t *testing.T) {
		client := createClientWithDefaultValues()
		connectTestClient(t, client)
		pool := client.pool.(*mockMpdRWPool)
		cmds := []commands.SingleCommand{
			commands.NewSingleCommand(commands.PLAY),
			commands.NewSingleCommand(commands.STOP)}
		pool.On("SendBatchCommand").Return(nil)
		err := client.SendBatchCommand(context.Background(), cmds)
		assert.NoError(t, err)
		client.cancelFunc()
	})
	t.Run("send single command. Error", func(t *testing.T) {
		client := createClientWithDefaultValues()
		connectTestClient(t, client)
		pool := client.pool.(*mockMpdRWPool)
		cmds := []commands.SingleCommand{
			commands.NewSingleCommand(commands.PLAY),
			commands.NewSingleCommand(commands.STOP)}
		pool.On("SendBatchCommand").Return(fmt.Errorf("error"))
		err := client.SendBatchCommand(context.Background(), cmds)
		assert.Error(t, err)
		assert.ErrorIs(t, err, CommandSendError)
		client.cancelFunc()
	})
	t.Run("send batch command when not connected", func(t *testing.T) {
		client := createClientWithDefaultValues()
		cmds := []commands.SingleCommand{
			commands.NewSingleCommand(commands.PLAY),
			commands.NewSingleCommand(commands.STOP)}
		err := client.SendBatchCommand(context.Background(), cmds)
		assert.Error(t, err)
		assert.ErrorIs(t, err, NotConnected)
	})
}

func connectTestClient(t *testing.T, client *Impl) {
	requestContext := context.Background()
	// checking connect call
	err := client.connect(requestContext, newMockMpdRWPoolFactory)
	assert.NoError(t, err)
	assert.NotNil(t, client.pool)
}
