package mpdrw

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/anpotashev/mpdgo/internal/commands"
	log "github.com/anpotashev/mpdgo/internal/logger"
	"github.com/google/uuid"
)

type Impl struct {
	rw          *bufio.ReadWriter
	readTimeout time.Duration
}

type Dialer func() (net.Conn, error)

func NewDialer(host string, port uint16) Dialer {
	return func() (net.Conn, error) {
		return net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	}
}

func (d Dialer) NewMpdRW(requestContext, ctx context.Context, password string, readTimeout time.Duration) (MpdRW, error) {
	return newMpdRW(requestContext, ctx, d, password, readTimeout)
}

func newMpdRW(requestContext, ctx context.Context, dialer Dialer, password string, readTimeout time.Duration) (*Impl, error) {
	log.DebugContext(requestContext, "Connecting to mpd")
	log.DebugContext(requestContext, "Dialing")
	conn, err := dialer()
	if err != nil {
		return nil, errors.Join(errors.Join(ErrIO, err), err)
	}
	log.DebugContext(requestContext, "Creating reader and writer")
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	impl := &Impl{
		rw:          bufio.NewReadWriter(r, w),
		readTimeout: readTimeout,
	}
	log.DebugContext(requestContext, "Starting a goroutine that closes the connection on ctx.Done")
	go func() {
		defer conn.Close()
		<-ctx.Done()
	}()
	log.DebugContext(requestContext, "Listening version")
	_, err = impl.readAnswerWithTimeout(requestContext)
	if err != nil {
		log.DebugContext(requestContext, "Error reading answer: %v", err)
		return nil, err
	}
	if password != "" {
		log.DebugContext(requestContext, "Authenticating with password")
		passwordCommand := commands.NewSingleCommand(commands.PASSWORD).AddParams(password)
		_, err := impl.SendSingleCommand(requestContext, passwordCommand)
		if err != nil {
			return nil, err
		}
	}
	return impl, nil
}

func (m *Impl) SendIdleCommand() ([]string, error) {
	idleCommandContext := context.Background()
	commandUUID, _ := uuid.NewUUID()
	idleCommandContext = context.WithValue(idleCommandContext, "command_id", commandUUID.String())
	log.DebugContext(idleCommandContext, "Sending idle command")
	command := commands.NewSingleCommand(commands.IDLE)
	log.DebugContext(idleCommandContext, "Writing the command")
	_, err := m.rw.WriteString(command.String())
	if err != nil {
		return nil, errors.Join(err, errors.Join(ErrIO, err))
	}
	log.DebugContext(idleCommandContext, "Flushing the writer")
	err = m.rw.Flush()
	if err != nil {
		log.ErrorContext(idleCommandContext, "Flushing the writer error", "err", err)
		return nil, errors.Join(err, errors.Join(ErrIO, err))
	}
	log.DebugContext(idleCommandContext, "Creating answer and error channels")
	answerChan := make(chan []string)
	errorChan := make(chan error)
	log.DebugContext(idleCommandContext, "Starting a goroutine that reads answer")
	//lint:ignore SA1012 ignore
	idleCommandContext, cancel := context.WithCancel(idleCommandContext)
	defer cancel()
	go m.readAnswer(idleCommandContext, answerChan, errorChan, nil)
	select {
	case answer := <-answerChan:
		log.DebugContext(idleCommandContext, "Got answer in the answer channel", "answer", log.Truncate(strings.Join(answer, "\n"), 100))
		return answer, nil
	case err := <-errorChan:
		log.DebugContext(idleCommandContext, "Got answer in the error channel", "err", err)
		return nil, err
	}
}

func (m *Impl) SendSingleCommand(requestContext context.Context, command commands.SingleCommand) ([]string, error) {
	return m.sendCommand(requestContext, &command)
}

func (m *Impl) SendBatchCommand(requestContext context.Context, command commands.BatchCommand) error {
	_, err := m.sendCommand(requestContext, command)
	return err
}

func (m *Impl) sendCommand(requestContext context.Context, command commands.MpdCommand) ([]string, error) {
	if requestContext == nil {
		requestContext = context.Background()
	}
	commandUUID, _ := uuid.NewUUID()
	requestContext = context.WithValue(requestContext, "command_id", commandUUID.String())
	log.DebugContext(requestContext, "Sending command", "command", command.String())
	_, err := m.rw.WriteString(command.String())
	if err != nil {
		return nil, errors.Join(errors.Join(ErrIO, err), err)
	}
	log.DebugContext(requestContext, "Flushing the writer")
	err = m.rw.Flush()
	if err != nil {
		return nil, errors.Join(errors.Join(ErrIO, err), err)
	}
	log.DebugContext(requestContext, "Waiting the answer")
	return m.readAnswerWithTimeout(requestContext)
}

func (m *Impl) readAnswerWithTimeout(requestContext context.Context) ([]string, error) {
	log.DebugContext(requestContext, "Creating answer and error channels")
	answerChan := make(chan []string)
	errorChan := make(chan error)
	log.DebugContext(requestContext, "Creation the timer")
	timer := time.NewTimer(m.readTimeout)
	log.DebugContext(requestContext, "Starting a goroutine that reads an answer")
	requestContext, cancel := context.WithCancel(requestContext)
	defer cancel()
	go m.readAnswer(requestContext, answerChan, errorChan, timer)
	select {
	case answer := <-answerChan:
		log.DebugContext(requestContext, "Received data from the answer channel", "answer", log.Truncate(strings.Join(answer, "\n"), 100))
		return answer, nil
	case err := <-errorChan:
		log.DebugContext(requestContext, "Received data from the error channel", "err", err)
		return nil, err
	case <-timer.C:
		log.DebugContext(requestContext, "Timeout")
		return nil, errors.Join(ErrIO, fmt.Errorf("timeout reading the answer"))
	}
}

func (m *Impl) readAnswer(requestContext context.Context, readChan chan []string, errorChan chan error, timer *time.Timer) {
	log.DebugContext(requestContext, "Starting reading the answer")
	var result []string
	for {
		line, err := m.rw.ReadString('\n')
		if err != nil {
			select {
			case errorChan <- errors.Join(ErrIO, err):
			default: // non-blocking send
			}
			return
		}
		line = strings.TrimSuffix(line, "\n")
		if len(line) == 0 {
			continue
		}
		isEnded, err := isAnswerEnded(line)
		if isEnded {
			if err != nil {
				select {
				case errorChan <- err:
				default: // non-blocking send
				}
				log.DebugContext(requestContext, "stop reading answers (error case)")
				return
			}
			select {
			case readChan <- result:
			case <-requestContext.Done():
			}
			log.DebugContext(requestContext, "stop reading answers (success case)")
			return
		}
		result = append(result, line)
		if timer != nil {
			timer.Reset(m.readTimeout)
		}
	}
}

func isAnswerEnded(line string) (bool, error) {
	if strings.HasPrefix(line, "OK") {
		return true, nil
	}
	if strings.HasPrefix(line, "ACK ") {
		return true, parseACKAnswer(line)
	}
	return false, nil
}
