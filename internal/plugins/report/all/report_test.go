// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	analyze "github.com/rangertaha/urlinsane/internal/plugins/analyze/all"
	"github.com/rangertaha/urlinsane/internal/plugins/report"
)

func registry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()
	lower := func(s string) (string, error) {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			return "", fmt.Errorf("empty key")
		}
		return s, nil
	}
	for _, d := range []graph.NodeTypeDef{
		{Name: "domain", Cap: graph.Nameable, Version: 1, Canonical: lower, Fields: []graph.FieldDef{
			{Name: "punycode", Kind: graph.KindString},
			{Name: "length", Kind: graph.KindInt},
		}},
		{Name: "package", Cap: graph.Nameable, Version: 1, Canonical: lower},
		{Name: "ip", Cap: graph.Observed, Version: 1, Canonical: lower},
	} {
		if _, err := r.AddType(d); err != nil {
			t.Fatalf("%s: %v", d.Name, err)
		}
	}
	for _, d := range []graph.RelDef{
		{Name: graph.VariantRel, Class: graph.Variant, Version: 1, Fields: []graph.FieldDef{
			{Name: "algorithm", Kind: graph.KindString},
			{Name: "distance", Kind: graph.KindInt},
		}},
		{Name: "RESOLVES_TO", Class: graph.Observation, Version: 1},
		{Name: "MX", Class: graph.Observation, Version: 1},
	} {
		if _, err := r.AddRel(d); err != nil {
			t.Fatalf("%s: %v", d.Name, err)
		}
	}
	return r
}

func by(op string) graph.Provenance { return graph.Provenance{Operator: op, Round: 1} }

// scan builds the shape a real run produces: a seed, variants, a shared IP, a
// mix of existence states, a declined candidate and a run-level truncation.
func scan(t *testing.T) *graph.Graph {
	t.Helper()
	g := graph.New(registry(t))
	seed, err := g.Seed("domain", "example.com")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	origin := graph.NodeRef{Type: "domain", Key: "example.com"}
	g.Apply(by("idn"), seed, graph.Delta{Props: []graph.PropSet{
		{Node: &origin, Field: "length", Value: graph.Int(11)},
	}})
	g.SetStatus(seed, "dns", graph.StatusOK)

	for _, k := range []string{"examp1e.com", "exarnple.com"} {
		v := graph.NodeRef{Type: "domain", Key: k}
		e := graph.EdgeRef{From: origin, Rel: graph.VariantRel, To: v}
		g.Apply(by("omission"), seed, graph.Delta{
			Edges: []graph.EdgeRef{e},
			Props: []graph.PropSet{
				{Edge: &e, Field: "distance", Value: graph.Int(1)},
				{Edge: &e, Field: "algorithm", Value: graph.String("omission")},
			},
		})
		id := node(t, g, "domain", k)
		g.Apply(by("dns"), id, graph.Delta{Edges: []graph.EdgeRef{{
			From: v, Rel: "RESOLVES_TO", To: graph.NodeRef{Type: "ip", Key: "1.2.3.4"},
		}}})
		g.SetStatus(id, "dns", graph.StatusOK)
	}

	// One confirmed-absent and one never-determined, so the three existence
	// values are all exercised.
	for k, s := range map[string]graph.Status{
		"exampie.com": graph.StatusEmpty,
		"exampl3.com": graph.StatusTimeout,
	} {
		v := graph.NodeRef{Type: "domain", Key: k}
		g.Apply(by("omission"), seed, graph.Delta{Edges: []graph.EdgeRef{
			{From: origin, Rel: graph.VariantRel, To: v},
		}})
		g.SetStatus(node(t, g, "domain", k), "dns", s)
	}

	// A live node of another type, so filters that combine existence with type
	// have something to exclude and the renderers see more than one type.
	pkg := graph.NodeRef{Type: "package", Key: "npm:example"}
	g.Apply(by("manifest"), seed, graph.Delta{Nodes: []graph.NodeRef{pkg}})
	g.SetStatus(node(t, g, "package", "npm:example"), "npm", graph.StatusOK)

	if err := g.Decline("domain", "3xample.com", 0, 0.02, graph.ReasonBelief, by("omission")); err != nil {
		t.Fatalf("decline: %v", err)
	}
	g.NoteTruncation(graph.ReasonRoundCap, 4, "round cap reached")

	if err := g.RunAnalyzers(t.Context(), analyze.All()); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	return g
}

