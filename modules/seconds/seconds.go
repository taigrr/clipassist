// Package seconds provides a clipboard matcher that detects 10-digit Unix
// timestamps (seconds since epoch) and shows a human-readable notification.
package seconds

import (
	"regexp"
	"strconv"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

var secondsRegex = regexp.MustCompile(`\b\d{10}\b`)

// Matchers returns a Matcher for 10-digit Unix timestamps.
func Matchers() []matchers.Matcher {
	return []matchers.Matcher{
		{
			Regex: secondsRegex,
			ID:    "seconds",
			F:     ConvertDate,
		},
	}
}

// ConvertDate parses a 10-digit Unix timestamp and sends a notification.
func ConvertDate(in string) {
	ts, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return
	}
	// Sanity check: between 2001-09-09 and 2286-11-20
	if ts < 1_000_000_000 || ts > 9_999_999_999 {
		return
	}
	t := time.Unix(ts, 0)
	beeep.Alert("Unix Seconds Converted", t.String(), "")
}
