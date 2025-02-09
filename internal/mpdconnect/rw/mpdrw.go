package rw

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/anpotashev/mpdgo/internal/commands"
	"net"
	"regexp"
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
	defaultTimeout   = time.Second
	maxCommandsCount = 100
)

var (
	WrongAnswerFromServerError = errors.New("wrong answer from mpd server")
	EmptyCommandsList          = errors.New("empty commands list")
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
		return "", WrongAnswerFromServerError
	}
	return strings.TrimPrefix(versionAnswer, "OK MPD "), nil
}

func (rw *MpdRWImpl) sendPassword(password string) error {
	if len(password) == 0 {
		return nil
	}
	passwordCommand := commands.NewSingleCommand(commands.PASSWORD)
	passwordCommand.AddParams(password)
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

func (rw *MpdRWImpl) parseError(answer string, command string) error {
	re := regexp.MustCompile(fmt.Sprintf(`ACK\ .*\{%s\} (.+)`, regexp.QuoteMeta(command)))
	matches := re.FindStringSubmatch(answer)
	if len(matches) > 1 {
		return newCommandError(command, matches[1])
	}
	return newCommandError(command, fmt.Sprintf("Unparseable answer: %s", answer))
}

func (rw *MpdRWImpl) SendCommand(command *commands.SingleCommand) ([]string, error) {
	_, err := rw.readerWriter.WriteString(command.String())
	if err != nil {
		return nil, err
	}
	if err = rw.readerWriter.Flush(); err != nil {
		return nil, err
	}
	timeout := defaultTimeout
	if command.String() == commands.NewSingleCommand(commands.IDLE).String() {
		timeout = time.Duration(0)
	}
	var result []string
	for {
		line, err := rw.readWithTimeout(timeout)
		if err != nil {
			return nil, err
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
	}
}

func isAnswerEnded(s string) (bool, error) {
	if s == "OK" {
		return true, nil
	}
	if strings.HasPrefix(s, "ACK") {
		re := regexp.MustCompile(`ACK .*\{(.*)\} (.+)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			command := matches[1]
			errorMessage := matches[2]
			return true, newCommandError(command, errorMessage)
		}
		return true, newCommandError("unexpected answer", s)
	}
	return false, nil
}

func (rw *MpdRWImpl) SendBatchCommand(command *commands.BatchCommand) error {
	_, err := rw.readerWriter.WriteString(command.String())
	if err != nil {
		return err
	}
	if err = rw.readerWriter.Flush(); err != nil {
		return err
	}
	for {
		line, err := rw.readWithTimeout(defaultTimeout)
		if err != nil {
			return err
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
