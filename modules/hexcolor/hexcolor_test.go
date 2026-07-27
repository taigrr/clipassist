package hexcolor

import (
	"math"
	"testing"
)

func TestMatchersReturnsValidMatcher(t *testing.T) {
	m := Matchers()
	if len(m) != 1 {
		t.Fatalf("expected 1 matcher, got %d", len(m))
	}
	if m[0].ID != "hexcolor" {
		t.Fatalf("expected ID 'hexcolor', got %q", m[0].ID)
	}
}

func TestHexColorRegex(t *testing.T) {
	tests := []struct {
		input string
		match bool
	}{
		{"#ff5733", true},
		{"#FF5733", true},
		{"#000000", true},
		{"ff5733", false},
		{"#fff", true},      // 3-digit shorthand — now matched
		{"#abc", true},      // 3-digit shorthand
		{"#ff57334", false}, // 7 hex digits — \b prevents match
		{"text #aabbcc!", true},
		{"#ff5733cc", true}, // 8-digit with alpha
		{"#FF5733CC", true}, // 8-digit uppercase
		{"#00000000", true}, // fully transparent black
	}
	for _, tt := range tests {
		got := hexColorRegex.FindString(tt.input)
		if tt.match && got == "" {
			t.Errorf("expected match for %q", tt.input)
		}
		if !tt.match && got != "" {
			t.Errorf("expected no match for %q, got %q", tt.input, got)
		}
	}
}

func TestParseHex3Digit(t *testing.T) {
	r, g, b, a, ok := ParseHex("#f0c")
	if !ok {
		t.Fatal("expected ok for #f0c")
	}
	if r != 0xff || g != 0x00 || b != 0xcc || a != 255 {
		t.Errorf("got RGBA(%d,%d,%d,%d), want (255,0,204,255)", r, g, b, a)
	}
}

func TestParseHex6Digit(t *testing.T) {
	r, g, b, a, ok := ParseHex("#ff5733")
	if !ok {
		t.Fatal("expected ok for #ff5733")
	}
	if r != 0xff || g != 0x57 || b != 0x33 || a != 255 {
		t.Errorf("got RGBA(%d,%d,%d,%d), want (255,87,51,255)", r, g, b, a)
	}
}

func TestParseHex8Digit(t *testing.T) {
	r, g, b, a, ok := ParseHex("#ff573380")
	if !ok {
		t.Fatal("expected ok for #ff573380")
	}
	if r != 0xff || g != 0x57 || b != 0x33 || a != 0x80 {
		t.Errorf("got RGBA(%d,%d,%d,%d), want (255,87,51,128)", r, g, b, a)
	}
}

func TestParseHexInvalid(t *testing.T) {
	tests := []string{"", "ff5733", "#gg5733", "#ff", "#f"}
	for _, in := range tests {
		_, _, _, _, ok := ParseHex(in)
		if ok {
			t.Errorf("expected not ok for %q", in)
		}
	}
}

func TestRGBToOKLCH(t *testing.T) {
	tests := []struct {
		r, g, b   uint8
		wantL     float64
		wantC     float64
		wantH     float64
		tolerance float64
	}{
		// Pure black
		{0, 0, 0, 0.0, 0.0, 0.0, 0.01},
		// Pure white
		{255, 255, 255, 1.0, 0.0, 0.0, 0.01},
		// Known red: L ≈ 0.628, C ≈ 0.258, H ≈ 29
		{255, 0, 0, 0.628, 0.258, 29.0, 0.02},
	}
	for _, tt := range tests {
		l, c, h := RGBToOKLCH(tt.r, tt.g, tt.b)
		if math.Abs(l-tt.wantL) > tt.tolerance {
			t.Errorf("RGB(%d,%d,%d) L=%.4f, want ~%.3f", tt.r, tt.g, tt.b, l, tt.wantL)
		}
		if math.Abs(c-tt.wantC) > tt.tolerance {
			t.Errorf("RGB(%d,%d,%d) C=%.4f, want ~%.3f", tt.r, tt.g, tt.b, c, tt.wantC)
		}
		// Skip hue check for achromatic colors (black/white)
		if tt.wantC > 0.01 && math.Abs(h-tt.wantH) > 3.0 {
			t.Errorf("RGB(%d,%d,%d) H=%.1f, want ~%.0f", tt.r, tt.g, tt.b, h, tt.wantH)
		}
	}
}

func TestNotifyDoesNotPanic(t *testing.T) {
	Notify("#ff5733")
	Notify("#fff")
	Notify("#ff573380")
	Notify("invalid")
}
