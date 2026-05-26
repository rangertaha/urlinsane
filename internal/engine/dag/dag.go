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

// Package dag turns a set of dependency-declaring nodes into ordered execution
// levels via a topological sort. It is pure (no I/O) so the engine's collector
// ordering can be unit-tested in isolation.
package dag

import (
	"fmt"
	"sort"
)

// Node is anything the DAG can order: it has a stable identity and a list of
// dependency ids that must run before it.
type Node interface {
	Id() string
	Dependencies() []string
}

// orderer is satisfied by nodes that carry a numeric priority. When present it
// is used only as a tie-breaker to order independent nodes within a level
// (ascending), giving deterministic, predictable output. Nodes that do not
// implement it sort by Id alone.
type orderer interface {
	Order() int
}

// Levels groups nodes into topological waves: every node in level i has all of
// its dependencies satisfied by levels < i, and nodes within a level are
// mutually independent (safe to run in any order, or concurrently). Within a
// level nodes are ordered by (Order, Id) for determinism.
//
// When a node depends on an id that is not in nodes, resolve is called to
// construct (auto-include) the missing dependency — i.e. dependency closure.
// Passing a nil resolve instead makes a missing dependency an error.
//
// Returns an error if the dependency graph contains a cycle.
func Levels[T Node](nodes []T, resolve func(id string) (T, bool)) ([][]T, error) {
	// Build the closed set of nodes: the given nodes plus, transitively, any
	// dependencies they reference that were not supplied.
	index := make(map[string]T)
	queued := make(map[string]bool)
	var queue []T
	enqueue := func(n T) {
		if !queued[n.Id()] {
			queued[n.Id()] = true
			queue = append(queue, n)
		}
	}
	for _, n := range nodes {
		enqueue(n)
	}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		index[n.Id()] = n
		for _, dep := range n.Dependencies() {
			if queued[dep] {
				continue
			}
			if resolve == nil {
				return nil, fmt.Errorf("dag: %q depends on missing %q", n.Id(), dep)
			}
			rn, ok := resolve(dep)
			if !ok {
				return nil, fmt.Errorf("dag: %q depends on unknown %q", n.Id(), dep)
			}
			enqueue(rn)
		}
	}

	// Kahn's algorithm: in-degree = number of (deduplicated) dependencies;
	// dependents maps a dep id to the ids that require it.
	indeg := make(map[string]int, len(index))
	dependents := make(map[string][]string, len(index))
	for id := range index {
		indeg[id] = 0
	}
	for id, n := range index {
		seen := make(map[string]bool)
		for _, dep := range n.Dependencies() {
			if seen[dep] {
				continue
			}
			seen[dep] = true
			indeg[id]++
			dependents[dep] = append(dependents[dep], id)
		}
	}

	var levels [][]T
	var frontier []string
	for id, d := range indeg {
		if d == 0 {
			frontier = append(frontier, id)
		}
	}

	processed := 0
	for len(frontier) > 0 {
		sortIDs(frontier, index)
		level := make([]T, 0, len(frontier))
		var next []string
		for _, id := range frontier {
			level = append(level, index[id])
			processed++
			for _, dep := range dependents[id] {
				indeg[dep]--
				if indeg[dep] == 0 {
					next = append(next, dep)
				}
			}
		}
		levels = append(levels, level)
		frontier = next
	}

	if processed < len(index) {
		var cyclic []string
		for id, d := range indeg {
			if d > 0 {
				cyclic = append(cyclic, id)
			}
		}
		sort.Strings(cyclic)
		return nil, fmt.Errorf("dag: dependency cycle among %v", cyclic)
	}

	return levels, nil
}

// Flatten returns the nodes of levels in execution order (level 0 first).
func Flatten[T Node](levels [][]T) []T {
	var out []T
	for _, level := range levels {
		out = append(out, level...)
	}
	return out
}

func sortIDs[T Node](ids []string, index map[string]T) {
	sort.Slice(ids, func(i, j int) bool {
		oi, oj := orderOf(index[ids[i]]), orderOf(index[ids[j]])
		if oi != oj {
			return oi < oj
		}
		return ids[i] < ids[j]
	})
}

func orderOf[T Node](n T) int {
	if o, ok := any(n).(orderer); ok {
		return o.Order()
	}
	return 0
}
