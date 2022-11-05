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
	return matchers
}

func Add(a ...Matcher) {
	matchLock.Lock()
	matchers = append(matchers, a...)
	matchLock.Unlock()
}

func Remove(id string) {
	toKeep := []Matcher{}
	matchLock.Lock()
	for _, m := range matchers {
		if m.ID != id {
			toKeep = append(toKeep, m)
		}
	}
	matchers = toKeep
	matchLock.Unlock()
}

func Run(in string) {
	matchLock.Lock()
	for _, matcher := range matchers {

		matches := matcher.Regex.FindAllString(in, -1)
		for _, m := range matches {
			if matcher.FullText {
				if m != in {
					break
				}
			}
			matcher.F(m)
		}
	}
	matchLock.Unlock()
}
