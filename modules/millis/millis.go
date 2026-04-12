// Package millis provides a clipboard matcher that detects 13-digit Unix
// millisecond timestamps and shows a human-readable notification.
package millis

import (
	"regexp"
	"strconv"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

var millisRegex = regexp.MustCompile(`\d{13}`)

// Matchers returns a Matcher for 13-digit Unix millisecond timestamps.
func Matchers() []matchers.Matcher {
	return []matchers.Matcher{
		{
			Regex: millisRegex,
			ID:    "millis",
			F:     ConvertDate,
		},
	}
}

// ConvertDate parses a 13-digit Unix millisecond timestamp and sends a notification.
func ConvertDate(in string) {
	ts, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return
	}
	// Sanity check: between 2001 and 2286 (same range as seconds module)
	if ts < 1_000_000_000_000 || ts > 9_999_999_999_999 {
		return
	}
	t := time.UnixMilli(ts)
	beeep.Alert("Millis Converted", t.String(), "")
}
