package mpdclient

import (
	"context"
	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/commands"
	"github.com/stretchr/testify/mock"
)

type mockMpdRWPool struct {
	mock.Mock
	observer.Observer[[]string]
}

func (m *mockMpdRWPool) SendSingleCommand(requestContext context.Context, command commands.SingleCommand) ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), nil
	}
	return nil, args.Error(1)
}

func (m *mockMpdRWPool) SendBatchCommand(requestContext context.Context, command commands.BatchCommand) error {
	return m.Called().Error(0)
}
