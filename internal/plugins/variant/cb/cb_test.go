// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package cb_test

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/plugins/variant/cb"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/varianttest"
)

// The declaration is what the scheduler binds and what the cache keys on, so a
// wrong id silently reroutes edges and a Version left at zero makes stale
// results look fresh.
func TestSpecDeclaresTheAlgorithm(t *testing.T) {
	s := cb.Spec([]string{"login", "secure"})
	if s.ID != "cb" {
		t.Errorf("ID = %q, want %q", s.ID, "cb")
	}
	if s.Title == "" {
		t.Error("Title is empty; --explain and the reports print it")
	}
	if s.Version < 1 {
		t.Errorf("Version = %d, want >= 1; zero cannot invalidate a cached result", s.Version)
	}
	if s.Gen == nil {
		t.Fatal("Gen is nil; the operator would produce nothing")
	}
}

// a brand plus a lure keyword, joined both ways and both orders.
func TestGeneratesKnownVariants(t *testing.T) {
	got := cb.Spec([]string{"login", "secure"}).Gen("example")
	want := []string{"example-login", "loginexample", "example-secure"}
	if miss := varianttest.Missing(got, want); len(miss) > 0 {
		t.Errorf("Gen(%q) missing %q\ngot %q", "example", miss, got)
	}
}

// A generator that returns its input claims the name is a typo of itself, and
// every consumer downstream would have to filter it out again.
func TestNeverReturnsItsInput(t *testing.T) {
	for _, in := range []string{"example", "example", "example.com"} {
		for _, g := range cb.Spec([]string{"login", "secure"}).Gen(in) {
			if g == in {
				t.Errorf("Gen(%q) returned its own input", in)
			}
		}
	}
}

// Duplicates inflate the candidate count and are charged to the scan budget
// twice.
func TestReturnsNoDuplicates(t *testing.T) {
	seen := map[string]bool{}
	for _, g := range cb.Spec([]string{"login", "secure"}).Gen("example") {
		if seen[g] {
			t.Errorf("Gen(%q) returned %q more than once", "example", g)
		}
		seen[g] = true
	}
}

// An empty key reaches the operators when a split yields nothing; it must not
// panic and must not invent a name out of nothing.
func TestEmptyInputIsQuiet(t *testing.T) {
	if got := cb.Spec([]string{"login", "secure"}).Gen(""); len(got) != 0 {
		t.Errorf("Gen(\"\") = %q, want nothing", got)
	}
}
