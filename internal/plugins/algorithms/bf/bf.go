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
package bf

// Typo-squatting through bit flipping involves manipulating domain names by changing one or more bits in the binary representation of characters. This can lead to visually similar but distinct characters being used in a domain name, which can trick users into visiting malicious sites.

// For instance, an attacker might register a domain like "exarnple.com," where the letter "m" has been flipped to a visually similar character, such as "rn." Users may not notice the subtle difference when typing quickly, leading them to inadvertently access the spoofed site.

// This technique exploits the similarities in how certain characters appear in certain fonts or how they are rendered on screens, making it a clever and deceptive method of capturing traffic meant for legitimate websites. It highlights the importance of vigilance and proactive measures in cybersecurity to protect against such tactics.

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "bf"
	NAME        = "Bit Flipping"
	DESCRIPTION = "Relies on random bit-errors to redirect connections"
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

func (n *Algo) email(typo internal.Typo) (typos []internal.Typo) {
	username, domain := typo.Original().Email()
	// fmt.Println(sub, prefix, suffix)

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

// // bitsquattingFunc relies on random bit- errors to redirect connections
// // intended for popular domains
// func bitsquattingFunc(tc Result) (results []Result) {
// 	// TOOO: need to improve.
// 	masks := []int{1, 2, 4, 8, 16, 32, 64, 128}
// 	charset := make(map[string][]string)
// 	for _, board := range tc.Keyboards {
// 		for _, alpha := range board.Language.Graphemes {
// 			for _, mask := range masks {
// 				new := int([]rune(alpha)[0]) ^ mask
// 				for _, a := range board.Language.Graphemes {
// 					if string(a) == string(new) {
// 						charset[string(alpha)] = append(charset[string(alpha)], string(new))
// 					}
// 				}
// 			}
// 		}
// 	}

// 	for d, dchar := range tc.Original.Domain {
// 		for _, char := range charset[string(dchar)] {

// 			dnew := tc.Original.Domain[:d] + string(char) + tc.Original.Domain[d+1:]
// 			dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
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
