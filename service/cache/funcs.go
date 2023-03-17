package cache

import (
	"github.com/3rooster/genericGoBox/syncMap"
)

type Config struct {
	ClearIntervalSec int `yaml:"clear_interval_sec"`
	DefaultTTLSec    int `yaml:"cache_ttl_sec"`
	Shards           int `yaml:"shards"`
}

func NewCacheGroup(cfg *Config) *CacheGroup {
	c := &CacheGroup{
		group:            map[uint64]*Cache{},
		stat:             stat{},
		clearIntervalSec: cfg.ClearIntervalSec,
		defaultTTLMs:     int64(cfg.DefaultTTLSec) * 1000,
		shards:           uint64(cfg.Shards),
	}
	if c.shards <= 0 || c.shards > 512 {
		c.shards = 8
	}
	if c.defaultTTLMs <= 0 {
		c.defaultTTLMs = 1000
	}
	if c.clearIntervalSec <= 0 {
		c.clearIntervalSec = 10
	}
	for i := uint64(0); i < c.shards; i++ {
		c.group[i] = NewCache(c.defaultTTLMs, false)
	}
	go c.runClear()
	return c
}

func NewCache(defaultTTLMS int64, runClean bool) *Cache {
	r := &Cache{
		store:        syncMap.Map[string, *CacheItem]{},
		stat:         stat{},
		defaultTTLMs: defaultTTLMS,
	}
	return r
}
