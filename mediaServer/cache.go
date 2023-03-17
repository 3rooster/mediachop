package mediaServer

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/cache"
)

var streams *streamInfoStore

// InitCache init cache
func InitCache() {
	streams = &streamInfoStore{cache: cache.NewCache(config.Cache)}
	streams.cache.SetLogger(zap.L().With(zap.String("cache", "stream_cache")))
}
