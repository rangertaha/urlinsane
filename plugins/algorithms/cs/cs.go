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
package cs

// // AlgoFunc typos are when two consecutive characters are swapped in the original domain name.
// // Example: www.examlpe.com
// func AlgoFunc(tc Result) (results []Result) {
// 	for i := range tc.Original.Domain {
// 		if i <= len(tc.Original.Domain)-2 {
// 			domain := fmt.Sprint(
// 				tc.Original.Domain[:i],
// 				string(tc.Original.Domain[i+1]),
// 				string(tc.Original.Domain[i]),
// 				tc.Original.Domain[i+2:],
// 			)
// 			if tc.Original.Domain != domain {
// 				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 			}
// 		}
// 	}
// 	return results
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const (
	CODE        = "cs"
	NAME        = "Character Swap"
	DESCRIPTION = "Character Swap Swapping two consecutive characters in a domain"
)

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
	return NAME
}
func (n *Algo) Description() string {
	return DESCRIPTION
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
