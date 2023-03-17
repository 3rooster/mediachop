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
