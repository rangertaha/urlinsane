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
package cn

import (
	"net"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
)

const (
	CODE        = "cn"
	DESCRIPTION = "DNS CNAME records"
)

type Ipaddr struct {
	conf internal.Config
}

func (n *Ipaddr) Id() string {
	return CODE
}

func (i *Ipaddr) Init(c internal.Config) {
	i.conf = c
}

func (n *Ipaddr) Description() string {
	return DESCRIPTION
}

func (n *Ipaddr) Headers() []string {
	return []string{"CNAME"}
}

func (i *Ipaddr) Exec(in internal.Typo) (out internal.Typo) {
	if name := in.Variant().Name(); name != "" {
		if cname, err := net.LookupCNAME(name); err == nil {
			in.Variant().Add("CNAME", cname)
			in.Variant().Live(true)
		}
	}
	return in
}

// Register the plugin
func init() {
	domains.Add(CODE, func() internal.Information {
		return &Ipaddr{}
	})
}
