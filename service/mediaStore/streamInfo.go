package mediaStore

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/cache"
	"sync"
)

type Stream struct {
	*cache.CacheGroup
	streamKey string
}

type streamStore struct {
	cache          *cache.CacheGroup
	streamInfoLock sync.Mutex
}

func (sc *streamStore) GetStreamInfo(mf *MediaFile) *Stream {
	ttlMs := 5 * 60 * int64(1000)
	if stream, o := sc.cache.TTL(mf.StreamKey(), ttlMs); o {
		cachedStream := stream.(*Stream)
		return cachedStream
	}
	sc.streamInfoLock.Lock()
	defer sc.streamInfoLock.Unlock()
	if stream, o := sc.cache.TTL(mf.StreamKey(), ttlMs); o {
		cachedStream := stream.(*Stream)
		return cachedStream
	}
	cachedStream := &Stream{
		CacheGroup: cache.NewCache(config.StreamCache),
		streamKey:  mf.StreamKey(),
	}
	cachedStream.SetLogger(
		zap.L().With(zap.String("cache", "stream"),
			zap.String("stream", mf.StreamKey())))
	sc.cache.SetEx(mf.StreamKey(), cachedStream, ttlMs)
	return cachedStream
}
