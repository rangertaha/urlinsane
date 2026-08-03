// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package pkg checks whether a package name exists on the public registries.
//
// This is the supply-chain surface. A dependency that resolves on one registry
// and nowhere else is what dependency confusion looks like, and it is only
// detectable because absence is recorded as absence rather than as failure.
package pkg

import (
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
)

// ID is the operator id, which is what --collect selects on.
const ID = "pkg"

// New builds the operator, or nothing when no source list is available.
//
// Returning nothing rather than an operator that always fails is the rule the
// whole plan rests on: --explain must not promise work that cannot happen.
func New(o observe.Options, list observe.SourceLister, prober observe.Prober) []graph.Operator {
	// Nothing configured for this kind is the same as no lister at all: the
	// operator could only ever report that it has no sources.
	if !observe.HasSources(list, "package") {
		return nil
	}
	return []graph.Operator{observe.NewSourceOp(o, ID, observe.TypePackage, "package", list, prober)}
}
