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
package txt

import (
	"net"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/models"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const (
	ORDER       = 2
	CODE        = "txt"
	NAME        = "TXT Records"
	DESCRIPTION = "DNS MX Records"
)

type Plugin struct {
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
	return []string{"TXT"}
}

func (i *Plugin) Exec(in internal.Typo) (out internal.Typo) {
	orig, vari := in.Get()
	records, err := net.LookupTXT(vari.Fqdn())
	if err == nil {
		in.SetMeta("TXT", strings.Join(records, " "))
		vari.Live = true
		for _, r := range records {
			vari.Dns = append(vari.Dns, models.DnsRecord{Type: "TXT", Value: r})
		}
	}

	in.Set(orig, vari)
	return in
}

func (i *Plugin) Close() {}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &Plugin{}
	})
}
