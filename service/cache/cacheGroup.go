package cache

import (
	"go.uber.org/zap"
	"hash/fnv"
	"time"
)

type Group struct {
	group            map[uint64]*Cache
	shards           uint64
	stat             stat
	clearIntervalSec int
	defaultTTLMs     int64
	logger           *zap.Logger
}

func (c *Group) getShardCache(key string) *Cache {
	h := fnv.New64()
	h.Write([]byte(key))
	return c.group[h.Sum64()%c.shards]
}

// Set Cache use default ttl
func (c *Group) Set(key string, value any) {
	c.getShardCache(key).SetEx(key, value, c.defaultTTLMs)
}

// SetEx set with ttl ms
func (c *Group) SetEx(key string, value any, ttlMs int64) {
	c.getShardCache(key).SetEx(key, value, ttlMs)
}

// TTL set key ttl
func (c *Group) TTL(key string, ttlMs int64) (data any, exist bool) {
	return c.getShardCache(key).TTL(key, ttlMs)
}

func (c *Group) GetCacheItem(key string) *Item {
	return c.getShardCache(key).GetCacheItem(key)
}

// Get get value of key
func (c *Group) Get(key string) (data any, expired bool) {
	return c.getShardCache(key).Get(key)
}

// Delete delete key
func (c *Group) Delete(key string) bool {
	return c.getShardCache(key).Delete(key)
}

// SetLogger set logger
func (c *Group) SetLogger(logger *zap.Logger) {
	c.logger = logger
}

// clear expired data
func (c *Group) runClear() {
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

func (c *Group) printStatToLog() {
	if c.logger == nil {
		c.logger = zap.L()
	}
	c.logger.With(
		zap.String("mod", "group_cache"),
		zap.Int64("hit", c.stat.Hit),
		zap.Int64("miss", c.stat.Miss),
		zap.Int("count", c.stat.CacheCount),
		zap.Int("expired_count", c.stat.ExpiredCount),
	).Info("Cache stat")
}
