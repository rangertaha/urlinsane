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
package di

// // adjacentCharacterInsertionFunc are created by inserting letters adjacent of each letter. For example, www.googhle.com
// // and www.goopgle.com
// func AlgoFunc(tc Result) (results []Result) {

// 	for i, char := range tc.Original.Domain {

// 		d1 := tc.Original.Domain[:i] + "." + string(char) + tc.Original.Domain[i+1:]
// 		dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 		results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
// 	}

// 	return
// }

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	algo "github.com/rangertaha/urlinsane/pkg/typo"
)

const (
	CODE        = "di"
	NAME        = "Dot Insertion"
	DESCRIPTION = "Inserting periods in the target domain name"
)

type Algo struct {
	config internal.Config
	// languages []internal.Language
	// keyboards []internal.Keyboard
	funcs map[int]func(internal.Typo) []internal.Typo
}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Init(conf internal.Config) {
	n.funcs = make(map[int]func(internal.Typo) []internal.Typo)
	// n.keyboards = conf.Keyboards()
	// n.languages = conf.Languages()
	n.config = conf

	// Supported targets
	n.funcs[internal.DOMAIN] = n.domain
	n.funcs[internal.PACKAGE] = n.name
	n.funcs[internal.EMAIL] = n.email
	n.funcs[internal.NAME] = n.name
}

func (n *Algo) Name() string {
	return NAME
}
func (n *Algo) Description() string {
	return DESCRIPTION
}

func (n *Algo) Exec(typo internal.Typo) []internal.Typo {
	return n.funcs[n.config.Type()](typo)
}
func (n *Algo) Exec(typo internal.Typo) (typos []internal.Typo) {
	orig, vari := typo.Get()

	for _, variant := range algo.BitFlipping(vari.Name) {
		if vari.Name != variant {

			new := typo.New(n, orig, domain.Parse(variant))
			typos = append(typos, new)
		}
	}

	return
}
func (n *Algo) domain(typo internal.Typo) (typos []internal.Typo) {
	sub, prefix, suffix := typo.Original().Domain()
	for _, variant := range algo.DotInsertion(prefix) {
		if prefix != variant {
			d := domain.New(sub, variant, suffix)
			new := typo.Clone(d.String())
			typos = append(typos, new)
		}
	}
	return
}

func (n *Algo) email(typo internal.Typo) (typos []internal.Typo) {
	username, domain := typo.Original().Email()
	for _, variant := range algo.DotInsertion(username) {
		if username != variant {
			new := typo.Clone(fmt.Sprintf("%s@%s", variant, domain))
			typos = append(typos, new)
		}
	}
	return
}

func (n *Algo) name(typo internal.Typo) (typos []internal.Typo) {
	original := n.config.Target().Name()
	for _, variant := range algo.DotInsertion(original) {
		if original != variant {
			typos = append(typos, typo.Clone(variant))
		}
	}
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Algo{}
	})
}
