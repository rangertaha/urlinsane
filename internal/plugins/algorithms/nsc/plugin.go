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

// Package nsc implements namespace/scope confusion, a package-registry squatting
// pattern. Scoped names (npm's @org/pkg) and org/repo paths can be impersonated
// by registering the unscoped name, a flattened org-pkg name, or the reverse —
// turning an unscoped name into a look-alike scoped one.
package nsc

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
)

type Plugin struct {
	algorithms.Plugin
}

// Exec generates scope/namespace variants depending on the name's shape.
func (p *Plugin) Exec(original *db.Domain) (domains []*db.Domain, err error) {
	name := original.Name
	var variants []string

	switch {
	case strings.HasPrefix(name, "@") && strings.Contains(name, "/"):
		// Scoped npm name: @org/pkg
		org, pkg, _ := strings.Cut(strings.TrimPrefix(name, "@"), "/")
		variants = append(variants, pkg, org+"-"+pkg, org+"/"+pkg, org+pkg)
	case strings.Contains(name, "/"):
		// Path-style name: org/pkg
		org, pkg, _ := strings.Cut(name, "/")
		variants = append(variants, pkg, org+"-"+pkg, "@"+org+"/"+pkg)
	case strings.Contains(name, "-"):
		// Flat name: org-pkg -> scoped / path look-alikes
		org, pkg, _ := strings.Cut(name, "-")
		variants = append(variants, "@"+org+"/"+pkg, org+"/"+pkg)
	}

	seen := map[string]bool{name: true, "": true}
	for _, v := range variants {
		if seen[v] {
			continue
		}
		seen[v] = true
		dist := fuzzy.Levenshtein(name, v)
		domains = append(domains, &db.Domain{Name: v, Levenshtein: dist, Algorithm: p.Algo()})
	}
	return
}

// Register the plugin
func init() {
	var CODE = "nsc"
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Plugin{
			Plugin: algorithms.Plugin{
				Code:     CODE,
				Title:    "Namespace Confusion",
				Summary:  "Confusing scoped/unscoped names (@org/pkg, org-pkg)",
				Entities: []entity.Type{entity.Package},
			},
		}
	})
}
