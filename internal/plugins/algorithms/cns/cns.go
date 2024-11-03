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
	"fmt"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "cns"
	NAME        = "Cardinal Swap"
	DESCRIPTION = "Swapping digial numbers and carninal numbers"
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
	for _, lang := range n.languages {

		for _, variant := range n.Func(lang.Cardinal(), prefix) {
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
	for _, lang := range n.languages {

		for _, variant := range n.Func(lang.Cardinal(), username) {
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
	for _, lang := range n.languages {
		for _, variant := range n.Func(lang.Cardinal(), original) {
			typos = append(typos, typo.Clone(variant))
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
		return &Algo{}
	})
}
