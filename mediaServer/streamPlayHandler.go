package mediaServer

import (
	"go.uber.org/zap"
	"mediachop/helpers/tm"
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
	w.Header().Set("Ext-Rcv-File-Time", strconv.FormatInt(cachedMf.RcvDateTimeMs, 10))
	w.Header().Set("Ext-Rcv-File-Date", cachedMf.RcvDateTime)
	w.Header().Set("Ext-Since-Rcv-File", strconv.FormatInt(tm.UnixMillionSeconds()-cachedMf.RcvDateTimeMs, 10))

	w.Header().Set("Ext-Published-Time", strconv.FormatInt(cachedMf.PublishedDateTimeMs, 10))
	w.Header().Set("Ext-Published-Date", cachedMf.PublishedDateTime)
	w.Header().Set("Ext-Since-Published", strconv.FormatInt(tm.UnixMillionSeconds()-cachedMf.PublishedDateTimeMs, 10))

	w.Header().Set("Ext-Publish-Cost", strconv.FormatInt(cachedMf.PublishCostMs, 10))
	_, err := w.Write(cachedMf.Content)
	if err != nil {
		logger.With(zap.Int64("cost", cs.CostMs())).
			Error("play failed on write to client, err=", zap.Error(err))
	} else {
		logger.With(zap.Int64("cost", cs.CostMs())).
			Info("play success")
	}
}
