package mediaStore

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/cache"
	"sync"
	"time"
)

type Stream struct {
	*cache.Cache
	streamKey string
}

type streamStore struct {
	cache            *cache.Cache
	streamInfoLock   sync.Mutex
	clearIntervalSec int
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
		Cache:     cache.NewCache(int64(config.StreamCache.DefaultTTLSec)*1000, false),
		streamKey: mf.StreamKey(),
	}
	cachedStream.SetLogger(
		zap.L().With(zap.String("cache", "stream"),
			zap.String("stream", mf.StreamKey())))
	sc.cache.SetEx(mf.StreamKey(), cachedStream, ttlMs)
	return cachedStream
}

func (sc *streamStore) runClean() {
	for {
		sc.cache.Range(func(key string, v *cache.Item) bool {
			ca := v.Data.(*cache.Cache)
			ca.Clear()
			ca.PrintStatToLog()
			return true
		})
		sc.cache.Clear()
		time.Sleep(time.Duration(sc.clearIntervalSec) * time.Second)
	}

}