func node(t *testing.T, g *graph.Graph, typ, key string) graph.NodeID {
	t.Helper()
	for _, n := range g.Nodes() {
		if n.Type.Name() == typ && n.Key == key {
			return n.ID
		}
	}
	t.Fatalf("node %s:%s not found", typ, key)
	return graph.NodeID{}
}

func render(t *testing.T, g *graph.Graph, o report.Options) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Render(&buf, g, o); err != nil {
		t.Fatalf("render %s: %v", o.Format, err)
	}
	return buf.String()
}

// --- determinism -------------------------------------------------------------

func TestReportFormatsAreByteIdenticalAcrossRuns(t *testing.T) {
	// The claim §11 makes: "what changed since last week" is a diff rather than
	// noise. That only holds if two identical scans render identically, which
	// map iteration alone is enough to break.
	for _, format := range []string{"table", "json", "csv", "dot"} {
		t.Run(format, func(t *testing.T) {
			first := render(t, scan(t), report.Options{Format: format, Target: "example.com"})
			for i := 0; i < 8; i++ {
				if got := render(t, scan(t), report.Options{Format: format, Target: "example.com"}); got != first {
					t.Fatalf("run %d differed:\n%s\n---\n%s", i, got, first)
				}
			}
		})
	}
}

func TestNDJSONIsStableToo(t *testing.T) {
	// ndjson makes no *ordering* guarantee across a live run, but rendering a
	// finished graph twice must still agree — otherwise the format is simply
	// nondeterministic, which is a different and worse thing.
	first := render(t, scan(t), report.Options{Format: "ndjson"})
	for i := 0; i < 5; i++ {
		if got := render(t, scan(t), report.Options{Format: "ndjson"}); got != first {
			t.Fatal("ndjson rendering varied between runs")
		}
	}
}

func TestElapsedIsNotInCanonicalOutput(t *testing.T) {
	// A duration differs on every run. Including it would defeat byte-comparison
	// exactly when a user is diffing two scans.
	g := scan(t)
	a := render(t, g, report.Options{Format: "json", Elapsed: 1})
	b := render(t, g, report.Options{Format: "json", Elapsed: 999999999})
	if a != b {
		t.Fatal("elapsed time leaked into the canonical JSON report")
	}
}

// --- the truncation ledger ---------------------------------------------------

func TestEveryFormatCarriesTheLedger(t *testing.T) {
	// §11: a truncated graph that reads as complete is a correctness bug. No
	// format gets to drop the ledger to keep its shape tidy.
	g := scan(t)
	for _, format := range []string{"table", "json", "csv", "dot", "ndjson"} {
		out := render(t, g, report.Options{Format: format})
		if !strings.Contains(out, "belief") && !strings.Contains(out, "declined") {
			t.Fatalf("%s output does not mention the declined candidate:\n%s", format, out)
		}
	}
}

func TestCSVCarriesDeclinedRows(t *testing.T) {
	out := render(t, scan(t), report.Options{Format: "csv"})
	rows, err := csv.NewReader(strings.NewReader(out)).ReadAll()
	if err != nil {
		t.Fatalf("csv: %v", err)
	}
	var declined int
	for _, r := range rows[1:] {
		if r[len(r)-1] != "" {
			declined++
			if r[1] != "3xample.com" {
				t.Fatalf("declined row key = %q", r[1])
			}
		}
	}
	if declined != 1 {
		t.Fatalf("declined rows = %d, want 1", declined)
	}
}

func TestTruncationForcesThePartialMarker(t *testing.T) {
	// The caller did not pass Partial. A run-level truncation means the graph is
	// a prefix regardless of what the caller believed, and deriving it here is
	// what keeps the marker honest when a budget bound expansion silently.
	r := report.Build(scan(t), report.Options{})
	if !r.Partial {
		t.Fatal("a graph with a run-level truncation was not marked partial")
	}
	if r.PartialWhy == "" {
		t.Fatal("partial report does not say why")
	}
	out := render(t, scan(t), report.Options{Format: "table"})
	if !strings.Contains(out, "PARTIAL") {
		t.Fatalf("table did not surface the partial marker:\n%s", out)
	}
}

