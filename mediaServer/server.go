package mediaServer

import (
	"go.uber.org/zap"
	"mediachop/config"
	"mediachop/service/mediaStore"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func Start(srvCfg *config.MediaServerConfig) {
	http.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		streamHandler(writer, request)
	}))
	zap.L().Info("starting server ", zap.Int("port", srvCfg.ListenPort))
	strPort := ":" + strconv.Itoa(srvCfg.ListenPort)
	zap.S().Fatal("ListenAndServe: ", http.ListenAndServe(strPort, nil))
}

// handles get requests.
func streamHandler(w http.ResponseWriter, r *http.Request) {
	logger := zap.L().With(zap.String("method", r.Method),
		zap.String("path", r.URL.String()))
	code, err, mf, sf := mediaStore.ParseMediaFileRequest(r.URL.Path, logger)
	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	switch r.Method {
	case http.MethodGet:
		playStream(w, r, mf, sf)
	case http.MethodPut:
		publishStream(w, r, mf, sf)
	default:
		w.WriteHeader(401)
		w.Write([]byte("not support"))
	}
}
