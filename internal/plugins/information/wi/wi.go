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
package wi

// func whoisLookupFunc(tr Result) (results []Result) {
// 	return
// }

import (
	"fmt"

	"github.com/likexian/whois"
	parser "github.com/likexian/whois-parser"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/models"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const (
	ORDER       = 2
	CODE        = "wi"
	NAME        = "Whois"
	DESCRIPTION = "Whois database search"
)

type Plugin struct {
	types []string
}

func (n *Plugin) Order() int {
	return ORDER
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Name() string {
	return NAME
}
func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Headers() []string {
	return []string{"WHOIS"}
}

func (n *Plugin) Exec(in internal.Typo) (out internal.Typo) {
	orig, vari := in.Get()
	raw, err := whois.Whois(vari.Fqdn())
	if err == nil {
		fmt.Println(raw)
	}

	result, err := parser.Parse(raw)
	if result.Domain != nil {

	}

	if err == nil {
		if result.Domain != nil {
			// Print the domain status
			fmt.Println(result.Domain.Status)

			// Print the domain created date
			fmt.Println(result.Domain.CreatedDate)

			// Print the domain expiration date
			fmt.Println(result.Domain.ExpirationDate)
		}

		if result.Registrar != nil {
			// Print the registrar name
			fmt.Println(result.Registrar.Name)
		}

		if result.Registrant != nil {
			// Print the registrant name
			fmt.Println(result.Registrant.Name)

			// Print the registrant email address
			fmt.Println(result.Registrant.Email)
		}

	}

	in.Set(orig, vari)
	return in
}

func (n *Plugin) Whois(name string) (rec models.WhoisRecord) {
	raw, err := whois.Whois(name)
	if err == nil {
		fmt.Println(raw)
	}

	result, err := parser.Parse(raw)
	if result.Domain != nil {

	}

	if err == nil {
		if result.Domain != nil {
			// Print the domain status
			fmt.Println(result.Domain.Status)

			// Print the domain created date
			fmt.Println(result.Domain.CreatedDate)

			// Print the domain expiration date
			fmt.Println(result.Domain.ExpirationDate)
		}

		if result.Registrar != nil {
			// Print the registrar name
			fmt.Println(result.Registrar.Name)
		}

		if result.Registrant != nil {
			// Print the registrant name
			fmt.Println(result.Registrant.Name)

			// Print the registrant email address
			fmt.Println(result.Registrant.Email)
		}

	}

	return
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &Plugin{}
	})
}
