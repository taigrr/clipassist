package millis

import (
	"testing"
)

func TestMatchersReturnsValidMatcher(t *testing.T) {
	m := Matchers()
	if len(m) != 1 {
		t.Fatalf("expected 1 matcher, got %d", len(m))
	}
	if m[0].ID != "millis" {
		t.Fatalf("expected ID 'millis', got %q", m[0].ID)
	}
	if m[0].Regex == nil {
		t.Fatal("regex should not be nil")
	}
}

func TestMatchersRegexMatches13Digits(t *testing.T) {
	m := Matchers()
	r := m[0].Regex

	tests := []struct {
		input string
		match bool
	}{
		{"1234567890123", true},
		{"0000000000000", true},
		{"123456789012", false},  // 12 digits
		{"12345678901234", true}, // 14 digits — contains a 13-digit substring
		{"abcdef", false},
	}

	for _, tt := range tests {
		got := r.FindString(tt.input)
		if tt.match && got == "" {
			t.Errorf("expected match for %q", tt.input)
		}
		if !tt.match && got != "" {
			t.Errorf("expected no match for %q, got %q", tt.input, got)
		}
	}
}

func TestConvertDateDoesNotPanic(t *testing.T) {
	// Valid millis timestamp (2023-01-01 ~)
	ConvertDate("1672531200000")
	// Invalid input
	ConvertDate("notanumber")
}
