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

	IsLive   bool                       `json:"live,omitempty"`
	Data     map[string]json.RawMessage `json:"data,omitempty"` // used for detailed JSON nested outputs
	Metadata map[string]string          `json:"-"`              // Used for simplified table based output

	// Internal use
	input       string // Name provded by user
	algo        internal.Algorithm
	levenshtein int
	active      bool
}

// type DnsRecord struct {
// 	Type  string `json:"type,omitempty"`
// 	Value string `json:"value,omitempty"`
// 	TTL   int    `json:"ttl,omitempty"`
// }

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
		FQDN:     name,
		PreName:  domainutil.Subdomain(name),
		Domain:   domainutil.DomainPrefix(name),
		SufName:  domainutil.DomainSuffix(name),
		Metadata: make(map[string]string),
		Data:     make(map[string]json.RawMessage),
		input:    name,
	}
}

func NewVariant(algo internal.Algorithm, names ...string) internal.Domain {
	name := strings.Join(names, ".")
	return &Domain{
		FQDN:     name,
		PreName:  domainutil.Subdomain(name),
		Domain:   domainutil.DomainPrefix(name),
		SufName:  domainutil.DomainSuffix(name),
		Metadata: make(map[string]string),
		Data:     make(map[string]json.RawMessage),
		algo:     algo,
		input:    name,
	}
}

func (t *Domain) Meta() map[string]string {
	return t.Metadata
}

func (t *Domain) SetMeta(key string, value string) {
	t.Metadata[key] = value
}

func (t *Domain) GetMeta(key string) (value string) {
	if value, ok := t.Metadata[key]; ok {
		return value
	}
	return
}

func (t *Domain) SetData(key string, value json.RawMessage) {
	t.Data[key] = value
}

func (t *Domain) GetData(key string) (value json.RawMessage) {
	if value, ok := t.Data[key]; ok {
		return value
	}
	return
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

func (d *Domain) String() string {
	// names := []string{d.PreName, d.Domain, d.SufName}
	// name = strings.Join(names, ".")
	// name = strings.ReplaceAll(name, "..", ".")
	// name = strings.Trim(name, ".")
	return d.input
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
		log.Error("Error:", err)
	}

	return string(jsonData)
}

func (d *Domain) Idn(names ...string) string {
	if len(names) > 0 {
		d.Punycode = names[0]
	}

	return d.Punycode
}
