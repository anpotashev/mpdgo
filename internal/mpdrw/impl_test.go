package mpdrw

import (
	"context"
	"fmt"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

var defaultConnectParams = struct {
	requestContext context.Context
	ctx            context.Context
	password       string
	readTimeout    time.Duration
	pingPeriod     time.Duration
}{
	requestContext: context.Background(),
	ctx:            context.Background(),
	password:       "12345",
	readTimeout:    time.Millisecond * 100,
}

const version = "1.2.3"

func TestNewMpdRW(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		happyPathConnect(t, mockConn)
	})
	t.Run("happy path with nil request context", func(t *testing.T) {
		requestContext := defaultConnectParams.requestContext
		defaultConnectParams.requestContext = nil
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		happyPathConnect(t, mockConn)
		defaultConnectParams.requestContext = requestContext
	})
	t.Run("timeout waiting version answer", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		var mockDialer Dialer
		mockDialer = func() (net.Conn, error) { return mockConn, nil }
		rw, err := mockDialer.NewMpdRW(defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.password, defaultConnectParams.readTimeout)
		assert.Error(t, err)
		assert.ErrorIs(t, err, IOError)
		assert.Nil(t, rw)
		assert.Empty(t, mockConn.readAllFromOutChan())
	})
	t.Run("error receiving version answer", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		mockConn.mockOnRead("ACK error")
		var mockDialer Dialer
		mockDialer = func() (net.Conn, error) { return mockConn, nil }
		rw, err := mockDialer.NewMpdRW(defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.password, defaultConnectParams.readTimeout)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ACKError)
		assert.Nil(t, rw)
		assert.Empty(t, mockConn.readAllFromOutChan())
	})
	t.Run("incorrect password", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		responses := []string{
			fmt.Sprintf("OK MPD %s", version),
			"ACK incorrect password",
		}
		mockConn.mockOnRead(responses...)
		var mockDialer Dialer
		mockDialer = func() (net.Conn, error) { return mockConn, nil }
		rw, err := mockDialer.NewMpdRW(defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.password, defaultConnectParams.readTimeout)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ACKError)
		assert.Nil(t, rw)
		expectedDataSentToWriter := fmt.Sprintf("password \"%s\"\n", defaultConnectParams.password)
		assert.Equal(t, expectedDataSentToWriter, mockConn.readAllFromOutChan())
	})
}

func TestImpl_SendSingleCommand(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		mockConn.mockOnRead("first", "", "second", "OK")
		cmd := commands.NewSingleCommand(commands.PING)
		response, err := rw.SendSingleCommand(defaultConnectParams.requestContext, cmd)
		assert.NoError(t, err)
		assert.Equal(t, []string{"first", "second"}, response)
		assert.Equal(t, cmd.String(), mockConn.readAllFromOutChan())
	})
	t.Run("received ACK error", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		mockConn.mockOnRead("first", "second", "ACK error sending command")
		cmd := commands.NewSingleCommand(commands.PING)
		response, err := rw.SendSingleCommand(defaultConnectParams.requestContext, cmd)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ACKError)
		assert.Nil(t, response)
		assert.Equal(t, cmd.String(), mockConn.readAllFromOutChan())
	})
	t.Run("timeout waiting response", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		mockConn.mockOnRead("first", "second")
		cmd := commands.NewSingleCommand(commands.PING)
		response, err := rw.SendSingleCommand(defaultConnectParams.requestContext, cmd)
		assert.Error(t, err)
		assert.ErrorIs(t, err, IOError)
		assert.Nil(t, response)
		assert.Equal(t, cmd.String(), mockConn.readAllFromOutChan())
	})
}

func TestImpl_SendMultipleCommands(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		mockConn.mockOnRead("first", "second", "OK")
		batchCommand := prepareBatchCommand()
		err := rw.SendBatchCommand(defaultConnectParams.requestContext, batchCommand)
		assert.NoError(t, err)
		assert.Equal(t, batchCommand.String(), mockConn.readAllFromOutChan())
	})

	t.Run("received ACK error", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		mockConn.mockOnRead("first", "second", "ACK error sending command")
		batchCommand := prepareBatchCommand()
		err := rw.SendBatchCommand(defaultConnectParams.requestContext, batchCommand)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ACKError)
		assert.Equal(t, batchCommand.String(), mockConn.readAllFromOutChan())
	})
	t.Run("timeout waiting response", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		mockConn.mockOnRead("first", "second")
		batchCommand := prepareBatchCommand()
		err := rw.SendBatchCommand(defaultConnectParams.requestContext, batchCommand)
		assert.Error(t, err)
		assert.ErrorIs(t, err, IOError)
		assert.Equal(t, batchCommand.String(), mockConn.readAllFromOutChan())
	})
}

