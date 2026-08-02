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
)

type DepConfusion struct{}

func (DepConfusion) Id() string { return "dep-confusion" }

func (DepConfusion) Exec(_ context.Context, a *graph.Analysis) ([]graph.Finding, error) {
	var out []graph.Finding
	for _, n := range a.Nodes() {
		if n.Type.Name() != "package" {
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
