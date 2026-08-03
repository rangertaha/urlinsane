// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package fsd_test

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/plugins/variant/fsd"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/varianttest"
)

// The declaration is what the scheduler binds and what the cache keys on, so a
// wrong id silently reroutes edges and a Version left at zero makes stale
// results look fresh.
func TestSpecDeclaresTheAlgorithm(t *testing.T) {
	s := fsd.Spec([]string{"github.io"})
	if s.ID != "fsd" {
		t.Errorf("ID = %q, want %q", s.ID, "fsd")
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

// the name re-hosted under a provider that delegates subdomains.
func TestGeneratesKnownVariants(t *testing.T) {
	got := fsd.Spec([]string{"github.io"}).Gen("example.com")
	want := []string{"example.github.io", "example.com.github.io"}
	if miss := varianttest.Missing(got, want); len(miss) > 0 {
		t.Errorf("Gen(%q) missing %q\ngot %q", "example.com", miss, got)
	}
}

// A generator that returns its input claims the name is a typo of itself, and
// every consumer downstream would have to filter it out again.
func TestNeverReturnsItsInput(t *testing.T) {
	for _, in := range []string{"example.com", "example", "example.com"} {
		for _, g := range fsd.Spec([]string{"github.io"}).Gen(in) {
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
	for _, g := range fsd.Spec([]string{"github.io"}).Gen("example.com") {
		if seen[g] {
			t.Errorf("Gen(%q) returned %q more than once", "example.com", g)
		}
		seen[g] = true
	}
}

// An empty key reaches the operators when a split yields nothing; it must not
// panic and must not invent a name out of nothing.
func TestEmptyInputIsQuiet(t *testing.T) {
	if got := fsd.Spec([]string{"github.io"}).Gen(""); len(got) != 0 {
		t.Errorf("Gen(\"\") = %q, want nothing", got)
	}
}
