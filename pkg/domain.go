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
package pkg

import (
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
)


// Domain ...
type Domain struct {
	subdomain string
	domain    string
	suffix    string
	meta      map[string]interface{}
	live      bool
}

func NewDomain(str string) (d *Domain) {
	str = strings.TrimSpace(str)
	d = &Domain{
		meta:      make(map[string]interface{}),
		subdomain: domainutil.Subdomain(str),
		domain:    domainutil.DomainPrefix(str),
		suffix:    domainutil.DomainSuffix(str),
	}
	if d.domain == "" {
		d.domain = str
	}
	return
}

func (d *Domain) Subdomain() string {
	return d.subdomain
}

func (d *Domain) Domain() string {
	return d.domain
}
func (d *Domain) Suffix() string {
	return d.suffix
}
func (d *Domain) Live() bool {
	return d.live
}

func (d *Domain) Meta() map[string]interface{} {
	return d.meta
}

func (d *Domain) Add(key string, value interface{}) {
	d.meta[key] = value
}

// Repr returns a printable representational string of the given domain
func (d *Domain) Repr() (domain string) {
	if d.subdomain != "" {
		domain = d.subdomain + "."
	}
	if d.domain != "" {
		domain = domain + d.domain
	}
	if d.suffix != "" {
		domain = domain + "." + d.suffix
	}
	return
}
