package watcher

import (
	"context"
	"sync"
	"testing"
	"time"

	"golang.design/x/clipboard"
)

func resetClipRing() {
	clipLock.Lock()
	defer clipLock.Unlock()

	for i := range clipRing {
		clipRing[i] = ""
	}
	current = 0
}

func TestClipRingInit(t *testing.T) {
	if len(clipRing) != clipRingSize {
		t.Errorf("expected clipRing length %d, got %d", clipRingSize, len(clipRing))
	}
}

func TestGetClipAtIndex_Empty(t *testing.T) {
	resetClipRing()

	got := GetClipAtIndex(0)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestGetClipAtIndex_Wraps(t *testing.T) {
	resetClipRing()
	storeClip("wrapped")

	got := GetClipAtIndex(clipRingSize + 1)
	if got != "wrapped" {
		t.Errorf("expected %q at wrapped index, got %q", "wrapped", got)
	}
}

func TestGetClipAtIndex_NegativeWraps(t *testing.T) {
	resetClipRing()
	storeClip("wrapped")

	got := GetClipAtIndex(-49)
	if got != "wrapped" {
		t.Errorf("expected %q at negative wrapped index, got %q", "wrapped", got)
	}
}

func TestClipRingStorage(t *testing.T) {
	resetClipRing()

	entries := []string{"first", "second", "third"}
	for _, entry := range entries {
		storeClip(entry)
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
	resetClipRing()

	for i := range clipRingSize + 5 {
		storeClip(string(rune('A' + i%26)))
	}

	clipLock.RLock()
	gotCurrent := current
	clipLock.RUnlock()

	if gotCurrent != 5 {
		t.Errorf("expected current to be 5 after overflow, got %d", gotCurrent)
	}

	// Index 1 should have the 51st entry (index 50 in 0-based, 50%26 = 24 = 'Y')
	got := GetClipAtIndex(1)
	expected := string(rune('A' + 50%26))
	if got != expected {
		t.Errorf("expected %q at index 1 after overflow, got %q", expected, got)
	}
}

func TestConcurrentStoreAndGet(t *testing.T) {
	resetClipRing()

	var waitGroup sync.WaitGroup
	for i := range 100 {
		waitGroup.Add(2)
		go func(value int) {
			defer waitGroup.Done()
			storeClip(string(rune('A' + value%26)))
		}(i)
		go func(index int) {
			defer waitGroup.Done()
			_ = GetClipAtIndex(index)
		}(i)
	}
	waitGroup.Wait()
}

func TestRunClipboardWatchReturnsWhenWatchCloses(t *testing.T) {
	watch := make(chan clipboard.Data)
	close(watch)

	done := make(chan struct{})
	go func() {
		runClipboardWatch(context.Background(), watch)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected runClipboardWatch to return when watch channel closes")
	}
}

func TestClipRingSize(t *testing.T) {
	if clipRingSize != 50 {
		t.Errorf("expected clipRingSize to be 50, got %d", clipRingSize)
	}
}
