package watcher

import (
	"context"

	"github.com/taigrr/clipassist/matchers"
	"golang.design/x/clipboard"
)

var (
	clipRing []string
	current  int
)

func init() {
	clipRing = make([]string, 50)
}

func Watch(ctx context.Context) {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	watch := clipboard.Watch(context.Background(), clipboard.FmtText)
	for {
		select {
		case clip := <-watch:
			sclip := string(clip)
			current++
			current %= 50
			clipRing[current] = sclip
			go matchers.Run(sclip)
		case <-ctx.Done():
			return
		}
	}
}
