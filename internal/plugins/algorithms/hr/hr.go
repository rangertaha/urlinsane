// Copyright 2024 Rangertaha. All Rights Reserved.
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
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
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
}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Init(conf internal.Config) {
	n.keyboards = conf.Keyboards()
	n.languages = conf.Languages()
	n.config = conf
}

func (n *Algo) Name() string {
	return NAME
}
func (n *Algo) Description() string {
	return DESCRIPTION
}

func (n *Algo) Exec(original *db.Domain) (domains []*db.Domain, err error) {
	for _, language := range n.languages {
		for _, variant := range algo.HomoglyphSwapping(original.Name, language.Homoglyphs()) {
			if original.Name != variant {
				domains = append(domains, &db.Domain{Name: variant})
				// acc.Add(domain.Variant(n, original.Prefix(), variant, original.Suffix()))
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
