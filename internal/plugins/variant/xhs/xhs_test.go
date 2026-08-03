// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package xhs_test

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/plugins/variant/varianttest"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/xhs"
)

// The declaration is what the scheduler binds and what the cache keys on, so a
// wrong id silently reroutes edges and a Version left at zero makes stale
// results look fresh.
func TestSpecDeclaresTheAlgorithm(t *testing.T) {
	s := xhs.Spec([][]string{{"foo", "phoo"}})
	if s.ID != "xhs" {
		t.Errorf("ID = %q, want %q", s.ID, "xhs")
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

// one sound spelled the way another language writes it.
func TestGeneratesKnownVariants(t *testing.T) {
	got := xhs.Spec([][]string{{"foo", "phoo"}}).Gen("foo")
	want := []string{"phoo"}
	if miss := varianttest.Missing(got, want); len(miss) > 0 {
		t.Errorf("Gen(%q) missing %q\ngot %q", "foo", miss, got)
	}
}

// A generator that returns its input claims the name is a typo of itself, and
// every consumer downstream would have to filter it out again.
func TestNeverReturnsItsInput(t *testing.T) {
	for _, in := range []string{"foo", "example", "example.com"} {
		for _, g := range xhs.Spec([][]string{{"foo", "phoo"}}).Gen(in) {
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
	for _, g := range xhs.Spec([][]string{{"foo", "phoo"}}).Gen("foo") {
		if seen[g] {
			t.Errorf("Gen(%q) returned %q more than once", "foo", g)
		}
		seen[g] = true
	}
}

// An empty key reaches the operators when a split yields nothing; it must not
// panic and must not invent a name out of nothing.
func TestEmptyInputIsQuiet(t *testing.T) {
	if got := xhs.Spec([][]string{{"foo", "phoo"}}).Gen(""); len(got) != 0 {
		t.Errorf("Gen(\"\") = %q, want nothing", got)
	}
}
