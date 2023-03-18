package mediaStore

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/cache"
	"sync"
	"time"
)

type MediaFileCache struct {
	*cache.Bucket
	streamKey string
}

type streamStore struct {
	fileCache      *cache.Bucket
	streamInfoLock sync.Mutex
	streamCfg      *cache.Config
	mediaFileCfg   *cache.Config
}

func (sc *streamStore) GetStreamInfo(mf *MediaFile) *MediaFileCache {
	if stream, o := sc.fileCache.TTL(mf.StreamKey, sc.streamCfg.DefaultTTLMs); o {
		cachedStream := stream.(*MediaFileCache)
		return cachedStream
	}
	sc.streamInfoLock.Lock()
	defer sc.streamInfoLock.Unlock()
	if stream, o := sc.fileCache.TTL(mf.StreamKey, sc.streamCfg.DefaultTTLMs); o {
		cachedStream := stream.(*MediaFileCache)
		return cachedStream
	}
	stream := &MediaFileCache{
		Bucket:    cache.NewBucket(sc.mediaFileCfg.DefaultTTLMs),
		streamKey: mf.StreamKey,
	}
	logger := zap.L().With(zap.String("fileCache", "stream_cache"),
		zap.String("stream", mf.StreamKey))
	logger.Info("new_stream_cache")
	stream.SetLogger(logger)
	sc.fileCache.SetEx(mf.StreamKey, stream, config.Cache.Stream.DefaultTTLMs)
	return stream
}

func (sc *streamStore) runClean() {
	for {
		sc.fileCache.Range(func(key string, v *cache.Item) bool {
			ca := v.Data.(*MediaFileCache)
			ca.Clear()
			ca.PrintStatToLog()
			return true
		})
		sc.fileCache.Clear()
		time.Sleep(time.Duration(sc.streamCfg.DefaultTTLMs) * time.Millisecond)
	}

}
