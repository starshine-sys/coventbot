package common

import (
	"time"

	"gitlab.com/1f320/x/duration"
)

func FormatTime(t time.Time) string {
	s, isBefore := duration.FormatAt(time.Now(), t)
	if isBefore {
		return s + " ago"
	}
	return s + " from now"
}
