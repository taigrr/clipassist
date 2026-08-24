package watcher

import (
	"context"
	"fmt"
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

func Watch(ctx context.Context) error {
	err := clipboard.Init()
	if err != nil {
		return fmt.Errorf("initialize clipboard: %w", err)
	}
	watch := clipboard.Watch(ctx, clipboard.FmtText)
	runClipboardWatch(ctx, watch)
	return nil
}

func runClipboardWatch(ctx context.Context, watch <-chan clipboard.Data) {
	for {
		select {
		case clip, ok := <-watch:
			if !ok {
				return
			}
			sclip := string(clip.Bytes)
			storeClip(sclip)
			go matchers.Run(sclip)
		case <-ctx.Done():
			return
		}
	}
}

func WriteToClip(text string) error {
	clipboard.Write(clipboard.FmtText, []byte(text))
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
