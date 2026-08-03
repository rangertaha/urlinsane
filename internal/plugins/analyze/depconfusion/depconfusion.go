// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package depconfusion flags the supply-chain case: a dependency that exists
// on one registry and nowhere else. Absence is the finding, which is why it can
// only be written against three-state existence — under a two-state model an
// unreachable registry and a free name are the same answer.
package depconfusion

import (
	"context"
	"fmt"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze"
)

type DepConfusion struct{}

func (DepConfusion) Id() string { return "dep-confusion" }

func (DepConfusion) Exec(_ context.Context, a *graph.Analysis) ([]graph.Finding, error) {
	var out []graph.Finding
	for _, n := range a.Nodes() {
		if n.Type.Name() != "package" {
			continue
		}
		// A generated variant is not a dependency. Being unpublished is the
		// normal condition of almost every one of them — that is what makes it
		// registrable — so reporting it as a supply-chain gap made every
		// variant of a package target a CRITICAL finding, hundreds per scan,
		// and `--fail-on` returned 2 for every package scan whatever level it
		// was given. The unregistered-lookalike case belongs to scoring, which
		// reports it only when the name is live.
		//
		// What is left is what the analyzer is actually about: a package this
		// run was told to care about — the seed, or one read out of a manifest —
		// that no public registry carries, where an attacker publishing it
		// would win resolution.
		if analyze.IsVariant(a, n.ID) {
			continue
		}
		if len(a.Outgoing(n.ID, "EXISTS_ON")) > 0 {
			continue
		}
		switch a.Existence(n.ID) {
		case graph.Absent:
			out = append(out, graph.Finding{
				Kind:     "dep-confusion",
				Severity: graph.SeverityCritical,
				Nodes:    []graph.NodeID{n.ID},
				Summary: fmt.Sprintf("%s is not published on its registry; a public package of that name would win resolution",
					n.Key),
			})
		case graph.Unknown:
			// Reporting this as a gap would be a guess. Saying so is the point
			// of having a third value.
			out = append(out, graph.Finding{
				Kind:     "dep-confusion-unknown",
				Severity: graph.SeverityInfo,
				Nodes:    []graph.NodeID{n.ID},
				Summary:  fmt.Sprintf("could not determine whether %s is published", n.Key),
			})
		}
	}
	return out, nil
}
