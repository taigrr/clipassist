package watcher

import (
	"context"
	"sync"

	"github.com/taigrr/clipassist/matchers"
	"golang.design/x/clipboard"
)

var (
	clipRing []string
	current  int
	clipLock sync.RWMutex

	writeClipboard = clipboard.Write
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
			sclip := string(clip.Bytes)
			storeClip(sclip)
			go matchers.Run(sclip)
		case <-ctx.Done():
			return
		}
	}
}

func WriteToClip(text string) error {
	_, err := writeClipboard(context.Background(), clipboard.FmtText, []byte(text))
	return err
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
