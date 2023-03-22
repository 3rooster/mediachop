package mediaStore

import (
	"errors"
	"go.uber.org/zap"
	"mediachop/helpers/tm"
	"strings"
)

type MediaFile struct {
	Path     string
	Event    string
	Stream   string
	FileName string

	StreamKey  string
	IsPlaylist bool
	IsInitFile bool

	RcvDateTime         string
	RcvDateTimeMs       int64
	PublishedDateTime   string
	PublishedDateTimeMs int64
	PublishCostMs       int64
	Content             []byte
}

func (m *MediaFile) CacheKey() string {
	return m.Path
}

func ParseStreamInfoFromPath(path string) (event, stream, fileName string, err error) {
	parts := strings.Split(path, "/")
	if parts[0] == "" {
		parts = parts[1:]
	}
	if len(parts) < 3 {
		err = errors.New("request path invalided")
		return
	}
	event = parts[0]
	stream = parts[1]
	fileName = strings.Join(parts[2:], "/")
	return
}

// ParseMediaFileRequest get requests.
func ParseMediaFileRequest(requestPath string, logger *zap.Logger) (rspCode int, err error, file *MediaFile, cache *MediaFileCache) {

	event, stream, fileName, err := ParseStreamInfoFromPath(requestPath)
	if err != nil {
		logger.Error("path not support")
		return 401, errors.New("path not support"), nil, nil
	}

	mf := &MediaFile{
		Path:          requestPath,
		Event:         event,
		Stream:        stream,
		FileName:      fileName,
		StreamKey:     event + "/" + stream,
		Content:       nil,
		RcvDateTimeMs: tm.UnixMillionSeconds(),
		RcvDateTime:   tm.NowDateTime(),
		IsPlaylist:    strings.HasSuffix(fileName, ".m3u8") || strings.HasSuffix(fileName, ".mpd"),
	}
	if !mf.IsPlaylist {
		mf.IsInitFile = strings.HasSuffix(fileName, "init.mp4")
	}
	sf := GetStreamInfo(mf)
	return 0, nil, mf, sf
}
