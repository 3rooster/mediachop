package mediaServer

import (
	"go.uber.org/zap"
	"mediachop/config"
	"net/http"
	"strconv"
	"strings"
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
	event, stream, fileName, err := parseStreamInfoFromPath(r.URL.Path)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("path not support"))
		zap.L().Error("path not support")
		return
	}
	mf := &mediaFileInfo{
		Path:          r.URL.Path,
		Event:         event,
		Stream:        stream,
		FileName:      fileName,
		IsSegmentFile: strings.HasSuffix(fileName, ".m3u8") || strings.HasSuffix(fileName, ".mpd"),
	}
	switch r.Method {
	case http.MethodGet:
		playStream(w, r, mf)
	case http.MethodPut:
		publishStream(w, r, mf)
	default:
		w.WriteHeader(401)
		w.Write([]byte("not support"))
	}
}
