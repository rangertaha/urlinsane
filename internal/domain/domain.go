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
	"os"
	"path/filepath"
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
	Metadata map[string]string          `json:"-"` // Used for simplified table based output
	// Input       string                     `json:"input,omitempty"` // Name provded by user
	Levenshtein int  `json:"ld,omitempty"`
	Involved    bool `json:"active,omitempty"`
	Cache       bool `json:"cached,omitempty"`
	// Directory   string `json:"dir,omitempty"`
	// Source      string `json:"src,omitempty"`  // Domain name provded by user
	// Variant     string `json:"typo,omitempty"` // Domain varient
	algo internal.Algorithm
}

func New(name string) internal.Domain {
	return &Domain{
		FQDN:     name,
		PreName:  domainutil.Subdomain(name),
		Domain:   domainutil.DomainPrefix(name),
		SufName:  domainutil.DomainSuffix(name),
		Metadata: make(map[string]string),
		Data:     make(map[string]json.RawMessage),
		// Source:   name,
	}
}

func Variant(algo internal.Algorithm, names ...string) internal.Domain {
	name := strings.Join(names, ".")
	name = strings.Trim(name, ".")

	return &Domain{
		FQDN:     name,
		PreName:  domainutil.Subdomain(name),
		Domain:   domainutil.DomainPrefix(name),
		SufName:  domainutil.DomainSuffix(name),
		Metadata: make(map[string]string),
		Data:     make(map[string]json.RawMessage),
		algo:     algo,
		// Source:   source,
	}
}

func NewVariant(algo internal.Algorithm, names ...string) internal.Domain {
	name := strings.Join(names, ".")
	name = strings.Trim(name, ".")

	return &Domain{
		FQDN:     name,
		PreName:  domainutil.Subdomain(name),
		Domain:   domainutil.DomainPrefix(name),
		SufName:  domainutil.DomainSuffix(name),
		Metadata: make(map[string]string),
		Data:     make(map[string]json.RawMessage),
		algo:     algo,
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

func (d *Domain) String() (n string) {
	return d.FQDN
}

// func (d *Domain) Dir(v ...string) string {
// 	if len(v) > 0 {
// 		d.Directory = v[0]
// 	}

// 	return d.Directory
// }

func (d *Domain) Live(v ...bool) bool {
	if len(v) > 0 {
		d.IsLive = v[0]
	}

	return d.IsLive
}

func (d *Domain) Cached(v ...bool) bool {
	if len(v) > 0 {
		d.Cache = v[0]
	}

	return d.Cache
}

func (d *Domain) Active(v ...bool) bool {
	if len(v) > 0 {
		d.Involved = v[0]
	}

	return d.Involved
}

// Ld returns the Levenshtein_distance
//
//	https://en.wikipedia.org/wiki/Levenshtein_distance
func (d *Domain) Ld(v ...int) int {
	if len(v) > 0 {
		d.Levenshtein = v[0]
	}

	return d.Levenshtein
}

func (d *Domain) Json(value ...string) (j string) {
	if len(value) == 0 {
		// Marshal the struct into JSON
		jsonData, err := json.Marshal(d)
		if err != nil {
			log.Error("Marshal:", err)
		}
		return string(jsonData)
	}
	if len(value) > 0 {
		data := value[0]
		// Unmarshal the JSON back into struct
		if err := json.Unmarshal([]byte(data), &d); err != nil {
			log.Error("Unmarshal:", err)
		}
		d.Cache = true
	}

	return
}

func (d *Domain) Idn(names ...string) string {
	if len(names) > 0 {
		d.Punycode = names[0]
	}

	return d.Punycode
}

func (d *Domain) Save(path string) (err error) {
	dir := filepath.Join(path, d.String())
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0750); err != nil {
			log.Error(err)
			return
		}
	}

	file := filepath.Join(dir, "scan.json")
	log.Debugf("creating %s", file)
	if err = os.WriteFile(file, []byte(d.Json()), 0644); err != nil {
		log.Error(err)
		return
	}
	return err
}
