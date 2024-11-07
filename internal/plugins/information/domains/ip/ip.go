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
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
<<<<<<< HEAD
	// "github.com/rangertaha/urlinsane/pkg/dns/resolver"
=======
>>>>>>> develop
)

const (
	ORDER       = 0 // We need this to run first
	CODE        = "ip"
	NAME        = "Ip Address"
	DESCRIPTION = "Domain IP addresses"
)

type Ipaddr struct {
	// resolver resolver.Client
	conf internal.Config
}

func (n *Ipaddr) Id() string {
	return CODE
}

func (n *Ipaddr) Order() int {
	return ORDER
}

func (i *Ipaddr) Init(c internal.Config) {
	i.conf = c
}

func (n *Ipaddr) Description() string {
	return DESCRIPTION
}

func (n *Ipaddr) Headers() []string {
	return []string{"A"}
}

func (i *Ipaddr) Exec(in internal.Typo) (out internal.Typo) {
	if name := in.Variant().Name(); name != "" {
		ips, err := net.LookupIP(name)
		if err == nil {
			var answers []string
			for _, ip := range ips {
				answers = append(answers, ip.String())
			}
			in.Variant().Add("A", strings.Join(answers, " "))
			in.Variant().Live(true)
		}

	}
	return in
}

func (i *Ipaddr) Close() {}

// Register the plugin
func init() {
	domains.Add(CODE, func() internal.Information {
		return &Ipaddr{}
	})
}
