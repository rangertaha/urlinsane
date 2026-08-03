// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

// BeliefModel scores a node for execution control only: frontier ordering,
// pruning and operator gating. It never contributes to a reported number, which
// is why analyzers cannot see belief at all.
//
// Belief is a pure function of the parent's belief and the node's props as of
// the current barrier. Nothing here reaches sideways into the graph, so the
// same run recomputes the same values in the same order.
type BeliefModel interface {
	// Initial is the seed's prior and starting state, the seed having no parent.
	Initial() (float64, State)
	// Step is the forward filtering step: the parent's state pushed through the
	// relation that admitted this node, then conditioned on its props. It
	// returns the child's scalar belief and the state its own children inherit.
	//
	// parent is whatever this model returned for the parent node, or nil if the
	// parent has none yet; a model must treat nil as its initial state.
	Step(parent State, rel string, v View) (float64, State)
}

// State is a model's latent state, carried between a parent and its children.
// The graph stores and forwards it without ever inspecting it.
//
// It exists because §10.1 specifies a hidden Markov model, and forward
// filtering in an HMM propagates a *distribution over latent states* — a
// vector. An earlier version of this interface passed the parent's scalar
// belief instead, which forced a model to reconstruct a plausible distribution
// from one number and collapse it again on the way out. That round trip is
// exact for a two-state model and lossy for three or more, so the interface
// silently answered §16's open question about state cardinality as "two" — not
// by anyone's decision, but as a consequence of a type signature. Numbers from
// a larger model would have looked entirely reasonable and been wrong.
//
// The scalar is still what the engine ranks and gates on; the state is the
// model's own business.
type State any

// uniformModel is the default. It reduces expansion to breadth-first and
// unranked — exactly today's behaviour — so the engine ships and runs correctly
// before any model exists, and a model that turns out poor can be dropped
// without invalidating a single result.
type uniformModel struct{}

func (uniformModel) Initial() (float64, State)                 { return 1, nil }
func (uniformModel) Step(State, string, View) (float64, State) { return 1, nil }

// SetBeliefModel installs the execution model.
func (g *Graph) SetBeliefModel(m BeliefModel) {
	if m == nil {
		m = uniformModel{}
	}
	g.model = m
}

// Belief returns a node's current belief.
func (g *Graph) Belief(id NodeID) float64 {
	if b, ok := g.belief[id]; ok {
		return b
	}
	b, _ := g.model.Initial()
	return b
}

// Parent returns the node's tree parent and the relation that admitted it.
func (g *Graph) Parent(id NodeID) (NodeID, string, bool) {
	p, ok := g.parent[id]
	if !ok {
		return NodeID{}, "", false
	}
	return p.node, p.rel, true
}

// considerParent offers a candidate parent. The tree parent is the minimum by
// (depth, relation, parent NodeID) over the candidates present when the node's
// creating round ends — min-so-far is not min-final, which is why this only
// records candidates and finalizeParents picks.
func (g *Graph) considerParent(child, from NodeID, rel string) {
	if child == g.seed || child == from {
		return
	}
	g.candidates[child] = append(g.candidates[child], parentRef{node: from, rel: rel})
}

type parentRef struct {
	node NodeID
	rel  string
}

// finalizeParents picks each unparented node's tree parent from the candidates
// gathered so far. Once fixed, a parent never moves: a parent revised after the
// fact would revise belief after the gating and pruning that belief already
// drove, and pruning cannot be undone.
func (g *Graph) finalizeParents() {
	for child, cands := range g.candidates {
		if _, done := g.parent[child]; done {
			continue
		}
		if _, live := g.nodes[child]; !live {
			continue
		}
		best := cands[0]
		for _, c := range cands[1:] {
			if betterParent(g, c, best) {
				best = c
			}
		}
		g.parent[child] = best
	}
	g.candidates = map[NodeID][]parentRef{}
}

