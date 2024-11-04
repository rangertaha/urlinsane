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
package dhus

import (
	"reflect"
	"sort"
	"testing"
)

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
			variants := algo.Func(test.original)
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
