package mpdclient

import (
	"errors"
	"github.com/anpotashev/mpdgo/internal/mpd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockMpdRWPoolFactory struct {
	mock.Mock
}

func (m *MockMpdRWPoolFactory) CreateAndConnect(host string, port uint16, password string) (*mpd.MpdRWPoolImpl, error) {
	args := m.Called(host, port, password)
	return args.Get(0).(*mpd.MpdRWPoolImpl), args.Error(1)
}

type MockMpdRwPool struct {
	mock.Mock
}

func (m *MockMpdRwPool) Disconnect() {
	m.Called()
}

func (m *MockMpdRwPool) SendCommand(command string) ([]string, error) {
	args := m.Called(command)
	return args.Get(0).([]string), args.Error(1)
}

const (
	testHost     = "localhost"
	testPort     = uint16(6600)
	testPassword = "12345678"
)

func TestNewMpdGo(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		NewMpdClient(testHost, testPort, testPassword)
	})
}

func TestMpdClientImpl_Connect(t *testing.T) {
	mockMpdRwPool := new(MockMpdRwPool)
	t.Run("already connected", func(t *testing.T) {
		client := &MpdClientImpl{
			host:             testHost,
			port:             testPort,
			password:         testPassword,
			currentMpdRWPool: mockMpdRwPool,
		}
		err := client.Connect()
		assert.ErrorIs(t, err, AlreadyConnected)
	})

	t.Run("successful connection", func(t *testing.T) {
		mockMpdRWPoolFactory := new(MockMpdRWPoolFactory)
		mockMpdRWPoolFactory.On("CreateAndConnect", testHost, testPort, testPassword).
			Return(&mpd.MpdRWPoolImpl{}, nil)
		client := &MpdClientImpl{
			host:     testHost,
			port:     testPort,
			password: testPassword,
			factory:  mockMpdRWPoolFactory,
		}
		err := client.Connect()
		assert.NoError(t, err)
		assert.NotNil(t, client.currentMpdRWPool)
		assert.True(t, client.IsConnected())
		mockMpdRWPoolFactory.AssertExpectations(t)
	})

	t.Run("connection error", func(t *testing.T) {
		mockMpdRWPoolFactory := new(MockMpdRWPoolFactory)
		expectedError := ConnectionError
		mockMpdRWPoolFactory.On("CreateAndConnect", testHost, testPort, testPassword).
			Return(&mpd.MpdRWPoolImpl{}, errors.New("connection error"))
		client := &MpdClientImpl{
			host:     testHost,
			port:     testPort,
			password: testPassword,
			factory:  mockMpdRWPoolFactory,
		}
		err := client.Connect()
		assert.ErrorIs(t, err, expectedError)
		assert.Nil(t, client.currentMpdRWPool)
		mockMpdRWPoolFactory.AssertExpectations(t)
	})
}

func TestMpdClientImpl_Disconnect(t *testing.T) {
	t.Run("already connected", func(t *testing.T) {
		client := &MpdClientImpl{
			host:             testHost,
			port:             testPort,
			password:         testPassword,
			currentMpdRWPool: nil,
		}
		err := client.Disconnect()
		assert.ErrorIs(t, err, NotConnected)
	})
	t.Run("successful disconnect", func(t *testing.T) {
		mockChan := make(chan struct{}, 1)
		mockMpdRwPool := new(MockMpdRwPool)
		mockMpdRwPool.On("Disconnect").Run(func(args mock.Arguments) {
			mockChan <- struct{}{}
		})
		client := &MpdClientImpl{
			host:             testHost,
			port:             testPort,
			password:         testPassword,
			currentMpdRWPool: mockMpdRwPool,
		}
		err := client.Disconnect()
		assert.NoError(t, err)
		select {
		case <-mockChan:
		case <-time.After(time.Second):
			t.Fatal("Disconnect() was not called in the goroutine")
		}
		mockMpdRwPool.AssertExpectations(t)
	})
}

func TestMpdClientImpl_IsConnected(t *testing.T) {

}

type TestCommand string

func (t TestCommand) String() string {
	return string(t)
}

func TestMpdClientImpl_SendCommand(t *testing.T) {
	command := TestCommand("command")
	t.Run("not connected", func(t *testing.T) {
		client := MpdClientImpl{
			host:             testHost,
			port:             testPort,
			password:         testPassword,
			currentMpdRWPool: nil,
		}
		result, err := client.SendCommand(command)
		assert.ErrorIs(t, err, NotConnected)
		assert.Nil(t, result)
	})
	t.Run("successful send command", func(t *testing.T) {
		expectedAnswer := []string{"string1", "string2"}
		mockMpdRwPool := new(MockMpdRwPool)
		mockMpdRwPool.On("SendCommand", command.String()).Return(expectedAnswer, nil)
		client := MpdClientImpl{
			host:             testHost,
			port:             testPort,
			password:         testPassword,
			currentMpdRWPool: mockMpdRwPool,
		}
		result, err := client.SendCommand(command)
		assert.NoError(t, err)
		assert.Equal(t, expectedAnswer, result)
		mockMpdRwPool.AssertExpectations(t)
	})
	t.Run("connection error", func(t *testing.T) {
		mockMpdRwPool := new(MockMpdRwPool)
		mockMpdRwPool.On("SendCommand", command.String()).Return([]string{}, mpd.ConnectionError)
		mockChan := make(chan struct{}, 1)
		mockMpdRwPool.On("Disconnect").Run(func(args mock.Arguments) {
			mockChan <- struct{}{}
		})
		client := MpdClientImpl{
			host:             testHost,
			port:             testPort,
			password:         testPassword,
			currentMpdRWPool: mockMpdRwPool,
		}
		result, err := client.SendCommand(command)
		assert.ErrorIs(t, err, ConnectionError)
		assert.Nil(t, result)
		select {
		case <-mockChan:
		case <-time.After(time.Second):
			t.Fatal("Disconnect() was not called in the goroutine")
		}
		mockMpdRwPool.AssertExpectations(t)
	})
	t.Run("connection error and error on disconnect() call", func(t *testing.T) {
		var client MpdClientImpl
		mockMpdRwPool := new(MockMpdRwPool)
		mockMpdRwPool.On("SendCommand", command.String()).Run(func(args mock.Arguments) {
			client.currentMpdRWPool = nil
		}).
			Return([]string{}, mpd.ConnectionError)
		client = MpdClientImpl{
			host:             testHost,
			port:             testPort,
			password:         testPassword,
			currentMpdRWPool: mockMpdRwPool,
		}
		result, err := client.SendCommand(command)
		assert.ErrorIs(t, err, NotConnected)
		assert.Nil(t, result)
		mockMpdRwPool.AssertExpectations(t)
	})
	t.Run("not connection error", func(t *testing.T) {
		mockMpdRwPool := new(MockMpdRwPool)
		expectedError := errors.New("other error")
		mockMpdRwPool.On("SendCommand", command.String()).Return([]string{}, expectedError)
		client := MpdClientImpl{
			host:             testHost,
			port:             testPort,
			password:         testPassword,
			currentMpdRWPool: mockMpdRwPool,
		}
		result, err := client.SendCommand(command)
		assert.ErrorIs(t, err, expectedError)
		assert.Nil(t, result)
		mockMpdRwPool.AssertExpectations(t)
		mockMpdRwPool.AssertNotCalled(t, "Disconnect")
	})
}
