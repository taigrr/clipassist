package jwt

import (
	"testing"
)

func TestDecodeSegment(t *testing.T) {
	// {"alg":"HS256","typ":"JWT"} base64url-encoded
	seg := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	m, err := DecodeSegment(seg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["alg"] != "HS256" {
		t.Errorf("expected alg=HS256, got %v", m["alg"])
	}
	if m["typ"] != "JWT" {
		t.Errorf("expected typ=JWT, got %v", m["typ"])
	}
}

func TestDecodeSegmentInvalid(t *testing.T) {
	_, err := DecodeSegment("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
	_, err = DecodeSegment("bm90LWpzb24=") // "not-json"
	if err == nil {
		t.Error("expected error for non-JSON payload")
	}
}

func TestJWTRegex(t *testing.T) {
	// Standard JWT
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	if !jwtRegex.MatchString(token) {
		t.Error("expected JWT regex to match valid token")
	}
	if jwtRegex.MatchString("not.a.jwt") {
		t.Error("expected JWT regex not to match non-JWT string")
	}
}

func TestMatchers(t *testing.T) {
	m := Matchers()
	if len(m) != 1 {
		t.Fatalf("expected 1 matcher, got %d", len(m))
	}
	if m[0].ID != "jwt" {
		t.Errorf("expected ID=jwt, got %s", m[0].ID)
	}
}
