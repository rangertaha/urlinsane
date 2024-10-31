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
package ons

// Ordinal Numeral Swap
// Ordinal numerals are the numbers that are used for counting something.
// For Example: first, second, third, fourth, fifth, sixth, seventh, eighth.
// Ordinal swapping replaces ordinal numerals with digit numbers and numbers for
// ordinal numerals. For example:
//
// Input: firstandsecondunited.com
//
// Output:
// ID     TYPE          TYPO
// -------------------------------------------
//  1      Ordinal Swap  1and2united.com
//  2      Ordinal Swap  1andsecondunited.com
//  3      Ordinal Swap  firstand2united.com
// -------------------------------------------
//  TOTAL  3
//
//
//
// Input: 1united23.com
//
// Output:

// ID     TYPE          TYPO
// -------------------------------------------------
//  1      Ordinal Swap  1unitedsecondthird.com
//  2      Ordinal Swap  1united2third.com
//  3      Ordinal Swap  firstunited2third.com
//  4      Ordinal Swap  firstunited23.com
//  5      Ordinal Swap  1unitedsecond3.com
//  6      Ordinal Swap  firstunitedsecond3.com
//  7      Ordinal Swap  firstunitedsecondthird.com
// -------------------------------------------------
//  TOTAL  7
//
// We can verify the number of permutations with some calculations.
// Assuming language plugins only have numbers and numerals upto 9, we can
// calculate the total number of variants using this formula:
// Total variants = 2^(number of numerals) - 1
//

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "ons"
	NAME        = "Ordinal Swap"
	DESCRIPTION = "Swapping digital numbers and ordinal numbers"
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

func (n *Algo) Exec(typo internal.Typo) (typos []internal.Typo) {
	for _, lang := range typo.Languages() {
		for _, variant := range n.Func(lang.Ordinal(), typo.Original().Repr()) {
			typos = append(typos, typo.New(variant))
		}
	}
	return
}

// Func swaps numbers and carninal numbers
func (n *Algo) Func(cardinals map[string]string, name string) []string {
	results := []string{}
	var fn func(map[string]string, string, bool) map[string]bool

	fn = func(data map[string]string, str string, reverse bool) (names map[string]bool) {
		names = make(map[string]bool)

		for num, word := range data {
			{
				var variant string
				if !reverse {
					variant = strings.Replace(str, word, num, -1)
				} else {
					variant = strings.Replace(str, num, word, -1)
				}

				if str != variant {
					if _, ok := names[variant]; !ok {
						names[variant] = true
						for k, v := range fn(cardinals, variant, reverse) {
							names[k] = v
						}

						fn(cardinals, variant, reverse)
					}
				}
			}
		}
		return names
	}

	for name := range fn(cardinals, name, false) {
		results = append(results, name)
	}
	for name := range fn(cardinals, name, true) {
		results = append(results, name)
	}

	return results
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Algo{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
