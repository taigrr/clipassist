package matchers

import (
	"regexp"
	"sync"
)

// A matcher is a regular expression tied to a function
// if the regular expression matches exactly to the input text,
// the function is called with the input.
// If the regular expression does not match, but the FullText boolean
// is not set to true, substrings will then be matched against the regexp.
// For each submatch, the function F is called.
type Matcher struct {
	Regex    *regexp.Regexp
	FullText bool
	F        func(string)
	ID       string
}

var (
	matchers  []Matcher
	matchLock sync.Mutex
)

func Get() []Matcher {
	matchLock.Lock()
	defer matchLock.Unlock()

	snapshot := make([]Matcher, len(matchers))
	copy(snapshot, matchers)
	return snapshot
}

func Add(a ...Matcher) {
	matchLock.Lock()
	defer matchLock.Unlock()

	matchers = append(matchers, a...)
}

func Remove(id string) {
	matchLock.Lock()
	defer matchLock.Unlock()

	toKeep := make([]Matcher, 0, len(matchers))
	for _, matcher := range matchers {
		if matcher.ID != id {
			toKeep = append(toKeep, matcher)
		}
	}
	matchers = toKeep
}

func Run(in string) {
	matchLock.Lock()
	snapshot := make([]Matcher, len(matchers))
	copy(snapshot, matchers)
	matchLock.Unlock()

	for _, matcher := range snapshot {
		if matcher.FullText {
			if matcher.Regex.MatchString(in) && matcher.Regex.FindString(in) == in {
				matcher.F(in)
			}
			continue
		}
		matches := matcher.Regex.FindAllString(in, -1)
		for _, m := range matches {
			matcher.F(m)
		}
	}
}
