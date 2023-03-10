package mediaCache

import (
	"go.uber.org/zap"
	"mediachop/helpers/tm"
	"sync"
	"time"
)

type cacheItem struct {
	CreateTimeMs  int64
	ExpiredTimeMs int64
	Data          any
}
type stat struct {
	Hit          int64
	Miss         int64
	SetTimes     int64
	CacheCount   int
	ExpiredCount int
}
type Cache struct {
	store            sync.Map
	stat             stat
	clearIntervalSec int
	defaultTTLMs     int64
}

func (c *Cache) run() {
	for {
		c.Clear()
		time.Sleep(time.Second * time.Duration(c.clearIntervalSec))
	}
}
func (c *Cache) Set(key string, value any) {
	c.SetEx(key, value, c.defaultTTLMs)

}
func (c *Cache) SetEx(key string, value any, ttlMs int64) {
	c.stat.SetTimes++
	c.store.Store(key, &cacheItem{
		CreateTimeMs:  tm.UnixMillionSeconds(),
		ExpiredTimeMs: tm.UnixMillionSeconds() + ttlMs,
		Data:          value,
	})

}

func (c *Cache) GetCacheItem(key string) *cacheItem {
	if item, o := c.store.Load(key); o {
		ci := item.(*cacheItem)
		return ci
	}
	return nil
}

func (c *Cache) Get(key string) (data any, expired bool) {
	if item, o := c.store.Load(key); o {
		ci := item.(*cacheItem)
		c.stat.Hit++
		return ci.Data, tm.UnixMillionSeconds() > ci.ExpiredTimeMs
	}
	c.stat.Miss++
	return nil, false
}

func (c *Cache) Clear() {
	deleteKeys := map[string]int{}
	cnt := 0
	c.store.Range(func(key, value any) bool {
		item := value.(*cacheItem)
		k := key.(string)
		if tm.UnixMillionSeconds() > item.ExpiredTimeMs {
			deleteKeys[k] = 1
		}
		cnt++
		return true
	})
	c.stat.CacheCount = cnt
	c.stat.ExpiredCount = len(deleteKeys)
	for k, _ := range deleteKeys {
		c.store.Delete(k)
	}
	c.printStatToLog()
}

func (c *Cache) Count() int {
	cnt := 0
	c.store.Range(func(key, value any) bool {
		cnt++
		return true
	})
	return cnt
}

func (c *Cache) GetStat() *stat {
	c.stat.CacheCount = c.Count()
	return &c.stat
}

func (c *Cache) printStatToLog() {
	zap.S().With(
		zap.String("m", "Cache"),
		zap.Int64("hit", c.stat.Hit),
		zap.Int64("miss", c.stat.Miss),
		zap.Int("count", c.stat.CacheCount),
		zap.Int("expired_count", c.stat.ExpiredCount),
	).Info("Cache stat")
}
