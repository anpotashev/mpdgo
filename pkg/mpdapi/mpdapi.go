package mpdapi

import (
	"context"
	"github.com/anpotashev/go-observer/pkg/observer"
	"github.com/anpotashev/mpdgo/internal/logger"
	"github.com/anpotashev/mpdgo/internal/mpdclient"
	"log/slog"
	"time"
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
	WithRequestContext(ctx context.Context) MpdApi
}

type Impl struct {
	mpdClient mpdclient.MpdClient
	observer.Observer[MpdEventType]
	ctx            context.Context
	requestContext context.Context
}

func NewMpdApi(ctx context.Context, host string, port uint16, password string, useCache bool, maxBatchCommandLength uint16, poolSize uint8, pingPeriod, pingTimeout time.Duration) (MpdApi, error) {
	mpdClient := mpdclient.NewMpdClientImpl(ctx, host, port, password, maxBatchCommandLength, poolSize, pingPeriod, pingTimeout)
	result := &Impl{mpdClient: mpdClient, ctx: ctx, Observer: observer.New[MpdEventType](), requestContext: context.Background()}
	result.initObserver()
	if useCache {
		return newWithCache(result), nil
	}
	return result, nil
}

func SetLogger(l *slog.Logger) {
	logger.Init(l)
}

func (api *Impl) WithRequestContext(ctx context.Context) MpdApi {
	return &Impl{
		mpdClient:      api.mpdClient,
		Observer:       api.Observer,
		ctx:            api.ctx,
		requestContext: ctx,
	}
}

func (api *Impl) Connect() error {
	return wrapPkgError(api.mpdClient.Connect(api.requestContext))
}

func (api *Impl) Disconnect() error {
	return wrapPkgError(api.mpdClient.Disconnect(api.requestContext))
}

func (api *Impl) IsConnected() bool {
	return api.mpdClient.IsConnected(api.requestContext)
}
