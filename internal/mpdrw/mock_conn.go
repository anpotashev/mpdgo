package mpdrw

import (
	"io"
	"net"
	"strings"
	"time"
)

// MockConn is a mock implementation of the [net.Conn] interface for connection
type MockConn struct {
	in  chan byte
	out chan byte
}

func (m *MockConn) mockOnRead(responses ...string) {
	joinedResponse := strings.Join(responses, "\n")
	if len(joinedResponse) > 0 {
		joinedResponse += "\n"
	}
	for _, b := range []byte(joinedResponse) {
		m.in <- b
	}
}

func (m *MockConn) readAllFromOutChan() string {
	var dataSentToWriter []byte
LABEL:
	for {
		select {
		case b := <-m.out:
			dataSentToWriter = append(dataSentToWriter, b)
		default:
			break LABEL
		}
	}
	return string(dataSentToWriter)
}

func (m *MockConn) Read(b []byte) (int, error) {
	first, ok := <-m.in
	if !ok {
		return 0, io.EOF
	}
	b[0] = first
	return 1, nil
}

func (m *MockConn) Write(b []byte) (int, error) {
	for _, elem := range b {
		m.out <- elem
	}
	return len(b), nil
}

func (m *MockConn) Close() error {
	close(m.in)
	close(m.out)
	return nil
}

func (m *MockConn) LocalAddr() net.Addr              { return nil }
func (m *MockConn) RemoteAddr() net.Addr             { return nil }
func (m *MockConn) SetDeadline(time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(time.Time) error { return nil }
