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

// Ordinal Numeral Substitution
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
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	algo "github.com/rangertaha/urlinsane/pkg/typo"
)

const (
	CODE        = "ons"
	NAME        = "Ordinal Numeral Substitution"
	DESCRIPTION = "Substituting digital numbers and ordinal numbers"
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
	original := n.config.Target().Name()
	for _, language := range n.languages {
		for _, variant := range algo.OrdinalSwap(prefix, language.Numerals()) {

			if original != variant {
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
	for _, language := range n.languages {
		for _, variant := range algo.OrdinalSwap(username, language.Numerals()) {
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
	for _, language := range n.languages {
		for _, variant := range algo.OrdinalSwap(original, language.Numerals()) {
			if original != variant {
				typos = append(typos, typo.Clone(variant))
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
