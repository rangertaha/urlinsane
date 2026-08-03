// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// A package key is registry-qualified and a repo key is forge-qualified. The
// qualifier says which namespace the name lives in, so an algorithm that varies
// it is not generating a squat of the name -- it is generating a name in a
// registry or on a forge that does not exist.
//
// This is a whole-class guard rather than a test of one algorithm. Three
// separate operators had it wrong (afx, nsc and sep, all via Whole:true) and
// nothing failed, because every one of them produced plausible-looking strings.
// Any new algorithm that reaches for the whole key fails here instead.
func TestNoAlgorithmVariesTheQualifier(t *testing.T) {
	for _, seed := range []struct{ typ, key, qualifier string }{
		{variant.TypePackage, "npm:lodash", "npm:"},
		{variant.TypePackage, "pypi:requests", "pypi:"},
		{variant.TypePackage, "npm:@acme/tool", "npm:"},
		{variant.TypeRepo, "github.com/acme/tool", "github.com/"},
	} {
		var total int
		for _, op := range All(testOptions()) {
			if !handles(op, seed.typ) {
				continue
			}
			_, _, d, _ := run(t, op, seed.typ, seed.key)
			for _, e := range d.Edges {
				if e.Rel != graph.VariantRel {
					continue
				}
				total++
				if !strings.HasPrefix(e.To.Key, seed.qualifier) {
					t.Errorf("%s on %s %q produced %q, which is outside %q",
						op.Id(), seed.typ, seed.key, e.To.Key, seed.qualifier)
				}
			}
		}
		// A guard that silently stopped generating anything would pass while
		// testing nothing.
		if total == 0 {
			t.Errorf("%s %q produced no variants at all", seed.typ, seed.key)
		}
	}
}

// handles reports whether the operator's trigger admits this node type, so the
// guard exercises the operators the scheduler would actually run and not the
// domain-only ones.
func handles(op graph.Operator, nodeType string) bool {
	types := op.Trigger().On.Types
	if len(types) == 0 {
		return true // bound by capability: every Nameable type
	}
	for _, t := range types {
		if t == nodeType {
			return true
		}
	}
	return false
}
