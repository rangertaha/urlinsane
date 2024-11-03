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
package vs

// AlgoFunc swaps vowels within the domain name except for the first letter.
// For example, www.google.com becomes www.gaagle.com.
// func AlgoFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for _, vchar := range keyboard.Language.Vowels {
// 			if strings.Contains(tc.Original.Domain, vchar) {
// 				for _, vvchar := range keyboard.Language.Vowels {
// 					new := strings.Replace(tc.Original.Domain, vchar, vvchar, -1)
// 					if new != tc.Original.Domain {
// 						dm := Domain{tc.Original.Subdomain, new, tc.Original.Suffix, Meta{}, false}
// 						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return
// }

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "vs"
	NAME        = "Vowel Swapping"
	DESCRIPTION = "Vowel Swapping"
)

type Algo struct {
	config    internal.Config
	languages []internal.Language
	keyboards []internal.Keyboard
	funcs     map[int]func(internal.Typo) []internal.Typo
}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Init(conf internal.Config) {
	n.funcs = make(map[int]func(internal.Typo) []internal.Typo)
	n.keyboards = conf.Keyboards()
	n.languages = conf.Languages()
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

func (n *Algo) domain(typo internal.Typo) (typos []internal.Typo) {
	sub, prefix, suffix := typo.Original().Domain()

	for _, variant := range n.Func(prefix) {
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

	for _, variant := range n.Func(username) {
		if username != variant {
			new := typo.Clone(fmt.Sprintf("%s@%s", variant, domain))

			typos = append(typos, new)
		}
	}
	return
}

func (n *Algo) name(typo internal.Typo) (typos []internal.Typo) {
	original := n.config.Target().Name()
	for _, variant := range n.Func(original) {
		if original != variant {
			typos = append(typos, typo.Clone(variant))
		}
	}
	return
}

func (n *Algo) Func(original string) (results []string) {
	// for i, char := range original {
	// 	for _, board := range n.keyboards {
	// 		for _, kchar := range board.Adjacent(string(char)) {
	// 			variant := fmt.Sprint(original[:i], kchar, original[i+1:])
	// 			results = append(results, variant)
	// 		}
	// 	}
	// }
	return results
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Algo{}
	})
}
