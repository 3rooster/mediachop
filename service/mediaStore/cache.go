package mediaStore

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/cache"
	"sync"
)

var store *streamStore

// Init init streams
func Init() {
	store = &streamStore{
		streams:           cache.NewBucket(config.Cache.Stream.DefaultTTLSec),
		streamInfoLock:    sync.Mutex{},
		clearIntervalSec:  10,
		defaultCacheTTLMS: config.Cache.Stream.DefaultTTLSec * 1000,
	}
	store.streams.SetLogger(zap.L().With(zap.String("streams", "stream_store")))
	go store.runClean()
}

func GetStreamInfo(mf *MediaFile) *Stream {
	return store.GetStreamInfo(mf)
}
