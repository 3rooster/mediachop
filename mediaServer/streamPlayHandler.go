package mediaServer

import (
	"go.uber.org/zap"
	"mediachop/service/cost"
	"net/http"
	"strconv"
)

func playStream(w http.ResponseWriter, r *http.Request, mf *mediaFileInfo) {
	logger := zap.L().With(
		zap.String("mod", "play"),
		zap.String("event", mf.Event),
		zap.String("stream", mf.Stream),
		zap.String("file", mf.FileName))
	cs := cost.NewCost()
	cachedData, _ := cache.Get(mf.CacheKey())
	if cachedData == nil {
		w.WriteHeader(404)
		logger.With(zap.Int64("cost", cs.CostMs())).
			Debug("not found")
		return
	}
	cachedMf := cachedData.(*mediaFileInfo)
	w.Header().Set("Ext-Publish-Time", strconv.FormatInt(cachedMf.RcvDateTimeMs, 10))
	w.Header().Set("Ext-Publish-Date", cachedMf.RcvDateTime)
	_, err := w.Write(cachedMf.Content)
	if err != nil {
		logger.With(zap.Int64("cost", cs.CostMs())).
			Error("play failed on write to client, err=", zap.Error(err))
	} else {
		logger.With(zap.Int64("cost", cs.CostMs())).
			Info("play success")
	}
}
