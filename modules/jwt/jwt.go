// Package jwt provides a clipboard matcher that detects JSON Web Tokens
// and shows a notification with the decoded header and payload.
package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

var jwtRegex = regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)

// Matchers returns a Matcher for JWT tokens.
func Matchers() []matchers.Matcher {
	return []matchers.Matcher{
		{
			Regex: jwtRegex,
			ID:    "jwt",
			F:     Notify,
		},
	}
}

// DecodeSegment decodes a base64url-encoded JWT segment.
func DecodeSegment(seg string) (map[string]any, error) {
	switch len(seg) % 4 {
	case 2:
		seg += "=="
	case 3:
		seg += "="
	}
	data, err := base64.URLEncoding.DecodeString(seg)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Notify parses a JWT and sends a notification with header/payload info.
func Notify(in string) {
	parts := strings.SplitN(in, ".", 3)
	if len(parts) != 3 {
		return
	}

	header, err := DecodeSegment(parts[0])
	if err != nil {
		return
	}
	payload, err := DecodeSegment(parts[1])
	if err != nil {
		return
	}

	alg, _ := header["alg"].(string)
	typ, _ := header["typ"].(string)

	msg := fmt.Sprintf("Type: %s, Alg: %s", typ, alg)

	if sub, ok := payload["sub"].(string); ok {
		msg += fmt.Sprintf("\nSub: %s", sub)
	}
	if iss, ok := payload["iss"].(string); ok {
		msg += fmt.Sprintf("\nIss: %s", iss)
	}
	if exp, ok := payload["exp"].(float64); ok {
		t := time.Unix(int64(exp), 0)
		if time.Now().After(t) {
			msg += fmt.Sprintf("\nExpired: %s", t.Format(time.RFC3339))
		} else {
			msg += fmt.Sprintf("\nExpires: %s", t.Format(time.RFC3339))
		}
	}

	beeep.Alert("JWT Token", msg, "")
}
