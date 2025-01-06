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
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
)

type Plugin struct {
	collectors.Plugin
}

func (p *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {
	ips, err := net.LookupIP(domain.Name)
	if err != nil {
		p.Log.Error("IP Lookup: ", err)
	}

	
	server := db.Device{Name: domain.Name}

	for _, ip := range ips {
		record := strings.TrimSpace(ip.String())
		record = strings.Trim(record, ".")

		if strings.Contains(ip.String(), ":") {
			domain.Dns = append(domain.Dns, &db.DnsRecord{Type: "AAAA", Value: record})
			server.IPs = append(server.IPs, &db.IP{Address: record, Type: "IPv6"})

		} else if strings.Contains(ip.String(), ".") {
			domain.Dns = append(domain.Dns, &db.DnsRecord{Type: "A", Value: record})
			server.IPs = append(server.IPs, &db.IP{Address: record, Type: "IPv4"})

		}
	}
	domain.Servers = append(domain.Servers, &server)
	return domain, err
}

// Register the plugin
func init() {
	var CODE = "ip"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:       0,
				Code:      CODE,
				Title:     "Ip Address",
				Summary:   "Domain IPv4 and IPv6 addresses",
				DependsOn: []string{},
			},
		}
	})
}
