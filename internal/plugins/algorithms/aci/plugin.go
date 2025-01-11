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
package aci

// Adjacent character insertion is where an attacker adds characters
// that are next to each other on a keyboard.

// For example, if a user intends to visit "example.com," a typo-squatter
// might register "examplw.com" or "exanple.com." These small alterations
// can trick users into clicking on the malicious sites, leading to phishing
// scams, malware downloads, or other harmful activities.
//
//               example.com  -> examplw.com
//                               exanple.com
//
// Adjacent character insertion exploits common typing errors, making it a
// particularly effective tactic, as users may not notice the difference,
// especially if they are typing quickly. It highlights the importance of
// vigilance and cybersecurity measures to protect against such deceptive
// practices.

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
	keyboards := p.Conf.Keyboards()
	prefix, name, suffix := dns.Split(original.Name)

	for _, keyboard := range keyboards {
		for _, variant := range typo.AdjacentCharacterInsertion(name, keyboard.Layouts()...) {
			if name != variant {
				variant = dns.Join(prefix, variant, suffix)
				dist := fuzzy.Levenshtein(original.Name, variant)
				domains = append(domains, &db.Domain{Name: variant, Levenshtein: dist, Algorithm: p.Algo()})
			}
		}
	}
	return
}

// Register the plugin
func init() {
	var CODE = "aci"
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Plugin{
			Plugin: algorithms.Plugin{
				Code:    CODE,
				Title:   "Adjacent Character Insertion",
				Summary: "Inserting adjacent character from the keyboard",
			},
		}
	})
}
