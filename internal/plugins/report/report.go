// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package report renders a finished scan.
//
// Every format is projected from one intermediate [Report], built once. Five
// renderers reading the graph directly would drift — the JSON would gain a
// field the table never showed, a filter would be applied in four places and
// three ways — and "the formats disagree about what the scan found" is not a
// defect a user can work around.
//
// Report formats are canonically ordered and byte-identical across runs, so
// "what changed since last week" is a diff rather than noise. ndjson is the one
// exception and says so: it is an event stream, and makes no ordering claim.
package report

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Options controls what is rendered and how.
type Options struct {
	// Format is one of table, json, ndjson, csv, dot.
	Format string
	// Filters select which nodes appear. An empty set selects everything.
	Filters []Filter
	// Partial marks a report whose expansion stopped early — an interrupt, a
	// deadline, a budget. §11 requires it: a truncated report that reads as
	// complete is a correctness bug.
	Partial bool
	// PartialWhy explains the early stop in one clause.
	PartialWhy string
	// Verbose includes per-node provenance and engine belief. Belief is
	// diagnostic only and never reaches a finding (§9).
	Verbose bool
	// Color enables ANSI styling in the table renderer.
	Color bool
	// Target and Scope describe the invocation, so a saved report is
	// self-describing.
	Target string
	Scope  []string
	// Plan is the compiled plan hash, which pins operators, config and models.
	Plan string
	// Rounds is how many dispatch generations ran.
	Rounds int
	// Elapsed is wall-clock scan time. It is deliberately excluded from every
	// canonical rendering — see Report.Elapsed.
	Elapsed time.Duration
}

// Report is the rendered view of a finished scan: the single intermediate every
// format projects from.
type Report struct {
	Target      string      `json:"target"`
	Scope       []string    `json:"scope,omitempty"`
	Plan        string      `json:"plan,omitempty"`
	Partial     bool        `json:"partial"`
	PartialWhy  string      `json:"partial_reason,omitempty"`
	Rounds      int         `json:"rounds"`
	Nodes       []NodeRow   `json:"nodes"`
	Edges       []EdgeRow   `json:"edges"`
	Findings    []Finding   `json:"findings"`
	Ledger      []Declined  `json:"declined"`
	Truncations []Truncated `json:"truncations"`
	Totals      Totals      `json:"totals"`

	// Elapsed is reported out-of-band, never inside a canonical rendering: a
	// duration differs on every run and would make byte-comparison useless
	// precisely when it matters most.
	Elapsed time.Duration `json:"-"`
}

// Totals summarize the scan. They count the *unfiltered* graph — a filtered
// report that reported filtered totals would misstate the scan's size, which is
// the one number a reader uses to judge coverage.
type Totals struct {
	Nodes     int            `json:"nodes"`
	Edges     int            `json:"edges"`
	Shown     int            `json:"shown"`
	Live      int            `json:"live"`
	Absent    int            `json:"absent"`
	Unknown   int            `json:"unknown"`
	Declined  int            `json:"declined"`
	ByType    map[string]int `json:"by_type"`
	typeOrder []string
}

// TypeOrder is the order node types were first seen, for a renderer that wants
// counts in a stable order rather than map order.
//
// A method rather than an exported field so it stays out of the JSON: the
// ordering is a presentation concern, and adding it to the document would
// change the serialized form for every consumer to serve the table.
func (t Totals) TypeOrder() []string { return t.typeOrder }

// NodeRow is one admitted node.
type NodeRow struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Key       string `json:"key"`
	Depth     int    `json:"depth"`
	Existence string `json:"existence"`
	Risk      string `json:"risk,omitempty"`
	InClosure bool   `json:"in_closure"`
	// Seed marks the target. Depth cannot stand in for it: variant edges cost
	// no depth, so every variant of the target sits at depth 0 alongside it.
	Seed      bool     `json:"seed,omitempty"`
	Props     []Prop   `json:"props,omitempty"`
	Statuses  []Status `json:"statuses,omitempty"`
	Scores    []Score  `json:"scores,omitempty"`
	Claims    []Claim  `json:"claims,omitempty"` // verbose only
	Belief    *float64 `json:"belief,omitempty"` // verbose only
	severity  graph.Severity
	existence graph.Existence
}

