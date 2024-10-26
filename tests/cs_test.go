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

package tests

import (
	"testing"

	"github.com/rangertaha/urlinsane/pkg/typo"
)

var characterSwapCases = []testpair{
	{[]string{"google.com"},
		map[string]bool{
			"ogogle.com": true,
			"googel.com": true,
			"goolge.com": true,
			"gogole.com": true,
		}, 5},
}

func TestCharacterSwap(t *testing.T) {
	for _, lang := range languages {
		count := 0
		for _, tcase := range characterSwapCases {
			conf := typo.BasicConfig{
				Domains:     tcase.domains,
				Keyboards:   []string{lang},
				Typos:       []string{"cs"},
				Funcs:       []string{""},
				Concurrency: 50,
				Format:      "text",
				Verbose:     false,
			}

			urli := typo.New(conf.Config())

			out := urli.Stream()

			for r := range out {
				_, ok := tcase.values[r.Variant.String()]
				if !ok {
					t.Errorf("Failed variant: %v for domains: %v, language: %v, algorithm %v", r.Variant.String(), tcase.domains, lang, r.Typo.Name)
				}
				count++
			}
			// TODO: Apply dup filter and uncomment
			// if count != tcase.total {
			// 	t.Errorf("Failed total number of records should be %v not %v", tcase.total, count)
			// }
			count = 0
		}
	}
}
