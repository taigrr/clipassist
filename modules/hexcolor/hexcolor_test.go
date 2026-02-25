package hexcolor

import (
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
		{"#fff", false},      // 3-digit shorthand — not matched
		{"#ff57334", false},   // 7 hex digits — \b prevents match
		{"text #aabbcc!", true},
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

func TestNotifyDoesNotPanic(t *testing.T) {
	Notify("#ff5733")
	Notify("invalid")
}
