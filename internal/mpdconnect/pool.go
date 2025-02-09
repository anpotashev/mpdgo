package mpdconnect

import (
	"github.com/anpotashev/mpdgo/internal/commands"
)

type MpdRWPoolFactory interface {
	CreateAndConnect(host string, port uint16, password string) (*MpdConnectImpl, error)
}

type MpdConnectFactoryImpl struct{}

type MpdConnect interface {
	Disconnect()
	SendCommand(command commands.MpdCommand) ([]string, error)
}

type MpdConnectImpl struct {
}

func (p *MpdConnectImpl) Disconnect() {

}

func (p *MpdConnectImpl) SendCommand(command commands.MpdCommand) ([]string, error) {
	return nil, nil
}

func (m MpdConnectFactoryImpl) CreateAndConnect(host string, port uint16, password string) (*MpdConnectImpl, error) {
	return &MpdConnectImpl{}, nil
}
