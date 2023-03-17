package mediaServer

import (
	"mediachop/config"
	"mediachop/service/mediaCache"
)

var cache *mediaCache.CacheGroup

// InitCache init mediaCache
func InitCache() {
	cache = mediaCache.NewCache(config.Cache)
}

var streamCache = mediaCache.NewCache(&mediaCache.Config{
	ClearIntervalSec: 60,
	DefaultTTLSec:    60,
	Shards:           8,
})
