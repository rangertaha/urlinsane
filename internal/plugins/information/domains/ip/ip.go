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

// // ipLookupFunc
// func ipLookupFunc(tr Result) (results []Result) {
// 	results = append(results, checkIP(tr))
// 	return
// }

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const (
	CODE        = "ip"
	NAME        = "Ip Address"
	DESCRIPTION = "Domain IP addresses"
)

type Ipaddr struct {
	types []string
}

func (n *Ipaddr) Id() string {
	return CODE
}

func (n *Ipaddr) Name() string {
	return NAME
}

func (n *Ipaddr) Description() string {
	return DESCRIPTION
}

func (n *Ipaddr) Headers() []string {
	return []string{"Online", "IPv4", "IPv6"}
}

func (n *Ipaddr) Exec(in internal.Typo) (out internal.Typo) {

	in.Variant().Add("Online", true)
	in.Variant().Add("IPv4", "100.0.0.0")
	in.Variant().Add("IPv6", "100.0.0.0")
	in.Variant().Add("JSON", "{}")
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &Ipaddr{}
	})
}
