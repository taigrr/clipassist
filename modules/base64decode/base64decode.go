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

// Match strings that look like base64: base64 alphabet with optional padding.
// Heuristics that reduce false positives are applied separately in IsCandidate.
var base64Regex = regexp.MustCompile(`^[A-Za-z0-9+/]+={0,2}$`)

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

// IsCandidate reports whether a string is plausibly base64 text rather than an
// arbitrary alphanumeric token.
func IsCandidate(in string) bool {
	if len(in) < 16 {
		return false
	}
	if !base64Regex.MatchString(in) {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, r := range in {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}

	return (hasUpper && hasLower) || hasDigit
}

// Decode attempts to decode a base64 string. Returns the decoded string and
// true if the result is valid UTF-8 text.
func Decode(in string) (string, bool) {
	if !IsCandidate(in) {
		return "", false
	}
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
		cut := maxPreviewLength
		for cut > 0 && !utf8.RuneStart(preview[cut]) {
			cut--
		}
		preview = preview[:cut] + "..."
	}

	_ = beeep.Alert("Base64 Decoded", fmt.Sprintf("(%d bytes)\n%s", len(decoded), preview), "")
}
