package mediaCache

import "mediachop/common/syncMap"

func NewCache(cfg *Config) *Cache {
	c := &Cache{
		store:            syncMap.Map[string, *cacheItem]{},
		stat:             stat{},
		clearIntervalSec: cfg.ClearIntervalSec,
		defaultTTLMs:     int64(cfg.DefaultTTLSec) * 1000,
	}
	go c.run()
	return c
}
