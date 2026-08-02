// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package plugins

import (
	"context"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

type stubOp struct{ id string }

func (o stubOp) Id() string             { return o.id }
func (o stubOp) Version() int           { return 1 }
func (o stubOp) Trigger() graph.Trigger { return graph.Trigger{} }
func (o stubOp) Emits() graph.Effects   { return graph.Effects{} }
func (o stubOp) Resource() string       { return "" }
func (o stubOp) Exec(context.Context, graph.View) (graph.Delta, graph.Outcome) {
	return graph.Delta{}, graph.Outcome{}
}

type stubAnalyzer struct{ id string }

func (a stubAnalyzer) Id() string { return a.id }
func (a stubAnalyzer) Exec(context.Context, *graph.Analysis) ([]graph.Finding, error) {
	return nil, nil
}

// stubSource stands in for the config file.
type stubSource map[string]map[string]any

func (s stubSource) Apply(id string) map[string]any { return s[id] }

func TestOperatorsBuildInIDOrder(t *testing.T) {
	reset()
	AddOperator("zeta", nil, func(Env) ([]graph.Operator, error) {
		return []graph.Operator{stubOp{"zeta"}}, nil
	})
	AddOperator("alpha", nil, func(Env) ([]graph.Operator, error) {
		return []graph.Operator{stubOp{"alpha"}}, nil
	})

	ops, err := Operators(nil, Env{})
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 2 || ops[0].Id() != "alpha" || ops[1].Id() != "zeta" {
		t.Fatalf("got %v, want alpha then zeta — the plan is hashed, so order cannot "+
			"depend on registration order", ops)
	}
}

// One plugin may contribute several operators, as a family does.
func TestOneOperatorPluginMayReturnSeveral(t *testing.T) {
	reset()
	AddOperator("family", nil, func(Env) ([]graph.Operator, error) {
		return []graph.Operator{stubOp{"a"}, stubOp{"b"}}, nil
	})
	ops, _ := Operators(nil, Env{})
	if len(ops) != 2 {
		t.Fatalf("got %d operators, want 2", len(ops))
	}
}

func TestSettingsReachThePlugin(t *testing.T) {
	reset()
	var got Settings
	AddOperator("dns", map[string]any{"timeout": 5, "retries": 2},
		func(e Env) ([]graph.Operator, error) {
			got = e.Settings
			return nil, nil
		})

	// No source: the plugin sees its declared defaults.
	if _, err := Operators(nil, Env{}); err != nil {
		t.Fatal(err)
	}
	if got.Int("timeout", 0) != 5 || got.Int("retries", 0) != 2 {
		t.Fatalf("defaults did not reach the plugin: %v", got)
	}

	// With a source: overrides apply per key, and the rest survive.
	src := stubSource{"dns": {"timeout": 30, "retries": 2}}
	if _, err := Operators(src, Env{}); err != nil {
		t.Fatal(err)
	}
	if got.Int("timeout", 0) != 30 {
		t.Errorf("override did not reach the plugin: %v", got)
	}
	if got.Int("retries", 0) != 2 {
		t.Errorf("overriding one key dropped another: %v", got)
	}
}

// A plugin that cannot build must stop the run, naming itself. A scan missing
// an operator it was configured to use is a smaller scan reported as complete.
func TestOperatorErrorNamesThePlugin(t *testing.T) {
	reset()
	AddOperator("broken", nil, func(Env) ([]graph.Operator, error) {
		return nil, context.DeadlineExceeded
	})
	_, err := Operators(nil, Env{})
	if err == nil {
		t.Fatal("a failing plugin did not stop the build")
	}
	if want := "broken"; !contains(err.Error(), want) {
		t.Errorf("error %q does not name the plugin", err)
	}
}

func TestDuplicateRegistrationPanics(t *testing.T) {
	reset()
	AddOperator("dup", nil, func(Env) ([]graph.Operator, error) { return nil, nil })
	defer func() {
		if recover() == nil {
			t.Error("registering an id twice did not panic; the second would " +
				"collide in the plan and the cache key")
		}
	}()
	AddOperator("dup", nil, func(Env) ([]graph.Operator, error) { return nil, nil })
}

func TestAnalyzersAndAlgorithms(t *testing.T) {
	reset()
	AddAnalyzer("scoring", nil, func(Env) (graph.Analyzer, error) {
		return stubAnalyzer{"scoring"}, nil
	})
	as, err := Analyzers(nil, Env{})
	if err != nil || len(as) != 1 || as[0].Id() != "scoring" {
		t.Fatalf("analyzers = %v, err %v", as, err)
	}

	AddSpec(variant.Spec{ID: "zz", Title: "Z"})
	AddSpec(variant.Spec{ID: "aa", Title: "A"})
	got, err := Algorithms(nil, Env{})
	if err != nil || len(got) != 2 || got[0].ID != "aa" {
		t.Fatalf("algorithms = %v, err %v; want aa first", got, err)
	}
}

// Settings coerce the types YAML produces rather than erroring on them.
func TestSettingsAccessorsCoerce(t *testing.T) {
	s := Settings{
		"n":    float64(30), // YAML numbers decode as float64
		"list": []any{"a", "b"},
		"flag": true,
		"name": "x",
	}
	if s.Int("n", 0) != 30 {
		t.Error("float64 not accepted as an int setting")
	}
	if got := s.Strings("list"); len(got) != 2 || got[0] != "a" {
		t.Errorf("[]any not accepted as a string list: %v", got)
	}
	if !s.Bool("flag", false) || s.String("name", "") != "x" {
		t.Error("bool or string accessor failed")
	}
	// Absent and wrong-typed keys fall back rather than zeroing.
	if s.Int("missing", 7) != 7 || s.String("n", "def") != "def" {
		t.Error("fallback not used for absent or mistyped keys")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
