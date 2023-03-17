package cache

import (
	"github.com/3rooster/genericGoBox/syncPool"
)

var cacheItemPool = syncPool.NewPool[*Item](func() any {
	return &Item{}
})

type Item struct {
	CreateTimeMs  int64
	ExpiredTimeMs int64
	Data          any
}

func (c *Item) reset() {
	c.Data = nil
	c.CreateTimeMs = 0
	c.ExpiredTimeMs = 0
}
