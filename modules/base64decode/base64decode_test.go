package base64decode

import (
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		input string
		want  string
		ok    bool
	}{
		{"SGVsbG8gV29ybGQh", "Hello World!", true},
		{"dGVzdGluZyAxMjM=", "testing 123", true},
		{"QUJDREVGR0hJSktMTQ==", "ABCDEFGHIJKLM", true},
		{"not-base64!!!", "", false},
		{"abcdefghijklmnop", "", false},
	}
	for _, tt := range tests {
		got, ok := Decode(tt.input)
		if ok != tt.ok {
			t.Errorf("Decode(%q) ok = %v, want %v", tt.input, ok, tt.ok)
		}
		if got != tt.want {
			t.Errorf("Decode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsCandidate(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"SGVsbG8gV29ybGQh", true},
		{"dGVzdGluZyAxMjM=", true},
		{"QUJDREVGR0hJSktMTQ==", true},
		{"short", false},
		{"hello world", false},
		{"abcdefghijklmnop", false},
		{"ABCDEFGHIJKLMNOP", false},
	}
	for _, tt := range tests {
		got := IsCandidate(tt.input)
		if got != tt.want {
			t.Errorf("IsCandidate(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestBase64Regex(t *testing.T) {
	tests := []struct {
		input string
		match bool
	}{
		{"SGVsbG8gV29ybGQh", true},
		{"dGVzdGluZyAxMjM=", true},
		{"QUJDREVGR0hJSktMTQ==", true},
		{"hello world", false},
		{"not-base64!!!", false},
	}
	for _, tt := range tests {
		got := base64Regex.MatchString(tt.input)
		if got != tt.match {
			t.Errorf("base64Regex.MatchString(%q) = %v, want %v", tt.input, got, tt.match)
		}
	}
}

func TestMatchers(t *testing.T) {
	m := Matchers()
	if len(m) != 1 {
		t.Fatalf("expected 1 matcher, got %d", len(m))
	}
	if m[0].ID != "base64" {
		t.Errorf("expected ID=base64, got %s", m[0].ID)
	}
	if !m[0].FullText {
		t.Error("expected FullText=true")
	}
}
