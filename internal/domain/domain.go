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
	"encoding/json"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/rangertaha/urlinsane/internal"
	log "github.com/sirupsen/logrus"
)

// Domain ...
type Domain struct {
	PreName  string `json:"prefix,omitempty"`
	Domain   string `json:"name,omitempty"`
	SufName  string `json:"suffix,omitempty"`
	FQDN     string `json:"fqdn"`
	Punycode string `json:"idn"`

	IsLive bool `json:"live,omitempty"`
	// IPv4       []IP   `json:"ipv4,omitempty"`
	// IPv6       []IP   `json:"ipv6,omitempty"`
	// Banner     Banner `json:"response,omitempty"`
	// Screenshot string `json:"screenshot,omitempty"`
	// Html       string `json:"html,omitempty"`
	// Ssdeep     string `json:"ssdeep,omitempty"`

	// Dns   []DnsRecord   `json:"dns,omitempty"`
	// Whois []WhoisRecord `json:"whois,omitempty"`

	algo        internal.Algorithm
	meta        map[string]interface{}
	levenshtein int
	active      bool
}

// type Domain struct {
// 	prefix string
// 	name   string
// 	suffix string

// 	algo        internal.Algorithm
// 	meta        map[string]interface{}
// 	levenshtein int
// 	live        bool
// 	active      bool
// }

func New(name string) internal.Domain {
	return &Domain{
		FQDN:    name,
		PreName: domainutil.Subdomain(name),
		Domain:  domainutil.DomainPrefix(name),
		SufName: domainutil.DomainSuffix(name),
		meta:    make(map[string]interface{}),
	}
}

func NewVariant(algo internal.Algorithm, names ...string) internal.Domain {
	name := strings.Join(names, ".")
	return &Domain{
		FQDN:    name,
		PreName: domainutil.Subdomain(name),
		Domain:  domainutil.DomainPrefix(name),
		SufName: domainutil.DomainSuffix(name),
		meta:    make(map[string]interface{}),
		algo:    algo,
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
		d.PreName = name
	}

	return d.PreName
}

func (d *Domain) Name(labels ...string) string {
	if len(labels) > 0 {
		name := strings.Join(labels, ".")
		d.Domain = name
	}

	return d.Domain
}

func (d *Domain) Suffix(labels ...string) string {
	if len(labels) > 0 {
		name := strings.Join(labels, ".")
		d.SufName = name
	}

	return d.SufName
}

func (d *Domain) Valid() bool {
	return d.Domain != ""
}

func (d *Domain) String(labels ...string) (name string) {
	names := []string{d.PreName, d.Domain, d.SufName}
	name = strings.Join(names, ".")
	name = strings.ReplaceAll(name, "..", ".")
	name = strings.Trim(name, ".")
	return
}

func (d *Domain) Live(v ...bool) bool {
	if len(v) > 0 {
		d.IsLive = v[0]
	}

	return d.IsLive
}

func (d *Domain) Active(v ...bool) bool {
	if len(v) > 0 {
		d.active = v[0]
	}

	return d.active
}

// Ld returns the Levenshtein_distance
//
//	https://en.wikipedia.org/wiki/Levenshtein_distance
func (d *Domain) Ld(v ...int) int {
	if len(v) > 0 {
		d.levenshtein = v[0]
	}

	return d.levenshtein
}

func (d *Domain) Json() string {
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Errorf("Error:", err)
	}

	return string(jsonData)
}

func (d *Domain) Idn(names ...string) string {
	if len(names) > 0 {
		d.Punycode = names[0]
	}

	return d.Punycode
}
