// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package email splits an address into its local part and its domain, so a
// scan of one address also varies the two halves independently.
package email

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
)

// Email splits an address into the account it names and the domain that hosts
// it: bob@example.com is both an account called bob and a domain called
// example.com, and a squatting scan wants variants of each.
func Email() graph.Operator {
	return decompose.New("decompose.email", 1, decompose.TypeEmail,
		graph.Effects{
			Nodes: []string{decompose.TypeUsername, decompose.TypeDomain},
			Rels:  []string{decompose.RelLocalPart, decompose.RelDomainOf},
		},
		splitEmail)
}

// splitEmail parses an already-canonical address. Canonicalization guarantees
// the "@" and both sides, so a failure here means the registry and this
// operator disagree — in which case emitting nothing is the safe answer.
func splitEmail(key string) graph.Delta {
	at := strings.LastIndex(key, "@")
	if at <= 0 || at == len(key)-1 {
		return graph.Delta{}
	}
	self := graph.NodeRef{Type: decompose.TypeEmail, Key: key}
	// Raw keys: the local part keeps its case here and the applier folds it,
	// because computing the canonical form — or worse, a NodeID — inside an
	// operator is how convergence quietly breaks.
	user := graph.NodeRef{Type: decompose.TypeUsername, Key: key[:at]}
	host := graph.NodeRef{Type: decompose.TypeDomain, Key: key[at+1:]}
	return graph.Delta{
		Nodes: []graph.NodeRef{user, host},
		Edges: []graph.EdgeRef{
			{From: self, Rel: decompose.RelLocalPart, To: user},
			{From: self, Rel: decompose.RelDomainOf, To: host},
		},
	}
}
