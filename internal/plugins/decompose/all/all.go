// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package all composes the shipped decomposers.
//
// Separate from the decompose library so each decomposer can import that
// library — for the shared operator shape and the schema vocabulary — without
// importing its siblings through it.
package all

import (
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose/email"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose/pkg"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose/repo"
)

// Operators returns the shipped decomposers.
//
// The order is the plan's listing order, not a run order: decomposers bind to
// node types and the scheduler decides what runs when (§4.1).
func Operators() []graph.Operator {
	return []graph.Operator{email.Email(), domain.Domain(), repo.Repo(), pkg.Package()}
}
