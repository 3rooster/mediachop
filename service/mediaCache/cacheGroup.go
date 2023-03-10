package mediaCache

import "sync"

type Config struct {
	ClearIntervalSec int
	DefaultTTLSec    int
}

func NewCache(cfg *Config) *Cache {
	c := &Cache{
		store:            sync.Map{},
		stat:             stat{},
		clearIntervalSec: cfg.ClearIntervalSec,
		defaultTTLMs:     int64(cfg.DefaultTTLSec) * 1000,
	}
	go c.run()
	return c
}
