package watcher

import (
	"context"
	"errors"
	"sync"

	"github.com/taigrr/clipassist/matchers"
	"golang.design/x/clipboard"
)

var (
	clipRing []string
	current  int
	clipLock sync.RWMutex
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
	watch := clipboard.Watch(ctx, clipboard.FmtText)
	for {
		select {
		case clip := <-watch:
			sclip := string(clip)
			storeClip(sclip)
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
	clipLock.RLock()
	defer clipLock.RUnlock()

	index = normalizeClipIndex(index)
	return clipRing[index]
}

func storeClip(text string) {
	clipLock.Lock()
	defer clipLock.Unlock()

	current++
	current %= clipRingSize
	clipRing[current] = text
}

func normalizeClipIndex(index int) int {
	index %= clipRingSize
	if index < 0 {
		index += clipRingSize
	}
	return index
}
