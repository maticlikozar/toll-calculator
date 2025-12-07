package request

import (
	"net"
	"strings"

	"toll/internal/errlog"
)

func ParseIp(addr string) (string, error) {
	if !strings.Contains(addr, ":") {
		addr += ":"
	}

	// Force brackets for IPv6 localhost.
	if addr == "::1" {
		addr = "[::1]:"
	}

	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", errlog.Error(err)
	}

	return ip, nil
}
