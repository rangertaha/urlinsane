// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package train

import (
	"context"
	"sort"

	"github.com/ipfs/go-cid"
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/model"
)

// nodeView adapts an admitted node to Featurable.
//
// It exists so training and inference featurize through the same function. It
// deliberately implements only what Featurable names, even though a finished
// graph could answer far more: a feature the trainer can compute and an
// operator cannot is a feature the model would never see at run time.
type nodeView struct {
	a  *graph.Analysis
	id graph.NodeID
	n  *graph.Node
	// out is the edge index for this graph, keyed by source then relation.
	//
	// Built once per Paths call. Answering EdgeProps from Analysis.Edges()
	// instead re-walks and re-sorts the whole edge list on every lookup, and
	// Features asks for six relations per node — so featurizing was six full
	// sweeps per node and Paths was quadratic. Measured over a 2000-node scan:
	// 24.1s before, 141ms after.
	out map[graph.NodeID]map[string][]Bag
}

func (v nodeView) Type() string { return v.n.Type.Name() }
func (v nodeView) Key() string  { return v.n.Key }
func (v nodeView) Depth() int   { return v.a.Depth(v.id) }

func (v nodeView) Prop(field string) (graph.Value, bool) {
	f, ok := v.n.Type.Field(field)
	if !ok {
		return graph.Value{}, false
	}
	return v.n.Props.Get(f)
}

func (v nodeView) EdgeProps(rel string) []Bag { return v.out[v.id][rel] }

// featurizer owns the per-graph index and hands out views over it.
//
// It exists so a nodeView cannot be built without its index. The first version
// let callers construct one directly, and a nodeView with a nil index reports
// every node as having no edges — silently, because "no edges" is a legitimate
// answer. That is the shape this repository keeps retiring: correctness by
// remembering to pass something. Here the constructor is the only way in.
type featurizer struct {
	a    *graph.Analysis
	out  map[graph.NodeID]map[string][]Bag
	byID map[graph.NodeID]*graph.Node
}

// newFeaturizer indexes every edge by its source and relation, in one pass.
func newFeaturizer(a *graph.Analysis) *featurizer {
	f := &featurizer{
		a:    a,
		out:  make(map[graph.NodeID]map[string][]Bag),
		byID: make(map[graph.NodeID]*graph.Node),
	}
	for _, n := range a.Nodes() {
		f.byID[n.ID] = n
	}
	for _, e := range a.Edges() {
		byRel := f.out[e.From]
		if byRel == nil {
			byRel = make(map[string][]Bag)
			f.out[e.From] = byRel
		}
		rel := e.Rel.Name()
		byRel[rel] = append(byRel[rel], edgeBag{e})
	}
	return f
}

// view returns the Featurable for one admitted node.
func (f *featurizer) view(id graph.NodeID) nodeView {
	return nodeView{a: f.a, id: id, n: f.byID[id], out: f.out}
}

// Symbols featurizes one node of a finished graph, the same way inference
// featurizes a View.
func (f *featurizer) Symbols(id graph.NodeID) []string { return Features(f.view(id)) }

// edgeBag reads an admitted edge's props, matching what graph.EdgeView exposes
// at run time.
type edgeBag struct{ e *graph.Edge }

func (b edgeBag) Prop(field string) (graph.Value, bool) {
	f, ok := b.e.Rel.Field(field)
	if !ok {
		return graph.Value{}, false
	}
	return b.e.Props.Get(f)
}

// Outcome is the label a node contributes to the emission alphabet.
//
// Three-valued, exactly as a scan reports it. Folding unknown into absent would
// train the model that a rate limit is evidence a name is free, which is the
// one inference this whole codebase is arranged to prevent — and it would do it
// invisibly, because the corpus would still look balanced.
func Outcome(a *graph.Analysis, id graph.NodeID) string {
	switch a.Existence(id) {
	case graph.Live:
		return "live"
	case graph.Absent:
		return "absent"
	case graph.Unknown:
		return "unknown"
	}
	return "" // untried: no observation ran, so there is nothing to report
}

