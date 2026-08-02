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

// BeliefModel scores a node for execution control only: frontier ordering,
// pruning and operator gating. It never contributes to a reported number, which
// is why analyzers cannot see belief at all.
//
// Belief is a pure function of the parent's belief and the node's props as of
// the current barrier. Nothing here reaches sideways into the graph, so the
// same run recomputes the same values in the same order.
type BeliefModel interface {
	// Initial is the seed's prior, the seed having no parent.
	Initial() float64
	// Step is the forward filtering step: the parent's belief pushed through
	// the relation that admitted this node, then conditioned on its props.
	Step(parent float64, rel string, v View) float64
}

// uniformModel is the default. It reduces expansion to breadth-first and
// unranked — exactly today's behaviour — so the engine ships and runs correctly
// before any model exists, and a model that turns out poor can be dropped
// without invalidating a single result.
type uniformModel struct{}

func (uniformModel) Initial() float64                   { return 1 }
func (uniformModel) Step(float64, string, View) float64 { return 1 }

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
	return g.model.Initial()
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

	// Nodes in depth order, so a parent's belief is always current before its
	// children read it.
	ids := make([]NodeID, 0, len(g.nodes))
	for _, id := range g.order {
		ids = append(ids, id)
	}
	sortByDepthThenKey(g, ids)

	for _, id := range ids {
		if id == g.seed {
			g.belief[id] = g.model.Initial()
			continue
		}
		p, ok := g.parent[id]
		if !ok {
			g.belief[id] = g.model.Initial()
			continue
		}
		v := g.fullView(id)
		g.belief[id] = g.model.Step(g.Belief(p.node), p.rel, v)
	}
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
