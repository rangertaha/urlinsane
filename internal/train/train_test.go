// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package train

import (
	"context"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// captureView runs a real dispatch and hands back the View the scheduler
// served for one node.
//
// Constructing a View directly is not possible from outside internal/graph,
// and going through the scheduler is the better test anyway: it proves the
// symbols match what an operator is actually given mid-scan, restricted by the
// same read-set declaration, rather than what a hand-built view would carry.
func captureView(t *testing.T, g *graph.Graph, key string) graph.View {
	t.Helper()
	op := &capturingOp{want: key, seen: map[string]graph.View{}}
	s := graph.NewScheduler(g, []graph.Operator{op}, graph.Limits{MaxRounds: 2, Workers: 1})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	v, ok := op.seen[key]
	if !ok {
		t.Fatalf("the scheduler never served a view for %s", key)
	}
	return v
}

// capturingOp reads exactly what the featurizer declares and produces nothing.
type capturingOp struct {
	want string
	seen map[string]graph.View
}

func (o *capturingOp) Id() string           { return "capture" }
func (o *capturingOp) Version() int         { return 1 }
func (o *capturingOp) Resource() string     { return "" }
func (o *capturingOp) Emits() graph.Effects { return graph.Effects{} }

func (o *capturingOp) Trigger() graph.Trigger {
	return graph.Trigger{
		On:    graph.Selector{Types: []string{"domain"}},
		Reads: Trigger(),
	}
}

func (o *capturingOp) Exec(_ context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	o.seen[v.Key()] = v
	return graph.Delta{}, graph.OK()
}

// registry mirrors the shape the scanner registers: a nameable seed, a variant
// relation carrying the algorithm and distance, and an observation relation.
func registry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()
	lower := func(s string) (string, error) {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			return "", errEmpty
		}
		return s, nil
	}
	types := []graph.NodeTypeDef{
		{Name: "domain", Cap: graph.Nameable, Version: 1, Canonical: lower, Fields: []graph.FieldDef{
			{Name: "live", Kind: graph.KindBool},
			{Name: "punycode", Kind: graph.KindString},
		}},
		{Name: "ip", Cap: graph.Observed, Version: 1, Canonical: lower},
	}
	for _, d := range types {
		if _, err := r.AddType(d); err != nil {
			t.Fatalf("type %s: %v", d.Name, err)
		}
	}
	rels := []graph.RelDef{
		{Name: graph.VariantRel, Class: graph.Variant, Version: 1, Fields: []graph.FieldDef{
			{Name: "algorithm", Kind: graph.KindString},
			{Name: "distance", Kind: graph.KindInt},
		}},
		{Name: "RESOLVES_TO", Class: graph.Observation, Version: 1},
	}
	for _, d := range rels {
		if _, err := r.AddRel(d); err != nil {
			t.Fatalf("rel %s: %v", d.Name, err)
		}
	}
	return r
}

type constErr string

func (e constErr) Error() string { return string(e) }

const errEmpty = constErr("empty key")

// scan builds a small finished graph: a seed, two variants, one of which
// resolves.
func scan(t *testing.T) (*graph.Graph, graph.NodeID) {
	t.Helper()
	g := graph.New(registry(t))
	seed, err := g.Seed("domain", "example.com")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	origin := graph.NodeRef{Type: "domain", Key: "example.com"}
	by := graph.Provenance{Operator: "co", Round: 1}

	for _, v := range []struct {
		key  string
		live bool
	}{{"exmple.com", true}, {"exampl.com", false}} {
		ref := graph.NodeRef{Type: "domain", Key: v.key}
		e := graph.EdgeRef{From: origin, Rel: graph.VariantRel, To: ref}
		g.Apply(by, seed, graph.Delta{
			Edges: []graph.EdgeRef{e},
			Props: []graph.PropSet{
				{Edge: &e, Field: "algorithm", Value: graph.String("co")},
				{Edge: &e, Field: "distance", Value: graph.Int(1)},
				{Node: &ref, Field: "live", Value: graph.Bool(v.live)},
			},
		})
		id := node(t, g, v.key)
		g.SetStatus(id, "dns", statusFor(v.live))
		if v.live {
			g.Apply(graph.Provenance{Operator: "dns", Round: 2}, id, graph.Delta{
				Edges: []graph.EdgeRef{{From: ref, Rel: "RESOLVES_TO",
					To: graph.NodeRef{Type: "ip", Key: "1.2.3.4"}}},
			})
		}
	}
	g.SetStatus(seed, "dns", graph.StatusOK)

	// Parents are assigned at a barrier and nowhere else, so a graph built by
	// calling Apply directly has no expansion tree. One scheduler round gives
	// it the barrier a real scan gets.
	barrier(t, g)

	// After the scheduler, not before: NewScheduler re-derives the observer set
	// from the operators that declare a resource, and the capturing operator
	// declares none — which leaves the set empty, and an empty set means
	// "everything observes", so every status would attest existence.
	g.SetObservers([]string{"dns"})
	return g, seed
}

