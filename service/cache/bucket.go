package cache

import (
	"github.com/3rooster/genericGoBox/syncMap"
	"go.uber.org/zap"
	"mediachop/helpers/tm"
)

type Bucket struct {
	store        syncMap.Map[string, *Item]
	stat         stat
	defaultTTLMs int64
	logger       *zap.Logger
}

func (c *Bucket) Set(key string, value any) {
	c.SetEx(key, value, c.defaultTTLMs)

}
func (c *Bucket) SetEx(key string, value any, ttlMs int64) {
	c.stat.SetTimes++
	item := cacheItemPool.Get()
	item.CreateTimeMs = tm.UnixMillionSeconds()
	item.ExpiredTimeMs = item.CreateTimeMs + ttlMs
	item.Data = value
	c.store.Store(key, item)
}

func (c *Bucket) TTL(key string, ttlMs int64) (data any, exist bool) {
	if item, o := c.store.Load(key); o {
		item.ExpiredTimeMs = tm.UnixMillionSeconds() + ttlMs
		return item.Data, true
	}
	return nil, false
}

func (c *Bucket) GetCacheItem(key string) *Item {
	if item, o := c.store.Load(key); o {
		return item
	}
	return nil
}

func (c *Bucket) Get(key string) (data any, expired bool) {
	if item, o := c.store.Load(key); o {
		c.stat.Hit++
		return item.Data, tm.UnixMillionSeconds() > item.ExpiredTimeMs
	}
	c.stat.Miss++
	return nil, false
}

func (c *Bucket) Delete(key string) bool {
	if v, o := c.store.Load(key); o {
		c.store.Delete(key)
		cacheItemPool.Put(v)
		return true
	}
	return false
}

func (c *Bucket) Clear() {
	deleteKeys := map[string]int{}
	cnt := 0
	c.store.Range(func(key string, value *Item) bool {
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

func (c *Bucket) Count() int {
	return c.store.Count()
}

func (c *Bucket) GetStat() *stat {
	c.stat.CacheCount = c.Count()
	return &c.stat
}

// SetLogger set logger
func (c *Bucket) SetLogger(logger *zap.Logger) {
	c.logger = logger
}

// Range range cache items
func (c *Bucket) Range(rangeFunc func(key string, v *Item) bool) {
	c.store.Range(func(key string, v *Item) bool {
		return rangeFunc(key, v)
	})
}

func (c *Bucket) PrintStatToLog() {
	if c.logger == nil {
		c.logger = zap.L()
	}
	c.logger.With(
		zap.String("mod", "CacheBucket"),
		zap.Int64("hit", c.stat.Hit),
		zap.Int64("miss", c.stat.Miss),
		zap.Int("count", c.stat.CacheCount),
		zap.Int("expired_count", c.stat.ExpiredCount),
	).Info("CacheBucket stat")
}
