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

package dag

import (
	"reflect"
	"testing"
)

type node struct {
	id    string
	deps  []string
	order int
}

func (n node) Id() string             { return n.id }
func (n node) Dependencies() []string { return n.deps }
func (n node) Order() int             { return n.order }

// ids extracts the node ids per level for easy comparison.
func ids(levels [][]node) [][]string {
	out := make([][]string, len(levels))
	for i, lvl := range levels {
		for _, n := range lvl {
			out[i] = append(out[i], n.id)
		}
	}
	return out
}

func TestLevels_Empty(t *testing.T) {
	levels, err := Levels([]node{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(levels) != 0 {
		t.Fatalf("want 0 levels, got %v", ids(levels))
	}
}

func TestLevels_Linear(t *testing.T) {
	// c <- b <- a (a runs first)
	nodes := []node{
		{id: "c", deps: []string{"b"}},
		{id: "b", deps: []string{"a"}},
		{id: "a"},
	}
	levels, err := Levels(nodes, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := [][]string{{"a"}, {"b"}, {"c"}}
	if got := ids(levels); !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestLevels_Diamond(t *testing.T) {
	// d depends on b,c ; b,c depend on a
	nodes := []node{
		{id: "d", deps: []string{"b", "c"}},
		{id: "b", deps: []string{"a"}},
		{id: "c", deps: []string{"a"}},
		{id: "a"},
	}
	levels, err := Levels(nodes, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := [][]string{{"a"}, {"b", "c"}, {"d"}}
	if got := ids(levels); !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

// TestLevels_RealCollectorGraph mirrors urlinsane's actual collector graph:
// ip has no deps; geo/ptr/wi/bn all depend on ip.
func TestLevels_RealCollectorGraph(t *testing.T) {
	nodes := []node{
		{id: "geo", deps: []string{"ip"}},
		{id: "ptr", deps: []string{"ip"}},
		{id: "wi", deps: []string{"ip"}},
		{id: "bn", deps: []string{"ip"}},
		{id: "ip"},
		{id: "cn"},
		{id: "mx"},
		{id: "ns"},
		{id: "txt"},
		{id: "idn"},
	}
	levels, err := Levels(nodes, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(levels) != 2 {
		t.Fatalf("want 2 levels, got %d: %v", len(levels), ids(levels))
	}
	// ip must be in level 0; geo/ptr/wi/bn strictly after it.
	level0 := map[string]bool{}
	for _, n := range levels[0] {
		level0[n.id] = true
	}
	if !level0["ip"] {
		t.Fatalf("ip must be in level 0, got %v", ids(levels)[0])
	}
	for _, dep := range []string{"geo", "ptr", "wi", "bn"} {
		if level0[dep] {
			t.Fatalf("%s must run after ip (not level 0)", dep)
		}
	}
}

func TestLevels_TieBreakByOrder(t *testing.T) {
	// Same level, distinct Order — lower Order first, then Id.
	nodes := []node{
		{id: "z", order: 1},
		{id: "a", order: 5},
		{id: "m", order: 1},
	}
	levels, err := Levels(nodes, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := [][]string{{"m", "z", "a"}} // order 1 (m,z by id), then order 5 (a)
	if got := ids(levels); !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestLevels_Cycle(t *testing.T) {
	nodes := []node{
		{id: "a", deps: []string{"b"}},
		{id: "b", deps: []string{"a"}},
	}
	if _, err := Levels(nodes, nil); err == nil {
		t.Fatal("expected a cycle error, got nil")
	}
}

func TestLevels_SelfCycle(t *testing.T) {
	nodes := []node{{id: "a", deps: []string{"a"}}}
	if _, err := Levels(nodes, nil); err == nil {
		t.Fatal("expected a self-cycle error, got nil")
	}
}

func TestLevels_MissingDepNoResolve(t *testing.T) {
	nodes := []node{{id: "wi", deps: []string{"ip"}}}
	if _, err := Levels(nodes, nil); err == nil {
		t.Fatal("expected a missing-dependency error, got nil")
	}
}

func TestLevels_MissingDepClosure(t *testing.T) {
	// Only wi is selected; its ip dependency is auto-included via resolve.
	registry := map[string]node{
		"ip": {id: "ip"},
		"wi": {id: "wi", deps: []string{"ip"}},
	}
	resolve := func(id string) (node, bool) {
		n, ok := registry[id]
		return n, ok
	}
	levels, err := Levels([]node{{id: "wi", deps: []string{"ip"}}}, resolve)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := [][]string{{"ip"}, {"wi"}}
	if got := ids(levels); !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestLevels_UnknownDepClosure(t *testing.T) {
	resolve := func(id string) (node, bool) { return node{}, false }
	if _, err := Levels([]node{{id: "wi", deps: []string{"ip"}}}, resolve); err == nil {
		t.Fatal("expected an unknown-dependency error, got nil")
	}
}

func TestLevels_Deterministic(t *testing.T) {
	nodes := []node{
		{id: "d", deps: []string{"b", "c"}},
		{id: "b", deps: []string{"a"}},
		{id: "c", deps: []string{"a"}},
		{id: "a"},
		{id: "e"},
	}
	var first [][]string
	for i := 0; i < 50; i++ {
		levels, err := Levels(nodes, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := ids(levels)
		if first == nil {
			first = got
			continue
		}
		if !reflect.DeepEqual(got, first) {
			t.Fatalf("non-deterministic: run %d gave %v, want %v", i, got, first)
		}
	}
}

func TestFlatten(t *testing.T) {
	nodes := []node{
		{id: "b", deps: []string{"a"}},
		{id: "a"},
	}
	levels, err := Levels(nodes, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	flat := Flatten(levels)
	if len(flat) != 2 || flat[0].id != "a" || flat[1].id != "b" {
		t.Fatalf("unexpected flatten: %v", flat)
	}
}
