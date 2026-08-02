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

package decompose

import (
	"strings"

	"golang.org/x/net/publicsuffix"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Domain pulls a domain's registry suffix out as its own node. The suffix is
// worth converging on: a hundred variants under one unusual TLD is a campaign,
// and that is invisible if the TLD only ever exists as a substring of a key.
func Domain() graph.Operator {
	return &decomposer{
		id:      "decompose.domain",
		version: 1,
		on:      TypeDomain,
		emits: graph.Effects{
			Nodes: []string{TypeTLD},
			Rels:  []string{RelTLDOf},
		},
		split: splitDomain,
	}
}

// splitDomain finds the suffix under which the name was registered. The public
// suffix list is used rather than the last label because "example.co.uk" is
// registered under co.uk, not under uk, and keying on uk would group it with
// names that share nothing with it.
//
// The key is already lowercase punycode, which is what the list expects.
func splitDomain(key string) graph.Delta {
	suffix, icann := publicsuffix.PublicSuffix(key)
	if !icann {
		// A private suffix (blogspot.com) or an unrecognised one. Neither is a
		// registry, and admitting them as tld nodes would put a squatter's own
		// hosting domain on equal footing with a real TLD; fall back to the
		// last label, which always is one.
		if i := strings.LastIndex(suffix, "."); i >= 0 {
			suffix = suffix[i+1:]
		}
	}
	// A bare TLD, or a single-label host: there is no suffix to separate out.
	if suffix == "" || suffix == key {
		return graph.Delta{}
	}
	self := graph.NodeRef{Type: TypeDomain, Key: key}
	tld := graph.NodeRef{Type: TypeTLD, Key: suffix}
	return graph.Delta{
		Nodes: []graph.NodeRef{tld},
		Edges: []graph.EdgeRef{{From: self, Rel: RelTLDOf, To: tld}},
	}
}
