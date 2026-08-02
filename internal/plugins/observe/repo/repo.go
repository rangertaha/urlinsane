// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package repo checks whether a repository exists on the public forges.
//
// The old collector bound to package and name entities; a repo is its own type
// now, so it binds where it belongs.
package repo

import (
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
)

// ID is the operator id, which is what --collect selects on.
const ID = "repo"

// New builds the operator, or nothing when no source list is available.
//
// Returning nothing rather than an operator that always fails is the rule the
// whole plan rests on: --explain must not promise work that cannot happen.
func New(o observe.Options, list observe.SourceLister, prober observe.Prober) []graph.Operator {
	if list == nil {
		return nil
	}
	return []graph.Operator{observe.NewSourceOp(o, ID, observe.TypeRepo, "repository", list, prober)}
}
