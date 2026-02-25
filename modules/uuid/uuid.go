// Package uuid provides a clipboard matcher that detects UUIDs and shows
// a desktop notification confirming the format and version.
package uuid

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

var uuidRegex = regexp.MustCompile(
	`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`,
)

// Matchers returns a Matcher that fires on UUID-shaped strings.
func Matchers() []matchers.Matcher {
	return []matchers.Matcher{
		{
			Regex: uuidRegex,
			ID:    "uuid",
			F:     Notify,
		},
	}
}

// Notify sends a desktop notification with the UUID version.
func Notify(in string) {
	version := "unknown"
	parts := strings.Split(in, "-")
	if len(parts) == 5 && len(parts[2]) == 4 {
		switch parts[2][0] {
		case '1':
			version = "v1 (time-based)"
		case '2':
			version = "v2 (DCE)"
		case '3':
			version = "v3 (MD5)"
		case '4':
			version = "v4 (random)"
		case '5':
			version = "v5 (SHA-1)"
		case '7':
			version = "v7 (time-ordered)"
		}
	}
	beeep.Alert("UUID Detected", fmt.Sprintf("%s\nVersion: %s", in, version), "")
}
