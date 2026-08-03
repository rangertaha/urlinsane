// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package cns_test

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/plugins/variant/cns"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/varianttest"
)

// The declaration is what the scheduler binds and what the cache keys on, so a
// wrong id silently reroutes edges and a Version left at zero makes stale
// results look fresh.
func TestSpecDeclaresTheAlgorithm(t *testing.T) {
	s := cns.Spec(varianttest.Langs())
	if s.ID != "cns" {
		t.Errorf("ID = %q, want %q", s.ID, "cns")
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

// digits and their cardinal words swapped both ways.
func TestGeneratesKnownVariants(t *testing.T) {
	got := cns.Spec(varianttest.Langs()).Gen("one2three")
	want := []string{"12three", "onetwothree"}
	if miss := varianttest.Missing(got, want); len(miss) > 0 {
		t.Errorf("Gen(%q) missing %q\ngot %q", "one2three", miss, got)
	}
}

// A generator that returns its input claims the name is a typo of itself, and
// every consumer downstream would have to filter it out again.
func TestNeverReturnsItsInput(t *testing.T) {
	for _, in := range []string{"one2three", "example", "example.com"} {
		for _, g := range cns.Spec(varianttest.Langs()).Gen(in) {
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
	for _, g := range cns.Spec(varianttest.Langs()).Gen("one2three") {
		if seen[g] {
			t.Errorf("Gen(%q) returned %q more than once", "one2three", g)
		}
		seen[g] = true
	}
}

// An empty key reaches the operators when a split yields nothing; it must not
// panic and must not invent a name out of nothing.
func TestEmptyInputIsQuiet(t *testing.T) {
	if got := cns.Spec(varianttest.Langs()).Gen(""); len(got) != 0 {
		t.Errorf("Gen(\"\") = %q, want nothing", got)
	}
}