func betterParent(g *Graph, a, b parentRef) bool {
	if da, db := g.depth[a.node], g.depth[b.node]; da != db {
		return da < db
	}
	if a.rel != b.rel {
		return a.rel < b.rel
	}
	return compareID(a.node[:], b.node[:]) < 0
}

// recomputeBelief runs at every barrier. Parents are fixed once; belief is not,
// because a node created in one round has only the props its creating operator
// set — observation operators on it run in the next round.
func (g *Graph) recomputeBelief() {
	g.finalizeParents()

	// Down the parent forest, so a parent's state is always current before its
	// children read it.
	//
	// Depth order does not give that. Depth counts observation hops, and
	// structural and variant edges cost none, so a parent and its child
	// routinely share one — the whole decomposition of a seed and every variant
	// of every part of it sit at depth 0 together. Within a depth the old sort
	// fell through to type and key, which have nothing to do with descent: for
	// example.com -> zzz.com -> aaa.com, all structural, "aaa.com" sorted first
	// and stepped from a parent that had not been recomputed yet. It read the
	// previous barrier's state, or nil on the first, and a three-generation
	// chain came out believing it was two. Forward filtering over a chain that
	// is not a chain means nothing.
	//
	// Siblings keep the old (depth, type, key) order, so a run whose ordering
	// was already correct produces exactly the values it did before.
	ids := make([]NodeID, 0, len(g.nodes))
	for _, id := range g.order {
		ids = append(ids, id)
	}
	sortByDepthThenKey(g, ids)

	children := make(map[NodeID][]NodeID, len(ids))
	roots := make([]NodeID, 0, 1)
	for _, id := range ids {
		p, ok := g.parent[id]
		if id == g.seed || !ok {
			roots = append(roots, id)
			continue
		}
		children[p.node] = append(children[p.node], id)
	}

	// Breadth-first rather than recursive: a chain is as long as the graph is
	// deep, and recursion here would put that on the stack.
	done := make(map[NodeID]bool, len(ids))
	queue := append(make([]NodeID, 0, len(ids)), roots...)
	for i := 0; i < len(queue); i++ {
		id := queue[i]
		if done[id] {
			continue
		}
		done[id] = true
		g.stepBelief(id)
		queue = append(queue, children[id]...)
	}

	// A node the walk never reached is in a parent cycle, which finalizeParents
	// is not supposed to be able to produce. Give it a value anyway rather than
	// leaving it with a stale one from an earlier barrier: a missing belief
	// gates admission and network calls just as firmly as a wrong one.
	for _, id := range ids {
		if !done[id] {
			done[id] = true
			g.stepBelief(id)
		}
	}
}

// stepBelief recomputes one node from its parent's current state.
func (g *Graph) stepBelief(id NodeID) {
	p, ok := g.parent[id]
	if id == g.seed || !ok {
		g.belief[id], g.bstate[id] = g.model.Initial()
		return
	}
	// The parent's state, not its scalar.
	g.belief[id], g.bstate[id] = g.model.Step(g.bstate[p.node], p.rel, g.fullView(id))
}

func sortByDepthThenKey(g *Graph, ids []NodeID) {
	sortSlice(ids, func(a, b NodeID) bool {
		if da, db := g.depth[a], g.depth[b]; da != db {
			return da < db
		}
		na, nb := g.nodes[a], g.nodes[b]
		if na.Type.name != nb.Type.name {
			return na.Type.name < nb.Type.name
		}
		return na.Key < nb.Key
	})
}

// fullView is the unfiltered view the belief model reads. The model is engine
// code rather than a plugin, so the read-set discipline that scopes operator
// views does not apply to it.
func (g *Graph) fullView(id NodeID) View {
	v := &view{g: g, id: id, fields: map[string]bool{}, rels: map[string]bool{}}
	if n, ok := g.nodes[id]; ok {
		for _, f := range n.Type.sch.fields {
			v.fields[f.name] = true
		}
	}
	for _, rel := range g.reg.Rels() {
		v.rels[rel.name] = true
	}
	return v
}
