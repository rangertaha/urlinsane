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
package hr

import (
	"reflect"
	"sort"
	"strings"
	"testing"
	"unicode"

		"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
)

func TestAlgo(t *testing.T) {
	tests := []struct {
		original string
		variants []string
	}{

		{
			original: "g.com",
			variants: []string{
				"q.com",
				"ɢ.com",
				"ɡ.com",
				"ԍ.com ",
				"ġ.com",
				"ğ.com",
				"ց.com",
				"ǵ.com ",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.original, func(t *testing.T) {
			algo := Algo{
				languages: languages.Languages("en"),
			}
			variants := algo.Func(test.original)
			sort.Strings(variants)

			if !reflect.DeepEqual(variants, test.variants) {
				t.Errorf("algo.Func(%s) = %s; want %s", test.original, variants, test.variants)
			}
		})
	}

	t.Run("hr id", func(t *testing.T) {
		algo := Algo{}
		if algo.Id() != CODE {
			t.Errorf("algo.Id() = '%s'; want '%s'", algo.Id(), CODE)
		}
	})

	t.Run("hr empty id", func(t *testing.T) {
		algo := Algo{}
		if algo.Id() == "" {
			t.Errorf("algo.Id() can not return an empty string")
		}
	})
	t.Run("hr lowercase id", func(t *testing.T) {
		algo := Algo{}
		for _, c := range algo.Id() {
			if unicode.IsUpper(c) {
				t.Errorf("algo.Id() must be lowercase")
			}
		}
	})
	t.Run("hr name", func(t *testing.T) {
		algo := Algo{}
		if algo.Name() != NAME {
			t.Errorf("algo.Name() = '%s'; want '%s'", algo.Name(), NAME)
		}
	})

	t.Run("hr empty name", func(t *testing.T) {
		algo := Algo{}
		if algo.Name() == "" {
			t.Errorf("algo.Name() can not return an empty string")
		}
	})

	t.Run("hr titlecase name", func(t *testing.T) {
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

	t.Run("hr description", func(t *testing.T) {
		algo := Algo{}
		if algo.Description() != DESCRIPTION {
			t.Errorf("algo.Description() = '%s'; want '%s'", algo.Name(), DESCRIPTION)
		}
	})

	t.Run("hr empty description", func(t *testing.T) {
		algo := Algo{}
		if algo.Description() == "" {
			t.Errorf("algo.Description() can not return an empty string")
		}
	})

	t.Run("hr capitalizing fist char of description", func(t *testing.T) {
		algo := Algo{}
		firstRune := []rune(algo.Description())[0]
		if !unicode.IsUpper(firstRune) {
			t.Errorf("algo.Name() must capitalize start of sentence: %s", algo.Description())
		}
	})
}
