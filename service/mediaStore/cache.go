package mediaStore

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/cache"
	"sync"
)

var store *streamStore

// Init init fileCache
func Init(streamCfg *cache.Config, mediaFileCfg *cache.Config) {
	store = &streamStore{
		fileCache:      cache.NewBucket(config.Cache.Stream.DefaultTTLMs),
		streamInfoLock: sync.Mutex{},
		streamCfg:      streamCfg,
		mediaFileCfg:   mediaFileCfg,
	}
	store.fileCache.SetLogger(zap.L().With(zap.String("fileCache", "stream_store")))
	go store.runClean()
}

func GetStreamInfo(mf *MediaFile) *MediaFileCache {
	return store.GetStreamInfo(mf)
}
