package mediaServer

import (
	"github.com/allegro/bigcache"
	"go.uber.org/zap"
	"mediachop/service/cost"
	"net/http"
)

func playStream(w http.ResponseWriter, r *http.Request, mf *mediaFileInfo) {
	logger := zap.L().With(
		zap.String("mod", "play"),
		zap.String("event", mf.Event),
		zap.String("stream", mf.Stream),
		zap.String("file", mf.FileName))
	cs := cost.NewCost()
	content, err := cache.Get(mf.CacheKey())
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			w.WriteHeader(404)
			logger.With(zap.Int64("cost", cs.CostMs())).
				Debug("not found")
			return
		}
		logger.With(zap.Int64("cost", cs.CostMs())).
			Error("play failed, err=", zap.Error(err))
		return
	}
	_, err = w.Write(content)
	if err != nil {
		logger.With(zap.Int64("cost", cs.CostMs())).
			Error("play failed on write to client, err=", zap.Error(err))
	} else {
		logger.With(zap.Int64("cost", cs.CostMs())).
			Info("play success")
	}
}
