// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package usr checks whether a handle is taken on the platforms that issue
// them. The squat here is of an identity rather than of infrastructure, so
// what matters is whether the handle is claimed and by whom.
package usr

import (
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
)

// ID is the operator id, which is what --collect selects on.
const ID = "usr"

// New builds the operator, or nothing when no source list is available.
//
// Returning nothing rather than an operator that always fails is the rule the
// whole plan rests on: --explain must not promise work that cannot happen.
func New(o observe.Options, list observe.SourceLister, prober observe.Prober) []graph.Operator {
	if list == nil {
		return nil
	}
	return []graph.Operator{observe.NewSourceOp(o, ID, observe.TypeUsername, "username", list, prober)}
}
