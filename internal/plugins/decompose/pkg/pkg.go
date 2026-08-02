// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package pkg splits a package reference into its owner scope and name.
package pkg

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
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
	return decompose.New("decompose.package", 1, decompose.TypePackage,
		graph.Effects{
			Nodes: []string{decompose.TypeUsername},
			Rels:  []string{decompose.RelOwner},
		},
		splitPackage)
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
	self := graph.NodeRef{Type: decompose.TypePackage, Key: key}
	user := graph.NodeRef{Type: decompose.TypeUsername, Key: owner}
	return graph.Delta{
		Nodes: []graph.NodeRef{user},
		Edges: []graph.EdgeRef{{From: self, Rel: decompose.RelOwner, To: user}},
	}
}
