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
	if stream, o := sc.cache.TTL(mf.StreamKey, ttlMs); o {
		cachedStream := stream.(*Stream)
		return cachedStream
	}
	sc.streamInfoLock.Lock()
	defer sc.streamInfoLock.Unlock()
	if stream, o := sc.cache.TTL(mf.StreamKey, ttlMs); o {
		cachedStream := stream.(*Stream)
		return cachedStream
	}
	cachedStream := &Stream{
		Cache:     cache.NewCache(int64(config.StreamCache.DefaultTTLSec)*1000, false),
		streamKey: mf.StreamKey,
	}
	logger := zap.L().With(zap.String("cache", "stream_cache"),
		zap.String("stream", mf.StreamKey))
	logger.Info("new_stream_cache")
	cachedStream.SetLogger(logger)
	sc.cache.SetEx(mf.StreamKey, cachedStream, ttlMs)
	return cachedStream
}

func (sc *streamStore) runClean() {
	for {
		sc.cache.Range(func(key string, v *cache.Item) bool {
			ca := v.Data.(*Stream)
			ca.Clear()
			ca.PrintStatToLog()
			return true
		})
		sc.cache.Clear()
		time.Sleep(time.Duration(sc.clearIntervalSec) * time.Second)
	}

}
