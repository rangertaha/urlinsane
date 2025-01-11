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
package cb

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

type Plugin struct {
	algorithms.Plugin
}

func (p *Plugin) Exec(original *db.Domain) (domains []*db.Domain, err error) {
	//
	// dist := fuzzy.Levenshtein(original.Name, variant)
	// prefix, name, suffix := dns.Split(original.Name)
	// variant = dns.Join(prefix, variant, suffix)

	return
}

// Register the plugin
func init() {
	var CODE = "com"
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Plugin{
			Plugin: algorithms.Plugin{
				Code:    CODE,
				Title:   "Combo Squatting",
				Summary: "Creating domain names by combining keywords",
			},
		}
	})
}
