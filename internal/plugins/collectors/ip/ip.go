// Copyright 2024 Rangertaha. All Rights Reserved.
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
	"github.com/rangertaha/urlinsane/internal/pkg"
	log "github.com/sirupsen/logrus"
)

type Plugin struct {
	// resolver resolver.Client
	conf internal.Config
	log  *log.Entry
}

func (i *Plugin) Init(c internal.Config) {
	i.log = log.WithFields(log.Fields{"plugin": CODE, "method": "Exec"})
	i.conf = c
}

func (i *Plugin) Exec(acc internal.Accumulator) (err error) {
	l := i.log.WithFields(log.Fields{"domain": acc.Domain().String()})

	dns := make(pkg.DnsRecords, 0)
	if err := acc.Unmarshal("DNS", &dns); err != nil {
		l.Error("Unmarshal DNS: ", err)
	}

	// Retrive data
	// dns := make(pkg.DnsRecords, 0)
	ipv4, ipv6 := i.getIp(acc.Domain().String())
	if len(ipv4) > 0 || len(ipv6) > 0 {
		dns.Add("A", 0, ipv4...)
		dns.Add("AAAA", 0, ipv6...)

		// Add simple table data
		acc.SetMeta("IPv4", dns.String("A"))
		acc.SetMeta("IPv6", dns.String("AAAA"))

		// Added to JSON data
		acc.SetJson("DNS", dns.Json())
		acc.Domain().Live(true)
	}

	return acc.Next()
}

func (i *Plugin) Close() {}

func (i *Plugin) getIp(d string) (v4, v6 []string) {
	var As []string
	var AAAAs []string

	ips, err := net.LookupIP(d)
	if err != nil {
		i.log.Error("IP Lookup: ", err)
	}
	for _, ip := range ips {
		if strings.Contains(ip.String(), ":") {
			AAAAs = append(AAAAs, ip.String())

		} else if strings.Contains(ip.String(), ".") {
			As = append(As, ip.String())
		}
	}
	return As, AAAAs
}
