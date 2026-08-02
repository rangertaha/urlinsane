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

package model

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// --- a stand-in view --------------------------------------------------------

// fakeView is the smallest thing that satisfies graph.View. The adapter reads
// nothing but props, so nothing else needs to be real.
type fakeView struct {
	typ, key string
	depth    int
	props    map[string]graph.Value
}

func (f fakeView) ID() graph.NodeID { return graph.NodeID{} }
func (f fakeView) Type() string     { return f.typ }
func (f fakeView) Key() string      { return f.key }
func (f fakeView) Depth() int       { return f.depth }
func (f fakeView) Prop(field string) (graph.Value, bool) {
	v, ok := f.props[field]
	return v, ok
}
func (f fakeView) Edges(string) []graph.EdgeView { return nil }
func (f fakeView) Ref() graph.NodeRef            { return graph.NodeRef{Type: f.typ, Key: f.key} }

// --- the uniform reduction --------------------------------------------------

// TestUniformBeliefIsExactlyOne is the property §10.5 rests on: before any
// model exists the engine must behave exactly as it does with graph's own
// default, which returns a literal 1. Anything else — even a constant 0.5 —
// would put every node under a BeliefAbove(0.5) gate and silently change what
// the engine does.
func TestUniformBeliefIsExactlyOne(t *testing.T) {
	b := UniformBelief()
	if got := b.Initial(); got != 1 {
		t.Fatalf("Initial() = %v, want exactly 1", got)
	}
	rels := []string{"", "VARIANT_OF", "RESOLVES_TO", "never-seen"}
	views := []graph.View{
		nil,
		fakeView{typ: "domain", key: "a.com"},
		fakeView{typ: "domain", key: "b.com", props: map[string]graph.Value{
			"punycode": graph.String("xn--"), "rank": graph.Int(3),
		}},
	}
	for _, rel := range rels {
		for i, v := range views {
			for _, parent := range []float64{0, 0.25, 1, -1, 2} {
				if got := b.Step(parent, rel, v); got != 1 {
					t.Fatalf("Step(%v, %q, view %d) = %v, want exactly 1", parent, rel, i, got)
				}
			}
		}
	}
}

// TestUniformBeliefWithAFeaturizerIsStillOne: a uniform model must stay
// uninformative no matter what the featurizer feeds it, or "ships before a
// model exists" becomes "ships until someone adds a featurizer".
func TestUniformBeliefWithAFeaturizerIsStillOne(t *testing.T) {
	b := NewBelief(Uniform(), PropFeatures("punycode", "rank"))
	v := fakeView{typ: "domain", key: "b.com", props: map[string]graph.Value{
		"punycode": graph.String("xn--"), "rank": graph.Int(3),
	}}
	if got := b.Step(0.3, "VARIANT_OF", v); got != 1 {
		t.Fatalf("Step = %v, want exactly 1", got)
	}
}

// TestUniformMultiStateModelIsUnranked is the weaker but more general form: a
// model whose tables are uniform gives every node the same belief regardless of
// relation or props, so the frontier ordering falls back to (depth, type, key)
// — breadth-first and unranked.
func TestUniformMultiStateModelIsUnranked(t *testing.T) {
	h, err := New(Spec{
		States:  []string{"dull", "interesting"},
		Focus:   []string{"interesting"},
		Rels:    []string{"VARIANT_OF", "RESOLVES_TO"},
		Symbols: []string{"resolves=true", "resolves=false"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	b := NewBelief(h, PropFeatures("resolves"))
	want := b.Initial()
	for _, rel := range []string{"VARIANT_OF", "RESOLVES_TO", "unknown"} {
		for _, val := range []bool{true, false} {
			v := fakeView{typ: "domain", key: "x", props: map[string]graph.Value{
				"resolves": graph.Bool(val),
			}}
			if got := b.Step(want, rel, v); math.Abs(got-want) > 1e-12 {
				t.Fatalf("uniform tables gave belief %v down %q, want the constant %v", got, rel, want)
			}
		}
	}
}

// --- end to end against the real engine -------------------------------------

// beliefRegistry is the smallest registry that produces a parent, a relation
// and an observed prop.
func beliefRegistry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()
	lower := func(s string) (string, error) {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			return "", fmt.Errorf("empty key")
		}
		return s, nil
	}
	if _, err := r.AddType(graph.NodeTypeDef{
		Name: "domain", Cap: graph.Nameable, Version: 1, Canonical: lower,
	}); err != nil {
		t.Fatalf("domain: %v", err)
	}
	if _, err := r.AddType(graph.NodeTypeDef{
		Name: "ip", Cap: graph.Observed, Version: 1, Canonical: lower,
		Fields: []graph.FieldDef{{Name: "asn", Kind: graph.KindString}},
	}); err != nil {
		t.Fatalf("ip: %v", err)
	}
	if _, err := r.AddRel(graph.RelDef{
		Name: "RESOLVES_TO", Class: graph.Observation, Version: 1,
	}); err != nil {
		t.Fatalf("RESOLVES_TO: %v", err)
	}
	return r
}

// resolveOp emits three ips from the seed, one of them carrying an asn prop.
type resolveOp struct{}

func (resolveOp) Id() string   { return "resolve" }
func (resolveOp) Version() int { return 1 }
func (resolveOp) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{"domain"}}}
}
func (resolveOp) Emits() graph.Effects {
	return graph.Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}, Props: []string{"asn"}}
}
func (resolveOp) Resource() string { return "dns" }
func (resolveOp) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	from := v.Ref()
	var d graph.Delta
	for _, ip := range []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"} {
		to := graph.NodeRef{Type: "ip", Key: ip}
		d.Nodes = append(d.Nodes, to)
		d.Edges = append(d.Edges, graph.EdgeRef{From: from, Rel: "RESOLVES_TO", To: to})
	}
	first := graph.NodeRef{Type: "ip", Key: "10.0.0.1"}
	d.Props = append(d.Props, graph.PropSet{
		Node: &first, Field: "asn", Value: graph.String("AS64500"),
	})
	return d, graph.OK()
}

