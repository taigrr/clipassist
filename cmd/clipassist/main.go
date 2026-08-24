// Command clipassist watches the system clipboard and runs registered matcher
// modules against copied text, surfacing helpful information via desktop
// notifications.
package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"github.com/taigrr/clipassist/matchers"
	"github.com/taigrr/clipassist/modules/base64decode"
	"github.com/taigrr/clipassist/modules/hexcolor"
	"github.com/taigrr/clipassist/modules/ipaddr"
	"github.com/taigrr/clipassist/modules/jwt"
	"github.com/taigrr/clipassist/modules/millis"
	"github.com/taigrr/clipassist/modules/seconds"
	"github.com/taigrr/clipassist/modules/uuid"
	"github.com/taigrr/clipassist/watcher"
)

// version is overridable via -ldflags; it otherwise falls back to the module
// version embedded by the Go toolchain.
var version = "devel"

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			version = v
		}
	}
}

func registerModules() {
	matchers.Add(millis.Matchers()...)
	matchers.Add(seconds.Matchers()...)
	matchers.Add(uuid.Matchers()...)
	matchers.Add(hexcolor.Matchers()...)
	matchers.Add(jwt.Matchers()...)
	matchers.Add(ipaddr.Matchers()...)
	matchers.Add(base64decode.Matchers()...)
}

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clipassist",
		Short: "Watch the clipboard and run matcher modules against copied text",
		Long: "clipassist runs in the background, watches the system clipboard, " +
			"and invokes registered matcher modules when copied text matches " +
			"their patterns.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			registerModules()

			ctx, stop := signal.NotifyContext(
				cmd.Context(), os.Interrupt, syscall.SIGTERM,
			)
			defer stop()

			return watcher.Watch(ctx)
		},
	}
}

func main() {
	if err := fang.Execute(
		context.Background(),
		newRootCmd(),
		fang.WithVersion(version),
	); err != nil {
		os.Exit(1)
	}
}