// EdgeRow is one admitted edge.
type EdgeRow struct {
	From  string `json:"from"`
	Rel   string `json:"rel"`
	To    string `json:"to"`
	Class string `json:"class"`
	Props []Prop `json:"props,omitempty"`
}

// Prop is a field assertion in structural order.
type Prop struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
	Kind  string `json:"kind"`
}

// Text renders a prop value for the line-oriented formats.
func (p Prop) Text() string { return text(p.Value) }

// Status is one operator's terminal outcome for a node.
type Status struct {
	Operator string `json:"operator"`
	Status   string `json:"status"`
}

// Score is a plugin model's judgement about an entity. Plugin scores are
// evidence about the thing; engine belief is not, and is verbose-only (§9).
type Score struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

// Claim is one operator's assertion about a field, winning or not. Retaining
// losers is the point: disagreement between two sources is signal.
type Claim struct {
	Field    string `json:"field"`
	Value    any    `json:"value"`
	Operator string `json:"operator"`
	Round    int    `json:"round"`
	Won      bool   `json:"won"`
}

// Finding is an analyzer's conclusion.
type Finding struct {
	Kind     string   `json:"kind"`
	Severity string   `json:"severity"`
	Summary  string   `json:"summary"`
	Nodes    []string `json:"nodes,omitempty"`
	Declined []string `json:"declined,omitempty"`
	Evidence []string `json:"evidence,omitempty"`
}

// Declined is a candidate the engine refused to admit.
type Declined struct {
	Type   string  `json:"type"`
	Key    string  `json:"key"`
	Depth  int     `json:"depth"`
	Reason string  `json:"reason"`
	By     string  `json:"by"`
	Round  int     `json:"round"`
	Belief float64 `json:"belief,omitempty"`
}

// Truncated is a run-level limit that bound expansion.
type Truncated struct {
	Reason string `json:"reason"`
	Round  int    `json:"round"`
	Detail string `json:"detail"`
}

