# ClipAssist

> **Note:** this project integrates with the system clipboard and the
> notifications system of your device, and as a result, requires CGO.

This tool serves as a framework for operating on selected text.
It provides functionality to run in the background of your system and listen in
on your clipboard operations.
When a piece of text is copied that matches a specified regex, a function can be
called against that piece of text.

## Built-in Modules

| Module | Detects | Action |
|--------|---------|--------|
| **millis** | 13-digit Unix millisecond timestamps | Shows human-readable date/time |
| **seconds** | 10-digit Unix second timestamps | Shows human-readable date/time |
| **uuid** | UUIDs (v1–v7) | Shows UUID version info |
| **hexcolor** | CSS hex colors (`#rgb`, `#rrggbb`, `#rrggbbaa`) | Shows RGB(A) breakdown + OKLCH |
| **jwt** | JSON Web Tokens | Shows header info, subject, issuer, expiry |
| **ipaddr** | IPv4 & IPv6 addresses | Shows address type (private/public/loopback) |
| **base64** | Base64-encoded strings | Decodes and shows preview of content |

## Writing Your Own Module

A module is simply a function that returns `[]matchers.Matcher`. Each matcher
pairs a regular expression with a callback:

```go
package mymodule

import (
    "regexp"
    "github.com/taigrr/clipassist/matchers"
)

func Matchers() []matchers.Matcher {
    return []matchers.Matcher{
        {
            Regex: regexp.MustCompile(`pattern`),
            ID:    "my-module",
            F:     func(matched string) {
                // Do something with the matched text
            },
        },
    }
}
```

Then register it in `cmd/clipassist/main.go`:

```go
matchers.Add(mymodule.Matchers()...)
```

Set `FullText: true` on a matcher if you only want it to fire when the
**entire** clipboard content matches (not substrings).

## Installation

```bash
go install github.com/taigrr/clipassist/cmd/clipassist@latest
```

A sample systemd unit file to be placed into
`/etc/systemd/user/clipassist.service` is provided. Be sure to change the path
to the binary to match your home folder.
