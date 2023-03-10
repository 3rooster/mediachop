package mediaServer

import (
	"github.com/allegro/bigcache"
	"go.uber.org/zap"
	"mediachop/config"
	"time"
)

type cacheLogger struct {
}

func (c *cacheLogger) Printf(format string, v ...interface{}) {
	zap.S().With(zap.String("m", "cache")).Infof(format, v...)
}

var cache *bigcache.BigCache

// InitCache init cache
func InitCache() {
	cache, _ = bigcache.NewBigCache(bigcache.Config{
		Shards:             config.Cache.Shards,
		LifeWindow:         time.Second * time.Duration(config.Cache.EntryTTLSec),
		CleanWindow:        time.Second * time.Duration(config.Cache.ClearIntervalSec),
		MaxEntriesInWindow: config.Cache.MaxEntries,
		MaxEntrySize:       32 * 1024 * 1024,
		Verbose:            false,
		Hasher:             nil,
		HardMaxCacheSize:   config.MediaServer.MaxCacheSize,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
		Logger:             &cacheLogger{},
	})

}
