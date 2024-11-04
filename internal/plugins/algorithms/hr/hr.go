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
package hr

// Homoglyph replacement in typosquatting involves substituting letters in a
// legitimate domain name, username, brand, or package name with visually similar
// characters (homoglyphs) from different alphabets or symbols. The aim is to
// create misleading, but seemingly identical or near-identical names that can
// trick users into thinking they're accessing a legitimate resource, when in
// fact, they’re visiting a malicious or spoofed one. This tactic is widely
// used in phishing attacks, brand impersonation, and other social engineering
// scams.

// homoglyphFunc when one or more characters that look similar to another
// character but are different are called homogylphs. An example is that the
// lower case l looks similar to the numeral one, e.g. l vs 1. For example,
// google.com becomes goog1e.com.

// INPUT:  g.com
//
// TYPE    TYPO
// ---------------
//  HR      ģ.com
//  HR      q.com
//  HR      ɢ.com
//  HR      ɡ.com
//  HR      ԍ.com
//  HR      ġ.com
//  HR      ğ.com
//  HR      ց.com
//  HR      ǵ.com
// ---------------
//  TOTAL   9

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	algo "github.com/rangertaha/urlinsane/pkg/typo"
)

const (
	CODE        = "hr"
	NAME        = "Homoglyphs Replacement"
	DESCRIPTION = "Replaces characters with characters that look similar"
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

	for _, language := range n.languages {
		for _, variant := range algo.HomoglyphSwapping(prefix, language.Homoglyphs()) {

			if prefix != variant {
				d := domain.New(sub, variant, suffix)

				new := typo.Clone(d.String())

				typos = append(typos, new)
			}
		}
	}
	return
}

func (n *Algo) email(typo internal.Typo) (typos []internal.Typo) {
	username, domain := typo.Original().Email()

	for _, language := range n.languages {
		for _, variant := range algo.HomoglyphSwapping(username, language.Homoglyphs()) {

			if username != variant {
				new := typo.Clone(fmt.Sprintf("%s@%s", variant, domain))

				typos = append(typos, new)
			}
		}
	}
	return
}

func (n *Algo) name(typo internal.Typo) (typos []internal.Typo) {
	original := n.config.Target().Name()
	for _, language := range n.languages {
		for _, variant := range algo.HomoglyphSwapping(original, language.Homoglyphs()) {

			if original != variant {
				typos = append(typos, typo.Clone(variant))
			}
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
