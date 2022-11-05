package watcher

import (
	"context"
	"errors"

	"github.com/taigrr/clipassist/matchers"
	"golang.design/x/clipboard"
)

var (
	clipRing []string
	current  int
)

const clipRingSize = 50

func init() {
	clipRing = make([]string, clipRingSize)
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

func WriteToClip(text string) error {
	success := clipboard.Write(clipboard.FmtText, []byte(text))
	if success == nil {
		return errors.New("could not write to clipboard")
	}
	return nil
}

func GetClipAtIndex(index int) string {
	index %= clipRingSize
	return clipRing[index]
}