func prepareBatchCommand() commands.BatchCommand {
	var singleCommands []commands.SingleCommand
	for range 5 {
		singleCommands = append(singleCommands, commands.NewSingleCommand(commands.PING))
	}
	batchCommands := commands.NewBatchCommands(singleCommands, 10)
	batchCommand := batchCommands[0]
	return batchCommand
}

func TestImpl_SendIdleCommand(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		idleChan := make(chan []string)
		errChan := make(chan error)
		go func() {
			idleEvents, err := rw.SendIdleCommand()
			if err != nil {
				errChan <- err
				return
			}
			if idleEvents != nil {
				idleChan <- idleEvents
			}
		}()
		var idleEvents []string
		var idleError error
		select {
		case events := <-idleChan:

			idleEvents = events
		case err := <-errChan:
			idleError = err
		case <-time.NewTimer(time.Millisecond * 100).C:
			break
		}
		assert.NoError(t, idleError)
		assert.Nil(t, idleEvents)
		mockConn.mockOnRead("first", "second", "OK")
		select {
		case events := <-idleChan:

			idleEvents = events
		case err := <-errChan:
			idleError = err
		case <-time.NewTimer(time.Millisecond * 100).C:
			break
		}
		assert.NoError(t, idleError)
		assert.Equal(t, []string{"first", "second"}, idleEvents)
		assert.Equal(t, mockConn.readAllFromOutChan(), commands.NewSingleCommand(commands.IDLE).String())
	})
	t.Run("connection closed", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		idleChan := make(chan []string)
		errChan := make(chan error)
		go func() {
			idleEvents, err := rw.SendIdleCommand()
			if err != nil {
				errChan <- err
				return
			}
			if idleEvents != nil {
				idleChan <- idleEvents
			}
		}()
		var idleEvents []string
		var idleError error
		close(mockConn.in)
		select {
		case events := <-idleChan:

			idleEvents = events
		case err := <-errChan:
			idleError = err
		case <-time.NewTimer(time.Millisecond * 100).C:
			break
		}
		assert.NotNil(t, idleError)
		assert.ErrorIs(t, idleError, IOError)
		assert.Nil(t, idleEvents)
		assert.Equal(t, mockConn.readAllFromOutChan(), commands.NewSingleCommand(commands.IDLE).String())
	})

	t.Run("received an ACK response", func(t *testing.T) {
		mockConn := &MockConn{
			in:  make(chan byte, 1024),
			out: make(chan byte, 1024),
		}
		rw := happyPathConnect(t, mockConn)
		idleChan := make(chan []string)
		errChan := make(chan error)
		go func() {
			idleEvents, err := rw.SendIdleCommand()
			if err != nil {
				errChan <- err
				return
			}
			if idleEvents != nil {
				idleChan <- idleEvents
			}
		}()
		var idleEvents []string
		var idleError error
		mockConn.mockOnRead("first", "second", "ACK error")
		select {
		case events := <-idleChan:

			idleEvents = events
		case err := <-errChan:
			idleError = err
		case <-time.NewTimer(time.Millisecond * 100).C:
			break
		}
		assert.Error(t, idleError)
		assert.ErrorIs(t, idleError, ACKError)
		assert.Nil(t, idleEvents)
		assert.Equal(t, mockConn.readAllFromOutChan(), commands.NewSingleCommand(commands.IDLE).String())
	})
}

func happyPathConnect(t *testing.T, mockConn *MockConn) MpdRW {
	responses := []string{
		fmt.Sprintf("OK MPD %s", version),
		"OK",
	}
	mockConn.mockOnRead(responses...)
	var mockDialer Dialer
	mockDialer = func() (net.Conn, error) { return mockConn, nil }
	rw, err := mockDialer.NewMpdRW(defaultConnectParams.requestContext, defaultConnectParams.ctx, defaultConnectParams.password, defaultConnectParams.readTimeout)
	assert.NoError(t, err)
	assert.NotNil(t, rw)

	expectedDataSentToWriter := fmt.Sprintf("password \"%s\"\n", defaultConnectParams.password)
	assert.Equal(t, expectedDataSentToWriter, mockConn.readAllFromOutChan())
	return rw
}
