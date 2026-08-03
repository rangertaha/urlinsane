// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package train

import (
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

func (v nodeView) EdgeProps(rel string) []Bag {
	var out []Bag
	for _, e := range v.a.Edges() {
		if e.From != v.id || e.Rel.Name() != rel {
			continue
		}
		out = append(out, edgeBag{e})
	}
	return out
}

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
// The graph must have completed at least one barrier, which every real scan
// has: parents are assigned there and nowhere else.
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

	// Deterministic: children in (type, key) order, so the corpus a scan
	// produces does not depend on admission order or map iteration.
	byID := map[graph.NodeID]*graph.Node{}
	for _, n := range nodes {
		byID[n.ID] = n
	}
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

	trace := func(id graph.NodeID) model.Trace {
		n := byID[id]
		return model.Trace{
			Rel:     rel[id],
			Props:   Features(nodeView{a: a, id: id, n: n}),
			Outcome: Outcome(a, id),
		}
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
