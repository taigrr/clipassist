package matchers

import (
	"regexp"
	"sync"
	"testing"
)

func resetMatchers() {
	matchLock.Lock()
	matchers = nil
	matchLock.Unlock()
}

func TestAddAndGet(t *testing.T) {
	resetMatchers()
	m := Matcher{
		Regex: regexp.MustCompile(`\d+`),
		ID:    "digits",
		F:     func(string) {},
	}
	Add(m)
	got := Get()
	if len(got) != 1 {
		t.Fatalf("expected 1 matcher, got %d", len(got))
	}
	if got[0].ID != "digits" {
		t.Fatalf("expected ID 'digits', got %q", got[0].ID)
	}
}

func TestRemove(t *testing.T) {
	resetMatchers()
	Add(
		Matcher{Regex: regexp.MustCompile(`a`), ID: "a", F: func(string) {}},
		Matcher{Regex: regexp.MustCompile(`b`), ID: "b", F: func(string) {}},
	)
	Remove("a")
	got := Get()
	if len(got) != 1 || got[0].ID != "b" {
		t.Fatalf("expected only 'b' remaining, got %v", got)
	}
}

func TestRemoveNonExistent(t *testing.T) {
	resetMatchers()
	Add(Matcher{Regex: regexp.MustCompile(`a`), ID: "a", F: func(string) {}})
	Remove("nonexistent")
	if len(Get()) != 1 {
		t.Fatal("removing non-existent ID should not affect matchers")
	}
}

func TestGetReturnsSnapshot(t *testing.T) {
	resetMatchers()
	Add(Matcher{Regex: regexp.MustCompile(`a`), ID: "a", F: func(string) {}})

	snapshot := Get()
	snapshot[0].ID = "changed"
	mutated := append(snapshot, Matcher{Regex: regexp.MustCompile(`b`), ID: "b", F: func(string) {}})
	if len(mutated) != 2 {
		t.Fatalf("expected mutated snapshot length 2, got %d", len(mutated))
	}

	got := Get()
	if len(got) != 1 {
		t.Fatalf("expected internal matcher list to stay length 1, got %d", len(got))
	}
	if got[0].ID != "a" {
		t.Fatalf("expected internal matcher ID to remain %q, got %q", "a", got[0].ID)
	}
}

func TestRunSubstringMatch(t *testing.T) {
	resetMatchers()
	var results []string
	var mu sync.Mutex
	Add(Matcher{
		Regex: regexp.MustCompile(`\d{3}`),
		ID:    "three-digits",
		F: func(s string) {
			mu.Lock()
			results = append(results, s)
			mu.Unlock()
		},
	})
	Run("abc123def456ghi")
	if len(results) != 2 {
		t.Fatalf("expected 2 matches, got %d: %v", len(results), results)
	}
	if results[0] != "123" || results[1] != "456" {
		t.Fatalf("unexpected matches: %v", results)
	}
}

func TestRunFullTextMatch(t *testing.T) {
	resetMatchers()
	called := false
	Add(Matcher{
		Regex:    regexp.MustCompile(`^\d+$`),
		FullText: true,
		ID:       "all-digits",
		F:        func(string) { called = true },
	})

	// Should not match — input has non-digit characters
	Run("abc123")
	if called {
		t.Fatal("FullText matcher should not match substring")
	}

	// Should match — entire input is digits
	Run("123456")
	if !called {
		t.Fatal("FullText matcher should match when regex matches full input")
	}
}

func TestRunFullTextNoMatch(t *testing.T) {
	resetMatchers()
	called := false
	Add(Matcher{
		Regex:    regexp.MustCompile(`hello`),
		FullText: true,
		ID:       "hello-full",
		F:        func(string) { called = true },
	})
	// "hello world" contains "hello" but is not equal to it
	Run("hello world")
	if called {
		t.Fatal("FullText matcher should not fire when match != full input")
	}
}

func TestRunNoMatchers(t *testing.T) {
	resetMatchers()
	// Should not panic
	Run("anything")
}

func TestRunCallbackCanAddMatcher(t *testing.T) {
	resetMatchers()
	// Ensure that calling Add from within a callback doesn't deadlock
	Add(Matcher{
		Regex: regexp.MustCompile(`trigger`),
		ID:    "trigger",
		F: func(string) {
			Add(Matcher{
				Regex: regexp.MustCompile(`new`),
				ID:    "new",
				F:     func(string) {},
			})
		},
	})
	Run("trigger") // should not deadlock
	if len(Get()) != 2 {
		t.Fatalf("expected 2 matchers after callback Add, got %d", len(Get()))
	}
}

func TestConcurrentAddAndRun(t *testing.T) {
	resetMatchers()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			Add(Matcher{
				Regex: regexp.MustCompile(`x`),
				ID:    "x",
				F:     func(string) {},
			})
		}()
		go func() {
			defer wg.Done()
			Run("x")
		}()
	}
	wg.Wait()
}
