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
	// "github.com/rangertaha/urlinsane/pkg/dns/resolver"
)

const (
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

func (i *Ipaddr) Init(c internal.Config) {
	i.conf = c
	// i.resolver = resolver.New(c.DnsServers(), 3, 1000, 50)
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
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		// 	// os.Exit(1)
		// }
		// for _, ip := range ips {
		// 	fmt.Printf("google.com. IN A %s\n", ip.String())
		// }
		if err == nil {
			var answers []string
			for _, ip := range ips {
				answers = append(answers, ip.String())
			}
			in.Variant().Add("A", strings.Join(answers, " "))
			in.Variant().Live(true)
		}

	}

	// i.resolver = resolver.New(i.conf.DnsServers(), 3, 1000, 50)
	// domains := []string{in.Variant().Name()}
	// results := i.resolver.Resolve(domains, resolver.TypeA)
	// var answers []string
	// for _, record := range results {
	// 	answers = append(answers, record.Answer)
	// }
	// in.Variant().Add("A", strings.Join(answers, "\n"))
	// // defer i.resolver.Close()
	return in
}

func (i *Ipaddr) Close() {
	// i.resolver.Close()
}

// Register the plugin
func init() {
	domains.Add(CODE, func() internal.Information {
		return &Ipaddr{}
	})
}