func runExpansion(t *testing.T, m graph.BeliefModel) map[string]float64 {
	t.Helper()
	g := graph.New(beliefRegistry(t))
	if m != nil {
		g.SetBeliefModel(m)
	}
	if _, err := g.Seed("domain", "Example.COM"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	s := graph.NewScheduler(g, []graph.Operator{resolveOp{}}, graph.Limits{})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}
	out := map[string]float64{}
	for _, n := range g.Nodes() {
		out[n.Type.Name()+":"+n.Key] = g.Belief(n.ID)
	}
	if len(out) != 4 {
		t.Fatalf("expansion produced %d nodes, want 4", len(out))
	}
	return out
}

// TestUniformModelMatchesTheEngineDefault runs a real expansion twice — once
// with graph's built-in default and once with this package's uniform model —
// and requires every belief to agree bit for bit. This is the concrete form of
// "a uniform model reduces exactly to breadth-first unranked expansion".
func TestUniformModelMatchesTheEngineDefault(t *testing.T) {
	def := runExpansion(t, nil)
	ours := runExpansion(t, UniformBelief())
	if len(def) != len(ours) {
		t.Fatalf("different node sets: %d vs %d", len(def), len(ours))
	}
	for k, want := range def {
		got, ok := ours[k]
		if !ok {
			t.Fatalf("%s missing under the uniform model", k)
		}
		if got != want {
			t.Fatalf("%s: uniform model gave %v, engine default gave %v", k, got, want)
		}
		if want != 1 {
			t.Fatalf("%s: engine default belief is %v, expected 1", k, want)
		}
	}
}

// TestTrainedModelIsDeterministicAcrossRuns: same model, same input, same
// belief, every run. It is the property gating and pruning are only defensible
// under, since both decisions are irreversible.
func TestTrainedModelIsDeterministicAcrossRuns(t *testing.T) {
	res, err := Train(syntheticCorpus(), trainConfig())
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	m := NewBelief(res.Model, PropFeatures("asn"))
	first := runExpansion(t, m)
	for i := 0; i < 5; i++ {
		got := runExpansion(t, m)
		for k, want := range first {
			if got[k] != want {
				t.Fatalf("run %d: %s = %v, want %v", i, k, got[k], want)
			}
		}
	}
}

// TestObservationsMoveBelief: the node that carries a prop must end up with a
// different belief from its siblings, or emissions are doing nothing and the
// model is a transition prior with extra steps.
func TestObservationsMoveBelief(t *testing.T) {
	h, err := New(Spec{
		States:  []string{"dull", "interesting"},
		Focus:   []string{"interesting"},
		Rels:    []string{"RESOLVES_TO"},
		Symbols: []string{"asn=AS64500"},
		Init:    []float64{0.5, 0.5},
		Trans:   map[string][][]float64{"RESOLVES_TO": {{0.9, 0.1}, {0.1, 0.9}}},
		Emit:    [][]float64{{0.05, 0.95}, {0.9, 0.1}},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	got := runExpansion(t, NewBelief(h, PropFeatures("asn")))
	observed, bare := got["ip:10.0.0.1"], got["ip:10.0.0.2"]
	if math.Abs(observed-bare) < 0.1 {
		t.Fatalf("observed node %v and bare sibling %v are indistinguishable", observed, bare)
	}
	// The bare sibling is a node with no props: its belief is the transition
	// prior alone, pushed from the seed's belief (§10.2).
	seed := got["domain:example.com"]
	closeTo(t, bare, h.Mass(h.Forward(h.Lift(seed), "RESOLVES_TO", nil)), "prior-only belief")
	closeTo(t, observed, h.Mass(h.Forward(h.Lift(seed), "RESOLVES_TO", []string{"asn=AS64500"})), "observed belief")
}

// --- featurizer -------------------------------------------------------------

func TestPropFeaturesRendersEveryKind(t *testing.T) {
	f := PropFeatures("s", "i", "fl", "b", "by", "missing")
	v := fakeView{props: map[string]graph.Value{
		"s":  graph.String("hi"),
		"i":  graph.Int(-7),
		"fl": graph.Float(1.5),
		"b":  graph.Bool(true),
		"by": graph.Bytes([]byte{0xde, 0xad}),
	}}
	got := f(v)
	want := []string{"b=true", "by=dead", "fl=1.5", "i=-7", "s=hi"} // sorted by field
	if len(got) != len(want) {
		t.Fatalf("features = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("features = %v, want %v", got, want)
		}
	}
}

func TestPropFeaturesIsStable(t *testing.T) {
	f := PropFeatures("b", "a", "c")
	v := fakeView{props: map[string]graph.Value{
		"a": graph.String("1"), "b": graph.String("2"), "c": graph.String("3"),
	}}
	first := strings.Join(f(v), ",")
	for i := 0; i < 50; i++ {
		if got := strings.Join(f(v), ","); got != first {
			t.Fatalf("run %d gave %q, want %q", i, got, first)
		}
	}
}

func TestNewBeliefToleratesANilModel(t *testing.T) {
	b := NewBelief(nil, nil)
	if b.Initial() != 1 {
		t.Fatal("a nil model should fall back to uniform")
	}
	if b.Model() == nil {
		t.Fatal("Model() should never be nil")
	}
}
