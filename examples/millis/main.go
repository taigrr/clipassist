package main

import (
	"context"

	"github.com/taigrr/clipassist/matchers"
	"github.com/taigrr/clipassist/modules/millis"
	"github.com/taigrr/clipassist/watcher"
)

func main() {
	ctx := context.Background()
	matchers.Add(millis.Matchers()...)
	watcher.Watch(ctx)
}
