// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package variant

import "testing"

// A package key is registry-qualified and a repo key is forge-qualified, so the
// whole key is not the name. Varying it produced "pm:lodash" -- a package in a
// registry that does not exist -- from every Whole:false algorithm, for every
// package and repo the tool was pointed at.
func TestQualifiersAreNotVaried(t *testing.T) {
	for _, tc := range []struct {
		typ, key, core, prefix string
	}{
		{TypePackage, "npm:lodash", "lodash", "npm:"},
		{TypePackage, "pypi:requests", "requests", "pypi:"},
		{TypePackage, "npm:@acme/tool", "@acme/tool", "npm:"},
		{TypeRepo, "github.com/acme/tool", "acme/tool", "github.com/"},
	} {
		p := DefaultSplit(tc.typ, tc.key)
		if p.Core != tc.core {
			t.Errorf("%s %q: core = %q, want %q", tc.typ, tc.key, p.Core, tc.core)
		}
		if p.Prefix != tc.prefix {
			t.Errorf("%s %q: prefix = %q, want %q", tc.typ, tc.key, p.Prefix, tc.prefix)
		}
		// The qualifier comes back on a varied core, unchanged.
		if got, want := p.Join("x"), tc.prefix+"x"; got != want {
			t.Errorf("%s %q: Join(x) = %q, want %q", tc.typ, tc.key, got, want)
		}
		// Round-tripping an unvaried core is the identity, or the algorithm's
		// own dedupe against the origin key stops working.
		if got := p.Join(p.Core); got != tc.key {
			t.Errorf("%s: Join(Core) = %q, want the original %q", tc.typ, got, tc.key)
		}
	}
}

// A key that never went through the decomposer has no qualifier to strip, and
// guessing one would be worse than varying it whole.
func TestUnqualifiedKeysAreVariedWhole(t *testing.T) {
	for _, tc := range []struct{ typ, key string }{
		{TypePackage, "lodash"},  // no registry
		{TypePackage, "npm:"},    // no name
		{TypePackage, ":lodash"}, // no registry
		{TypeRepo, "acme"},       // no path
		{TypeUsername, "acme"},   // never qualified
	} {
		p := DefaultSplit(tc.typ, tc.key)
		if p.Core != tc.key {
			t.Errorf("%s %q: core = %q, want the whole key", tc.typ, tc.key, p.Core)
		}
	}
}

// An emptied core is not a package. Omission of a one-character name must not
// leave a bare "npm:", which would canonicalize as a package with no name.
func TestAnEmptiedCoreDoesNotLeaveABareQualifier(t *testing.T) {
	for _, tc := range []struct{ typ, key string }{
		{TypePackage, "npm:a"},
		{TypeRepo, "github.com/a"},
	} {
		if got := DefaultSplit(tc.typ, tc.key).Join(""); got != "" {
			t.Errorf("%s %q: Join(\"\") = %q, want empty", tc.typ, tc.key, got)
		}
	}
}
