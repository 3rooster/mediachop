package mediaCache

import (
	"github.com/3rooster/genericGoBox/syncPool"
	"go.uber.org/zap"
	"mediachop/common/syncMap"
	"mediachop/helpers/tm"
	"time"
)

type cacheItem struct {
	CreateTimeMs  int64
	ExpiredTimeMs int64
	Data          any
}

func (c *cacheItem) reset() {
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
type Cache struct {
	store            syncMap.Map[string, *cacheItem]
	stat             stat
	clearIntervalSec int
	defaultTTLMs     int64
}

var cacheItemPool = syncPool.Pool[*cacheItem]{}

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
	item := cacheItemPool.Get()
	item.CreateTimeMs = tm.UnixMillionSeconds()
	item.ExpiredTimeMs = item.CreateTimeMs + ttlMs
	item.Data = value
	c.store.Store(key, item)
}

func (c *Cache) GetCacheItem(key string) *cacheItem {
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
	c.store.Range(func(key string, value *cacheItem) bool {
		item := value
		k := key
		if tm.UnixMillionSeconds() > item.ExpiredTimeMs {
			deleteKeys[k] = 1
		}
		cnt++
		return true
	})
	c.stat.CacheCount = cnt
	c.stat.ExpiredCount = len(deleteKeys)
	for k, _ := range deleteKeys {
		v, _ := c.store.Load(k)
		c.store.Delete(k)
		if v != nil {
			cacheItemPool.Put(v)
		}
	}
	c.printStatToLog()
}

func (c *Cache) Count() int {
	cnt := 0
	c.store.Range(func(key string, value *cacheItem) bool {
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
		zap.String("mod", "Cache"),
		zap.Int64("hit", c.stat.Hit),
		zap.Int64("miss", c.stat.Miss),
		zap.Int("count", c.stat.CacheCount),
		zap.Int("expired_count", c.stat.ExpiredCount),
	).Info("Cache stat")
}

func (c *Cache) clearHitAndMissStat() {
	c.stat.Hit = 0
	c.stat.Miss = 0
}
