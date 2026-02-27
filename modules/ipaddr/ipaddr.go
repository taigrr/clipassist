// Package ipaddr provides a clipboard matcher that detects IPv4 and IPv6
// addresses and shows a notification with address type information.
package ipaddr

import (
	"fmt"
	"net"
	"regexp"

	"github.com/gen2brain/beeep"
	"github.com/taigrr/clipassist/matchers"
)

var ipRegex = regexp.MustCompile(`(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?)|(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(?:[0-9a-fA-F]{1,4}:){1,7}:|::(?:[0-9a-fA-F]{1,4}:){0,5}[0-9a-fA-F]{1,4}`)

// Matchers returns a Matcher for IP addresses.
func Matchers() []matchers.Matcher {
	return []matchers.Matcher{
		{
			Regex: ipRegex,
			ID:    "ipaddr",
			F:     Notify,
		},
	}
}

// Classify returns a human-readable classification of the IP address.
func Classify(in string) string {
	ip := net.ParseIP(in)
	if ip == nil {
		return ""
	}

	var version string
	if ip.To4() != nil {
		version = "IPv4"
	} else {
		version = "IPv6"
	}

	var kind string
	switch {
	case ip.IsLoopback():
		kind = "Loopback"
	case ip.IsPrivate():
		kind = "Private"
	case ip.IsMulticast():
		kind = "Multicast"
	case ip.IsLinkLocalUnicast():
		kind = "Link-Local"
	case ip.IsUnspecified():
		kind = "Unspecified"
	default:
		kind = "Public"
	}

	return fmt.Sprintf("%s — %s", version, kind)
}

// Notify classifies an IP address and sends a notification.
func Notify(in string) {
	info := Classify(in)
	if info == "" {
		return
	}
	beeep.Alert("IP Address", fmt.Sprintf("%s\n%s", in, info), "")
}
