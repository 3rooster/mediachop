package mediaServer

import (
	"go.uber.org/zap"
	"io"
	"mediachop/helpers/tm"
	"mediachop/service/cost"
	"mediachop/service/mediaStore"
	"net/http"
)

func publishStream(w http.ResponseWriter, r *http.Request, mf *mediaStore.MediaFile, sf *mediaStore.Stream) {
	logger := zap.L().With(
		zap.String("mod", "publish"),
		zap.String("event", mf.Event),
		zap.String("stream", mf.Stream),
		zap.String("file", mf.FileName))
	cs := cost.NewCost()
	bn, e := io.Copy(mf.Content, r.Body)
	if e != nil {
		w.WriteHeader(502)
		w.Write([]byte(e.Error()))
		logger.With(zap.Int64("cost", cs.CostMs())).
			Error("error on read content, err=", zap.Error(e))
		return
	}
	mf.PublishedDateTimeMs = tm.UnixMillionSeconds()
	mf.PublishedDateTime = tm.NowDateTime()
	mf.PublishCostMs = mf.PublishedDateTimeMs - mf.RcvDateTimeMs

	sf.Set(mf.CacheKey(), mf)
	logger.With(zap.Int64("cost", cs.CostMs()),
		zap.Int64("bytes", bn)).
		Info("publish success")
}