// Paths walks a finished graph into root-to-leaf expansion paths.
//
// An HMM is defined over sequences, and the expansion *tree* — not the graph —
// is what provides them: Graph.Parent gives every admitted node exactly one
// parent, fixed at a barrier and never revised, so seed → node is a sequence.
// The graph itself is cyclic (domain → ip → PTR → domain) and has no sequences
// in it at all.
//
// One path per leaf. An interior node appears in every path that passes through
// it, which is correct: it was observed once per descendant lineage, and
// Baum-Welch weights a transition by how often it was traversed.
//
// The graph must have completed at least one barrier, which every live scan
// has and no rehydrated one does — parents are assigned there, and the store
// does not persist them. Call Finalize after loading from the store.
func Paths(g *graph.Graph, seed graph.NodeID) []model.Path {
	a := g.Analyze()

	// Parent and relation for every node, plus the child lists.
	parent := map[graph.NodeID]graph.NodeID{}
	rel := map[graph.NodeID]string{}
	children := map[graph.NodeID][]graph.NodeID{}
	nodes := a.Nodes()

	for _, n := range nodes {
		p, r, ok := g.Parent(n.ID)
		if !ok {
			continue
		}
		parent[n.ID] = p
		rel[n.ID] = r
		children[p] = append(children[p], n.ID)
	}

	f := newFeaturizer(a)
	byID := f.byID

	// Deterministic: children in (type, key) order, so the corpus a scan
	// produces does not depend on admission order or map iteration.
	for p := range children {
		kids := children[p]
		sort.Slice(kids, func(i, j int) bool {
			a, b := byID[kids[i]], byID[kids[j]]
			if a.Type.Name() != b.Type.Name() {
				return a.Type.Name() < b.Type.Name()
			}
			return a.Key < b.Key
		})
		children[p] = kids
	}

	// Memoized: an interior node appears in every path through it, and
	// featurizing it once per path rather than once per node made the seed of
	// a thousand-leaf scan cost a thousand featurizations.
	cache := make(map[graph.NodeID]model.Trace, len(nodes))
	trace := func(id graph.NodeID) model.Trace {
		if t, ok := cache[id]; ok {
			return t
		}
		t := model.Trace{
			Rel:     rel[id],
			Props:   f.Symbols(id),
			Outcome: Outcome(a, id),
		}
		cache[id] = t
		return t
	}

	// The parent tree is finalized at a barrier, so a graph that has never run
	// one has no parents at all. Walking it would yield a single-step path from
	// the seed, which trains without complaint and means nothing — so refuse it
	// here and let Fit's empty-corpus error say so.
	if len(nodes) > 1 && len(parent) == 0 {
		return nil
	}

	var out []model.Path
	var walk func(id graph.NodeID, prefix []model.Trace)
	walk = func(id graph.NodeID, prefix []model.Trace) {
		if byID[id] == nil {
			return
		}
		step := append(append([]model.Trace(nil), prefix...), trace(id))
		kids := children[id]
		if len(kids) == 0 {
			out = append(out, model.Path{Steps: step})
			return
		}
		for _, k := range kids {
			walk(k, step)
		}
	}
	walk(seed, nil)
	return out
}

// Finalize gives a graph the barrier its expansion tree comes from.
//
// Parents are assigned at a barrier and are *not* persisted: the side tables
// carry depth and closure but not the tree, because the tree is derived state
// the engine rebuilds. A graph replayed out of the store therefore arrives with
// every edge and no parents at all, and Paths refuses it.
//
// A scheduler with no operators is the smallest way to get one: it runs barrier
// 0, finds no eligible work, and stops.
//
// The observer set is saved and restored around it, and that is not defensive
// tidying. NewScheduler *rebuilds* the set from the operators it is given, and
// with no operators the set comes out empty — at which point Graph.observes
// answers true for everything, and a pure variant generator's "I produced this"
// counts as evidence the name exists. Measured before the restore was added:
// exampl.com, whose only lookup returned an authoritative absence, reported
// Outcome "live". Every corpus built from a rehydrated scan was labelled that
// way, and the AUC computed from it was meaningless.
//
// The set is scan state, persisted with the graph precisely so a re-render
// agrees with the run that produced it (store side tables, FormatVersion 2).
// Losing it here would undo that on the one path — rehydration — where it
// matters most.
func Finalize(g *graph.Graph) error {
	observers := g.Observers()
	err := graph.NewScheduler(g, nil, graph.Limits{MaxRounds: 1}).Run(context.Background())
	g.SetObservers(observers)
	return err
}

// CorpusOf builds a training corpus from one scan.
//
// roots are the scan roots the graph came from, copied into the model's
// provenance so a trained model always points back at the data it was fitted
// on. Pass none for an in-memory graph; the model then records that it was
// fitted on something unaddressed, which is exactly what it was.
func CorpusOf(g *graph.Graph, seed graph.NodeID, roots ...cid.Cid) model.Corpus {
	return model.Corpus{Paths: Paths(g, seed), CIDs: roots}
}

// Alphabet is every symbol and relation a corpus contains, sorted.
//
// The alphabet is part of the model's identity and therefore of its CID, so it
// is derived from the corpus rather than declared: a hand-written alphabet
// drifts from the featurizer the first time a symbol is added, and the symptom
// is a model that maps a live observation to out-of-vocabulary.
func Alphabet(c model.Corpus) (symbols, rels []string) {
	symSet, relSet := map[string]bool{}, map[string]bool{}
	for _, p := range c.Paths {
		for i, s := range p.Steps {
			for _, sym := range s.Symbols() {
				symSet[sym] = true
			}
			// Steps[0] is the seed and its Rel is ignored — it has no parent.
			if i > 0 && s.Rel != "" {
				relSet[s.Rel] = true
			}
		}
	}
	return sortedKeys(symSet), sortedKeys(relSet)
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
