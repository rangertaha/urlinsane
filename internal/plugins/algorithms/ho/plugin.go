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
package ho

// Missing Dashes
//
// Func omits a dash from the name. For example:
//
// Original: www.one-two-three.com
//
// Variants: www.onetwo-three.com
//           www.one-twothree.com

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/pkg/dns"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

type Plugin struct {
	algorithms.Plugin
}

func (p *Plugin) Exec(original *db.Domain) (domains []*db.Domain, err error) {
	algo := db.Algorithm{Code: p.Code, Name: p.Title}
	prefix, name, suffix := dns.Split(original.Name)

	for _, variant := range typo.HyphenOmission(name) {
		if name != variant {
			variant = dns.Join(prefix, variant, suffix)
			dist := fuzzy.Levenshtein(original.Name, variant)
			domains = append(domains, &db.Domain{Name: variant, Levenshtein: dist, Algorithm: algo})
		}
	}

	return
}

// Register the plugin
func init() {
	var CODE = "ho"
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Plugin{
			Plugin: algorithms.Plugin{
				Code:    CODE,
				Title:   "Hyphen Omission",
				Summary: "Created by removing hyphens from the name",
			},
		}
	})
}
