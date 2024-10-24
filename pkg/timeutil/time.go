package timeutil

import (
	"strings"
	"time"
)

func Format(ts time.Time) string {
	s := ts.Round(time.Millisecond).Format(time.RFC3339Nano)
	dot := strings.Index(s, ".")
	if dot == -1 {
		s = strings.Replace(s, "Z", ".000Z", 1)
		dot = strings.Index(s, ".")
	}
	zone := strings.Index(s, "Z")

	switch zone - dot {
	case 3:
		s = strings.Replace(s, "Z", "0Z", 1)
	case 2:
		s = strings.Replace(s, "Z", "00Z", 1)
	}

	return s
}

func Parse(s string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, s)
}
