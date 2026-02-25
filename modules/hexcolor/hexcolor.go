// Package hexcolor provides a clipboard matcher that detects CSS hex color
// codes and shows an informational notification with the RGB breakdown.
package hexcolor

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

var hexColorRegex = regexp.MustCompile(`#[0-9a-fA-F]{6}\b`)

// Matchers returns a Matcher for 6-digit hex color codes.
func Matchers() []matchers.Matcher {
	return []matchers.Matcher{
		{
			Regex: hexColorRegex,
			ID:    "hexcolor",
			F:     Notify,
		},
	}
}

// Notify parses a hex color and sends a notification with RGB values.
func Notify(in string) {
	if len(in) != 7 {
		return
	}
	hex := in[1:]
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return
	}
	beeep.Alert("Hex Color", fmt.Sprintf("%s\nRGB(%d, %d, %d)", in, r, g, b), "")
}
