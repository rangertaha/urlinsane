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
	"github.com/rangertaha/urlinsane/internal/models"
)

// // Domain ...
// type Domain struct {
// 	Subdomain string
// 	Prefix    string
// 	Suffix    string
// }

func New(prefix, name, suffix string) (d models.Domain) {
	name = fmt.Sprintf("%s.%s.%s", prefix, name, suffix)
	name = strings.ReplaceAll(name, "..", ".")
	name = strings.Trim(name, ".")

	domain := models.Domain{
		Prefix: domainutil.Subdomain(name),
		Name:   domainutil.DomainPrefix(name),
		Suffix: domainutil.DomainSuffix(name),
	}
	domain.Fqdn()
	return domain
}

func Parse(name string) (d models.Domain) {
	domain := models.Domain{
		FQDN: name,
		Prefix: domainutil.Subdomain(name),
		Name:   domainutil.DomainPrefix(name),
		Suffix: domainutil.DomainSuffix(name),
	}
	return domain
}

// func (d *Domain) String() (name string) {
// 	name = fmt.Sprintf("%s.%s.%s", d.Subdomain, d.Prefix, d.Suffix)
// 	name = strings.ReplaceAll(name, "..", ".")
// 	name = strings.Trim(name, ".")
// 	return
// }
