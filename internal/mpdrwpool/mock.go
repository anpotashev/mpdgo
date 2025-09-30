package mpdrwpool

import (
	"context"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/stretchr/testify/mock"
)

type mockMpdRW struct {
	mock.Mock
}

func (m *mockMpdRW) SendIdleCommand() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}
func (m *mockMpdRW) SendSingleCommand(requestContext context.Context, command commands.SingleCommand) ([]string, error) {
	args := m.Called(requestContext, command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}
func (m *mockMpdRW) SendBatchCommand(requestContext context.Context, command commands.BatchCommand) error {
	args := m.Called(requestContext, command)
	return args.Error(0)
}
