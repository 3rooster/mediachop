package mediaServer

import (
	"github.com/allegro/bigcache"
	"go.uber.org/zap"
	"net/http"
)

func playStream(w http.ResponseWriter, r *http.Request, mf *mediaFileInfo) {

	content, err := cache.Get(mf.CacheKey())
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			w.WriteHeader(404)
			zap.L().With(zap.String("event", mf.Event),
				zap.String("stream", mf.Stream),
				zap.String("file", mf.FileName)).Debug("not found")
			return
		}
		zap.L().With(zap.String("event", mf.Event),
			zap.String("stream", mf.Stream),
			zap.String("file", mf.FileName)).Error("play failed, err=", zap.Error(err))
		return
	}
	_, err = w.Write(content)
	if err != nil {
		zap.L().With(zap.String("event", mf.Event),
			zap.String("stream", mf.Stream),
			zap.String("file", mf.FileName)).Error("play failed on write to client, err=", zap.Error(err))
	} else {
		zap.L().With(zap.String("event", mf.Event),
			zap.String("stream", mf.Stream),
			zap.String("file", mf.FileName)).Debug("publish success")
	}
}
