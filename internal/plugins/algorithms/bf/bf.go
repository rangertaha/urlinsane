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
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)


const (
	CODE        = "bf"
	NAME        = "Bit Flipping"
	DESCRIPTION = "Relies on random bit-errors to redirect connections"
)

type Algo struct {

}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Name() string {
	return NAME
}

func (n *Algo) Description() string {
	return DESCRIPTION
}

func (n *Algo) Exec(in internal.Typo) (out []internal.Typo) {
	out = append(out, in)
	return
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
		return &Algo{
			
		}
	})
}
