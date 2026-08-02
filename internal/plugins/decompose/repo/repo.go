// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package repo splits a repository reference into forge, owner and name —
// each of which can be squatted independently.
package repo

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
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
	return decompose.New("decompose.repo", 1, decompose.TypeRepo,
		graph.Effects{
			Nodes: []string{decompose.TypePlatform, decompose.TypeUsername},
			Rels:  []string{decompose.RelHostedOn, decompose.RelOwner},
		},
		splitRepo)
}

// splitRepo parses the canonical host/owner/name form. Every other spelling was
// already collapsed into it by canonRepo, so there is exactly one shape here.
func splitRepo(key string) graph.Delta {
	parts := strings.Split(key, "/")
	if len(parts) != 3 {
		return graph.Delta{}
	}
	self := graph.NodeRef{Type: decompose.TypeRepo, Key: key}
	platform := graph.NodeRef{Type: decompose.TypePlatform, Key: parts[0]}
	owner := graph.NodeRef{Type: decompose.TypeUsername, Key: parts[1]}
	return graph.Delta{
		Nodes: []graph.NodeRef{platform, owner},
		Edges: []graph.EdgeRef{
			{From: self, Rel: decompose.RelHostedOn, To: platform},
			{From: self, Rel: decompose.RelOwner, To: owner},
		},
	}
}
