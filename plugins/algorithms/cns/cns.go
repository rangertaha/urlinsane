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
package cns

// Cardinal Numeral Swap
// A numeral is a word describing a number and a number is expressed with
// digits. Numeral swapping replaces numerals with numbers and numbers for
// numerals. For example:
//
// Original: onetwothree.com
//
// Variants: 1twothree.com
//           one2three.com
//           onetwo3.com
//           one23.com
//           12three.com
//           1two3.com
//           123.com
//
// Assuming language plugins only have numbers and numerals upto 9, we can
// calculate the total number of variants using this formula:
// Total variants = 2^(number of numerals) - 1
//
// numerals. For example:
//
// Original: onetwothree.com
//
// Variants: 1twothree.com
//           one2three.com
//           onetwo3.com
//           one23.com
//           12three.com
//           1two3.com
//           123.com
//
// Assuming language plugins only have numbers and numerals upto 9, we can
// calculate the total number of variants using this formula:
// Total variants = 2^(number of numerals) - 1

import (
	"strings"

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const (
	CODE        = "cns"
	NAME        = "Cardinal Swap"
	DESCRIPTION = "Swapping numbers and carninal numbers"
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

func (n *Algo) Exec(typo urlinsane.Typo) (typos []urlinsane.Typo) {
	for _, lang := range typo.Languages() {
		for _, variant := range n.Func(lang.Cardinal(), typo.Original().Repr()) {
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
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &Algo{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
