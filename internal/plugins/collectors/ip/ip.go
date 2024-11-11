// Copyright (C) 2024 Rangertaha
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package ip

import (
	"net"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	log "github.com/sirupsen/logrus"
)

const (
	ORDER       = 1 // We need this to run first
	CODE        = "ip"
	NAME        = "Ip Address"
	DESCRIPTION = "Domain IP addresses"
)

type Plugin struct {
	// resolver resolver.Client
	conf internal.Config
	db   internal.Database
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Order() int {
	return ORDER
}

func (i *Plugin) Init(c internal.Config) {
	i.db = c.Database()
	i.conf = c
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Headers() []string {
	return []string{"IPv4", "IPv6"}
}

func (i *Plugin) Exec(domain internal.Domain, acc internal.Accumulator) (err error) {
	ipv4, _ := i.db.Read(domain.String(), "IPv4")
	ipv6, _ := i.db.Read(domain.String(), "IPv6")
	if ipv4 != "" {
		domain.SetMeta("IPv4", ipv4)
		domain.SetMeta("IPv6", ipv6)
		domain.Live(true)
		acc.Add(domain)

		return
	}

	ipv4, ipv6 = i.getIp(domain.String())
	domain.SetMeta("IPv4", ipv4)
	domain.SetMeta("IPv6", ipv6)
	domain.Live(true)
	acc.Add(domain)

	_ = i.db.Write(ipv4, domain.String(), "IPv4")
	err = i.db.Write(ipv6, domain.String(), "IPv6")

	return
}

func (i *Plugin) Close() {}

func (i *Plugin) getIp(d string) (v4, v6 string) {
	var As []string
	var AAAAs []string

	ips, err := net.LookupIP(d)
	if err != nil {
		log.Error("IP Lookup: ", err)
	}
	for _, ip := range ips {
		if strings.Contains(ip.String(), ":") {
			AAAAs = append(AAAAs, ip.String())

		} else if strings.Contains(ip.String(), ".") {
			As = append(As, ip.String())
		}
	}
	return strings.Join(As, " "), strings.Join(AAAAs, " ")
}

// Register the plugin
func init() {
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{}
	})
}
