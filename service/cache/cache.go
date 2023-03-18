package cache

import (
	"go.uber.org/zap"
	"hash/fnv"
	"time"
)

type Cache struct {
	group           map[uint64]*Bucket
	shards          uint64
	stat            stat
	clearIntervalMs int64
	defaultTTLMs    int64
	logger          *zap.Logger
}

func (c *Cache) getBucket(key string) *Bucket {
	h := fnv.New64()
	h.Write([]byte(key))
	return c.group[h.Sum64()%c.shards]
}

// Set CacheBucket use default ttl
func (c *Cache) Set(key string, value any) {
	c.getBucket(key).SetEx(key, value, c.defaultTTLMs)
}

// SetEx set with ttl ms
func (c *Cache) SetEx(key string, value any, ttlMs int64) {
	c.getBucket(key).SetEx(key, value, ttlMs)
}

// TTL set key ttl
func (c *Cache) TTL(key string, ttlMs int64) (data any, exist bool) {
	return c.getBucket(key).TTL(key, ttlMs)
}

func (c *Cache) GetCacheItem(key string) *Item {
	return c.getBucket(key).GetCacheItem(key)
}

// Get get value of key
func (c *Cache) Get(key string) (data any, expired bool) {
	return c.getBucket(key).Get(key)
}

// Delete delete key
func (c *Cache) Delete(key string) bool {
	return c.getBucket(key).Delete(key)
}

// SetLogger set logger
func (c *Cache) SetLogger(logger *zap.Logger) {
	c.logger = logger
}

// clear expired data
func (c *Cache) runClear() {
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
		time.Sleep(time.Millisecond * time.Duration(c.clearIntervalMs))
	}
}

func (c *Cache) printStatToLog() {
	if c.logger == nil {
		c.logger = zap.L()
	}
	c.logger.With(
		zap.String("mod", "group_cache"),
		zap.Int64("hit", c.stat.Hit),
		zap.Int64("miss", c.stat.Miss),
		zap.Int("count", c.stat.CacheCount),
		zap.Int("expired_count", c.stat.ExpiredCount),
	).Info("CacheBucket stat")
}
