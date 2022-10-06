package mediaServer

import (
	"go.uber.org/zap"
	"io"
	"net/http"
)

func publishStream(w http.ResponseWriter, r *http.Request, mf *mediaFileInfo) {
	content, e := io.ReadAll(r.Body)
	if e != nil {
		w.WriteHeader(502)
		w.Write([]byte(e.Error()))
		zap.L().With(zap.String("event", mf.Event),
			zap.String("stream", mf.Stream),
			zap.String("file", mf.FileName)).Error("error on read content, err=", zap.Error(e))
		return
	}
	cache.Set(mf.CacheKey(), content)
	zap.L().With(zap.String("event", mf.Event),
		zap.String("stream", mf.Stream),
		zap.String("file", mf.FileName)).Debug("publish success")
}
