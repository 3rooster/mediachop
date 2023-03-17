package mediaCache

type Config struct {
	ClearIntervalSec int
	DefaultTTLSec    int
	Shards           int
}

func NewCache(cfg *Config) *CacheGroup {
	c := &CacheGroup{
		group:            map[uint64]*Cache{},
		stat:             stat{},
		clearIntervalSec: cfg.ClearIntervalSec,
		defaultTTLMs:     int64(cfg.DefaultTTLSec) * 1000,
	}
	go c.runClear()
	return c
}
