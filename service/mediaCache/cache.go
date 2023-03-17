package mediaCache

import (
	"github.com/3rooster/genericGoBox/syncMap"
	"github.com/3rooster/genericGoBox/syncPool"
	"mediachop/helpers/tm"
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
type cache struct {
	store        syncMap.Map[string, *cacheItem]
	stat         stat
	defaultTTLMs int64
}

var cacheItemPool = syncPool.Pool[*cacheItem]{}

func (c *cache) Set(key string, value any) {
	c.SetEx(key, value, c.defaultTTLMs)

}
func (c *cache) SetEx(key string, value any, ttlMs int64) {
	c.stat.SetTimes++
	item := cacheItemPool.Get()
	item.CreateTimeMs = tm.UnixMillionSeconds()
	item.ExpiredTimeMs = item.CreateTimeMs + ttlMs
	item.Data = value
	c.store.Store(key, item)
}

func (c *cache) GetCacheItem(key string) *cacheItem {
	if item, o := c.store.Load(key); o {
		return item
	}
	return nil
}

func (c *cache) Get(key string) (data any, expired bool) {
	if item, o := c.store.Load(key); o {
		c.stat.Hit++
		return item.Data, tm.UnixMillionSeconds() > item.ExpiredTimeMs
	}
	c.stat.Miss++
	return nil, false
}

func (c *cache) Delete(key string) bool {
	if v, o := c.store.Load(key); o {
		c.store.Delete(key)
		cacheItemPool.Put(v)
		return true
	}
	return false
}

func (c *cache) Clear() {
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

func (c *cache) Count() int {
	return c.store.Count()
}

func (c *cache) GetStat() *stat {
	c.stat.CacheCount = c.Count()
	return &c.stat
}

func (s *stat) clearHitAndMissStat() {
	s.Hit = 0
	s.Miss = 0
}
