package mediaServer

import "testing"

func TestParseStreamInfoFromPath(t *testing.T) {
	urlPath := "/event/stream/file.ts"
	event := "event"
	stream := "stream"
	file := "file.ts"
	ev, s, f, e := parseStreamInfoFromPath(urlPath)
	noError(t, e)
	assertEqual(t, event, ev)
	assertEqual(t, stream, s)
	assertEqual(t, file, f)
}

func TestParseStreamInfoFromPathWithNoPrefixSlash(t *testing.T) {
	urlPath := "event/stream/file.ts"
	event := "event"
	stream := "stream"
	file := "file.ts"
	ev, s, f, e := parseStreamInfoFromPath(urlPath)
	noError(t, e)
	assertEqual(t, event, ev)
	assertEqual(t, stream, s)
	assertEqual(t, file, f)
}

func TestParseStreamInfoFromPathWithExtraDir(t *testing.T) {
	urlPath := "/event/stream/720p/file.ts"
	event := "event"
	stream := "stream"
	file := "720p/file.ts"
	ev, s, f, e := parseStreamInfoFromPath(urlPath)
	noError(t, e)
	assertEqual(t, event, ev)
	assertEqual(t, stream, s)
	assertEqual(t, file, f)
}
