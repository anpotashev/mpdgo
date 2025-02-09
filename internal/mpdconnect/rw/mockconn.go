package rw

import (
	"github.com/stretchr/testify/mock"
	"net"
	"time"
)

type MockConn struct {
	mock.Mock
}

func (m *MockConn) Read(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Write(b []byte) (int, error) {
	args := m.Called(b)
	if err := args.Error(0); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (m *MockConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConn) LocalAddr() net.Addr              { return nil }
func (m *MockConn) RemoteAddr() net.Addr             { return nil }
func (m *MockConn) SetDeadline(time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(time.Time) error { return nil }
