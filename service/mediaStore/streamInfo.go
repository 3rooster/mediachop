package mediaStore

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/cache"
	"sync"
	"time"
)

type Stream struct {
	*cache.Bucket
	streamKey string
}

type streamStore struct {
	streams           *cache.Bucket
	streamInfoLock    sync.Mutex
	clearIntervalSec  int
	defaultCacheTTLMS int64
}

func (sc *streamStore) GetStreamInfo(mf *MediaFile) *Stream {
	if stream, o := sc.streams.TTL(mf.StreamKey, sc.defaultCacheTTLMS); o {
		cachedStream := stream.(*Stream)
		return cachedStream
	}
	sc.streamInfoLock.Lock()
	defer sc.streamInfoLock.Unlock()
	if stream, o := sc.streams.TTL(mf.StreamKey, sc.defaultCacheTTLMS); o {
		cachedStream := stream.(*Stream)
		return cachedStream
	}
	cachedStream := &Stream{
		Bucket:    cache.NewBucket(config.Cache.MediaFile.DefaultTTLSec * 1000),
		streamKey: mf.StreamKey,
	}
	logger := zap.L().With(zap.String("streams", "stream_cache"),
		zap.String("stream", mf.StreamKey))
	logger.Info("new_stream_cache")
	cachedStream.SetLogger(logger)
	sc.streams.SetEx(mf.StreamKey, cachedStream, config.Cache.Stream.DefaultTTLSec*1000)
	return cachedStream
}

func (sc *streamStore) runClean() {
	for {
		sc.streams.Range(func(key string, v *cache.Item) bool {
			ca := v.Data.(*Stream)
			ca.Clear()
			ca.PrintStatToLog()
			return true
		})
		sc.streams.Clear()
		time.Sleep(time.Duration(sc.clearIntervalSec) * time.Second)
	}

}
