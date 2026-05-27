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

// Package afx implements affix (combosquatting) generation, a package-registry
// squatting pattern. Attackers bracket a popular name with a plausible
// language/role token — py-requests, requests2, requests-js — so the result
// still reads as "the requests package." Longer affixes raise the edit distance,
// so widen --distance to surface them all.
package afx

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
)

type Plugin struct {
	algorithms.Plugin
}

// prefixes and suffixes are common combosquatting affixes seen in real
// registry attacks (ecosystem hints, versions, and role qualifiers).
var (
	prefixes = []string{"py", "py-", "python-", "node-", "js-", "go-", "lib", "lib-", "the-", "get-"}
	suffixes = []string{"2", "3", "js", "-js", "py", "-py", "-python", "-cli", "-dev",
		"-core", "-utils", "-api", "-sdk", "-lib", ".js", "-ng", "-master", "-official"}
)

// Exec brackets the name with each common ecosystem/role affix.
func (p *Plugin) Exec(original *db.Domain) (domains []*db.Domain, err error) {
	name := original.Name
	seen := map[string]bool{name: true}

	add := func(variant string) {
		if variant == "" || seen[variant] {
			return
		}
		seen[variant] = true
		dist := fuzzy.Levenshtein(name, variant)
		domains = append(domains, &db.Domain{Name: variant, Levenshtein: dist, Algorithm: p.Algo()})
	}

	for _, pre := range prefixes {
		add(pre + name)
	}
	for _, suf := range suffixes {
		add(name + suf)
	}
	return
}

// Register the plugin
func init() {
	var CODE = "afx"
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Plugin{
			Plugin: algorithms.Plugin{
				Code:     CODE,
				Title:    "Affix Squatting",
				Summary:  "Adding an ecosystem/role prefix or suffix (py-, -js, 2)",
				Entities: []entity.Type{entity.Package},
			},
		}
	})
}
