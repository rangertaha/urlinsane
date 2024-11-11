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
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	algo "github.com/rangertaha/urlinsane/pkg/typo"
)

const (
	CODE        = "cns"
	NAME        = "Cardinal Substitution"
	DESCRIPTION = "Swapping digial numbers and carninal numbers"
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

func (n *Algo) Exec(original internal.Domain, acc internal.Accumulator) (err error) {
	for _, lang := range n.languages {
		for _, variant := range algo.CardinalSwap(original.Name(), lang.Numerals()) {
			if original.Name() != variant {
				acc.Add(domain.NewVariant(n, original.Prefix(), variant, original.Suffix()))
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
