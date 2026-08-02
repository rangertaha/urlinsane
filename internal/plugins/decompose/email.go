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

// Email splits an address into the account it names and the domain that hosts
// it: bob@example.com is both an account called bob and a domain called
// example.com, and a squatting scan wants variants of each.
func Email() graph.Operator {
	return &decomposer{
		id:      "decompose.email",
		version: 1,
		on:      TypeEmail,
		emits: graph.Effects{
			Nodes: []string{TypeUsername, TypeDomain},
			Rels:  []string{RelLocalPart, RelDomainOf},
		},
		split: splitEmail,
	}
}

// splitEmail parses an already-canonical address. Canonicalization guarantees
// the "@" and both sides, so a failure here means the registry and this
// operator disagree — in which case emitting nothing is the safe answer.
func splitEmail(key string) graph.Delta {
	at := strings.LastIndex(key, "@")
	if at <= 0 || at == len(key)-1 {
		return graph.Delta{}
	}
	self := graph.NodeRef{Type: TypeEmail, Key: key}
	// Raw keys: the local part keeps its case here and the applier folds it,
	// because computing the canonical form — or worse, a NodeID — inside an
	// operator is how convergence quietly breaks.
	user := graph.NodeRef{Type: TypeUsername, Key: key[:at]}
	host := graph.NodeRef{Type: TypeDomain, Key: key[at+1:]}
	return graph.Delta{
		Nodes: []graph.NodeRef{user, host},
		Edges: []graph.EdgeRef{
			{From: self, Rel: RelLocalPart, To: user},
			{From: self, Rel: RelDomainOf, To: host},
		},
	}
}
