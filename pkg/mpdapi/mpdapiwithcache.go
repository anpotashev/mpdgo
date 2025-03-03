package mpdapi

import (
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	treeCN     = "tree"
	playlistCN = "playlist"
	statusCN   = "status"
)

// all cache keys to invalidate on disconnect
var cacheKeys = []string{treeCN, playlistCN, statusCN}

var clearCacheByEventMap = map[MpdEventType][]string{
	ON_CONNECT:                 {},
	ON_DISCONNECT:              {treeCN, playlistCN, statusCN},
	ON_DATABASE_CHANGED:        {treeCN},
	ON_UPDATE_CHANGED:          {},
	ON_STORED_PLAYLIST_CHANGED: {},
	ON_PLAYLIST_CHANGED:        {playlistCN},
	ON_PLAYER_CHANGED:          {},
	ON_MIXER_CHANGED:           {},
	ON_OUTPUT_CHANGED:          {},
	ON_OPTIONS_CHANGED:         {},
	ON_PARTITION_CHANGED:       {},
	ON_STICKER_CHANGED:         {},
	ON_SUBSCRIPTION_CHANGED:    {},
	ON_MESSAGE_CHANGED:         {},
}

// map idleCode to cacheKeys
var idleEventCacheKey = map[string][]string{
	"update":   {treeCN},
	"playlist": {playlistCN},
}

type ImplWithCache struct {
	MpdApi
	cache *cache.Cache
}

func newWithCache(api *Impl) MpdApi {
	c := cache.New(cache.NoExpiration, cache.NoExpiration)
	ch := api.Subscribe(time.Millisecond * 100)
	go func() {
		for event := range ch {
			if cacheNames, ok := clearCacheByEventMap[event]; ok {
				for _, cacheName := range cacheNames {
					c.Delete(cacheName)
				}
			}
			//switch event {
			//case ON_DISCONNECT:
			//	onDisconnect(c)
			//}
		}
	}()
	return &ImplWithCache{MpdApi: api, cache: c}
}

//func onDisconnect(c *cache.Cache) {
//	for _, key := range cacheKeys {
//		c.Delete(key)
//	}
//}

func (api *ImplWithCache) Tree() (TreeItem, error) {
	value, found := api.cache.Get(treeCN)
	if found {
		return value.(TreeItem), nil
	}
	result, err := api.MpdApi.Tree()
	if err != nil {
		return nil, err
	}
	api.cache.Set(treeCN, result, cache.NoExpiration)
	return result, err
}
