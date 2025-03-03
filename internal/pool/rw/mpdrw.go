package rw

import (
	"bufio"
	"context"
	"fmt"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/rs/zerolog/log"
	"net"
	"strings"
	"time"
)

type MpdRW interface {
	SendCommand(command *commands.SingleCommand) ([]string, error)
	SendBatchCommand(command *commands.BatchCommand) error
}

type MpdRWImpl struct {
	conn         net.Conn
	readerWriter *bufio.ReadWriter
	version      string
	cancel       context.CancelFunc
}

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type RealDialer struct{}

func (d *RealDialer) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

const (
	//  defaultTimeout - стандартное время ожидания ответа, после отправки команды
	// используется для всех команд, кроме idle
	defaultTimeout = time.Second
)

func NewMpdRW(ctx context.Context, host string, port uint16, password string) (*MpdRWImpl, error) {
	return newMpdRwImpl(&RealDialer{}, ctx, host, port, password)
}

func newMpdRwImpl(dialer Dialer, ctx context.Context, host string, port uint16, password string) (*MpdRWImpl, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	rw := bufio.NewReadWriter(r, w)
	ctx1, cancel := context.WithCancel(ctx)
	goroutineStarted := make(chan struct{})
	go func() {
		goroutineStarted <- struct{}{}
		defer conn.Close()
		<-ctx1.Done()
	}()
	<-goroutineStarted
	result := &MpdRWImpl{
		conn:         conn,
		readerWriter: rw,
		cancel:       cancel,
	}
	version, err := result.getVersion()
	if err != nil {
		cancel()
		return nil, err
	}
	result.version = version
	if err := result.sendPassword(password); err != nil {
		return nil, err
	}
	return result, nil
}

func (rw *MpdRWImpl) getVersion() (string, error) {
	versionAnswer, err := rw.readWithTimeout(defaultTimeout)
	if err != nil {
		return "", err
	}
	versionAnswer = strings.TrimSuffix(versionAnswer, "\n")
	if !strings.HasPrefix(versionAnswer, "OK MPD ") {
		return "", ServerError
	}
	return strings.TrimPrefix(versionAnswer, "OK MPD "), nil
}

func (rw *MpdRWImpl) sendPassword(password string) error {
	if len(password) == 0 {
		return nil
	}
	passwordCommand := commands.NewSingleCommand(commands.PASSWORD).
		AddParams(password)
	_, err := rw.SendCommand(passwordCommand)
	return err
}

func (rw *MpdRWImpl) readWithTimeout(timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	if time.Duration(0) == timeout {
		deadline = time.Time{}
	}
	err := rw.conn.SetReadDeadline(deadline)
	if err != nil {
		return "", err
	}
	return rw.readerWriter.ReadString('\n')
}

func (rw *MpdRWImpl) SendCommand(command *commands.SingleCommand) ([]string, error) {
	log.Debug().Str("command", command.String()).Msg("Sending command")
	_, err := rw.readerWriter.WriteString(command.String())
	if err != nil {
		return nil, wrapIntoServerError(err)
	}
	if err = rw.readerWriter.Flush(); err != nil {
		return nil, wrapIntoServerError(err)
	}
	timeout := defaultTimeout
	if command.String() == commands.NewSingleCommand(commands.IDLE).String() {
		timeout = time.Duration(0)
	}
	var result []string
	for {
		line, err := rw.readWithTimeout(timeout)
		if err != nil {
			return nil, wrapIntoServerError(err)
		}
		line = strings.TrimSuffix(line, "\n")
		if len(line) == 0 {
			continue
		}
		isEnded, err := isAnswerEnded(line)
		if isEnded {
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		result = append(result, line)
	}
}

func isAnswerEnded(s string) (bool, error) {
	if s == "OK" {
		return true, nil
	}
	if strings.HasPrefix(s, "ACK") {
		return true, parseACKAnswer(s)
	}
	return false, nil
}

func (rw *MpdRWImpl) SendBatchCommand(command *commands.BatchCommand) error {
	_, err := rw.readerWriter.WriteString(command.String())
	if err != nil {
		return wrapIntoServerError(err)
	}
	if err = rw.readerWriter.Flush(); err != nil {
		return wrapIntoServerError(err)
	}
	for {
		line, err := rw.readWithTimeout(defaultTimeout)
		if err != nil {
			return wrapIntoServerError(err)
		}
		line = strings.TrimSuffix(line, "\n")
		if len(line) == 0 {
			continue
		}
		isEnded, err := isAnswerEnded(line)
		if isEnded {
			return err
		}
	}
}
