package uuid

import (
	"testing"
)

func TestMatchersReturnsValidMatcher(t *testing.T) {
	m := Matchers()
	if len(m) != 1 {
		t.Fatalf("expected 1 matcher, got %d", len(m))
	}
	if m[0].ID != "uuid" {
		t.Fatalf("expected ID 'uuid', got %q", m[0].ID)
	}
}

func TestUUIDRegexMatches(t *testing.T) {
	tests := []struct {
		input string
		match bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"not-a-uuid", false},
		{"550e8400-e29b-41d4-a716", false},
		{"ffffffff-ffff-ffff-ffff-ffffffffffff", true},
		{"text 550e8400-e29b-41d4-a716-446655440000 more", true},
	}
	for _, tt := range tests {
		got := uuidRegex.FindString(tt.input)
		if tt.match && got == "" {
			t.Errorf("expected match for %q", tt.input)
		}
		if !tt.match && got != "" {
			t.Errorf("expected no match for %q, got %q", tt.input, got)
		}
	}
}

func TestNotifyDoesNotPanic(t *testing.T) {
	Notify("550e8400-e29b-41d4-a716-446655440000")
	Notify("invalid")
}
