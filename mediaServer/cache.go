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
	zap.S().Infof(format, v...)
}

var cache *bigcache.BigCache

// InitCache init cache
func InitCache() {
	cache, _ = bigcache.NewBigCache(bigcache.Config{
		Shards:             8,
		LifeWindow:         time.Minute,
		CleanWindow:        time.Minute,
		MaxEntriesInWindow: 1024,
		MaxEntrySize:       32 * 1024 * 1024,
		Verbose:            false,
		Hasher:             nil,
		HardMaxCacheSize:   config.MediaServer.MaxCacheSize,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
		Logger:             &cacheLogger{},
	})

}
