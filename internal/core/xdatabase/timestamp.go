package xdatabase

import (
	"fmt"
	"time"
)

func nowDateTimeStamp(t time.Time) string {
	return fmt.Sprintf("%04d/%02d/%02d-%02d:%02d:%02d.%09d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

func nowJSTimestamp(t time.Time) int64 {
	base := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	return t.Sub(base).Milliseconds()
}
