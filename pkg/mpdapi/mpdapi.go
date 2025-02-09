package mpdapi

import (
	"context"
	"github.com/anpotashev/mpdgo/internal/client"
)

type MpdApi interface {
	Player
	CurrentPlaylist
	StoredPlaylists
	Settings
	Outputs
	Tree
	Connect() error
	Disconnect() error
	IsConnected() bool
}

type MpdApiImpl struct {
	mpdClient client.MpdClient
}

func NewMpdApi(ctx context.Context, host string, port uint16, password string) (*MpdApiImpl, error) {
	mpdClient, err := client.NewMpdClient(ctx, host, port, password)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	return &MpdApiImpl{mpdClient: mpdClient}, nil
}

func (api *MpdApiImpl) Connect() error {
	return wrapPkgError(api.mpdClient.Connect())
}

func (api *MpdApiImpl) Disconnect() error {
	return wrapPkgError(api.mpdClient.Disconnect())
}

func (api *MpdApiImpl) IsConnected() bool {
	return api.mpdClient.IsConnected()
}
