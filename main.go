package main

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/gen2brain/beeep"
	"golang.design/x/clipboard"
)

type Matcher struct {
	regex *regexp.Regexp
	f     func(string)
}

var matchers []Matcher

func init() {
}

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	watch := clipboard.Watch(context.Background(), clipboard.FmtText)
	matchers = LoadMatchers()
	for clip := range watch {
		RunMatchers(string(clip))
	}
}

func LoadMatchers() []Matcher {
	var m []Matcher
	millis := regexp.MustCompile(`\d\d\d\d\d\d\d\d\d\d\d\d\d`)
	m = append(m, Matcher{regex: millis, f: ConvertDate})
	return m
}

func RunMatchers(in string) {
	for _, m := range matchers {
		go func(matcher Matcher) {
			if matcher.regex.MatchString(in) {
				matcher.f(in)
			}
		}(m)
	}
}

func ConvertDate(in string) {
	ts, err := strconv.Atoi(in)
	if err != nil {
		return
	}
	t := time.UnixMilli(int64(ts))
	beeep.Alert("Millis Converted", t.String(), "")
	clipboard.Write(clipboard.FmtText, []byte(t.String()))
}
