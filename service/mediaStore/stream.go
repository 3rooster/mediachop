package mediaStore

import (
	"errors"
	"strings"
)

type MediaFile struct {
	Path     string
	Event    string
	Stream   string
	FileName string

	StreamKey  string
	IsPlaylist bool

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
