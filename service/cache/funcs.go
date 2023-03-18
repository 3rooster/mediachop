package cache

import (
	"github.com/3rooster/genericGoBox/syncMap"
)

type Config struct {
	ClearIntervalSec int64 `yaml:"clear_interval_sec"`
	DefaultTTLSec    int64 `yaml:"cache_ttl_sec"`
	Shards           int   `yaml:"shards"`
}

func NewCacheGroup(cfg *Config) *Cache {
	c := &Cache{
		group:            map[uint64]*Bucket{},
		stat:             stat{},
		clearIntervalSec: cfg.ClearIntervalSec,
		defaultTTLMs:     cfg.DefaultTTLSec * 1000,
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
		c.group[i] = NewBucket(c.defaultTTLMs)
	}
	go c.runClear()
	return c
}

func NewBucket(defaultTTLMS int64) *Bucket {
	r := &Bucket{
		store:        syncMap.Map[string, *Item]{},
		stat:         stat{},
		defaultTTLMs: defaultTTLMS,
	}
	return r
}
