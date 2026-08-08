// Package base64decode provides a clipboard matcher that detects base64-encoded
// strings and shows a notification with the decoded content.
package base64decode

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

// Match strings that look like base64: base64 or base64url alphabet with optional padding.
// Heuristics that reduce false positives are applied separately in IsCandidate.
var base64Regex = regexp.MustCompile(`^[A-Za-z0-9+/_-]+={0,2}$`)

const maxPreviewLength = 200

var decoders = []*base64.Encoding{
	base64.StdEncoding,
	base64.RawStdEncoding,
	base64.URLEncoding,
	base64.RawURLEncoding,
}

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
// true if the result is valid UTF-8 that looks like text (printable characters
// or whitespace, per isText). Requiring the decoded bytes to look like text
// rather than merely valid UTF-8 avoids firing on url-safe-alphabet inputs that
// aren't really base64, such as UUIDs and other hyphenated identifiers, which
// otherwise decode to binary garbage.
func Decode(in string) (string, bool) {
	if !IsCandidate(in) {
		return "", false
	}
	var (
		data []byte
		ok   bool
	)
	for _, decoder := range decoders {
		decoded, err := decoder.DecodeString(in)
		if err == nil {
			data = decoded
			ok = true
			break
		}
	}
	if !ok {
		return "", false
	}
	if !utf8.Valid(data) || !isText(data) {
		return "", false
	}
	return string(data), true
}

// isText reports whether data consists entirely of printable characters
// (per unicode.IsPrint) or whitespace (per unicode.IsSpace, which covers tabs,
// newlines, and Unicode spaces such as NBSP), i.e. it plausibly represents
// human-readable text rather than decoded binary content.
func isText(data []byte) bool {
	for _, r := range string(data) {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
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
