// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package campaign finds variants that cluster: shared addresses, shared
// nameservers, shared registrant. One squatted lookalike is an incident; forty
// on one nameserver is a campaign, and only the graph can see the difference.
package campaign

import (
	"context"
	"fmt"
	"sort"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze"
)

type Campaign struct {
	// Rels are the shared-infrastructure relations to cluster on.
	Rels []string
	// Min is the smallest cluster worth reporting. Two is the honest floor:
	// one variant on an address is not a campaign.
	Min int
}

func (Campaign) Id() string { return "campaign" }

func (c Campaign) Exec(_ context.Context, a *graph.Analysis) ([]graph.Finding, error) {
	rels := c.Rels
	if len(rels) == 0 {
		rels = []string{"RESOLVES_TO", "NS", "REGISTERED_BY"}
	}
	min := c.Min
	if min < 2 {
		min = 2
	}

	var out []graph.Finding
	for _, hub := range a.Nodes() {
		for _, rel := range rels {
			in := a.Incoming(hub.ID, rel)
			if len(in) < min {
				continue
			}
			// Only variants count. A hub shared by the origin and one variant
			// is ordinary hosting, not a campaign.
			var members []graph.NodeID
			var names []string
			for _, e := range in {
				if !analyze.IsVariant(a, e.From) {
					continue
				}
				members = append(members, e.From)
				if n, ok := a.Node(e.From); ok {
					names = append(names, n.Key)
				}
			}
			if len(members) < min {
				continue
			}
			sort.Strings(names)
			out = append(out, graph.Finding{
				Kind:     "campaign",
				Severity: graph.SeverityHigh,
				Nodes:    append(members, hub.ID),
				Summary: fmt.Sprintf("%d variants share %s %s: %v",
					len(members), rel, hub.Key, names),
			})
		}
	}
	return out, nil
}

// Scoring rates a variant on the signals that separate a parked typo from an
// operational one. It is rule-based and entirely independent of the execution
// model, which never contributes to a reported number.
