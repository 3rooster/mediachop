package mediaStore

import (
	"go.uber.org/zap"
	"mediachop/service/cache"
	"sync"
)

var store *streamStore

// Init init cache
func Init() {
	store = &streamStore{
		cache:            cache.NewCache(60*1000, false),
		streamInfoLock:   sync.Mutex{},
		clearIntervalSec: 10,
	}
	store.cache.SetLogger(zap.L().With(zap.String("cache", "stream_store")))
	go store.runClean()
}

func GetStreamInfo(mf *MediaFile) *Stream {
	return store.GetStreamInfo(mf)
}
