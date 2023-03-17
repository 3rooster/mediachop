package cache

import (
	"github.com/3rooster/genericGoBox/syncMap"
	"github.com/3rooster/genericGoBox/syncPool"
	"go.uber.org/zap"
	"mediachop/helpers/tm"
)

var cacheItemPool = syncPool.NewPool[*CacheItem](func() any {
	return &CacheItem{}
})

type CacheItem struct {
	CreateTimeMs  int64
	ExpiredTimeMs int64
	Data          any
}

func (c *CacheItem) reset() {
	c.Data = nil
	c.CreateTimeMs = 0
	c.ExpiredTimeMs = 0
}

type stat struct {
	Hit          int64
	Miss         int64
	SetTimes     int64
	CacheCount   int
	ExpiredCount int
}

func (s *stat) clearHitAndMissStat() {
	s.Hit = 0
	s.Miss = 0
}

type Cache struct {
	store        syncMap.Map[string, *CacheItem]
	stat         stat
	defaultTTLMs int64
	logger       *zap.Logger
}

func (c *Cache) Set(key string, value any) {
	c.SetEx(key, value, c.defaultTTLMs)

}
func (c *Cache) SetEx(key string, value any, ttlMs int64) {
	c.stat.SetTimes++
	item := cacheItemPool.Get()
	item.CreateTimeMs = tm.UnixMillionSeconds()
	item.ExpiredTimeMs = item.CreateTimeMs + ttlMs
	item.Data = value
	c.store.Store(key, item)
}

func (c *Cache) TTL(key string, ttlMs int64) (data any, exist bool) {
	if item, o := c.store.Load(key); o {
		item.ExpiredTimeMs = tm.UnixMillionSeconds() + ttlMs
		return item.Data, true
	}
	return nil, false
}

func (c *Cache) GetCacheItem(key string) *CacheItem {
	if item, o := c.store.Load(key); o {
		return item
	}
	return nil
}

func (c *Cache) Get(key string) (data any, expired bool) {
	if item, o := c.store.Load(key); o {
		c.stat.Hit++
		return item.Data, tm.UnixMillionSeconds() > item.ExpiredTimeMs
	}
	c.stat.Miss++
	return nil, false
}

func (c *Cache) Delete(key string) bool {
	if v, o := c.store.Load(key); o {
		c.store.Delete(key)
		cacheItemPool.Put(v)
		return true
	}
	return false
}

func (c *Cache) Clear() {
	deleteKeys := map[string]int{}
	cnt := 0
	c.store.Range(func(key string, value *CacheItem) bool {
		item := value
		k := key
		if tm.UnixMillionSeconds() > item.ExpiredTimeMs {
			deleteKeys[k] = 1
		}
		cnt++
		return true
	})
	c.stat.ExpiredCount = len(deleteKeys)
	c.stat.CacheCount = cnt - c.stat.ExpiredCount
	for k, _ := range deleteKeys {
		v, _ := c.store.Load(k)
		c.store.Delete(k)
		if v != nil {
			cacheItemPool.Put(v)
		}
	}

}

func (c *Cache) Count() int {
	return c.store.Count()
}

func (c *Cache) GetStat() *stat {
	c.stat.CacheCount = c.Count()
	return &c.stat
}

// SetLogger set logger
func (c *Cache) SetLogger(logger *zap.Logger) {
	c.logger = logger
}

// Range range cache items
func (c *Cache) Range(rangeFunc func(key string, v *CacheItem) bool) {
	c.store.Range(func(key string, v *CacheItem) bool {
		return rangeFunc(key, v)
	})
}

func (c *Cache) PrintStatToLog() {
	if c.logger == nil {
		c.logger = zap.L()
	}
	c.logger.With(
		zap.String("mod", "Cache"),
		zap.Int64("hit", c.stat.Hit),
		zap.Int64("miss", c.stat.Miss),
		zap.Int("count", c.stat.CacheCount),
		zap.Int("expired_count", c.stat.ExpiredCount),
	).Info("Cache stat")
}
