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

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Repo splits a repository reference into the platform hosting it and the
// account owning it. Both are targets in their own right: typosquatting a
// repository starts with typosquatting its owner's account name.
//
// It deliberately does not emit MANIFEST. Reading a repository's dependency
// manifest is a network call, which makes that an observation edge — and that
// classification is the only thing keeping "typo github.com/acme/tool" from
// generating variants of several hundred declared dependencies (§8).
func Repo() graph.Operator {
	return &decomposer{
		id:      "decompose.repo",
		version: 1,
		on:      TypeRepo,
		emits: graph.Effects{
			Nodes: []string{TypePlatform, TypeUsername},
			Rels:  []string{RelHostedOn, RelOwner},
		},
		split: splitRepo,
	}
}

// splitRepo parses the canonical host/owner/name form. Every other spelling was
// already collapsed into it by canonRepo, so there is exactly one shape here.
func splitRepo(key string) graph.Delta {
	parts := strings.Split(key, "/")
	if len(parts) != 3 {
		return graph.Delta{}
	}
	self := graph.NodeRef{Type: TypeRepo, Key: key}
	platform := graph.NodeRef{Type: TypePlatform, Key: parts[0]}
	owner := graph.NodeRef{Type: TypeUsername, Key: parts[1]}
	return graph.Delta{
		Nodes: []graph.NodeRef{platform, owner},
		Edges: []graph.EdgeRef{
			{From: self, Rel: RelHostedOn, To: platform},
			{From: self, Rel: RelOwner, To: owner},
		},
	}
}
