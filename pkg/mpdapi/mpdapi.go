package mpdapi

import (
	"context"
	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/client"
)

type MpdApi interface {
	Player
	CurrentPlaylist
	StoredPlaylists
	Settings
	Outputs
	Tree
	observer.Subscriber[MpdEventType]
	Connect() error
	Disconnect() error
	IsConnected() bool
}

type Impl struct {
	mpdClient client.MpdClient
	observer.Observer[MpdEventType]
	ctx context.Context
}

func NewMpdApi(ctx context.Context, host string, port uint16, password string, useCache bool) (MpdApi, error) {
	mpdClient, err := client.NewMpdClient(ctx, host, port, password)
	if err != nil {
		return nil, wrapPkgError(err)
	}
	result := &Impl{mpdClient: mpdClient, ctx: ctx, Observer: observer.New[MpdEventType]()}
	result.initObserver()
	if useCache {
		return newWithCache(result), nil
	}
	return result, nil
}

func (api *Impl) Connect() error {
	return wrapPkgError(api.mpdClient.Connect())
}

func (api *Impl) Disconnect() error {
	return wrapPkgError(api.mpdClient.Disconnect())
}

func (api *Impl) IsConnected() bool {
	return api.mpdClient.IsConnected()
}
