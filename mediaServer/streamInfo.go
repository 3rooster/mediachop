package mediaServer

import (
	"mediachop/config"
	"mediachop/service/cache"
	"sync"
)

type streamInfo struct {
	cache *cache.CacheGroup
}

type streamInfoStore struct {
	cache *cache.CacheGroup
}

var streamInfoLock = sync.Mutex{}

func (sc *streamInfoStore) getStreamInfo(mf *mediaFileInfo) *streamInfo {
	ttlMs := 60 * int64(1000)
	if stream, o := sc.cache.TTL(mf.StreamKey(), ttlMs); o {
		cachedStream := stream.(*streamInfo)
		return cachedStream
	}
	streamInfoLock.Lock()
	defer streamInfoLock.Unlock()
	if stream, o := sc.cache.TTL(mf.StreamKey(), ttlMs); o {
		cachedStream := stream.(*streamInfo)
		return cachedStream
	}
	cachedStream := &streamInfo{
		cache: cache.NewCache(config.StreamCache),
	}
	sc.cache.SetEx(mf.StreamKey(), cachedStream, ttlMs)
	return cachedStream
}
