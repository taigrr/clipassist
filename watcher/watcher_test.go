package watcher

import (
	"testing"
)

func TestClipRingInit(t *testing.T) {
	if len(clipRing) != clipRingSize {
		t.Errorf("expected clipRing length %d, got %d", clipRingSize, len(clipRing))
	}
}

func TestGetClipAtIndex_Empty(t *testing.T) {
	// Reset ring for test isolation
	for i := range clipRing {
		clipRing[i] = ""
	}
	current = 0

	got := GetClipAtIndex(0)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestGetClipAtIndex_Wraps(t *testing.T) {
	for i := range clipRing {
		clipRing[i] = ""
	}
	clipRing[0] = "wrapped"

	// clipRingSize should wrap to 0
	got := GetClipAtIndex(clipRingSize)
	if got != "wrapped" {
		t.Errorf("expected %q at wrapped index, got %q", "wrapped", got)
	}
}

func TestClipRingStorage(t *testing.T) {
	// Reset ring
	for i := range clipRing {
		clipRing[i] = ""
	}
	current = 0

	// Simulate clipboard entries being added (mimicking Watch loop logic)
	entries := []string{"first", "second", "third"}
	for _, entry := range entries {
		current++
		current %= clipRingSize
		clipRing[current] = entry
	}

	if GetClipAtIndex(1) != "first" {
		t.Errorf("expected %q at index 1, got %q", "first", GetClipAtIndex(1))
	}
	if GetClipAtIndex(2) != "second" {
		t.Errorf("expected %q at index 2, got %q", "second", GetClipAtIndex(2))
	}
	if GetClipAtIndex(3) != "third" {
		t.Errorf("expected %q at index 3, got %q", "third", GetClipAtIndex(3))
	}
}

func TestClipRingOverflow(t *testing.T) {
	// Reset ring
	for i := range clipRing {
		clipRing[i] = ""
	}
	current = 0

	// Fill the ring past capacity
	for i := range clipRingSize + 5 {
		current++
		current %= clipRingSize
		clipRing[current] = string(rune('A' + i%26))
	}

	// Verify ring wrapped — index 1 should have been overwritten
	// After 55 entries: current = 55 % 50 = 5
	// Entries 51-55 wrote to indices 1-5
	if current != 5 {
		t.Errorf("expected current to be 5 after overflow, got %d", current)
	}

	// Index 1 should have the 51st entry (index 50 in 0-based, 50%26 = 24 = 'Y')
	got := GetClipAtIndex(1)
	expected := string(rune('A' + 50%26))
	if got != expected {
		t.Errorf("expected %q at index 1 after overflow, got %q", expected, got)
	}
}

func TestClipRingSize(t *testing.T) {
	if clipRingSize != 50 {
		t.Errorf("expected clipRingSize to be 50, got %d", clipRingSize)
	}
}
