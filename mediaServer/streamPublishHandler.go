package mediaServer

import (
	"go.uber.org/zap"
	"io"
	"mediachop/service/cost"
	"net/http"
)

func publishStream(w http.ResponseWriter, r *http.Request, mf *mediaFileInfo) {
	logger := zap.L().With(zap.String("event", mf.Event),
		zap.String("stream", mf.Stream),
		zap.String("file", mf.FileName))
	cs := cost.NewCost()
	content, e := io.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(502)
		w.Write([]byte(e.Error()))
		logger.With(zap.Int64("cost", cs.CostMs())).
			Error("error on read content, err=", zap.Error(e))
		return
	}
	cache.Set(mf.CacheKey(), content)
	logger.With(zap.Int64("cost", cs.CostMs())).
		Info("publish success")
}
