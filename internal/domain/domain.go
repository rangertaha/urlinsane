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
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/rangertaha/urlinsane/internal"
)

// Domain ...
type Domain struct {
	prefix string
	name   string
	suffix string

	algo internal.Algorithm
	meta map[string]interface{}
	live bool
}

func New(name string) internal.Domain {
	return &Domain{
		prefix: domainutil.Subdomain(name),
		name:   domainutil.DomainPrefix(name),
		suffix: domainutil.DomainSuffix(name),
		meta:   make(map[string]interface{}),
	}
}

func NewVariant(algo internal.Algorithm, names ...string) internal.Domain {
	name := strings.Join(names, ".")
	return &Domain{
		prefix: domainutil.Subdomain(name),
		name:   domainutil.DomainPrefix(name),
		suffix: domainutil.DomainSuffix(name),
		meta:   make(map[string]interface{}),
		algo:   algo,
	}
}

func (t *Domain) Meta() map[string]interface{} {
	return t.meta
}

func (t *Domain) SetMeta(key string, value interface{}) {
	t.meta[key] = value
}

func (t *Domain) GetMeta(key string) (value interface{}) {
	if value, ok := t.meta[key]; ok {
		return value
	}
	return nil
}

func (t *Domain) Algorithm() internal.Algorithm {
	return t.algo
}

func (d *Domain) Prefix(labels ...string) string {
	if len(labels) > 0 {
		name := strings.Join(labels, ".")
		d.prefix = name
	}

	return d.prefix
}

func (d *Domain) Name(labels ...string) string {
	if len(labels) > 0 {
		name := strings.Join(labels, ".")
		d.name = name
	}

	return d.name
}

func (d *Domain) Suffix(labels ...string) string {
	if len(labels) > 0 {
		name := strings.Join(labels, ".")
		d.suffix = name
	}

	return d.suffix
}

func (d *Domain) Valid() bool {
	return d.name != ""
}

func (d *Domain) String(labels ...string) (name string) {
	names := []string{d.prefix, d.name, d.suffix}
	name = strings.Join(names, ".")
	name = strings.ReplaceAll(name, "..", ".")
	name = strings.Trim(name, ".")
	return
}

func (d *Domain) Live(v ...bool) (ip bool) {
	if len(v) > 0 {
		d.live = v[0]
	}

	return d.live
}
