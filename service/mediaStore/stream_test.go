package mediaStore

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseStreamInfoFromPath(t *testing.T) {
	urlPath := "/event/stream/file.ts"
	event := "event"
	stream := "stream"
	file := "file.ts"
	ev, s, f, e := ParseStreamInfoFromPath(urlPath)
	require.NoError(t, e)
	require.Equal(t, event, ev)
	require.Equal(t, stream, s)
	require.Equal(t, file, f)
}

func TestParseStreamInfoFromPathWithNoPrefixSlash(t *testing.T) {
	urlPath := "event/stream/file.ts"
	event := "event"
	stream := "stream"
	file := "file.ts"
	ev, s, f, e := ParseStreamInfoFromPath(urlPath)
	require.NoError(t, e)
	require.Equal(t, event, ev)
	require.Equal(t, stream, s)
	require.Equal(t, file, f)
}

func TestParseStreamInfoFromPathWithExtraDir(t *testing.T) {
	urlPath := "/event/stream/720p/file.ts"
	event := "event"
	stream := "stream"
	file := "720p/file.ts"
	ev, s, f, e := ParseStreamInfoFromPath(urlPath)
	require.NoError(t, e)
	require.Equal(t, event, ev)
	require.Equal(t, stream, s)
	require.Equal(t, file, f)
}
