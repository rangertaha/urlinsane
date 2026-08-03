package train

import (
	"fmt"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// The edge index must answer exactly what a direct scan of the edge list does.
//
// An oracle test: `old` below is the pre-index implementation, kept as the
// reference definition of the answer rather than as a second implementation to
// maintain. The index was introduced for speed (24.1s -> 141ms over 2000 nodes),
// and the risk of that kind of change is a silent difference in what it returns
// — a missing relation, a dropped duplicate, a different order — none of which
// would fail anything, because "no edges" and "these edges in another order" are
// both legitimate answers.
//
// Checks every node against every declared relation plus two that match nothing.
func TestEdgeIndexAgreesWithADirectScan(t *testing.T) {
	g, _ := wideScan(t, 200)
	a := g.Analyze()
	f := newFeaturizer(a)

	// The pre-refactor implementation, verbatim.
	old := func(id graph.NodeID, rel string) []Bag {
		var out []Bag
		for _, e := range a.Edges() {
			if e.From != id || e.Rel.Name() != rel {
				continue
			}
			out = append(out, edgeBag{e})
		}
		return out
	}

	rels := append([]string{}, Rels...)
	rels = append(rels, "NOPE", "")
	for _, n := range a.Nodes() {
		for _, rel := range rels {
			want, got := old(n.ID, rel), f.view(n.ID).EdgeProps(rel)
			if len(want) != len(got) {
				t.Fatalf("node %s rel %q: old %d edges, new %d", n.Key, rel, len(want), len(got))
			}
			for i := range want {
				if fmt.Sprint(want[i]) != fmt.Sprint(got[i]) {
					t.Fatalf("node %s rel %q edge %d differs", n.Key, rel, i)
				}
			}
		}
	}

	// And the symbols themselves must be identical for every node.
	for _, n := range a.Nodes() {
		got := f.Symbols(n.ID)
		want := Features(oldView{a: a, id: n.ID, n: n, old: old})
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("node %s: new %q, old %q", n.Key, got, want)
		}
	}
}

type oldView struct {
	a   *graph.Analysis
	id  graph.NodeID
	n   *graph.Node
	old func(graph.NodeID, string) []Bag
}

func (v oldView) Type() string { return v.n.Type.Name() }
func (v oldView) Key() string  { return v.n.Key }
func (v oldView) Depth() int   { return v.a.Depth(v.id) }
func (v oldView) Prop(f string) (graph.Value, bool) {
	return nodeView{a: v.a, id: v.id, n: v.n}.Prop(f)
}
func (v oldView) EdgeProps(rel string) []Bag { return v.old(v.id, rel) }