func TestCompleteScanIsNotMarkedPartial(t *testing.T) {
	g := graph.New(registry(t))
	if _, err := g.Seed("domain", "example.com"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if report.Build(g, report.Options{}).Partial {
		t.Fatal("an untruncated scan was marked partial")
	}
}

// --- filters -----------------------------------------------------------------

func TestFilterByExistence(t *testing.T) {
	g := scan(t)
	for _, tc := range []struct {
		filter string
		want   []string
	}{
		{"live", []string{"example.com", "examp1e.com", "exarnple.com", "npm:example"}},
		{"absent", []string{"exampie.com"}},
		{"unknown", []string{"exampl3.com"}},
		// The IP is reached by resolution but never observed itself, so it is
		// untried rather than live. Inheriting liveness from an in-edge would
		// claim an observation nobody made.
		{"untried", []string{"1.2.3.4"}},
	} {
		t.Run(tc.filter, func(t *testing.T) {
			f, err := report.ParseFilters([]string{tc.filter})
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			r := report.Build(g, report.Options{Filters: f})
			got := map[string]bool{}
			for _, n := range r.Nodes {
				got[n.Key] = true
			}
			if len(got) != len(tc.want) {
				t.Fatalf("%s selected %d nodes %v, want %d %v",
					tc.filter, len(got), keys(got), len(tc.want), tc.want)
			}
			for _, w := range tc.want {
				if !got[w] {
					t.Fatalf("%s did not select %s (got %v)", tc.filter, w, keys(got))
				}
			}
		})
	}
}

func TestExistenceFiltersAreAlternatives(t *testing.T) {
	// A node cannot be both live and absent. Treating every filter as a
	// conjunction would make the obvious two-value query return nothing at all.
	f, err := report.ParseFilters([]string{"live,absent"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r := report.Build(scan(t), report.Options{Filters: f})
	var live, absent bool
	for _, n := range r.Nodes {
		switch n.Existence {
		case "live":
			live = true
		case "absent":
			absent = true
		}
	}
	if !live || !absent {
		t.Fatalf("live,absent returned live=%v absent=%v; both must appear", live, absent)
	}
}

func TestFiltersAcrossFamiliesNarrow(t *testing.T) {
	f, err := report.ParseFilters([]string{"live", "type=domain"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r := report.Build(scan(t), report.Options{Filters: f})
	if len(r.Nodes) == 0 {
		t.Fatal("live+type=domain selected nothing; the assertions below would pass vacuously")
	}
	for _, n := range r.Nodes {
		if n.Type != "domain" || n.Existence != "live" {
			t.Fatalf("node %s:%s (%s) passed a live+type=domain filter", n.Type, n.Key, n.Existence)
		}
	}
	// And the conjunction must actually exclude: the fixture has live non-domains
	// and non-live domains, so a filter that dropped either half would still
	// satisfy the loop above.
	if len(r.Nodes) >= r.Totals.Live {
		t.Fatalf("shown %d >= %d live; the type filter did not narrow", len(r.Nodes), r.Totals.Live)
	}
}

func TestFilterByRisk(t *testing.T) {
	f, err := report.ParseFilters([]string{"risk>=high"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r := report.Build(scan(t), report.Options{Filters: f})
	if len(r.Nodes) == 0 {
		t.Fatal("risk>=high selected nothing; the fixture has a campaign finding")
	}
	for _, n := range r.Nodes {
		if s, _ := graph.ParseSeverity(n.Risk); s < graph.SeverityHigh {
			t.Fatalf("%s has risk %q, below the filter", n.Key, n.Risk)
		}
	}
}

func TestFilterDropsDanglingEdges(t *testing.T) {
	// An edge to a filtered-out node would make dot draw a node the table never
	// listed and json reference an id it never defined.
	f, _ := report.ParseFilters([]string{"type=domain"})
	r := report.Build(scan(t), report.Options{Filters: f})
	shown := map[string]bool{}
	for _, n := range r.Nodes {
		shown[n.Type+":"+n.Key] = true
	}
	for _, e := range r.Edges {
		if !shown[e.From] || !shown[e.To] {
			t.Fatalf("edge %s -> %s dangles past the filter", e.From, e.To)
		}
	}
}

func TestFiltersNeverHideFindingsOrLedger(t *testing.T) {
	// Letting --filter hide a finding would hide exactly what --fail-on gates
	// on, and hiding the ledger would report a truncated scan as complete.
	full := report.Build(scan(t), report.Options{})
	f, _ := report.ParseFilters([]string{"type=ip"})
	narrowed := report.Build(scan(t), report.Options{Filters: f})

	if len(narrowed.Findings) != len(full.Findings) {
		t.Fatalf("filter changed finding count %d -> %d",
			len(full.Findings), len(narrowed.Findings))
	}
	if len(narrowed.Ledger) != len(full.Ledger) {
		t.Fatalf("filter changed ledger size %d -> %d", len(full.Ledger), len(narrowed.Ledger))
	}
}

func TestTotalsCountTheUnfilteredGraph(t *testing.T) {
	// report.Totals are how a reader judges coverage. Reporting filtered totals would
	// misstate the scan's size precisely when a filter is narrowing it.
	f, _ := report.ParseFilters([]string{"type=ip"})
	r := report.Build(scan(t), report.Options{Filters: f})
	if r.Totals.Nodes <= r.Totals.Shown {
		t.Fatalf("totals %d nodes with %d shown; totals must count the whole graph",
			r.Totals.Nodes, r.Totals.Shown)
	}
	if r.Totals.Live == 0 || r.Totals.Absent == 0 || r.Totals.Unknown == 0 {
		t.Fatalf("existence totals %+v lost values the filter excluded", r.Totals)
	}
}

func TestBadFilterIsAnError(t *testing.T) {
	// A silently ignored filter is worse than a rejected one: the user believes
	// they narrowed the report and reads a full one as filtered.
	for _, bad := range []string{"registered", "risk>enormous", "depth<=abc", "type>domain", ""} {
		if bad == "" {
			continue
		}
		if _, err := report.ParseFilter(bad); err == nil {
			t.Fatalf("%q was accepted as a filter", bad)
		}
	}
}

// --- format mechanics --------------------------------------------------------

func TestJSONIsValidAndComplete(t *testing.T) {
	var r report.Report
	if err := json.Unmarshal([]byte(render(t, scan(t), report.Options{Format: "json"})), &r); err != nil {
		t.Fatalf("json: %v", err)
	}
	if len(r.Nodes) == 0 || len(r.Findings) == 0 || len(r.Ledger) == 0 {
		t.Fatalf("round-tripped report is missing sections: %+v", r.Totals)
	}
	// Numbers must be numbers, not quoted strings, or every consumer has to
	// parse them back out.
	for _, n := range r.Nodes {
		for _, p := range n.Props {
			if p.Kind == "int" {
				if _, ok := p.Value.(float64); !ok {
					t.Fatalf("int prop %s encoded as %T", p.Name, p.Value)
				}
			}
		}
	}
}

func TestNDJSONIsOneTaggedObjectPerLine(t *testing.T) {
	out := render(t, scan(t), report.Options{Format: "ndjson"})
	kinds := map[string]int{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("line is not JSON: %q", line)
		}
		kind, ok := m["kind"].(string)
		if !ok {
			t.Fatalf("line has no kind discriminator: %q", line)
		}
		kinds[kind]++
	}
	for _, want := range []string{"run", "node", "edge", "finding", "declined", "truncation", "totals"} {
		if kinds[want] == 0 {
			t.Fatalf("ndjson emitted no %s events: %v", want, kinds)
		}
	}
}

// A Finding carries its own "kind" -- which analyzer concluded what -- and the
// record discriminator was assigned straight over it, so every finding line
// said only that it was a finding. The analyzer that raised it, which is the
// whole content of the field, was gone.
func TestNDJSONKeepsAFindingsOwnKind(t *testing.T) {
	out := render(t, scan(t), report.Options{Format: "ndjson"})

	var seen bool
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("line is not JSON: %q", line)
		}
		if m["kind"] != "finding" {
			continue
		}
		seen = true
		own, ok := m["finding_kind"].(string)
		if !ok || own == "" {
			t.Fatalf("finding line lost the analyzer discriminator: %q", line)
		}
	}
	if !seen {
		t.Fatal("no finding lines in the fixture")
	}
}

func TestDOTQuotesHostileKeys(t *testing.T) {
	// Variant algorithms generate exactly the strings that break naive quoting;
	// an unescaped quote produces a file graphviz refuses to parse.
	g := graph.New(registry(t))
	seed, err := g.Seed("domain", `ex"ample\.com`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	_ = seed
	out := render(t, g, report.Options{Format: "dot"})
	if strings.Contains(out, `"ex"ample`) {
		t.Fatalf("unescaped quote in dot output:\n%s", out)
	}
	if !strings.Contains(out, `ex\"ample\\.com`) {
		t.Fatalf("dot did not escape the key:\n%s", out)
	}
}

func TestDOTShowsConvergence(t *testing.T) {
	// The reason the data model puts infrastructure in nodes: two variants
	// pointing at one IP is visible as convergence, and no row-per-domain
	// rendering can show it.
	out := render(t, scan(t), report.Options{Format: "dot"})
	var arrows int
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, `-> "ip:1.2.3.4"`) {
			arrows++
		}
	}
	if arrows < 2 {
		t.Fatalf("dot shows %d edges into the shared IP, want 2:\n%s", arrows, out)
	}
}

func TestTableRespectsColorOption(t *testing.T) {
	g := scan(t)
	plain := render(t, g, report.Options{Format: "table", Target: "example.com"})
	if strings.Contains(plain, "\033[") {
		t.Fatal("ANSI codes in an uncolored table; NO_COLOR and pipes would break")
	}
	if !strings.Contains(render(t, g, report.Options{Format: "table", Color: true}), "\033[") {
		t.Fatal("color was requested but nothing was styled")
	}
}

func TestVerboseAddsBeliefAndClaims(t *testing.T) {
	// Belief is diagnostic only. It may appear under --verbose and must never
	// reach a finding (§9).
	quiet := report.Build(scan(t), report.Options{})
	loud := report.Build(scan(t), report.Options{Verbose: true})
	for _, n := range quiet.Nodes {
		if n.Belief != nil {
			t.Fatalf("%s exposed engine belief without --verbose", n.Key)
		}
	}
	var sawBelief bool
	for _, n := range loud.Nodes {
		if n.Belief != nil {
			sawBelief = true
		}
	}
	if !sawBelief {
		t.Fatal("--verbose did not surface belief")
	}
}

func TestUnknownFormatIsRejected(t *testing.T) {
	var buf bytes.Buffer
	err := Render(&buf, scan(t), report.Options{Format: "yaml"})
	if err == nil {
		t.Fatal("an unknown format was accepted")
	}
	if !strings.Contains(err.Error(), "table") {
		t.Fatalf("error does not list the valid formats: %v", err)
	}
}

// --- CLI-facing helpers ------------------------------------------------------

func TestMaxDrivesFailOn(t *testing.T) {
	r := report.Build(scan(t), report.Options{})
	if r.Max() < graph.SeverityHigh {
		t.Fatalf("Max() = %v; the fixture has a campaign finding at high", r.Max())
	}
	empty := report.Build(graph.New(registry(t)), report.Options{})
	if empty.Max() != 0 {
		t.Fatalf("Max() = %v on a scan with no findings, want none", empty.Max())
	}
}

func TestFormatForExtension(t *testing.T) {
	// --save is a separate sink from -o: the format follows the extension, so
	// `--save out.json` writes a report whatever stdout is doing.
	for path, want := range map[string]string{
		"out.json": "json", "out.ndjson": "ndjson", "out.csv": "csv",
		"out.dot": "dot", "out.gv": "dot", "out.txt": "table", "OUT.JSON": "json",
	} {
		got, ok := report.FormatFor(path)
		if !ok || got != want {
			t.Fatalf("FormatFor(%q) = %q,%v; want %q", path, got, ok, want)
		}
	}
	for _, path := range []string{"out", "out.xlsx", ""} {
		if _, ok := report.FormatFor(path); ok {
			t.Fatalf("FormatFor(%q) claimed a format", path)
		}
	}
}

func TestWriteRendersTwoSinksFromOneBuild(t *testing.T) {
	// A run writes stdout and --save from the same report.Report, so the two sinks
	// cannot disagree about what the scan found.
	r := report.Build(scan(t), report.Options{Target: "example.com"})
	var table, js bytes.Buffer
	if err := Write(&table, r, report.Options{Format: "table"}); err != nil {
		t.Fatalf("table: %v", err)
	}
	if err := Write(&js, r, report.Options{Format: "json"}); err != nil {
		t.Fatalf("json: %v", err)
	}
	var decoded report.Report
	if err := json.Unmarshal(js.Bytes(), &decoded); err != nil {
		t.Fatalf("json: %v", err)
	}
	for _, n := range decoded.Nodes {
		if !strings.Contains(table.String(), n.Key) {
			t.Fatalf("%s is in the json sink but not the table sink", n.Key)
		}
	}
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestDOTLabelsUseRealNewlines(t *testing.T) {
	// graphviz's line break is \n inside the quoted string. Writing the escape
	// into the label in Go source produces a literal backslash-n, which quote
	// then escapes again and graphviz renders verbatim on the node.
	out := render(t, scan(t), report.Options{Format: "dot"})
	if strings.Contains(out, `\\n`) {
		t.Fatalf("dot labels contain a double-escaped newline:\n%s", out)
	}
	if !strings.Contains(out, `example.com\ndomain`) {
		t.Fatalf("dot node label is not type-annotated:\n%s", out)
	}
}

func TestOnlyTheSeedIsMarkedAsSeed(t *testing.T) {
	// Depth cannot stand in for this: variant edges cost no depth, so every
	// variant of the target sits at depth 0 right beside it.
	r := report.Build(scan(t), report.Options{})
	var seeds []string
	var zeroDepth int
	for _, n := range r.Nodes {
		if n.Seed {
			seeds = append(seeds, n.Key)
		}
		if n.Depth == 0 {
			zeroDepth++
		}
	}
	if len(seeds) != 1 || seeds[0] != "example.com" {
		t.Fatalf("seeds = %v, want exactly [example.com]", seeds)
	}
	if zeroDepth < 2 {
		t.Fatal("fixture no longer has variants at depth 0; this test proves nothing")
	}
}

func TestTableFlattensMultilineValues(t *testing.T) {
	// A domain's TXT set is several records joined into one prop, because props
	// are single-valued and the records have no entity to hang off. tabwriter
	// measures column widths per line, so one embedded newline destroys the
	// alignment of every row below it.
	g := graph.New(registry(t))
	seed, err := g.Seed("domain", "example.com")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	ref := graph.NodeRef{Type: "domain", Key: "example.com"}
	g.Apply(by("dns"), seed, graph.Delta{Props: []graph.PropSet{
		{Node: &ref, Field: "punycode", Value: graph.String("v=spf1 -all\ngoogle-site-verification=abc")},
	}})

	out := render(t, g, report.Options{Format: "table"})
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if strings.HasPrefix(line, "google-site-verification") {
			t.Fatalf("a multi-line prop value spilled onto its own table row:\n%s", out)
		}
	}

	// The other formats keep the value intact — only the table flattens.
	var r report.Report
	if err := json.Unmarshal([]byte(render(t, g, report.Options{Format: "json"})), &r); err != nil {
		t.Fatalf("json: %v", err)
	}
	for _, n := range r.Nodes {
		for _, p := range n.Props {
			if p.Name == "punycode" && !strings.Contains(p.Value.(string), "\n") {
				t.Fatal("json lost the newline the table flattened")
			}
		}
	}
}

// --- what every format owes the reader ---------------------------------------

// The class test for the five renderers.
//
// The package's claim is that all five project one intermediate "so they cannot
// disagree". They can, and did, because each one decides for itself which parts
// of that intermediate to surface — so the disagreements show up one renderer at
// a time and only under a filter.
//
// Two run-level facts have to survive into every format whatever the filter
// says. `report.Build` never filters findings, because hiding one behind
// `--filter live` would hide exactly what `--fail-on` gates on; and a scan that
// stopped early must never export as a complete one. CSV surfaced findings only
// as a column on node rows, so filtering the nodes away deleted the findings
// with them, and CSV had no partial marker at all — an analyst exporting a
// truncated scan to a spreadsheet got a file indistinguishable from a whole one.
func TestEveryFormatSurvivesFiltering(t *testing.T) {
	g := scan(t)

	// A filter that keeps almost nothing: the findings in this fixture are on
	// live variant nodes, so selecting only absent ones strands them.
	absent, err := report.ParseFilters([]string{"absent"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	opts := func(format string) report.Options {
		return report.Options{
			Format:  format,
			Target:  "example.com",
			Filters: absent,
			Partial: true, PartialWhy: "interrupted",
		}
	}

	built := report.Build(g, report.Options{Filters: absent, Partial: true})
	if len(built.Findings) == 0 {
		t.Fatal("fixture is wrong: no findings survive Build, so this proves nothing")
	}
	var stranded int
	shown := map[string]bool{}
	for _, n := range built.Nodes {
		shown[n.Type+":"+n.Key] = true
	}
	for _, f := range built.Findings {
		for _, n := range f.Nodes {
			if !shown[n] {
				stranded++
			}
		}
	}
	if stranded == 0 {
		t.Fatal("fixture is wrong: every finding still has a visible node row")
	}

	for _, format := range report.Formats() {
		t.Run(format, func(t *testing.T) {
			out := render(t, g, opts(format))

			// The partial marker, however the format spells it.
			if !strings.Contains(strings.ToLower(out), "partial") {
				t.Errorf("%s does not say the scan was partial:\n%s", format, out)
			}

			// At least one finding kind, even though its nodes were filtered out.
			var carried bool
			for _, f := range built.Findings {
				if strings.Contains(out, f.Kind) {
					carried = true
					break
				}
			}
			if !carried {
				t.Errorf("%s dropped every finding when the nodes were filtered away; "+
					"findings are deliberately never filtered:\n%s", format, out)
			}
		})
	}
}

// An unfiltered report has no stranded findings, so no format may claim it has.
//
// dot built its "is this node shown" set from NodeRow.ID (a hex node id) but
// looked it up with Finding.Nodes entries ("type:key"), so the test never
// matched and every finding read as stranded: a plain `-o dot` drew a
// "findings on filtered nodes" cluster listing findings that were right there
// in the graph. csv had the same code with the right key, which is how the two
// drifted apart.
func TestNoFormatClaimsStrandedFindingsWhenNothingIsFiltered(t *testing.T) {
	g := scan(t)

	// The helper both renderers now share must agree there is nothing stranded.
	if s := report.StrandedFindings(report.Build(g, report.Options{})); len(s) != 0 {
		t.Fatalf("StrandedFindings = %d on an unfiltered report, want 0: %+v", len(s), s)
	}

	for _, format := range []string{"dot", "csv"} {
		out := render(t, g, report.Options{Format: format})
		if strings.Contains(out, "filtered") {
			t.Errorf("%s output mentions filtered findings with no filter applied:\n%s", format, out)
		}
	}
}

// And when a filter really does strand one, it must still be surfaced — the
// fix must not have bought silence.
func TestStrandedFindingsAreStillReportedWhenAFilterHidesTheNode(t *testing.T) {
	g := scan(t)
	// The campaign finding sits on live variants; keep only absent nodes.
	o := report.Options{Format: "dot", Filters: mustParse(t, "absent")}

	if s := report.StrandedFindings(report.Build(g, o)); len(s) == 0 {
		t.Fatal("filtering the findings' nodes away stranded nothing")
	}
	if out := render(t, g, o); !strings.Contains(out, "filtered") {
		t.Errorf("dot did not surface a stranded finding:\n%s", out)
	}
}

func mustParse(t *testing.T, specs ...string) []report.Filter {
	t.Helper()
	f, err := report.ParseFilters(specs)
	if err != nil {
		t.Fatal(err)
	}
	return f
}
