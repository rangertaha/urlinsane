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

// Package decomposes a registry-qualified package reference such as
// "npm:lodash". Most of its work is done before it runs: qualifying the key by
// registry is canonPackage's job, and it is what keeps npm's lodash and PyPI's
// lodash from converging into one node.
//
// What is left to decompose is the namespace of a scoped name — "npm:@acme/ui"
// belongs to the account acme, and dependency-confusion analysis turns on
// exactly that ownership.
func Package() graph.Operator {
	return &decomposer{
		id:      "decompose.package",
		version: 1,
		on:      TypePackage,
		emits: graph.Effects{
			Nodes: []string{TypeUsername},
			Rels:  []string{RelOwner},
		},
		split: splitPackage,
	}
}

// splitPackage extracts the owning namespace, if the name has one. An
// unqualified name such as "npm:lodash" has nothing to decompose, which is an
// answer rather than a failure.
//
// The "@scope/name" form is treated as syntax rather than as an npm feature: it
// is the shape a scope takes wherever a registry has scopes, and a registry
// without them simply never produces a name of that shape.
func splitPackage(key string) graph.Delta {
	_, name, ok := strings.Cut(key, ":")
	if !ok {
		return graph.Delta{}
	}
	scope, _, ok := strings.Cut(name, "/")
	if !ok || !strings.HasPrefix(scope, "@") {
		return graph.Delta{}
	}
	owner := strings.TrimPrefix(scope, "@")
	if owner == "" {
		return graph.Delta{}
	}
	self := graph.NodeRef{Type: TypePackage, Key: key}
	user := graph.NodeRef{Type: TypeUsername, Key: owner}
	return graph.Delta{
		Nodes: []graph.NodeRef{user},
		Edges: []graph.EdgeRef{{From: self, Rel: RelOwner, To: user}},
	}
}
