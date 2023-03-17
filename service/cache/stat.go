package cache

import "go.uber.org/zap"

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

func (s *stat) Print(logger *zap.Logger) {
	logger.With(
		zap.Int64("hit", s.Hit),
		zap.Int64("miss", s.Miss),
		zap.Int("count", s.CacheCount),
		zap.Int("expired_count", s.ExpiredCount),
	).Info("Cache stat")
}
