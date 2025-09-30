package mpdapi

import (
	"github.com/anpotashev/mpdgo/internal/logger"
	"github.com/anpotashev/mpdgo/internal/mpdclient"
	"time"
)

type MpdEventType uint8

const (
	UNKNOWN MpdEventType = iota
	ON_CONNECT
	ON_DISCONNECT
	ON_DATABASE_CHANGED
	ON_UPDATE_CHANGED
	ON_STORED_PLAYLIST_CHANGED
	ON_PLAYLIST_CHANGED
	ON_PLAYER_CHANGED
	ON_MIXER_CHANGED
	ON_OUTPUT_CHANGED
	ON_OPTIONS_CHANGED
	ON_PARTITION_CHANGED
	ON_STICKER_CHANGED
	ON_SUBSCRIPTION_CHANGED
	ON_MESSAGE_CHANGED
)

var eventsMap = map[string]MpdEventType{
	mpdclient.OnConnect:    ON_CONNECT,
	mpdclient.OnDisconnect: ON_DISCONNECT,
	"database":             ON_DATABASE_CHANGED,
	"update":               ON_UPDATE_CHANGED,
	"stored_playlist":      ON_STORED_PLAYLIST_CHANGED,
	"playlist":             ON_PLAYLIST_CHANGED,
	"player":               ON_PLAYER_CHANGED,
	"mixer":                ON_MIXER_CHANGED,
	"output":               ON_OUTPUT_CHANGED,
	"options":              ON_OPTIONS_CHANGED,
	"partition":            ON_PARTITION_CHANGED,
	"sticker":              ON_STICKER_CHANGED,
	"subscription":         ON_SUBSCRIPTION_CHANGED,
	"message":              ON_MESSAGE_CHANGED,
}

//func (api *Impl) Subscribe(timeout time.Duration) chan MpdEventType {
//	return api.observer.Subscribe(timeout)
//}
//
//func (api *Impl) Unsubscribe(ch chan MpdEventType) {
//	api.observer.Unsubscribe(ch)
//}

func (api *Impl) initObserver() {
	ch := api.mpdClient.Subscribe(100 * time.Millisecond)
	go func() {
		for {
			select {
			case event := <-ch:
				logger.Info("got event", "event", event)
				eventType := getEventType(event)
				logger.Info("got event", "eventType", eventType)
				if eventType != UNKNOWN {
					logger.Info("Notifying about event", "eventType", eventType)
					api.Notify(eventType)
				}
			case <-api.ctx.Done():
				return
			}
		}
	}()
}

func getEventType(event string) MpdEventType {
	logger.Info("Got event", "event type", event)
	if result, ok := eventsMap[event]; ok {
		return result
	}
	return UNKNOWN
}
