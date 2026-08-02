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

// Package dag lays out a directed graph for presentation. Nothing here decides
// execution order — the engine dispatches on data patterns, not topology.
//
// The type-flow graph is legitimately cyclic (domain → ip → PTR → domain), so a
// plain topological sort is not defined on it. Condensing strongly connected
// components first gives a DAG that is, which is what lets `--explain` render
// the plan in readable layers without pretending the cycles are not there.
package dag

import "sort"

// Graph is an adjacency list keyed by node name.
type Graph map[string][]string

// Component is one strongly connected component. A component with a single
// member and no self-edge is an ordinary node; anything larger is a genuine
// cycle in the type flow.
type Component struct {
	Members []string
	Layer   int
	Cyclic  bool
}

// Condense returns the strongly connected components of g, each assigned a
// layer such that every edge runs from a lower layer to an equal or higher one.
// Components are returned in layer order, and members within a component are
// sorted, so the rendering is stable across runs.
func Condense(g Graph) []Component {
	nodes := make([]string, 0, len(g))
	seen := map[string]bool{}
	for n, outs := range g {
		if !seen[n] {
			seen[n] = true
			nodes = append(nodes, n)
		}
		for _, m := range outs {
			if !seen[m] {
				seen[m] = true
				nodes = append(nodes, m)
			}
		}
	}
	sort.Strings(nodes)

	comps := tarjan(g, nodes)

	// Map each node to its component index, then build the condensation.
	idx := map[string]int{}
	for i, c := range comps {
		for _, m := range c.Members {
			idx[m] = i
		}
	}
	cond := make(map[int]map[int]bool, len(comps))
	indeg := make([]int, len(comps))
	for n, outs := range g {
		for _, m := range outs {
			a, b := idx[n], idx[m]
			if a == b {
				continue
			}
			if cond[a] == nil {
				cond[a] = map[int]bool{}
			}
			if !cond[a][b] {
				cond[a][b] = true
				indeg[b]++
			}
		}
	}

	// Kahn over the condensation. It cannot cycle: that is the point of
	// condensing first.
	layer := make([]int, len(comps))
	var frontier []int
	for i := range comps {
		if indeg[i] == 0 {
			frontier = append(frontier, i)
		}
	}
	sort.Ints(frontier)
	for depth := 0; len(frontier) > 0; depth++ {
		var next []int
		for _, i := range frontier {
			layer[i] = depth
			for j := range cond[i] {
				indeg[j]--
				if indeg[j] == 0 {
					next = append(next, j)
				}
			}
		}
		sort.Ints(next)
		frontier = next
	}

	out := make([]Component, len(comps))
	for i, c := range comps {
		c.Layer = layer[i]
		out[i] = c
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Layer != out[j].Layer {
			return out[i].Layer < out[j].Layer
		}
		return out[i].Members[0] < out[j].Members[0]
	})
	return out
}

// tarjan returns the strongly connected components of g over nodes.
func tarjan(g Graph, nodes []string) []Component {
	var (
		index   = map[string]int{}
		low     = map[string]int{}
		onStack = map[string]bool{}
		stack   []string
		counter int
		out     []Component
	)

	var strongconnect func(v string)
	strongconnect = func(v string) {
		index[v] = counter
		low[v] = counter
		counter++
		stack = append(stack, v)
		onStack[v] = true

		outs := append([]string(nil), g[v]...)
		sort.Strings(outs)
		for _, w := range outs {
			switch {
			case func() bool { _, ok := index[w]; return !ok }():
				strongconnect(w)
				if low[w] < low[v] {
					low[v] = low[w]
				}
			case onStack[w]:
				if index[w] < low[v] {
					low[v] = index[w]
				}
			}
		}

		if low[v] == index[v] {
			var members []string
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				members = append(members, w)
				if w == v {
					break
				}
			}
			sort.Strings(members)
			cyclic := len(members) > 1
			if !cyclic {
				for _, w := range g[v] {
					if w == v {
						cyclic = true
						break
					}
				}
			}
			out = append(out, Component{Members: members, Cyclic: cyclic})
		}
	}

	for _, n := range nodes {
		if _, ok := index[n]; !ok {
			strongconnect(n)
		}
	}
	return out
}

// Reachable returns every node reachable from any of roots, inclusive.
func Reachable(g Graph, roots ...string) map[string]bool {
	out := map[string]bool{}
	var walk func(string)
	walk = func(n string) {
		if out[n] {
			return
		}
		out[n] = true
		for _, m := range g[n] {
			walk(m)
		}
	}
	for _, r := range roots {
		walk(r)
	}
	return out
}
