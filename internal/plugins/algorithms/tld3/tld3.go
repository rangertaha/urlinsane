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
package tld3

// // wrongThirdLevelDomainFunc uses an alternate, valid third level domain.
// func wrongThirdLevelDomainFunc(tc Result) (results []Result) {
// 	labels := strings.Split(tc.Original.Suffix, ".")
// 	length := len(labels)
// 	for _, suffix := range datasets.TLD {
// 		suffixLbl := strings.Split(suffix, ".")
// 		suffixLen := len(suffixLbl)
// 		if length == suffixLen && length == 3 {
// 			if suffixLbl[1] == labels[1] && suffixLbl[2] == labels[2] {
// 				if suffix != tc.Original.Suffix {
// 					dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "tld3"
	NAME        = "Wrong TLD3"
	DESCRIPTION = "Wrong third level domain (TLD3)"
)

type Algo struct {
	ctype     int
	config    internal.Config
	languages []internal.Language
	keyboards []internal.Keyboard
}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Init(conf internal.Config) {
	n.keyboards = conf.Keyboards()
	n.languages = conf.Languages()
	n.ctype = conf.Type()
	n.config = conf
}

func (n *Algo) Name() string {
	return NAME
}
func (n *Algo) Description() string {
	return DESCRIPTION
}

func (n *Algo) Exec(typo internal.Typo) (typos []internal.Typo) {
	if n.config.Type() == internal.DOMAIN {
		return n.domain(typo)
	}

	if n.config.Type() == internal.PACKAGE {
		return n.code(typo)
	}

	if n.config.Type() == internal.NAME {
		return n.name(typo)
	}
	return
}

func (n *Algo) domain(typo internal.Typo) (typos []internal.Typo) {
	sub, prefix, suffix := typo.Original().Domain()
	// fmt.Println(sub, prefix, suffix)

	for _, variant := range n.Func(prefix) {
		if prefix != variant {
			d := domain.New(sub, variant, suffix)
			// fmt.Println(sub, variant, suffix)

			new := typo.Clone(d.String())

			typos = append(typos, new)
		}
	}
	return
}

func (n *Algo) code(typo internal.Typo) (typos []internal.Typo) {
	original := n.config.Target().Name()
	for _, variant := range n.Func(original) {
		if original != variant {
			typos = append(typos, typo.Clone(variant))
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
