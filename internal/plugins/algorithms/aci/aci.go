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
package aci

// Adjacent character insertion is where an attacker adds characters
// that are next to each other on a keyboard.

// For example, if a user intends to visit "example.com," a typo-squatter
// might register "examplw.com" or "exanple.com." These small alterations
// can trick users into clicking on the malicious sites, leading to phishing
// scams, malware downloads, or other harmful activities.

// Adjacent character insertion exploits common typing errors, making it a
// particularly effective tactic, as users may not notice the difference,
// especially if they are typing quickly. It highlights the importance of
// vigilance and cybersecurity measures to protect against such deceptive
// practices.

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "aci"
	NAME        = "Adjacent Character Insertion"
	DESCRIPTION = "Inserting adjacent character from the keyboard"
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

// func (n *Algo) Exec(in internal.Typo) (out []internal.Typo) {
// 	original := in.Original().Repr()

// 	for i, char := range original {
// 		for _, board := range in.Keyboards() {
// 			for _, key := range board.Adjacent(string(char)) {
// 				variant := original[:i] + string(key) + string(char) + original[i+1:]
// 				vt := in.New(variant)
// 				out = append(out, vt)
// 			}
// 		}
// 	}
// 	return
// }

func (n *Algo) Func(original string) (results []string) {
	for i, char := range original {
		for _, board := range n.keyboards {
			for _, kchar := range board.Adjacent(string(char)) {
				variant := fmt.Sprint(original[:i], kchar, original[i+1:])
				results = append(results, variant)
			}
		}
	}
	return results
}

// func AlgoFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for i, char := range tc.Original.Domain {
// 			for _, key := range keyboard.Adjacent(string(char)) {
// 				d1 := tc.Original.Domain[:i] + string(key) + string(char) + tc.Original.Domain[i+1:]
// 				dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})

// 				d2 := tc.Original.Domain[:i] + string(char) + string(key) + tc.Original.Domain[i+1:]
// 				dm2 := Domain{tc.Original.Subdomain, d2, tc.Original.Suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm2, Typo: tc.Typo, Data: tc.Data})
// 			}
// 		}
// 	}
// 	return
// }

// Register the plugin
func init() {
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Algo{}
	})
}
