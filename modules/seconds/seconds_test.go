package seconds

import (
	"testing"
)

func TestMatchersReturnsValidMatcher(t *testing.T) {
	m := Matchers()
	if len(m) != 1 {
		t.Fatalf("expected 1 matcher, got %d", len(m))
	}
	if m[0].ID != "seconds" {
		t.Fatalf("expected ID 'seconds', got %q", m[0].ID)
	}
}

func TestSecondsRegex(t *testing.T) {
	tests := []struct {
		input string
		match bool
	}{
		{"1672531200", true},
		{"123456789", false},  // 9 digits
		{"12345678901", true}, // contains 10-digit substring
		{"abcdef", false},
	}
	for _, tt := range tests {
		got := secondsRegex.FindString(tt.input)
		if tt.match && got == "" {
			t.Errorf("expected match for %q", tt.input)
		}
		if !tt.match && got != "" {
			t.Errorf("expected no match for %q, got %q", tt.input, got)
		}
	}
}

func TestConvertDateDoesNotPanic(t *testing.T) {
	ConvertDate("1672531200")
	ConvertDate("notanumber")
	ConvertDate("0000000000") // out of range
}
