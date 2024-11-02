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
)

const (
	CODE        = "hr"
	NAME        = "Homoglyphs Replacement"
	DESCRIPTION = "Replaces characters with characters that look similar"
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

// AlgoFunc typos are when two consecutive characters are swapped in the original domain name.
// Example: www.examlpe.com
func (n *Algo) Func(original string) (results []string) {
	for i, char := range original {
		for _, language := range n.languages {
			for _, kchar := range language.SimilarChars(string(char)) {
				variant := fmt.Sprint(original[:i], kchar, original[i+1:])
				results = append(results, variant)
			}
		}
	}

	return
}

// func homoglyphFunc(tc Result) (results []Result) {
// 	for i, char := range tc.Original.Domain {
// 		// Check the alphabet of the language associated with the keyboard for
// 		// Algo
// 		for _, keyboard := range tc.Keyboards {
// 			for _, kchar := range keyboard.Language.SimilarChars(string(char)) {
// 				domain := fmt.Sprint(tc.Original.Domain[:i], kchar, tc.Original.Domain[i+1:])
// 				if tc.Original.Domain != domain {
// 					dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}
// 		// Check languages given with the (-l --language) CLI options for Algo.
// 		for _, language := range tc.Languages {
// 			for _, lchar := range language.SimilarChars(string(char)) {
// 				domain := fmt.Sprint(tc.Original.Domain[:i], lchar, tc.Original.Domain[i+1:])
// 				if tc.Original.Domain != domain {
// 					dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}

// 	}
// 	return results
// }

// Register the plugin
func init() {
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Algo{}
	})
}
