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
package typo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/models"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
)

type Typo struct {
	Algorithm internal.Algorithm `json:"-"`
	// Algorithm  string        `json:"algorithm"`
	Original   models.Domain `json:"original"`
	Variant    models.Domain `json:"variant"`
	Distance   int           `json:"distance"`
	Similarity int           `json:"similarity"`

	meta map[string]interface{}
}

func New() Typo {
	return Typo{
		meta: make(map[string]interface{}),
	}
}

func (t *Typo) Metatable() map[string]interface{} {
	return t.meta
}

func (t *Typo) SetMeta(key string, value interface{}) {
	t.meta[key] = value
}

func (t *Typo) GetMeta(key string) (value interface{}) {
	if value, ok := t.meta[key]; ok {
		return value
	}
	return nil
}

func (t *Typo) Algo() internal.Algorithm {
	return t.Algorithm
}

func (t *Typo) Set(origin, variant models.Domain) {
	t.Original = origin
	t.Variant = variant
}

func (t *Typo) Get() (origin, variant models.Domain) {
	fmt.Println(t.Original.Name, len(t.Original.Name), "-", t.Variant.Name, len(t.Variant.Name), t.Valid())
	return t.Original, t.Variant
}

// "Origin" refers to the starting point or source of something, while "derive"

func (t *Typo) Derived(labels ...string) models.Domain {
	if len(labels) > 0 {
		name := strings.Join(labels, ".")
		t.Variant = models.Domain{
			Prefix: domainutil.Subdomain(name),
			Name:   domainutil.DomainPrefix(name),
			Suffix: domainutil.DomainSuffix(name),
		}
		// domainutil.SplitDomain()
	}

	return t.Variant
}

func (t *Typo) Origin(labels ...string) models.Domain {
	if len(labels) > 0 {
		name := strings.Join(labels, ".")
		t.Original = models.Domain{
			Prefix: domainutil.Subdomain(name),
			Name:   domainutil.DomainPrefix(name),
			Suffix: domainutil.DomainSuffix(name),
		}
	}

	return t.Original
}

func (t *Typo) New(algo internal.Algorithm, origin, variant models.Domain) internal.Typo {
	dist := fuzzy.Levenshtein(origin.Name, variant.Name)

	return &Typo{
		Algorithm: algo,
		Original:  origin,
		Variant:   variant,
		Distance:  dist,
		meta:      make(map[string]interface{}),
	}
}

func (t *Typo) String() string {
	return t.Variant.Fqdn()
}

func (t *Typo) Threat() int {
	return 0
}

func (t *Typo) Live() bool {
	return t.Variant.Live
}

func (t *Typo) Valid() bool {
	return t.Variant.Name != ""
}

func (t *Typo) Dist() int {
	return t.Distance
}

func (t *Typo) Json() string {
	// Marshal the struct into JSON
	jsonData, err := json.Marshal(t)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonData)
}

// type Typo struct {
// 	algorithm internal.Algorithm
// 	original  internal.Target
// 	variant   internal.Target
// 	dist      int
// }

// func (t *Typo) Algorithm() internal.Algorithm {
// 	return t.algorithm
// }

// func (t *Typo) Original() internal.Target {
// 	return t.original
// }

// func (t *Typo) Variant() internal.Target {
// 	return t.variant
// }

// func (t *Typo) Active() bool {
// 	return t.variant.Live()
// }

// func (t *Typo) Ld() int {
// 	return t.dist
// }

// func (t *Typo) Clone(name string) internal.Typo {
// 	dist := fuzzy.Levenshtein(t.original.Name(), name)

// 	return &Typo{
// 		algorithm: t.algorithm,
// 		original:  t.original,
// 		variant:   target.New(name),
// 		dist:      dist,
// 	}
// }

// func (t *Typo) String() (val string) {
// 	if t.variant != nil {
// 		return t.variant.Name()
// 	}

// 	return ""
// }
