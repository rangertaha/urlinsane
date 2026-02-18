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
package cns

// Cardinal Numeral Swap
// Cardinal numerals are the numbers that are used for counting something.
// For Example: one, two, three, four, five, six, seven, eight, nine, ten.
// Cardinal swapping replaces cardinal numerals with numbers and numbers for
// cardinal numerals. For example:
//
// Input: 123.com
//
// Output:
//  ID     TYPE           TYPO
// ---------------------------------------
//  7      Cardinal Swap  one2three.com
//  1      Cardinal Swap  one23.com
//  2      Cardinal Swap  1two3.com
//  3      Cardinal Swap  1twothree.com
//  4      Cardinal Swap  onetwothree.com
//  5      Cardinal Swap  onetwo3.com
//  6      Cardinal Swap  12three.com
// ---------------------------------------
//  TOTAL  7
//
//
//
// Input: onetwothree.com
//
// Output:
// ID     TYPE           TYPO
// -------------------------------------
//  1      Cardinal Swap  one2three.com
//  2      Cardinal Swap  1twothree.com
//  3      Cardinal Swap  12three.com
//  4      Cardinal Swap  123.com
//  5      Cardinal Swap  1two3.com
//  6      Cardinal Swap  onetwo3.com
//  7      Cardinal Swap  one23.com
// -------------------------------------
//  TOTAL  7
//
// We can verify the number of permutations with some calculations.
// Assuming language plugins only have numbers and numerals upto 9, we can
// calculate the total number of variants using this formula:
// Total variants = 2^(number of numerals) - 1
//

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

type Plugin struct {
	algorithms.Plugin
}

func (p *Plugin) Exec(original *db.Domain) (domains []*db.Domain, err error) {
	languages := p.Conf.Languages()
	// prefix, name, suffix := dns.Split(original.Name)
	// variant = dns.Join(prefix, variant, suffix)

	for _, lang := range languages {
		for _, variant := range typo.CardinalSwap(original.Name, lang.Numerals()) {
			if original.Name != variant {
				dist := fuzzy.Levenshtein(original.Name, variant)
				domains = append(domains, &db.Domain{Name: variant, Levenshtein: dist, Algorithm: p.Algo()})
			}
		}
	}

	return
}

// Register the plugin
func init() {
	var CODE = "cns"
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Plugin{
			Plugin: algorithms.Plugin{
				Code:    CODE,
				Title:   "Cardinal Substitution",
				Summary: "Swapping digial numbers and carninal numbers",
			},
		}
	})
}
