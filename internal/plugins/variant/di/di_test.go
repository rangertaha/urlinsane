// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package di_test

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/plugins/variant/di"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/varianttest"
)

// The declaration is what the scheduler binds and what the cache keys on, so a
// wrong id silently reroutes edges and a Version left at zero makes stale
// results look fresh.
func TestSpecDeclaresTheAlgorithm(t *testing.T) {
	s := di.Spec()
	if s.ID != "di" {
		t.Errorf("ID = %q, want %q", s.ID, "di")
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

// a dot inserted at each interior gap.
func TestGeneratesKnownVariants(t *testing.T) {
	got := di.Spec().Gen("abc")
	want := []string{"a.bc", "ab.c"}
	if miss := varianttest.Missing(got, want); len(miss) > 0 {
		t.Errorf("Gen(%q) missing %q\ngot %q", "abc", miss, got)
	}
}

// A generator that returns its input claims the name is a typo of itself, and
// every consumer downstream would have to filter it out again.
func TestNeverReturnsItsInput(t *testing.T) {
	for _, in := range []string{"abc", "example", "example.com"} {
		for _, g := range di.Spec().Gen(in) {
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
	for _, g := range di.Spec().Gen("abc") {
		if seen[g] {
			t.Errorf("Gen(%q) returned %q more than once", "abc", g)
		}
		seen[g] = true
	}
}

// An empty key reaches the operators when a split yields nothing; it must not
// panic and must not invent a name out of nothing.
func TestEmptyInputIsQuiet(t *testing.T) {
	if got := di.Spec().Gen(""); len(got) != 0 {
		t.Errorf("Gen(\"\") = %q, want nothing", got)
	}
}
