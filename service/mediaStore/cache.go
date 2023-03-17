package mediaStore

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/cache"
	"sync"
)

var store *streamStore

// Init init cache
func Init() {
	store = &streamStore{
		cache:          cache.NewCache(config.Cache),
		streamInfoLock: sync.Mutex{},
	}
	store.cache.SetLogger(zap.L().With(zap.String("cache", "stream_store")))
}

func GetStreamInfo(mf *MediaFile) *Stream {
	return store.GetStreamInfo(mf)
}
