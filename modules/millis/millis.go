package millis

import (
	"regexp"
	"strconv"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

func Matchers() []matchers.Matcher {
	var m []matchers.Matcher
	millis := regexp.MustCompile(`\d\d\d\d\d\d\d\d\d\d\d\d\d`)
	m = append(m, matchers.Matcher{
		F:     ConvertDate,
		ID:    "millis",
		Regex: millis,
	})
	return m
}

func ConvertDate(in string) {
	ts, err := strconv.Atoi(in)
	if err != nil {
		return
	}
	t := time.UnixMilli(int64(ts))
	beeep.Alert("Millis Converted", t.String(), "")
}
