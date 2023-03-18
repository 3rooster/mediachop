package cache

import (
	"github.com/3rooster/genericGoBox/syncMap"
)

var Default = NewCache(&Config{
	ClearIntervalMs: 10,
	DefaultTTLMs:    5,
	Shards:          8,
})

type Config struct {
	ClearIntervalMs int64 `yaml:"clear_interval_ms"`
	DefaultTTLMs    int64 `yaml:"cache_ttl_ms"`
	Shards          int   `yaml:"shards"`
}

func InitDefault(cfg *Config) {
	Default = NewCache(cfg)
}

func NewCache(cfg *Config) *Cache {
	c := &Cache{
		group:           map[uint64]*Bucket{},
		stat:            stat{},
		clearIntervalMs: cfg.ClearIntervalMs,
		defaultTTLMs:    cfg.DefaultTTLMs * 1000,
		shards:          uint64(cfg.Shards),
	}
	if c.shards <= 0 || c.shards > 512 {
		c.shards = 8
	}
	if c.defaultTTLMs <= 0 {
		c.defaultTTLMs = 5 * 1000
	}
	if c.clearIntervalMs <= 0 {
		c.clearIntervalMs = 5 * 1000
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
