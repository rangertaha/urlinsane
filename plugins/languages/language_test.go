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
package languages

import (
	"reflect"
	"testing"
)

// | Number | Cardinal | Ordinal |
// |--------| -------- |---------|
// | 0      | zero     |         |
// | 1      | one      |first    |
// | 2      | two      |second   |
// | 3      | three    |third    |
// | 4      | four     |fourth   |
// | 5      | five     |fifth    |

var numerals map[string][]string = map[string][]string{
	"0": {"zero"},
	"1": {"one", "first"},
	"2": {"two", "second"},
	"3": {"three", "third"},
	"4": {"four", "fourth"},
	"5": {"five", "fifth"},
	"6": {"six", "sixth"},
	"7": {"seven", "seventh"},
	"8": {"eight", "eighth"},
	"9": {"nine", "ninth"},
}

type TestCase struct {
	numerals map[string][]string
	position int
	returns  map[string]string
}

func TestNumeralMap(t *testing.T) {
	tests := []TestCase{
		{
			numerals: numerals,
			position: 0,
			returns: map[string]string{
				"0":     "zero",
				"1":     "one",
				"2":     "two",
				"3":     "three",
				"4":     "four",
				"5":     "five",
				"6":     "six",
				"7":     "seven",
				"8":     "eight",
				"9":     "nine",
				"zero":  "0",
				"one":   "1",
				"two":   "2",
				"three": "3",
				"four":  "4",
				"five":  "5",
				"six":   "6",
				"seven": "7",
				"eight": "8",
				"nine":  "9",
			},
		},
	}

	for _, test := range tests {
		t.Run("numerals", func(t *testing.T) {
			returned := NumeralMap(test.numerals, 0)

			if !reflect.DeepEqual(returned, test.returns) {
				t.Errorf("languages.NumeralMap(<numerals>, 0) = %s; want %s", returned, test.returns)
			}
		})
	}

}


func TestOrdinalMap(t *testing.T) {
	tests := []TestCase{
		{
			numerals: numerals,
			position: 0,
			returns: map[string]string{
				"1":     "first",
				"2":     "second",
				"3":     "third",
				"4":     "fourth",
				"5":     "fifth",
				"6":     "sixth",
				"7":     "seventh",
				"8":     "eighth",
				"9":     "ninth",
				"first":   "1",
				"second":   "2",
				"third": "3",
				"fourth":  "4",
				"fifth":  "5",
				"sixth":   "6",
				"seventh": "7",
				"eighth": "8",
				"ninth":  "9",
			},
		},
	}

	for _, test := range tests {
		t.Run("numerals", func(t *testing.T) {
			returned := NumeralMap(test.numerals, 1)

			if !reflect.DeepEqual(returned, test.returns) {
				t.Errorf("languages.NumeralMap(<numerals>, 1) = %s; want %s", returned, test.returns)
			}
		})
	}

}
