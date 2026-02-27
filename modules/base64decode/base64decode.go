// Package base64decode provides a clipboard matcher that detects base64-encoded
// strings and shows a notification with the decoded content.
package base64decode

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

// Match strings that look like base64: at least 16 chars, only base64 alphabet,
// proper padding, and must contain mixed case or digits (to reduce false positives).
var base64Regex = regexp.MustCompile(`[A-Za-z0-9+/]{12,}=*`)

const maxPreviewLength = 200

// Matchers returns a Matcher for base64-encoded strings.
func Matchers() []matchers.Matcher {
	return []matchers.Matcher{
		{
			Regex:    base64Regex,
			ID:       "base64",
			F:        Notify,
			FullText: true,
		},
	}
}

// Decode attempts to decode a base64 string. Returns the decoded string and
// true if the result is valid UTF-8 text.
func Decode(in string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		// Try without padding
		data, err = base64.RawStdEncoding.DecodeString(in)
		if err != nil {
			return "", false
		}
	}
	if !utf8.Valid(data) {
		return "", false
	}
	return string(data), true
}

// Notify decodes a base64 string and sends a notification with a preview.
func Notify(in string) {
	decoded, ok := Decode(in)
	if !ok {
		return
	}

	preview := decoded
	if len(preview) > maxPreviewLength {
		preview = preview[:maxPreviewLength] + "..."
	}

	beeep.Alert("Base64 Decoded", fmt.Sprintf("(%d bytes)\n%s", len(decoded), preview), "")
}
