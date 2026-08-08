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
		{"aGVsbG8td29ybGQ_", "hello-world?", true},
		{"dXJsLXNhZmV-fn5-", "url-safe~~~~", true},
		{"not-base64!!!", "", false},
		{"abcdefghijklmnop", "", false},
		// A UUID matches the url-safe alphabet and length but decodes to
		// binary garbage, so it must be rejected rather than firing a
		// spurious notification.
		{"550e8400-e29b-41d4-a716-446655440000", "", false},
		// Mixed std (+) and url-safe (_) alphabets can't decode under any
		// single encoding.
		{"abc+def_ghij1234", "", false},
		// Decodes to valid UTF-8 containing a control character (ESC), so
		// it must be rejected by the text guard even though utf8.Valid is
		// true. This directly exercises the isText branch.
		{"YWJjG2RlZmdoaWpr", "", false},
		// Decodes to text containing a non-breaking space (U+00A0), which
		// is accepted because it is Unicode whitespace.
		{"aGVsbG/CoHdvcmxkIQ==", "hello\u00a0world!", true},
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
		{"aGVsbG8td29ybGQ_", true},
		{"dXJsLXNhZmV-fn5-", true},
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
		{"aGVsbG8td29ybGQ_", true},
		{"dXJsLXNhZmV-fn5-", true},
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