// Build projects a graph into the intermediate every renderer reads.
func Build(g *graph.Graph, o Options) Report {
	a := g.Analyze()
	r := Report{
		Target:     o.Target,
		Scope:      o.Scope,
		Plan:       o.Plan,
		Partial:    o.Partial,
		PartialWhy: o.PartialWhy,
		Rounds:     o.Rounds,
		Elapsed:    o.Elapsed,
		Totals:     Totals{ByType: map[string]int{}},
	}

	all := g.Nodes() // already sorted by (type, key)
	kept := map[graph.NodeID]bool{}
	for _, n := range all {
		row := buildNode(g, a, n, o.Verbose)

		r.Totals.Nodes++
		if _, seen := r.Totals.ByType[row.Type]; !seen {
			r.Totals.typeOrder = append(r.Totals.typeOrder, row.Type)
		}
		r.Totals.ByType[row.Type]++
		switch row.existence {
		case graph.Live:
			r.Totals.Live++
		case graph.Absent:
			r.Totals.Absent++
		case graph.Unknown:
			r.Totals.Unknown++
		}

		if !keep(o.Filters, row) {
			continue
		}
		kept[n.ID] = true
		r.Nodes = append(r.Nodes, row)
	}
	r.Totals.Shown = len(r.Nodes)

	for _, e := range g.Edges() { // already sorted by (from, rel, to)
		r.Totals.Edges++
		// An edge to a filtered-out node would dangle: dot would draw a node
		// the table never listed, and json would reference an absent id.
		if !kept[e.From] || !kept[e.To] {
			continue
		}
		r.Edges = append(r.Edges, buildEdge(g, e))
	}

	// Findings and the ledger are never filtered. §11 requires the ledger
	// unconditionally, and letting `--filter live` hide a finding would hide
	// exactly what `--fail-on` gates on.
	for _, f := range g.Findings() { // already sorted, most severe first
		r.Findings = append(r.Findings, buildFinding(g, f))
	}
	for _, d := range g.Ledger() { // already sorted by (type, key, reason)
		r.Ledger = append(r.Ledger, Declined{
			Type: d.Type, Key: d.Key, Depth: d.Depth, Reason: d.Reason.String(),
			By: d.By.Operator, Round: d.By.Round, Belief: d.Belief,
		})
		r.Totals.Declined++
	}
	for _, t := range g.Truncations() {
		r.Truncations = append(r.Truncations, Truncated{
			Reason: t.Reason.String(), Round: t.Round, Detail: t.Detail,
		})
	}
	// Truncation of either kind means the graph is a prefix, whatever the
	// caller believed. Deriving it here rather than trusting the flag is what
	// keeps the marker honest when a budget bound expansion silently — the
	// caller never sees that happen, and §8 is explicit that a truncated graph
	// which reads as complete is a correctness bug.
	//
	// The ledger counts, not just run-level limits: a declined candidate is a
	// node a fuller scan would have had. That includes belief-gated declines,
	// which are working as configured and still leave the report incomplete.
	if !r.Partial {
		switch {
		case len(r.Truncations) > 0:
			r.Partial, r.PartialWhy = true, r.Truncations[0].Reason
		case len(r.Ledger) > 0:
			r.Partial = true
			r.PartialWhy = fmt.Sprintf("%d candidates declined (%s)",
				len(r.Ledger), r.Ledger[0].Reason)
		}
	}
	return r
}

func buildNode(g *graph.Graph, a *graph.Analysis, n *graph.Node, verbose bool) NodeRow {
	ex := a.Existence(n.ID)
	row := NodeRow{
		ID:        n.ID.String(),
		Type:      n.Type.Name(),
		Key:       n.Key,
		Depth:     a.Depth(n.ID),
		Existence: ex.String(),
		InClosure: a.InClosure(n.ID),
		Seed:      n.ID == g.SeedID(),
		severity:  g.Risk(n.ID),
		existence: ex,
	}
	if row.severity != 0 {
		row.Risk = row.severity.String()
	}
	n.Props.Each(func(f graph.Field, v graph.Value) {
		row.Props = append(row.Props, prop(f.Name(), v))
	})

	for _, s := range a.Statuses(n.ID) {
		row.Statuses = append(row.Statuses, Status{Operator: s.Operator, Status: s.Status.String()})
	}
	for _, s := range a.Scores(n.ID) {
		row.Scores = append(row.Scores, Score{Key: s.Key, Value: s.Value})
	}
	if verbose {
		for _, c := range a.Assertions(n.ID) {
			row.Claims = append(row.Claims, Claim{
				Field: c.Field, Value: value(c.Value), Operator: c.By.Operator,
				Round: c.By.Round, Won: c.Won,
			})
		}
		b := g.Belief(n.ID)
		row.Belief = &b
	}
	return row
}

func buildEdge(g *graph.Graph, e *graph.Edge) EdgeRow {
	row := EdgeRow{Rel: e.Rel.Name(), Class: e.Rel.Class().String()}
	if n, ok := g.Node(e.From); ok {
		row.From = ref(n)
	}
	if n, ok := g.Node(e.To); ok {
		row.To = ref(n)
	}
	e.Props.Each(func(f graph.Field, v graph.Value) {
		row.Props = append(row.Props, prop(f.Name(), v))
	})
	return row
}

