package mediaCache

import (
	"github.com/3rooster/genericGoBox/syncMap"
)

type Config struct {
	ClearIntervalSec int
	DefaultTTLSec    int
	Shards           int
}

func NewCache(cfg *Config) *CacheGroup {
	c := &CacheGroup{
		group:            map[uint64]*cache{},
		stat:             stat{},
		clearIntervalSec: cfg.ClearIntervalSec,
		defaultTTLMs:     int64(cfg.DefaultTTLSec) * 1000,
	}
	if cfg.Shards <= 0 {
		c.shards = 8
	} else {
		c.shards = uint64(cfg.Shards)
	}
	for i := uint64(0); i < c.shards; i++ {
		c.group[i] = &cache{
			store:        syncMap.Map[string, *cacheItem]{},
			stat:         stat{},
			defaultTTLMs: c.defaultTTLMs,
		}
	}
	go c.runClear()
	return c
}
