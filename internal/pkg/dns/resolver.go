// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package dns

import (
	"context"
	"math/rand"
	"net"
	"strings"
)

// Resolver is the DNS resolver used by the collectors. It defaults to the
// system resolver but can be pointed at custom servers with SetResolver
// (wired from the --nameservers flag).
var Resolver = net.DefaultResolver

// SetResolver points Resolver at the given DNS servers (host or host:port;
// port 53 is assumed when omitted). Queries are spread across the servers.
// An empty list restores the system resolver.
func SetResolver(servers []string) {
	var addrs []string
	for _, s := range servers {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if !strings.Contains(s, ":") {
			s += ":53"
		}
		addrs = append(addrs, s)
	}
	if len(addrs) == 0 {
		Resolver = net.DefaultResolver
		return
	}
	Resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, network, addrs[rand.Intn(len(addrs))])
		},
	}
}