func buildFinding(g *graph.Graph, f graph.Finding) Finding {
	out := Finding{Kind: f.Kind, Severity: f.Severity.String(), Summary: f.Summary}
	for _, id := range f.Nodes {
		if n, ok := g.Node(id); ok {
			out.Nodes = append(out.Nodes, ref(n))
		}
	}
	for _, d := range f.Declined {
		out.Declined = append(out.Declined, d.Type+":"+d.Key)
	}
	for _, p := range f.Evidence {
		out.Evidence = append(out.Evidence, fmt.Sprintf("%s@%d", p.Operator, p.Round))
	}
	// A finding's node list is an unordered set; rendering it in whatever order
	// the analyzer happened to append would make the report non-reproducible.
	sort.Strings(out.Nodes)
	sort.Strings(out.Declined)
	sort.Strings(out.Evidence)
	return out
}

// ref is the stable human-facing name of a node: type-qualified, because two
// types can share a key and an unqualified `acme` in a report is ambiguous.
func ref(n *graph.Node) string { return n.Type.Name() + ":" + n.Key }

func prop(name string, v graph.Value) Prop {
	return Prop{Name: name, Value: value(v), Kind: v.Kind().String()}
}

// value converts a graph value to a Go native, so the JSON encoder emits a
// number as a number rather than a quoted string.
func value(v graph.Value) any {
	switch v.Kind() {
	case graph.KindString:
		return v.Str()
	case graph.KindInt:
		return v.Num()
	case graph.KindFloat:
		return v.Real()
	case graph.KindBool:
		return v.Flag()
	case graph.KindTime:
		return v.Time().Format(time.RFC3339)
	case graph.KindBytes:
		return fmt.Sprintf("%x", v.Raw())
	}
	return nil
}

func text(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	}
	return fmt.Sprint(v)
}

// Max returns the highest severity in the report. It is what --fail-on tests.
func (r Report) Max() graph.Severity {
	var max graph.Severity
	for _, f := range r.Findings {
		if s, ok := graph.ParseSeverity(f.Severity); ok && s > max {
			max = s
		}
	}
	return max
}

// Formats lists every renderable format, for --help and validation.
func Formats() []string { return []string{"table", "json", "ndjson", "csv", "dot"} }

// FormatFor resolves a format from a filename extension, for --save.
func FormatFor(path string) (string, bool) {
	i := strings.LastIndex(path, ".")
	if i < 0 {
		return "", false
	}
	switch ext := strings.ToLower(path[i+1:]); ext {
	case "json", "ndjson", "csv", "dot":
		return ext, true
	case "txt", "text":
		return "table", true
	case "gv":
		return "dot", true
	}
	return "", false
}

// Stranded is a finding whose node the filter removed from the report.
type Stranded struct {
	// Node is the "type:key" reference, as Finding.Nodes carries it.
	Node string
	// Finding is the finding that node carried.
	Finding Finding
}

// StrandedFindings returns the findings whose nodes are not in the report's
// node list, so a renderer can surface them separately.
//
// Findings are never filtered — hiding one behind --filter would hide exactly
// what --fail-on gates on — but every format draws them as an attribute of a
// node, so a finding whose node was filtered out has nothing to attach to and
// vanishes unless the renderer says so itself.
//
// This exists because two renderers worked that out independently and one got
// it wrong. Finding.Nodes holds "type:key" references; csv built its lookup set
// with those and dot built its with NodeRow.ID, a hex node id, so dot's test
// never matched and EVERY finding read as stranded — an unfiltered `-o dot`
// drew a "findings on filtered nodes" cluster listing findings that were right
// there in the graph. A renderer calls this rather than deriving the set again;
// the key is built once, here.
func StrandedFindings(r Report) []Stranded {
	shown := make(map[string]bool, len(r.Nodes))
	for _, n := range r.Nodes {
		shown[n.Type+":"+n.Key] = true
	}
	var out []Stranded
	for _, f := range r.Findings {
		for _, n := range f.Nodes {
			if !shown[n] {
				out = append(out, Stranded{Node: n, Finding: f})
			}
		}
	}
	return out
}
