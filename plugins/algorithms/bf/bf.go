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

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "bf"
// const (
// 	CODE        = ""
// 	NAME        = ""
// 	DESCRIPTION = ""
// )


type Algo struct {
	types []string
}

func (n *Algo) Id() string {
	return CODE
}
func (n *Algo) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *Algo) Name() string {
	return "Bit Flipping"
}

func (n *Algo) Description() string {
	return "Relies on random bit-errors to redirect connections"
}

func (n *Algo) Fields() []string {
	return []string{}
}

func (n *Algo) Headers() []string {
	return []string{}
}

func (n *Algo) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &Algo{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
