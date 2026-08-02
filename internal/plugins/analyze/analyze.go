// Copyright 2024 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package analyze holds whole-graph analyzers. They run once, after expansion
// stops, over a stable graph.
//
// Every analysis here is only expressible because infrastructure is nodes
// rather than fields hanging off a domain: "which variants share this IP" is an
// in-edge query, and it has no answer in a model where each domain carries its
// own private copy of its addresses.
package analyze

import (
	"context"
	"fmt"
	"sort"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// All returns the standard analyzer set.
func All() []graph.Analyzer {
	return []graph.Analyzer{Campaign{}, Scoring{}, DepConfusion{}}
}

// Campaign clusters variants that share infrastructure. Two typo variants of a
// brand resolving to the same address, sitting behind the same nameserver, or
// registered by the same party are one actor, and that is a materially
// different finding from two unrelated squats.
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
				if !isVariant(a, e.From) {
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

// isVariant reports whether a node was reached by a VARIANT_OF edge.
func isVariant(a *graph.Analysis, id graph.NodeID) bool {
	return len(a.Incoming(id, graph.VariantRel)) > 0
}

// Scoring rates a variant on the signals that separate a parked typo from an
// operational one. It is rule-based and entirely independent of the execution
// model, which never contributes to a reported number.
type Scoring struct{}

func (Scoring) Id() string { return "scoring" }

func (Scoring) Exec(_ context.Context, a *graph.Analysis) ([]graph.Finding, error) {
	var out []graph.Finding
	for _, n := range a.Nodes() {
		if !isVariant(a, n.ID) {
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
		if d := editDistance(a, n.ID); d > 0 && d <= 1 {
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

// editDistance reads the distance off the VARIANT_OF edge that produced this
// node. The algorithm already computed it; recomputing here would risk
// disagreeing with the edge.
func editDistance(a *graph.Analysis, id graph.NodeID) int64 {
	for _, e := range a.Incoming(id, graph.VariantRel) {
		if f, ok := e.Rel.Field("distance"); ok {
			if v, set := e.Props.Get(f); set {
				return v.Num()
			}
		}
	}
	return 0
}

// DepConfusion finds internal package names absent from the public registry
// they would be resolved against. Unlike the other analyses this one is about
// something that does *not* exist, which is why it depends on the status
// vocabulary distinguishing "confirmed absent" from "could not determine".
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
