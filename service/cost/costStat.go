package cost

import "time"

type Cost struct {
	StartTimeMs int64
}

func NewCost() *Cost {
	return &Cost{
		StartTimeMs: time.Now().UnixMilli(),
	}
}

func (c *Cost) CostMs() int64 {
	return time.Now().UnixMilli() - c.StartTimeMs
}
