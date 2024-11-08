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
	"log"
	"net"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/models"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const (
	ORDER       = 0 // We need this to run first
	CODE        = "ip"
	NAME        = "Ip Address"
	DESCRIPTION = "Domain IP addresses"
)

type Plugin struct {
	// resolver resolver.Client
	conf internal.Config
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Order() int {
	return ORDER
}

func (i *Plugin) Init(c internal.Config) {
	i.conf = c
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Headers() []string {
	return []string{"IPv4", "IPv6"}
}

func (i *Plugin) Exec(in internal.Typo) (out internal.Typo) {
	orig, vari := in.Get()
	var As []string
	var AAAAs []string

	vari.IPv4 = []models.IP{}
	vari.IPv6 = []models.IP{}
	vari.Dns = []models.DnsRecord{}
	ips, err := net.LookupIP(vari.Fqdn())
	if err != nil {
		log.Print(err)
	}
	for _, ip := range ips {
		dr := models.DnsRecord{}
		vari.Live = true
		if strings.Contains(ip.String(), ":") {
			AAAAs = append(AAAAs, ip.String())
			dr.Type = "AAAA"
			dr.Value = ip.String()
			vari.IPv6 = append(vari.IPv6, models.IP{Address: ip.String()})
			vari.Dns = append(vari.Dns, dr)

		} else if strings.Contains(ip.String(), ".") {
			As = append(As, ip.String())
			dr.Type = "A"
			dr.Value = ip.String()
			vari.IPv4 = append(vari.IPv6, models.IP{Address: ip.String()})
			vari.Dns = append(vari.Dns, dr)
		}
	}

	in.Set(orig, vari)
	in.SetMeta("IPv4", strings.Join(As, " "))
	in.SetMeta("IPv6", strings.Join(AAAAs, " "))

	return in
}

func (i *Plugin) Close() {}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &Plugin{}
	})
}
