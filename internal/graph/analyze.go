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

package graph

import (
	"context"
	"sort"
)

// Severity is an ordered, named level. It is deliberately not a bare int:
// `--fail-on high` has to mean the same thing in every release, and an integer
// scale drifts silently the first time someone inserts a level.
type Severity uint8

const (
	SeverityInfo Severity = iota + 1
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	case SeverityCritical:
		return "critical"
	}
	return "none"
}

// ParseSeverity resolves a level name, for --fail-on.
func ParseSeverity(s string) (Severity, bool) {
	for _, v := range []Severity{SeverityInfo, SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical} {
		if v.String() == s {
			return v, true
		}
	}
	return 0, false
}

// LedgerRef names a declined candidate.
type LedgerRef struct {
	Type string
	Key  string
}

// Finding is an analyzer's conclusion.
type Finding struct {
	Kind     string
	Severity Severity
	Nodes    []NodeID    // admitted nodes this concerns
	Declined []LedgerRef // ledger rows this concerns
	Summary  string
	Evidence []Provenance
}

// Analyzer runs once, over the whole graph, after expansion stops — for any
// reason, including an interrupt. There is exactly one lifetime.
type Analyzer interface {
	Id() string
	Exec(ctx context.Context, a *Analysis) ([]Finding, error)
}

// Analysis is the read-only surface analyzers get. It exposes nodes, edges,
// status, provenance, plugin scores and the truncation ledger — but not the
// engine's belief.
//
// Withholding belief is what makes "the execution model never contributes to a
// reported number" true by construction rather than by convention: an analyzer
// that could read it could launder it into a Severity and no reviewer would
// spot it. A plugin score is different — it describes the entity, not the
// traversal — so those are visible.
type Analysis struct {
	g *Graph
}

// Analyze returns the read-only analysis surface.
func (g *Graph) Analyze() *Analysis { return &Analysis{g: g} }

func (a *Analysis) Nodes() []*Node                   { return a.g.Nodes() }
func (a *Analysis) Edges() []*Edge                   { return a.g.Edges() }
func (a *Analysis) Node(id NodeID) (*Node, bool)     { return a.g.Node(id) }
func (a *Analysis) Depth(id NodeID) int              { return a.g.Depth(id) }
func (a *Analysis) InClosure(id NodeID) bool         { return a.g.InClosure(id) }
func (a *Analysis) Assertions(id NodeID) []Assertion { return a.g.Assertions(id) }
func (a *Analysis) Ledger() []LedgerRow              { return a.g.Ledger() }
func (a *Analysis) Truncations() []RunTruncation     { return a.g.Truncations() }
func (a *Analysis) Rejections() []Rejection          { return a.g.Rejections() }
func (a *Analysis) Status(id NodeID, op string) (Status, bool) {
	return a.g.Status(id, op)
}
func (a *Analysis) Score(id NodeID, key string) (float64, bool) { return a.g.Score(id, key) }

// StatusRow pairs an operator with its terminal outcome for one node.
type StatusRow struct {
	Operator string
	Status   Status
}

// Statuses enumerates a node's recorded outcomes, sorted by operator.
// Status(id, op) answers a question you already know to ask; a report has to
// enumerate, and iterating the side table directly would expose map order.
func (a *Analysis) Statuses(id NodeID) []StatusRow {
	var out []StatusRow
	for k, s := range a.g.status {
		if k.node == id {
			out = append(out, StatusRow{Operator: k.op, Status: s})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Operator < out[j].Operator })
	return out
}

// ScoreRow pairs a plugin model's key with its score.
type ScoreRow struct {
	Key   string
	Value float64
}

// Scores enumerates a node's plugin-model scores, sorted by key. Engine belief
// is deliberately absent — see the type comment.
func (a *Analysis) Scores(id NodeID) []ScoreRow {
	var out []ScoreRow
	for k, v := range a.g.scores[id] {
		out = append(out, ScoreRow{Key: k, Value: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// Existence is the three-valued rollup of a node's observation statuses.
type Existence uint8

const (
	// Live — at least one observation operator returned ok.
	Live Existence = iota + 1
	// Absent — none did, and at least one authoritatively determined absence.
	Absent
	// Unknown — every attempt failed, timed out or was skipped.
	Unknown
)

func (e Existence) String() string {
	switch e {
	case Live:
		return "live"
	case Absent:
		return "absent"
	case Unknown:
		return "unknown"
	}
	return "untried"
}

// Existence rolls a node's statuses up. Three values, not two: collapsing
// "confirmed absent" into "unregistered" discards the distinction the whole
// status vocabulary exists to preserve — a variant nobody could resolve is not
// a variant proven free.
func (a *Analysis) Existence(id NodeID) Existence {
	var sawEmpty, sawAny bool
	for k, s := range a.g.status {
		if k.node != id {
			continue
		}
		sawAny = true
		switch s {
		case StatusOK:
			return Live
		case StatusEmpty:
			sawEmpty = true
		}
	}
	switch {
	case sawEmpty:
		return Absent
	case sawAny:
		return Unknown
	}
	return 0
}

// Outgoing returns a node's outgoing edges of a relation.
func (a *Analysis) Outgoing(id NodeID, rel string) []*Edge { return a.g.outEdges(id, rel) }

// Incoming returns edges pointing at a node, optionally filtered by relation.
// Analyzers need this to cluster: "which variants share this IP" is an
// in-edge query, and it is only answerable because infrastructure is nodes.
func (a *Analysis) Incoming(id NodeID, rel string) []*Edge {
	var out []*Edge
	for _, eid := range a.g.eord {
		e := a.g.edges[eid]
		if e.To != id {
			continue
		}
		if rel != "" && e.Rel.name != rel {
			continue
		}
		out = append(out, e)
	}
	return out
}

// Findings stores an analyzer's conclusions in the side table.
func (g *Graph) AddFindings(f ...Finding) { g.findings = append(g.findings, f...) }

// Findings returns every finding, most severe first, then stably by kind and
// summary so two runs render identically.
func (g *Graph) Findings() []Finding {
	out := append([]Finding(nil), g.findings...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return out[i].Severity > out[j].Severity
		}
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].Summary < out[j].Summary
	})
	return out
}

// Risk is the maximum severity among findings referencing a node. It is what
// --filter risk>N and --fail-on select on, and the only user-facing score in
// the system.
func (g *Graph) Risk(id NodeID) Severity {
	var max Severity
	for _, f := range g.findings {
		for _, n := range f.Nodes {
			if n == id && f.Severity > max {
				max = f.Severity
			}
		}
	}
	return max
}

// SetScore records a plugin model's judgement about an entity. Scores are side
// table state, not props: as props they would make every node's CID depend on a
// model version and break cross-run diffing the moment anything is retrained.
func (g *Graph) SetScore(id NodeID, key string, v float64) {
	if g.scores[id] == nil {
		g.scores[id] = map[string]float64{}
	}
	g.scores[id][key] = v
}

// Score returns a plugin model's score for an entity.
func (g *Graph) Score(id NodeID, key string) (float64, bool) {
	v, ok := g.scores[id][key]
	return v, ok
}

// RunAnalyzers runs every analyzer once and records their findings.
func (g *Graph) RunAnalyzers(ctx context.Context, analyzers []Analyzer) error {
	sorted := append([]Analyzer(nil), analyzers...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Id() < sorted[j].Id() })
	a := g.Analyze()
	for _, an := range sorted {
		f, err := an.Exec(ctx, a)
		if err != nil {
			return err
		}
		g.AddFindings(f...)
	}
	return nil
}
