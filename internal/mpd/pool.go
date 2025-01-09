package mpd

type MpdRWPoolFactory interface {
	CreateAndConnect(host string, port uint16, password string) (*MpdRWPoolImpl, error)
}

type MpdRWPoolFactoryImpl struct{}

type MpdRWPool interface {
	Disconnect()
	SendCommand(command string) ([]string, error)
}

type MpdRWPoolImpl struct {
}

func (p *MpdRWPoolImpl) Disconnect() {

}

func (p *MpdRWPoolImpl) SendCommand(command string) ([]string, error) {
	return nil, nil
}

func (m MpdRWPoolFactoryImpl) CreateAndConnect(host string, port uint16, password string) (*MpdRWPoolImpl, error) {
	return &MpdRWPoolImpl{}, nil
}