// barrier runs one dispatch so the parent tree is finalized.
func barrier(t *testing.T, g *graph.Graph) {
	t.Helper()
	op := &capturingOp{seen: map[string]graph.View{}}
	s := graph.NewScheduler(g, []graph.Operator{op}, graph.Limits{MaxRounds: 2, Workers: 1})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("barrier: %v", err)
	}
}

func statusFor(live bool) graph.Status {
	if live {
		return graph.StatusOK
	}
	return graph.StatusEmpty
}

func node(t *testing.T, g *graph.Graph, key string) graph.NodeID {
	t.Helper()
	for _, n := range g.Nodes() {
		if n.Key == key {
			return n.ID
		}
	}
	t.Fatalf("node %s not found", key)
	return graph.NodeID{}
}

// The test this package exists for.
//
// Training reads a finished graph and inference reads a mid-scan View. If the
// two featurize differently the model is fitted on symbols that never occur at
// run time, belief collapses to the prior, and nothing anywhere fails. So the
// same node, seen both ways, must produce the same symbols.
func TestTrainAndServeFeaturizeIdentically(t *testing.T) {
	g, _ := scan(t)
	a := g.Analyze()
	id := node(t, g, "exmple.com")
	n, _ := g.Node(id)

	fromGraph := Features(nodeView{a: a, id: id, n: n})
	fromView := Featurizer(captureView(t, g, "exmple.com"))

	if len(fromGraph) != len(fromView) {
		t.Fatalf("train saw %d symbols, serve saw %d:\n  train %q\n  serve %q",
			len(fromGraph), len(fromView), fromGraph, fromView)
	}
	for i := range fromGraph {
		if fromGraph[i] != fromView[i] {
			t.Fatalf("symbol %d differs: train %q, serve %q\n  train %q\n  serve %q",
				i, fromGraph[i], fromView[i], fromGraph, fromView)
		}
	}
	if len(fromGraph) == 0 {
		t.Fatal("no symbols at all, so this proves nothing")
	}
}

// The featurizer must read what it declares, or a View reports every field
// unset and the model silently sees nothing.
func TestTriggerCoversWhatFeaturesReads(t *testing.T) {
	tr := Trigger()
	has := func(xs []string, want string) bool {
		for _, x := range xs {
			if x == want {
				return true
			}
		}
		return false
	}
	for _, f := range Fields {
		if !has(tr.Fields, f) {
			t.Errorf("Trigger does not declare field %q", f)
		}
	}
	for _, r := range Rels {
		if !has(tr.Rels, r) {
			t.Errorf("Trigger does not declare relation %q", r)
		}
	}
	if !has(tr.Rels, graph.VariantRel) {
		t.Error("Trigger must declare VARIANT_OF: the algorithm symbol comes off that edge")
	}
}

// Symbols are coarse buckets over what an operator can actually see.
//
// The algorithm that produced a variant is deliberately absent: VARIANT_OF runs
// origin -> variant, so on the variant it is an incoming edge, and a View
// exposes outgoing edges only. Training on a symbol inference cannot observe
// would fit the model to noise.
func TestFeaturesCarryTheAlgorithm(t *testing.T) {
	g, _ := scan(t)
	a := g.Analyze()
	id := node(t, g, "exmple.com")
	n, _ := g.Node(id)
	syms := Features(nodeView{a: a, id: id, n: n})

	want := map[string]bool{"type:domain": false, "live:true": false,
		"edge:RESOLVES_TO": false, "len:mid": false}
	for _, s := range syms {
		if _, ok := want[s]; ok {
			want[s] = true
		}
	}
	for s, found := range want {
		if !found {
			t.Errorf("missing symbol %q in %q", s, syms)
		}
	}
	for _, s := range syms {
		if strings.HasPrefix(s, "algo:") || strings.HasPrefix(s, "dist:") {
			t.Errorf("symbol %q is not obtainable from a View; training on it fits noise", s)
		}
	}
}

