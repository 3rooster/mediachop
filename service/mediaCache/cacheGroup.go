package mediaCache

import (
	"go.uber.org/zap"
	"hash/fnv"
	"time"
)

type CacheGroup struct {
	group            map[uint64]*Cache
	shards           uint64
	stat             stat
	clearIntervalSec int
	defaultTTLMs     int64
}

func (c *CacheGroup) getShardCache(key string) *Cache {
	h := fnv.New64()
	h.Write([]byte(key))
	return c.group[h.Sum64()%c.shards]
}

func (c *CacheGroup) Set(key string, value any) {
	c.getShardCache(key).Set(key, value)

}
func (c *CacheGroup) SetEx(key string, value any, ttlMs int64) {
	c.getShardCache(key).SetEx(key, value, ttlMs)
}

func (c *CacheGroup) GetCacheItem(key string) *cacheItem {
	return c.getShardCache(key).GetCacheItem(key)
}

func (c *CacheGroup) Get(key string) (data any, expired bool) {
	return c.getShardCache(key).Get(key)
}

func (c *CacheGroup) Delete(key string) bool {
	return c.getShardCache(key).Delete(key)
}

func (c *CacheGroup) runClear() {
	for {
		for _, cacheInstance := range c.group {
			cacheInstance.Clear()
			s := cacheInstance.GetStat()
			c.stat.SetTimes += s.SetTimes
			c.stat.Hit += s.Hit
			c.stat.Miss += s.Miss
			c.stat.ExpiredCount += s.ExpiredCount
			c.stat.CacheCount += s.CacheCount
			cacheInstance.stat.clearHitAndMissStat()
		}
		c.printStatToLog()
		c.stat.clearHitAndMissStat()
		time.Sleep(time.Second * time.Duration(c.clearIntervalSec))
	}
}

func (c *CacheGroup) printStatToLog() {
	zap.S().With(
		zap.String("mod", "Cache"),
		zap.Int64("hit", c.stat.Hit),
		zap.Int64("miss", c.stat.Miss),
		zap.Int("count", c.stat.CacheCount),
		zap.Int("expired_count", c.stat.ExpiredCount),
	).Info("Cache stat")
}
