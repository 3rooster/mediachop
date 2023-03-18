package mediaServer

import (
	"bytes"
	"go.uber.org/zap"
	"mediachop/helpers/tm"
	"mediachop/service/cost"
	"mediachop/service/mediaStore"
	"net/http"
	"strconv"
)

func playStream(w http.ResponseWriter, r *http.Request, mf *mediaStore.MediaFile, sf *mediaStore.Stream) {
	logger := zap.L().With(
		zap.String("mod", "play"),
		zap.String("event", mf.Event),
		zap.String("stream", mf.Stream),
		zap.String("file", mf.FileName))
	cs := cost.NewCost()
	cachedData, _ := sf.Get(mf.CacheKey())
	if cachedData == nil {
		w.WriteHeader(404)
		logger.With(zap.Int64("cost", cs.CostMs())).
			Debug("not found")
		return
	}
	cachedMf := cachedData.(*mediaStore.MediaFile)
	w.Header().Set("Ext-Rcv-File-Time", strconv.FormatInt(cachedMf.RcvDateTimeMs, 10))
	w.Header().Set("Ext-Rcv-File-Date", cachedMf.RcvDateTime)
	w.Header().Set("Ext-Since-Rcv-File", strconv.FormatInt(tm.UnixMillionSeconds()-cachedMf.RcvDateTimeMs, 10))

	w.Header().Set("Ext-Published-Time", strconv.FormatInt(cachedMf.PublishedDateTimeMs, 10))
	w.Header().Set("Ext-Published-Date", cachedMf.PublishedDateTime)
	w.Header().Set("Ext-Since-Published", strconv.FormatInt(tm.UnixMillionSeconds()-cachedMf.PublishedDateTimeMs, 10))

	w.Header().Set("Ext-Publish-Cost", strconv.FormatInt(cachedMf.PublishCostMs, 10))
	if cachedMf.IsPlaylist {
		w.Header().Set("Cache-Control", "no-store")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}
	br := bytes.NewReader(cachedMf.Content)
	bn, err := br.WriteTo(w)
	if err != nil {
		logger.With(zap.Int64("cost", cs.CostMs())).
			Error("play failed on write to client, err=", zap.Error(err))
	} else {
		logger.With(
			zap.Int64("cost", cs.CostMs()),
			zap.Int64("bytes", bn),
		).Info("play success")
	}
}
