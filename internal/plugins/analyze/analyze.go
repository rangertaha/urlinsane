// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package analyze is the small shared library the analyzers are built on.
//
// It holds the graph queries more than one analyzer needs and nothing else.
// The analyzers themselves are packages under this one; the composed list is
// in analyze/all, which imports them — so that this package can be imported by
// them without a cycle.
package analyze

import (
	"github.com/rangertaha/urlinsane/internal/graph"
)

// isVariant reports whether a node was reached by a VARIANT_OF edge.
func IsVariant(a *graph.Analysis, id graph.NodeID) bool {
	return len(a.Incoming(id, graph.VariantRel)) > 0
}

// editDistance reads the distance off the VARIANT_OF edge that produced this
// node. The algorithm already computed it; recomputing here would risk
// disagreeing with the edge.
func EditDistance(a *graph.Analysis, id graph.NodeID) int64 {
	for _, e := range a.Incoming(id, graph.VariantRel) {
		if f, ok := e.Rel.Field("distance"); ok {
			if v, set := e.Props.Get(f); set {
				return v.Num()
			}
		}
	}
	return 0
}
