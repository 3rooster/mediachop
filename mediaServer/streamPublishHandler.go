package mediaServer

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"mediachop/helpers/tm"
	"mediachop/service/cost"
	"mediachop/service/mediaStore"
	"net/http"
)

func publishStream(w http.ResponseWriter, r *http.Request, mf *mediaStore.MediaFile, sf *mediaStore.MediaFileCache) {
	logger := zap.L().With(
		zap.String("mod", "publish"),
		zap.String("event", mf.Event),
		zap.String("stream", mf.Stream),
		zap.String("file", mf.FileName))
	cs := cost.NewCost()

	bw := bytes.NewBuffer(make([]byte, 20*1024*1024))
	bw.Reset()
	bn, e := io.Copy(bw, r.Body)
	if e != nil {
		w.WriteHeader(502)
		w.Write([]byte(e.Error()))
		logger.With(zap.Int64("cost", cs.CostMs())).
			Error("error on read content, err=", zap.Error(e))
		return
	}
	mf.Content = bw.Bytes()
	mf.PublishedDateTimeMs = tm.UnixMillionSeconds()
	mf.PublishedDateTime = tm.NowDateTime()
	mf.PublishCostMs = mf.PublishedDateTimeMs - mf.RcvDateTimeMs
	if mf.IsInitFile {
		sf.SetEx(mf.CacheKey(), mf, 7*24*3600*1000)
	} else {
		sf.Set(mf.CacheKey(), mf)
	}
	logger.With(zap.Int64("cost", cs.CostMs()),
		zap.Int64("bytes", bn)).
		Info("publish success")
}
