package tm

import (
	"time"
)

const DefaultTimeFormat = "2006-01-02 15:04:05"

// NowDateTime date time
func NowDateTime() string {
	return time.Now().Format(DefaultTimeFormat)
}

// UnixSeconds unix timestamp
func UnixSeconds() int64 {
	return time.Now().Unix()
}

// UnixMillionSeconds ms
func UnixMillionSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// FromUnixSecToTime unix to time
func FromUnixSecToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// UnixTimeToDatetimeStr unix sec to datetime
func UnixTimeToDatetimeStr(unixSec int64) string {
	return FromUnixSecToTime(unixSec).Format(DefaultTimeFormat)
}
