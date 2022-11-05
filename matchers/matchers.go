package matchers

import (
	"regexp"
	"sync"
)

type Matcher struct {
	Regex *regexp.Regexp
	F     func(string)
	ID    string
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
		if matcher.Regex.MatchString(in) {
			matcher.F(in)
		}
	}
	matchLock.Unlock()
}