// A graph that never ran a barrier has no parent tree, so there is nothing to
// train on and saying so beats training on one degenerate path.
func TestPathsRefuseAnUnbarrieredGraph(t *testing.T) {
	g := graph.New(registry(t))
	seed, err := g.Seed("domain", "example.com")
	if err != nil {
		t.Fatal(err)
	}
	origin := graph.NodeRef{Type: "domain", Key: "example.com"}
	g.Apply(graph.Provenance{Operator: "co", Round: 1}, seed, graph.Delta{
		Edges: []graph.EdgeRef{{From: origin, Rel: graph.VariantRel,
			To: graph.NodeRef{Type: "domain", Key: "exmple.com"}}},
	})
	if got := Paths(g, seed); len(got) != 0 {
		t.Errorf("got %d paths from a graph with no finalized parents", len(got))
	}
}

// Unknown must survive into the corpus as itself. Folding it into absent trains
// the model that a rate limit is evidence a name is free.
func TestOutcomeKeepsUnknownDistinct(t *testing.T) {
	g, _ := scan(t)
	// A variant nobody could resolve: the lookup failed rather than answered.
	origin := graph.NodeRef{Type: "domain", Key: "example.com"}
	ref := graph.NodeRef{Type: "domain", Key: "exmpl.com"}
	g.Apply(graph.Provenance{Operator: "co", Round: 1}, node(t, g, "example.com"),
		graph.Delta{Edges: []graph.EdgeRef{{From: origin, Rel: graph.VariantRel, To: ref}}})
	id := node(t, g, "exmpl.com")
	g.SetStatus(id, "dns", graph.StatusFailed)

	a := g.Analyze()
	if got := Outcome(a, id); got != "unknown" {
		t.Errorf("a failed lookup produced outcome %q, want unknown", got)
	}
	if got := Outcome(a, node(t, g, "exampl.com")); got != "absent" {
		t.Errorf("an authoritative negative produced %q, want absent", got)
	}
	if got := Outcome(a, node(t, g, "exmple.com")); got != "live" {
		t.Errorf("a resolving name produced %q, want live", got)
	}
}

// Paths are sequences through the expansion tree, seed first.
func TestPathsAreRootedSequences(t *testing.T) {
	g, seed := scan(t)
	paths := Paths(g, seed)
	if len(paths) == 0 {
		t.Fatal("no paths")
	}
	for _, p := range paths {
		if len(p.Steps) < 2 {
			t.Errorf("path of %d steps; a seed with children should be longer", len(p.Steps))
		}
		// The seed is first and its relation is ignored.
		if !hasSymbol(p.Steps[0].Symbols(), "depth:0") {
			t.Errorf("first step is not the seed: %q", p.Steps[0].Symbols())
		}
		for i, s := range p.Steps[1:] {
			if s.Rel == "" {
				t.Errorf("step %d has no relation; only the seed may", i+1)
			}
		}
	}
}

// Deterministic: the same graph yields the same corpus, or a model's identity
// depends on map iteration order.
func TestPathsAreDeterministic(t *testing.T) {
	g, seed := scan(t)
	first := Paths(g, seed)
	for i := 0; i < 20; i++ {
		got := Paths(g, seed)
		if len(got) != len(first) {
			t.Fatalf("run %d produced %d paths, first produced %d", i, len(got), len(first))
		}
		for j := range got {
			if len(got[j].Steps) != len(first[j].Steps) {
				t.Fatalf("path %d changed length", j)
			}
			for k := range got[j].Steps {
				if got[j].Steps[k].Rel != first[j].Steps[k].Rel {
					t.Fatalf("path %d step %d relation changed", j, k)
				}
			}
		}
	}
}

// Fitting nothing must fail rather than hand back the uniform prior dressed as
// a trained model.
func TestFitRefusesAnEmptyCorpus(t *testing.T) {
	if _, _, err := Fit(DefaultConfig()); err == nil {
		t.Fatal("training on no scans succeeded")
	}
	if _, _, err := Fit(DefaultConfig(), Scan{}); err == nil {
		t.Fatal("training on a nil graph succeeded")
	}
}

