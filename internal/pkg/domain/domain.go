// // Copyright (C) 2024 Rangertaha
// //
// // This program is free software: you can redistribute it and/or modify
// // it under the terms of the GNU General Public License as published by
// // the Free Software Foundation, either version 3 of the License, or
// // (at your option) any later version.
// //
// // This program is distributed in the hope that it will be useful,
// // but WITHOUT ANY WARRANTY; without even the implied warranty of
// // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// // GNU General Public License for more details.
// //
// // You should have received a copy of the GNU General Public License
// // along with this program.  If not, see <http://www.gnu.org/licenses/>.
package domain

import (
	"fmt"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
)

// Domain ...
type Domain struct {
	Subdomain string
	Prefix    string
	Suffix    string
}

func New(sub, prefix, suffix string) (d *Domain) {

	// if len(names) == 3 {
	// 	domain := fmt.Sprintf("%s.%s.%s", names[0], names[1], names[2])
	// 	d.Subdomain = domainutil.Subdomain(domain)
	// 	d.Prefix = domainutil.DomainPrefix(domain)
	// 	d.Suffix = domainutil.DomainSuffix(domain)
	// }
	// if len(names) == 2 {
	// 	domain := fmt.Sprintf("%s.%s", names[0], names[1])
	// 	d.Subdomain = domainutil.Subdomain(domain)
	// 	d.Prefix = domainutil.DomainPrefix(domain)
	// 	d.Suffix = domainutil.DomainSuffix(domain)
	// }

	// if len(names) == 1 {
	// 	domain := names[0]
	// 	d.Subdomain = domainutil.Subdomain(domain)
	// 	d.Prefix = domainutil.DomainPrefix(domain)
	// 	d.Suffix = domainutil.DomainSuffix(domain)
	// }
	name := fmt.Sprintf("%s.%s.%s", sub, prefix, suffix)
	name = strings.ReplaceAll(name, "..", ".")
	name = strings.Trim(name, ".")

	return &Domain{
		Subdomain: domainutil.Subdomain(name),
		Prefix:    domainutil.DomainPrefix(name),
		Suffix:    domainutil.DomainSuffix(name),
	}
}

func Parse(name string) (d *Domain) {
	return &Domain{
		Subdomain: domainutil.Subdomain(name),
		Prefix:    domainutil.DomainPrefix(name),
		Suffix:    domainutil.DomainSuffix(name),
	}
}

func (d *Domain) String() (name string) {
	name = fmt.Sprintf("%s.%s.%s", d.Subdomain, d.Prefix, d.Suffix)
	name = strings.ReplaceAll(name, "..", ".")
	name = strings.Trim(name, ".")
	return
}
