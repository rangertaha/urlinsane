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
	"github.com/rangertaha/urlinsane/internal"
)

// Domain ...
type Domain struct {
	prefix string
	name   string
	suffix string
}

//	type Domain interface {
//		Prefix(...string) string
//		Name(...string) string
//		Suffix(...string) string
//		String() string
//		Valid() bool
//		Live() bool
//	}
func New(name string) internal.Domain {
	return &Domain{
		prefix: domainutil.Subdomain(name),
		name:   domainutil.DomainPrefix(name),
		suffix: domainutil.DomainSuffix(name),
	}
}

// func New(prefix, name, suffix string) internal.Domain {
// 	name = fmt.Sprintf("%s.%s.%s", prefix, name, suffix)
// 	name = strings.ReplaceAll(name, "..", ".")
// 	name = strings.Trim(name, ".")

// 	domain := &Domain{
// 		prefix: domainutil.Subdomain(name),
// 		name:   domainutil.DomainPrefix(name),
// 		suffix: domainutil.DomainSuffix(name),
// 	}
// 	// domain.Fqdn()
// 	return domain
// }

func (d *Domain) Prefix(labels ...string) (name string) {
	return d.prefix
}

func (d *Domain) Name(labels ...string) (name string) {
	return d.name
}

func (d *Domain) Suffix(labels ...string) (name string) {
	return d.suffix
}

func (d *Domain) Valid() bool {
	return d.name != ""
}

func (d *Domain) String(labels ...string) (name string) {
	name = fmt.Sprintf("%s.%s.%s", d.prefix, d.name, d.suffix)
	name = strings.ReplaceAll(name, "..", ".")
	name = strings.Trim(name, ".")
	return
}

func (d *Domain) Live() (ip bool) {
	return true
}
