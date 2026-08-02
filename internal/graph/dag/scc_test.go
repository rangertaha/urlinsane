// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package dag

import (
	"reflect"
	"testing"
)

func TestCondenseLinear(t *testing.T) {
	g := Graph{"a": {"b"}, "b": {"c"}, "c": nil}
	got := Condense(g)
	want := [][]string{{"a"}, {"b"}, {"c"}}
	if len(got) != 3 {
		t.Fatalf("components = %d, want 3", len(got))
	}
	for i, c := range got {
		if !reflect.DeepEqual(c.Members, want[i]) {
			t.Fatalf("component %d = %v, want %v", i, c.Members, want[i])
		}
		if c.Layer != i {
			t.Fatalf("component %v at layer %d, want %d", c.Members, c.Layer, i)
		}
		if c.Cyclic {
			t.Fatalf("component %v marked cyclic", c.Members)
		}
	}
}

func TestCondenseCollapsesCycle(t *testing.T) {
	// The type flow this exists for: domain -> ip -> domain is correct data,
	// not a defect, and cannot be topologically sorted without condensing.
	g := Graph{
		"domain":     {"ip"},
		"ip":         {"domain", "asn"},
		"asn":        nil,
		"registrant": nil,
		"domain2":    {"domain"},
	}
	got := Condense(g)

	var cycle *Component
	for i := range got {
		if len(got[i].Members) > 1 {
			cycle = &got[i]
		}
	}
	if cycle == nil {
		t.Fatalf("no component collapsed the domain/ip cycle: %+v", got)
	}
	if !reflect.DeepEqual(cycle.Members, []string{"domain", "ip"}) {
		t.Fatalf("cycle members = %v, want [domain ip]", cycle.Members)
	}
	if !cycle.Cyclic {
		t.Fatal("multi-member component not marked cyclic")
	}

	// Every edge must run forward across layers, which is the property that
	// makes the rendering readable.
	layer := map[string]int{}
	for _, c := range got {
		for _, m := range c.Members {
			layer[m] = c.Layer
		}
	}
	for from, outs := range g {
		for _, to := range outs {
			if layer[to] < layer[from] {
				t.Fatalf("edge %s->%s runs backwards: layers %d -> %d", from, to, layer[from], layer[to])
			}
		}
	}
}

func TestCondenseSelfLoopIsCyclic(t *testing.T) {
	// A variant operator emits its own type, so domain -> domain is normal.
	g := Graph{"domain": {"domain"}}
	got := Condense(g)
	if len(got) != 1 || !got[0].Cyclic {
		t.Fatalf("self-loop not reported as cyclic: %+v", got)
	}
}

func TestCondenseIsDeterministic(t *testing.T) {
	g := Graph{
		"a": {"b", "c"}, "b": {"d"}, "c": {"d"}, "d": {"a"},
		"e": {"f"}, "f": nil, "g": nil,
	}
	first := Condense(g)
	for i := 0; i < 20; i++ {
		got := Condense(g)
		if !reflect.DeepEqual(got, first) {
			t.Fatalf("condensation varied between runs:\n%+v\nvs\n%+v", got, first)
		}
	}
}

func TestReachable(t *testing.T) {
	g := Graph{"domain": {"ip"}, "ip": {"asn"}, "asn": nil, "package": {"registry"}}
	got := Reachable(g, "domain")
	for _, want := range []string{"domain", "ip", "asn"} {
		if !got[want] {
			t.Fatalf("%s not reachable from domain", want)
		}
	}
	if got["package"] || got["registry"] {
		t.Fatal("package flow reachable from a domain seed")
	}
}

func TestReachableFollowsCycles(t *testing.T) {
	g := Graph{"domain": {"ip"}, "ip": {"domain"}}
	got := Reachable(g, "domain")
	if len(got) != 2 {
		t.Fatalf("reachable = %v, want both nodes without looping forever", got)
	}
}

func TestCondenseEmpty(t *testing.T) {
	if got := Condense(Graph{}); len(got) != 0 {
		t.Fatalf("empty graph produced %d components", len(got))
	}
}
