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


import (
	"reflect"
	"strings"
	"testing"
	"unicode"
)

func TestAlgo(t *testing.T) {
	tests := []struct {
		original string
		variants []string
	}{
		{
			original: "facebook.com.io.uk",
			variants: []string{
				"facebookcom.io.uk",
				"facebook.comio.uk",
				"facebook.com.iouk",
			},
		},
		{
			original: "google.com.uk",
			variants: []string{
				"googlecom.uk",
				"google.comuk",
			},
		},
		{
			original: "www.google.com",
			variants: []string{
				"wwwgoogle.com",
				"www.googlecom",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.original, func(t *testing.T) {
			algo := Algo{}
			variants := algo.Func(test.original, ".")

			if !reflect.DeepEqual(variants, test.variants) {
				t.Errorf("algo.Func(%s, '.') = %s; want %s", test.original, variants, test.variants)
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
	t.Run("md lowercase id", func(t *testing.T) {
		algo := Algo{}
		for _, c := range algo.Id() {
			if unicode.IsUpper(c) {
				t.Errorf("algo.Id() must be lowercase")
			}
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

	t.Run("md titlecase name", func(t *testing.T) {
		algo := Algo{}
		strs := ""
		for _, name := range strings.Split(algo.Name(), " ") {
			strs = strs + string(name[0])
		}
		for _, c := range strs {
			if !unicode.IsUpper(c) {
				t.Errorf("algo.Name() must be in titlecase: %s, Found: (%s)", algo.Name(), string(c))
			}
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

	t.Run("md capitalizing fist char of description", func(t *testing.T) {
		algo := Algo{}
		firstRune := []rune(algo.Description())[0]
		if !unicode.IsUpper(firstRune) {
			t.Errorf("algo.Name() must capitalize start of sentence: %s", algo.Description())
		}
	})
}
