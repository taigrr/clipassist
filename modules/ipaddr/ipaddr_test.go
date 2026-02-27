package ipaddr

import (
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"192.168.1.1", "IPv4 — Private"},
		{"10.0.0.1", "IPv4 — Private"},
		{"127.0.0.1", "IPv4 — Loopback"},
		{"8.8.8.8", "IPv4 — Public"},
		{"::1", "IPv6 — Loopback"},
		{"not-an-ip", ""},
	}
	for _, tt := range tests {
		got := Classify(tt.input)
		if got != tt.want {
			t.Errorf("Classify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIPRegex(t *testing.T) {
	tests := []struct {
		input string
		match bool
	}{
		{"192.168.1.1", true},
		{"255.255.255.255", true},
		{"8.8.8.8", true},
		{"999.999.999.999", false},
		{"hello", false},
	}
	for _, tt := range tests {
		got := ipRegex.MatchString(tt.input)
		if got != tt.match {
			t.Errorf("ipRegex.MatchString(%q) = %v, want %v", tt.input, got, tt.match)
		}
	}
}

func TestMatchers(t *testing.T) {
	m := Matchers()
	if len(m) != 1 {
		t.Fatalf("expected 1 matcher, got %d", len(m))
	}
	if m[0].ID != "ipaddr" {
		t.Errorf("expected ID=ipaddr, got %s", m[0].ID)
	}
}
