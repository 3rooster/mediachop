package mediaServer

import (
	"errors"
	"strings"
)

type mediaFileInfo struct {
	Path          string
	Event         string
	Stream        string
	FileName      string
	IsSegmentFile bool

	RcvDateTime   string
	RcvDateTimeMs int64
	Content       []byte
}

func (m *mediaFileInfo) CacheKey() string {
	return m.Path
}

func parseStreamInfoFromPath(path string) (event, stream, fileName string, err error) {
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