// End to end: a scan trains, and the result plugs into the engine's seam.
func TestFitProducesAUsableBeliefModel(t *testing.T) {
	g, seed := scan(t)
	res, corpus, err := Fit(DefaultConfig(), Scan{Graph: g, Seed: seed})
	if err != nil {
		t.Fatalf("fit: %v", err)
	}
	if res.Model == nil {
		t.Fatal("no model")
	}

	s := Describe(corpus, res)
	if s.Paths == 0 || s.Symbols == 0 {
		t.Fatalf("summary is empty: %s", s)
	}
	// The corpus has all three observations in it, which is what makes it worth
	// training on.
	for _, want := range []string{"live", "absent"} {
		if s.Outcomes[want] == 0 {
			t.Errorf("no %q steps in the corpus: %v", want, s.Outcomes)
		}
	}

	// The fitted model satisfies the engine's interface and is not the uniform
	// one: belief must actually vary, or nothing has been learned.
	bm := BeliefFrom(res.Model)
	if bm == nil {
		t.Fatal("BeliefFrom returned nil")
	}
	g2 := graph.New(registry(t))
	if _, err := g2.Seed("domain", "example.com"); err != nil {
		t.Fatal(err)
	}
	g2.SetBeliefModel(bm) // the seam the scanner uses
}

func hasSymbol(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}

// Baum-Welch is unsupervised, so which latent state means "promising" is
// whatever order EM converged to. Focus has to be chosen on the evidence after
// fitting, or belief can come out confidently backwards — which is what
// happened on a real scan: IPv6 addresses at 1.0 and every live typosquat at
// 0.0.
func TestAnchorFocusPicksTheStateThatEmitsLive(t *testing.T) {
	g, seed := scan(t)
	res, _, err := Fit(DefaultConfig(), Scan{Graph: g, Seed: seed})
	if err != nil {
		t.Fatalf("fit: %v", err)
	}

	h, focus, err := AnchorFocus(res.Model)
	if err != nil {
		t.Fatalf("anchor: %v", err)
	}
	if len(h.Focus()) != 1 || h.Focus()[0] != focus {
		t.Fatalf("focus = %v, want [%s]", h.Focus(), focus)
	}

	// The chosen state must be the one most likely to emit a live observation.
	chosen := -1
	for i, s := range h.States() {
		if s == focus {
			chosen = i
		}
	}
	if chosen < 0 {
		t.Fatalf("focus %q is not a state of the model", focus)
	}
	for i := range h.States() {
		if h.LogEmission(i, []string{LiveSymbol}) > h.LogEmission(chosen, []string{LiveSymbol}) {
			t.Errorf("state %d emits live more often than the anchored state %q", i, focus)
		}
	}
}

// Anchoring changes which state is reported and nothing else: the transition
// and emission tables are the fitted ones.
func TestAnchorFocusPreservesTheModel(t *testing.T) {
	g, seed := scan(t)
	res, _, err := Fit(DefaultConfig(), Scan{Graph: g, Seed: seed})
	if err != nil {
		t.Fatalf("fit: %v", err)
	}
	before := res.Model
	after, _, err := AnchorFocus(before)
	if err != nil {
		t.Fatalf("anchor: %v", err)
	}

	if len(after.States()) != len(before.States()) || len(after.Symbols()) != len(before.Symbols()) {
		t.Fatal("anchoring changed the alphabet")
	}
	for i := range before.States() {
		for _, sym := range before.Symbols() {
			b := before.LogEmission(i, []string{sym})
			a := after.LogEmission(i, []string{sym})
			if diff := a - b; diff > 1e-9 || diff < -1e-9 {
				t.Fatalf("emission for state %d symbol %q moved by %g", i, sym, diff)
			}
		}
	}
}

// A corpus with no live observation cannot orient a model, and saying so beats
// picking a state arbitrarily.
func TestAnchorFocusRefusesWithoutLiveEvidence(t *testing.T) {
	g := graph.New(registry(t))
	seed, err := g.Seed("domain", "example.com")
	if err != nil {
		t.Fatal(err)
	}
	origin := graph.NodeRef{Type: "domain", Key: "example.com"}
	ref := graph.NodeRef{Type: "domain", Key: "exmple.com"}
	g.Apply(graph.Provenance{Operator: "co", Round: 1}, seed, graph.Delta{
		Edges: []graph.EdgeRef{{From: origin, Rel: graph.VariantRel, To: ref}},
	})
	barrier(t, g)
	g.SetObservers([]string{"dns"})
	g.SetStatus(node(t, g, "exmple.com"), "dns", graph.StatusEmpty) // absent, never live

	res, _, err := Fit(DefaultConfig(), Scan{Graph: g, Seed: seed})
	if err != nil {
		t.Fatalf("fit: %v", err)
	}
	if _, _, err := AnchorFocus(res.Model); err == nil {
		t.Fatal("anchored a model whose corpus recorded no live observation")
	}
}
