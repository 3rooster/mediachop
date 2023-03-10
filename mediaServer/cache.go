package mediaServer

import (
	"mediachop/config"
	"mediachop/service/mediaCache"
)

var cache *mediaCache.Cache

// InitCache init mediaCache
func InitCache() {
	cache = mediaCache.NewCache(&mediaCache.Config{
		ClearIntervalSec: config.Cache.ClearIntervalSec,
		DefaultTTLSec:    config.Cache.CacheTTLSec,
	})
}
