package rw

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"os"
	"testing"
	"time"
)

func (m *MockDialer) Dial(network, address string) (net.Conn, error) {
	return m.conn, m.err
}

type MockDialer struct {
	conn net.Conn
	err  error
}

const (
	host     = "localhost"
	port     = 6600
	password = "pa55w0rd"
	version  = "1.2.3"
)

func TestMpdRwImpl(t *testing.T) {
	t.Run("happy pass connect and MpdRWImp.cancel()", func(t *testing.T) {
		ctx := context.Background()
		mpdRW := mockAndCreateMpdImpl(t, ctx)
		closeChan := make(chan interface{})
		mpdRW.conn.(*MockConn).On("Close", mock.Anything).Run(
			func(args mock.Arguments) {
				closeChan <- struct{}{}
			}).Return(nil)
		mpdRW.cancel()
		select {
		case <-closeChan:
			assert.True(t, true)
		case <-time.After(time.Second):
			t.Fatal("Timeout: Close was not called within 1 second")
		}
	})
	t.Run("happy pass connect and cancel on parent context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		mpdRW := mockAndCreateMpdImpl(t, ctx)
		closeChan := make(chan interface{})
		mpdRW.conn.(*MockConn).On("Close", mock.Anything).Run(
			func(args mock.Arguments) {
				closeChan <- struct{}{}
			}).Return(nil)
		cancel()
		select {
		case <-closeChan:
			assert.True(t, true)
		case <-time.After(time.Second):
			t.Fatal("Timeout: Close was not called within 1 second")
		}
	})
	t.Run("error on getting error during reading", func(t *testing.T) {
		ctx := context.Background()
		mockConn := new(MockConn)
		dialer := &MockDialer{
			conn: mockConn,
			err:  nil,
		}
		mpdAnswer := []byte("Wrong answer\n")
		mockConn.On("Read", mock.Anything).
			Return(len(mpdAnswer), nil).Run(func(args mock.Arguments) {
			copy(args.Get(0).([]byte), mpdAnswer)
		})
		closeChan := make(chan interface{})
		mockConn.On("Close", mock.Anything).Run(
			func(args mock.Arguments) {
				closeChan <- struct{}{}
			}).Return(nil)
		mpdRw, err := newMpdRwImpl(dialer, ctx, host, port, password)
		assert.ErrorIs(t, err, ServerError)
		assert.Nil(t, mpdRw)
		select {
		case <-closeChan:
			assert.True(t, true)
		case <-time.After(time.Second):
			t.Fatal("Timeout: Close was not called within 1 second")
		}
	})
	t.Run("timeout waiting version answer", func(t *testing.T) {
		ctx := context.Background()
		mockConn := new(MockConn)
		dialer := &MockDialer{
			conn: mockConn,
			err:  nil,
		}
		mockConn.On("Read", mock.Anything).
			Return(0, os.ErrDeadlineExceeded)
		mockConn.On("Close", mock.Anything).Return(nil)
		mpdRw, err := newMpdRwImpl(dialer, ctx, host, port, password)
		assert.ErrorIs(t, err, os.ErrDeadlineExceeded)
		assert.Nil(t, mpdRw)
	})
	t.Run("wrong password", func(t *testing.T) {
		errMsg := "incorrect password"
		ctx := context.Background()
		mockConn := new(MockConn)
		dialer := &MockDialer{
			conn: mockConn,
			err:  nil,
		}
		mpdAnswer := []byte(fmt.Sprintf("OK MPD %s\nACK [3@0] {password} %s\n", version, errMsg))
		mockConn.On("Read", mock.Anything).
			Return(len(mpdAnswer), nil).Run(func(args mock.Arguments) {
			copy(args.Get(0).([]byte), mpdAnswer)
		})
		mockConn.On("Write", mock.Anything).
			Return(nil)
		mockConn.On("Close", mock.Anything).Return(nil)
		mpdRw, err := newMpdRwImpl(dialer, ctx, host, port, password)
		expectedError := &CommandError{
			command:      "password",
			errorMessage: errMsg,
		}
		assert.EqualError(t, err, expectedError.Error())
		assert.Nil(t, mpdRw)
	})
	t.Run("error reading answer on sending password", func(t *testing.T) {
		ctx := context.Background()
		mockConn := new(MockConn)
		dialer := &MockDialer{
			conn: mockConn,
			err:  nil,
		}
		mpdAnswer := []byte(fmt.Sprintf("OK MPD %s\n", version))
		mockConn.On("Read", mock.Anything).
			Return(len(mpdAnswer), nil).Run(func(args mock.Arguments) {
			copy(args.Get(0).([]byte), mpdAnswer)
		}).Once()
		mockConn.On("Read", mock.Anything).
			Return(0, os.ErrDeadlineExceeded).Once()
		mockConn.On("Write", mock.Anything).
			Return(nil)
		mockConn.On("Close", mock.Anything).Return(nil)
		mpdRw, err := newMpdRwImpl(dialer, ctx, host, port, password)
		assert.ErrorIs(t, err, os.ErrDeadlineExceeded)
		assert.Nil(t, mpdRw)
	})
	t.Run("error opening connect on realdealer", func(t *testing.T) {
		ctx := context.Background()
		_, err := NewMpdRW(ctx, "wrong-host.111", 1, password)
		assert.ErrorContains(t, err, "no such host")
	})
}

func mockAndCreateMpdImpl(t *testing.T, ctx context.Context) *MpdRWImpl {
	mockConn := new(MockConn)
	dialer := &MockDialer{
		conn: mockConn,
		err:  nil,
	}
	mpdAnswer := []byte(fmt.Sprintf("OK MPD %s\nOK\n", version))
	mockConn.On("Read", mock.Anything).
		Return(len(mpdAnswer), nil).Run(func(args mock.Arguments) {
		copy(args.Get(0).([]byte), mpdAnswer)
	})
	mockConn.On("Write", mock.Anything).
		Return(nil)
	mpdRw, err := newMpdRwImpl(dialer, ctx, host, port, password)
	assert.NoError(t, err)
	assert.NotNil(t, mpdRw)
	assert.Equal(t, mpdRw.version, version)
	mockConn.AssertExpectations(t)
	return mpdRw
}
