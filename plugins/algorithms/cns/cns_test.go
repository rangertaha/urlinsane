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

// Numeral Numeral Swap
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

import (
	"reflect"
	"sort"
	"testing"
)

var cardinals = map[string]string{
	"0":  "zero",
	"1":  "one",
	"2":  "two",
	"3":  "three",
	"4":  "four",
	"5":  "five",
	"6":  "six",
	"7":  "seven",
	"8":  "eight",
	"9":  "nine",
	"10": "ten",
}


func TestAlgo(t *testing.T) {
	tests := []struct {
		original string
		variants []string
	}{
		{
			original: "123.com",
			variants: []string{
				"12three.com",
				"1two3.com",
				"1twothree.com",
				"one23.com",
				"one2three.com",
				"onetwo3.com",
				"onetwothree.com",
			},
		},
		{
			original: "onetwothree.com",
			variants: []string{
				"123.com",
				"12three.com",
				"1two3.com",
				"1twothree.com",
				"one23.com",
				"one2three.com",
				"onetwo3.com",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.original, func(t *testing.T) {
			algo := Algo{}
			variants := algo.Func(cardinals, test.original)
			sort.Strings(variants)

			if !reflect.DeepEqual(variants, test.variants) {
				t.Errorf("algo.Func(<numerals>, %s) = %s; want %s", test.original, variants, test.variants)
			}
		})
	}

	t.Run("md id", func(t *testing.T) {
		algo := Algo{}
		if algo.Id() != CODE {
			t.Errorf("algo.Id() = '%s'; want '%s'", algo.Id(), CODE)
		}
	})

	t.Run("md empty id", func(t *testing.T) {
		algo := Algo{}
		if algo.Id() == "" {
			t.Errorf("algo.Id() can not return an empty string")
		}
	})

	t.Run("md name", func(t *testing.T) {
		algo := Algo{}
		if algo.Name() != NAME {
			t.Errorf("algo.Name() = '%s'; want '%s'", algo.Name(), NAME)
		}
	})

	t.Run("md empty name", func(t *testing.T) {
		algo := Algo{}
		if algo.Name() == "" {
			t.Errorf("algo.Name() can not return an empty string")
		}
	})

	t.Run("md description", func(t *testing.T) {
		algo := Algo{}
		if algo.Description() != DESCRIPTION {
			t.Errorf("algo.Description() = '%s'; want '%s'", algo.Name(), DESCRIPTION)
		}
	})

	t.Run("md empty description", func(t *testing.T) {
		algo := Algo{}
		if algo.Description() == "" {
			t.Errorf("algo.Description() can not return an empty string")
		}
	})

}
