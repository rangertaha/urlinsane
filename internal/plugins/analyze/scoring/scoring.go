// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package scoring ranks a variant by how plausible the confusion is —
// edit distance, whether it resolves, what it resolves to.
package scoring

import (
	"context"
	"fmt"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze"
)

type Scoring struct{}

func (Scoring) Id() string { return "scoring" }

func (Scoring) Exec(_ context.Context, a *graph.Analysis) ([]graph.Finding, error) {
	var out []graph.Finding
	for _, n := range a.Nodes() {
		if !analyze.IsVariant(a, n.ID) {
			continue
		}
		if a.Existence(n.ID) != graph.Live {
			continue
		}

		var signals []string
		score := 0
		if len(a.Outgoing(n.ID, "RESOLVES_TO")) > 0 {
			score += 2
			signals = append(signals, "resolves")
		}
		if len(a.Outgoing(n.ID, "MX")) > 0 {
			// Mail is the signal that separates a parked name from one built
			// to receive credentials.
			score += 3
			signals = append(signals, "accepts mail")
		}
		if d := analyze.EditDistance(a, n.ID); d > 0 && d <= 1 {
			score += 2
			signals = append(signals, "one edit from the target")
		}

		sev := graph.SeverityLow
		switch {
		case score >= 6:
			sev = graph.SeverityCritical
		case score >= 4:
			sev = graph.SeverityHigh
		case score >= 2:
			sev = graph.SeverityMedium
		}
		out = append(out, graph.Finding{
			Kind:     "live-variant",
			Severity: sev,
			Nodes:    []graph.NodeID{n.ID},
			Summary:  fmt.Sprintf("%s is live: %v", n.Key, signals),
		})
	}
	return out, nil
}

// DepConfusion finds internal package names absent from the public registry
// they would be resolved against. Unlike the other analyses this one is about
// something that does *not* exist, which is why it depends on the status
// vocabulary distinguishing "confirmed absent" from "could not determine".
