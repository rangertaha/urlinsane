// Copyright 2024 Rangertaha. All Rights Reserved.
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
package co

// Character Omission
// Created by leaving out a character in the name.
//
//For example:
//
// Input: google.com
//
// Output:
// ID     TYPE    TYPO
// --------------------------
// 1      CO      gogle.com
// 2      CO      googlecom
// 5      CO      google.cm
// 6      CO      google.co
// 7      CO      oogle.com
// 8      CO      goole.com
// 9      CO      googe.com
// 3      CO      googl.com
// 4      CO      google.om
// --------------------------
// TOTAL  9
//
//
// Input: abcd
//
// Output:
// ID     TYPE    TYPO
// ---------------------
//  3      CO      abd
//  4      CO      abc
//  1      CO      bcd
//  2      CO      acd
// ---------------------
//  TOTAL  4

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

const (
	CODE        = "co"
	NAME        = "Character Omission"
	DESCRIPTION = "Omitting a character from the name"
)

type Plugin struct {
	algorithms.Plugin
}

func (p *Plugin) Exec(original *db.Domain) (domains []*db.Domain, err error) {
	algo := db.Algorithm{Code: p.Code, Name: p.Title}

	for _, variant := range typo.CharacterOmission(original.Name) {
		if original.Name != variant {
			dist := fuzzy.Levenshtein(original.Name, variant)
			domains = append(domains, &db.Domain{Name: variant, Levenshtein: dist, Algorithm: algo})
		}
	}

	return domains, err
}

// Register the plugin
func init() {
	var CODE = "cns"
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Plugin{
			Plugin: algorithms.Plugin{
				Code:    CODE,
				Title:   "Cardinal Substitution",
				Summary: "Swapping digial numbers and carninal numbers",
			},
		}
	})
}
