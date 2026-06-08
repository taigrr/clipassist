// Package hexcolor provides a clipboard matcher that detects CSS hex color
// codes (#rgb, #rrggbb, #rrggbbaa) and shows an informational notification
// with the RGB(A) breakdown and OKLCH approximation.
package hexcolor

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

// hexColorRegex matches 3-digit (#rgb), 6-digit (#rrggbb), and 8-digit
// (#rrggbbaa) CSS hex color codes. The negative lookahead-style \b prevents
// partial matches on longer hex strings.
var hexColorRegex = regexp.MustCompile(`#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{3})\b`)

// Matchers returns a Matcher for hex color codes.
func Matchers() []matchers.Matcher {
	return []matchers.Matcher{
		{
			Regex: hexColorRegex,
			ID:    "hexcolor",
			F:     Notify,
		},
	}
}

// ParseHex parses a hex color string (#rgb, #rrggbb, or #rrggbbaa) and returns
// r, g, b in 0–255, alpha in 0–255 (255 = opaque), and whether parsing succeeded.
func ParseHex(in string) (r, g, b, a uint8, ok bool) {
	if len(in) < 1 || in[0] != '#' {
		return 0, 0, 0, 0, false
	}
	hex := in[1:]
	switch len(hex) {
	case 3:
		// #rgb → expand each digit: f → ff
		rv, err := strconv.ParseUint(string([]byte{hex[0], hex[0]}), 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		gv, err := strconv.ParseUint(string([]byte{hex[1], hex[1]}), 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		bv, err := strconv.ParseUint(string([]byte{hex[2], hex[2]}), 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		return uint8(rv), uint8(gv), uint8(bv), 255, true
	case 6:
		rv, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		gv, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		bv, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		return uint8(rv), uint8(gv), uint8(bv), 255, true
	case 8:
		rv, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		gv, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		bv, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		av, err := strconv.ParseUint(hex[6:8], 16, 8)
		if err != nil {
			return 0, 0, 0, 0, false
		}
		return uint8(rv), uint8(gv), uint8(bv), uint8(av), true
	default:
		return 0, 0, 0, 0, false
	}
}

// RGBToOKLCH converts sRGB (0–255) to OKLCH (L 0–1, C 0–~0.4, H 0–360).
func RGBToOKLCH(r, g, b uint8) (l, c, h float64) {
	// sRGB to linear
	lr := srgbToLinear(float64(r) / 255.0)
	lg := srgbToLinear(float64(g) / 255.0)
	lb := srgbToLinear(float64(b) / 255.0)

	// Linear RGB to LMS (using OKLab matrix)
	ll := 0.4122214708*lr + 0.5363325363*lg + 0.0514459929*lb
	m := 0.2119034982*lr + 0.6806995451*lg + 0.1073969566*lb
	s := 0.0883024619*lr + 0.2817188376*lg + 0.6299787005*lb

	// Cube root
	ll = math.Cbrt(ll)
	m = math.Cbrt(m)
	s = math.Cbrt(s)

	// LMS to OKLab
	labL := 0.2104542553*ll + 0.7936177850*m - 0.0040720468*s
	labA := 1.9779984951*ll - 2.4285922050*m + 0.4505937099*s
	labB := 0.0259040371*ll + 0.7827717662*m - 0.8086757660*s

	// OKLab to OKLCH
	l = labL
	c = math.Sqrt(labA*labA + labB*labB)
	h = math.Atan2(labB, labA) * 180.0 / math.Pi
	if h < 0 {
		h += 360
	}
	return l, c, h
}

// srgbToLinear converts a single sRGB channel (0–1) to linear light.
func srgbToLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// Notify parses a hex color and sends a notification with RGB and OKLCH values.
func Notify(in string) {
	r, g, b, a, ok := ParseHex(in)
	if !ok {
		return
	}

	l, c, h := RGBToOKLCH(r, g, b)

	msg := fmt.Sprintf("%s\nRGB(%d, %d, %d)", in, r, g, b)
	if a < 255 {
		msg += fmt.Sprintf(" α=%d%%", int(float64(a)/255.0*100+0.5))
	}
	msg += fmt.Sprintf("\nOKLCH(%.3f, %.3f, %.0f°)", l, c, h)

	_ = beeep.Alert("Hex Color", msg, "")
}
