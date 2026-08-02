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

package graph

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// testRegistry mirrors the shape of the real registry: a nameable domain, an
// observed ip, and one relation of each class.
func testRegistry(t *testing.T) *Registry {
	t.Helper()
	r := NewRegistry()

	lower := func(s string) (string, error) {
		s = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(s, ".")))
		if s == "" {
			return "", fmt.Errorf("empty key")
		}
		return s, nil
	}

	if _, err := r.AddType(NodeTypeDef{
		Name: "domain", Cap: Nameable, Version: 1, Canonical: lower,
		Fields: []FieldDef{
			{Name: "punycode", Kind: KindString},
			{Name: "rank", Kind: KindInt},
			{Name: "created", Kind: KindTime, Merge: Precedence("rdap", "whois")},
		},
	}); err != nil {
		t.Fatalf("domain: %v", err)
	}
	if _, err := r.AddType(NodeTypeDef{
		Name: "tld", Cap: Observed, Version: 1, Canonical: lower,
	}); err != nil {
		t.Fatalf("tld: %v", err)
	}
	if _, err := r.AddType(NodeTypeDef{
		Name: "ip", Cap: Observed, Version: 1, Canonical: lower,
		Fields: []FieldDef{{Name: "asn", Kind: KindString}},
	}); err != nil {
		t.Fatalf("ip: %v", err)
	}

	for _, d := range []RelDef{
		{Name: "TLD_OF", Class: Structural, Version: 1},
		{Name: "VARIANT_OF", Class: Variant, Version: 1, Fields: []FieldDef{
			{Name: "algorithm", Kind: KindString},
			{Name: "distance", Kind: KindInt},
		}},
		{Name: "RESOLVES_TO", Class: Observation, Version: 1},
	} {
		if _, err := r.AddRel(d); err != nil {
			t.Fatalf("%s: %v", d.Name, err)
		}
	}
	return r
}

func seeded(t *testing.T) (*Graph, NodeID) {
	t.Helper()
	g := New(testRegistry(t))
	id, err := g.Seed("domain", "Example.COM")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	return g, id
}

func op(name string) Provenance { return Provenance{Operator: name, Round: 1} }

// --- canonicalization and identity -----------------------------------------

func TestCanonicalizationConverges(t *testing.T) {
	g, seed := seeded(t)
	if n, _ := g.Node(seed); n.Key != "example.com" {
		t.Fatalf("seed key = %q, want example.com", n.Key)
	}

	// Two operators reach the same entity by different spellings. Convergence
	// on one node is the property the whole design rests on.
	res := g.Apply(op("a"), seed, Delta{Nodes: []NodeRef{{Type: "domain", Key: "Foo.COM"}}})
	res2 := g.Apply(op("b"), seed, Delta{Nodes: []NodeRef{{Type: "domain", Key: "foo.com."}}})
	if len(res.Nodes) != 1 || len(res2.Nodes) != 1 {
		t.Fatalf("want one node each, got %d and %d", len(res.Nodes), len(res2.Nodes))
	}
	if res.Nodes[0] != res2.Nodes[0] {
		t.Fatal("Foo.COM and foo.com. produced different NodeIDs")
	}
	if got := len(g.Nodes()); got != 2 {
		t.Fatalf("graph has %d nodes, want 2 (seed + foo.com)", got)
	}
}

func TestCanonicalizationFailureIsRefused(t *testing.T) {
	g, seed := seeded(t)
	res := g.Apply(op("a"), seed, Delta{Nodes: []NodeRef{{Type: "domain", Key: "   "}}})
	if len(res.Nodes) != 0 {
		t.Fatal("admitted a node whose key could not be canonicalized")
	}
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != RejectCanonical {
		t.Fatalf("rejections = %+v, want one RejectCanonical", res.Rejected)
	}
}

func TestNodeIDStableAsPropsAccumulate(t *testing.T) {
	g, seed := seeded(t)
	n, _ := g.Node(seed)
	before := n.ID
	_, cidBefore, err := n.Addressed()
	if err != nil {
		t.Fatalf("addressed: %v", err)
	}

	g.Apply(op("a"), seed, Delta{Props: []PropSet{
		{Node: &NodeRef{Type: "domain", Key: "example.com"}, Field: "rank", Value: Int(42)},
	}})

	n, _ = g.Node(seed)
	if n.ID != before {
		t.Fatal("NodeID changed when props accumulated; identity must be stable")
	}
	_, cidAfter, err := n.Addressed()
	if err != nil {
		t.Fatalf("addressed: %v", err)
	}
	if cidBefore.Equals(cidAfter) {
		t.Fatal("CID did not change when props changed; it must address content")
	}
}

// --- content addressing ------------------------------------------------------

func TestCIDStableAcrossAssertionOrder(t *testing.T) {
	// Same facts, asserted in opposite orders by different operators. The
	// content address must not depend on arrival order.
	build := func(forward bool) string {
		g, seed := seeded(t)
		ref := NodeRef{Type: "domain", Key: "example.com"}
		a := PropSet{Node: &ref, Field: "punycode", Value: String("example.com")}
		b := PropSet{Node: &ref, Field: "rank", Value: Int(7)}
		if forward {
			g.Apply(op("x"), seed, Delta{Props: []PropSet{a, b}})
		} else {
			g.Apply(op("x"), seed, Delta{Props: []PropSet{b, a}})
		}
		n, _ := g.Node(seed)
		_, c, err := n.Addressed()
		if err != nil {
			t.Fatalf("addressed: %v", err)
		}
		return c.String()
	}
	if f, r := build(true), build(false); f != r {
		t.Fatalf("CID depends on assertion order: %s vs %s", f, r)
	}
}

func TestCIDStableAcrossRuns(t *testing.T) {
	// The reproducibility property in miniature: two independent graphs built
	// identically must produce identical content addresses.
	var first string
	for i := 0; i < 3; i++ {
		g, seed := seeded(t)
		ref := NodeRef{Type: "domain", Key: "example.com"}
		g.Apply(op("x"), seed, Delta{Props: []PropSet{
			{Node: &ref, Field: "punycode", Value: String("xn--fsq.com")},
			{Node: &ref, Field: "rank", Value: Int(1)},
		}})
		n, _ := g.Node(seed)
		_, c, err := n.Addressed()
		if err != nil {
			t.Fatalf("addressed: %v", err)
		}
		if i == 0 {
			first = c.String()
		} else if c.String() != first {
			t.Fatalf("run %d CID %s != %s", i, c, first)
		}
	}
}

func TestDistinctPropsDistinctCID(t *testing.T) {
	mk := func(rank int64) string {
		g, seed := seeded(t)
		ref := NodeRef{Type: "domain", Key: "example.com"}
		g.Apply(op("x"), seed, Delta{Props: []PropSet{{Node: &ref, Field: "rank", Value: Int(rank)}}})
		n, _ := g.Node(seed)
		_, c, _ := n.Addressed()
		return c.String()
	}
	if mk(1) == mk(2) {
		t.Fatal("different props produced the same CID")
	}
}

func TestEdgeAddressed(t *testing.T) {
	g, seed := seeded(t)
	res := g.Apply(op("algo"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "examp1e.com"},
	}}})
	if len(res.Edges) != 1 {
		t.Fatalf("edges = %d, want 1", len(res.Edges))
	}
	e, _ := g.Edge(res.Edges[0])
	block, c, err := e.Addressed()
	if err != nil {
		t.Fatalf("addressed: %v", err)
	}
	if len(block) == 0 || !c.Defined() {
		t.Fatal("edge produced no addressed form")
	}
}

// --- props -------------------------------------------------------------------

func TestPropsPositionalOrder(t *testing.T) {
	g, seed := seeded(t)
	ref := NodeRef{Type: "domain", Key: "example.com"}
	// Assert the last field first; Each must still yield declaration order.
	g.Apply(op("x"), seed, Delta{Props: []PropSet{
		{Node: &ref, Field: "rank", Value: Int(3)},
		{Node: &ref, Field: "punycode", Value: String("p")},
	}})
	n, _ := g.Node(seed)
	var got []string
	n.Props.Each(func(f Field, _ Value) { got = append(got, f.Name()) })
	if len(got) != 2 || got[0] != "punycode" || got[1] != "rank" {
		t.Fatalf("Each order = %v, want [punycode rank]", got)
	}
}

func TestPropKindMismatchRejected(t *testing.T) {
	g, seed := seeded(t)
	ref := NodeRef{Type: "domain", Key: "example.com"}
	res := g.Apply(op("x"), seed, Delta{Props: []PropSet{
		{Node: &ref, Field: "rank", Value: String("not an int")},
	}})
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != RejectKindMismatch {
		t.Fatalf("rejections = %+v, want one RejectKindMismatch", res.Rejected)
	}
}

func TestUnknownFieldRejected(t *testing.T) {
	g, seed := seeded(t)
	ref := NodeRef{Type: "domain", Key: "example.com"}
	res := g.Apply(op("x"), seed, Delta{Props: []PropSet{
		{Node: &ref, Field: "nonesuch", Value: Int(1)},
	}})
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != RejectUnknownField {
		t.Fatalf("rejections = %+v, want one RejectUnknownField", res.Rejected)
	}
}

func TestMergePrecedenceIndependentOfOrder(t *testing.T) {
	// whois and rdap both assert `created`; precedence says rdap wins whichever
	// answered first. Under concurrent dispatch, arrival order is network
	// timing, so a last-write-wins rule would be nondeterministic.
	run := func(firstOp, secondOp string) string {
		g, seed := seeded(t)
		ref := NodeRef{Type: "domain", Key: "example.com"}
		g.Apply(op(firstOp), seed, Delta{Props: []PropSet{
			{Node: &ref, Field: "created", Value: Time(time.Unix(1000, 0))},
		}})
		g.Apply(op(secondOp), seed, Delta{Props: []PropSet{
			{Node: &ref, Field: "created", Value: Time(time.Unix(2000, 0))},
		}})
		n, _ := g.Node(seed)
		f, _ := n.Type.Field("created")
		winner, _ := n.Props.Setter(f)
		return winner
	}
	if got := run("whois", "rdap"); got != "rdap" {
		t.Fatalf("whois-then-rdap winner = %q, want rdap", got)
	}
	if got := run("rdap", "whois"); got != "rdap" {
		t.Fatalf("rdap-then-whois winner = %q, want rdap", got)
	}
}

func TestLosingAssertionIsRetained(t *testing.T) {
	g, seed := seeded(t)
	ref := NodeRef{Type: "domain", Key: "example.com"}
	g.Apply(op("rdap"), seed, Delta{Props: []PropSet{{Node: &ref, Field: "created", Value: Time(time.Unix(1, 0))}}})
	g.Apply(op("whois"), seed, Delta{Props: []PropSet{{Node: &ref, Field: "created", Value: Time(time.Unix(2, 0))}}})

	as := g.Assertions(seed)
	if len(as) != 2 {
		t.Fatalf("assertions = %d, want 2 — disagreement is signal and must be kept", len(as))
	}
	var won int
	for _, a := range as {
		if a.Won {
			won++
		}
	}
	if won != 1 {
		t.Fatalf("%d winning assertions, want exactly 1", won)
	}
}

// --- applier invariants ------------------------------------------------------

func TestVariantFromOutsideClosureRejected(t *testing.T) {
	g, seed := seeded(t)

	// An IP discovered by observation, then a PTR domain hanging off it. That
	// domain is Nameable but outside the seed closure.
	g.Apply(op("dns"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "RESOLVES_TO",
		To:   NodeRef{Type: "ip", Key: "1.2.3.4"},
	}}})
	res := g.Apply(op("ptr"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "ip", Key: "1.2.3.4"},
		Rel:  "RESOLVES_TO",
		To:   NodeRef{Type: "domain", Key: "parked.example.net"},
	}}})
	if len(res.Edges) != 1 {
		t.Fatalf("ptr edge not admitted: %+v", res.Rejected)
	}

	// Now a variant operator tries to root on it. This is the explosion path
	// the closure rule exists to close.
	before := len(g.Nodes())
	res = g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{
		{From: NodeRef{Type: "domain", Key: "parked.example.net"}, Rel: "VARIANT_OF",
			To: NodeRef{Type: "domain", Key: "prked.example.net"}},
		{From: NodeRef{Type: "domain", Key: "parked.example.net"}, Rel: "VARIANT_OF",
			To: NodeRef{Type: "domain", Key: "paked.example.net"}},
	}})
	if len(res.Edges) != 0 {
		t.Fatal("admitted a VARIANT_OF edge rooted outside the seed closure")
	}
	if len(res.Rejected) != 2 || res.Rejected[0].Kind != RejectClosure {
		t.Fatalf("rejections = %+v, want two RejectClosure", res.Rejected)
	}

	// And the targets must not be admitted either. Rejecting the edge while
	// keeping its endpoint would leave the variants in the graph unrooted: the
	// invariant exists to bound expansion, and a check that admits every
	// candidate before refusing its edge bounds nothing. A real algorithm emits
	// thousands of these per origin.
	if got := len(g.Nodes()); got != before {
		t.Fatalf("nodes %d -> %d: a rejected variant edge admitted its target anyway", before, got)
	}
}

func TestVariantInsideClosureAccepted(t *testing.T) {
	g, seed := seeded(t)
	res := g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exmple.com"},
	}}})
	if len(res.Edges) != 1 || len(res.Rejected) != 0 {
		t.Fatalf("edges=%d rejected=%+v, want 1 and none", len(res.Edges), res.Rejected)
	}
}

func TestClosureFollowsStructuralOnly(t *testing.T) {
	g, seed := seeded(t)
	g.Apply(op("decompose"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "TLD_OF",
		To:   NodeRef{Type: "tld", Key: "com"},
	}}})
	g.Apply(op("dns"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "RESOLVES_TO",
		To:   NodeRef{Type: "ip", Key: "1.2.3.4"},
	}}})

	tldID := newNodeID("tld", "com")
	ipID := newNodeID("ip", "1.2.3.4")
	if !g.InClosure(tldID) {
		t.Fatal("structural edge did not extend the seed closure")
	}
	if g.InClosure(ipID) {
		t.Fatal("observation edge extended the seed closure")
	}
}

func TestDepthCountsObservationHopsOnly(t *testing.T) {
	g, seed := seeded(t)
	g.Apply(op("decompose"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "TLD_OF",
		To:   NodeRef{Type: "tld", Key: "com"},
	}}})
	g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exmple.com"},
	}}})
	variant := newNodeID("domain", "exmple.com")
	g.Apply(op("dns"), variant, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "exmple.com"},
		Rel:  "RESOLVES_TO",
		To:   NodeRef{Type: "ip", Key: "1.2.3.4"},
	}}})

	if d := g.Depth(newNodeID("tld", "com")); d != 0 {
		t.Fatalf("structural hop cost depth: tld at %d, want 0", d)
	}
	if d := g.Depth(variant); d != 0 {
		t.Fatalf("variant hop cost depth: variant at %d, want 0", d)
	}
	if d := g.Depth(newNodeID("ip", "1.2.3.4")); d != 1 {
		t.Fatalf("ip at depth %d, want 1", d)
	}
}

func TestDeclinedCandidateStaysDeclined(t *testing.T) {
	// "Pruning is irreversible" is only true if the ledger denies re-admission.
	// Otherwise whether a node exists depends on how many operators found it.
	g, seed := seeded(t)
	if err := g.Decline("domain", "Pruned.COM", 1, 0.01, ReasonBelief, op("engine")); err != nil {
		t.Fatalf("decline: %v", err)
	}
	res := g.Apply(op("other"), seed, Delta{Nodes: []NodeRef{{Type: "domain", Key: "pruned.com"}}})
	if len(res.Nodes) != 0 {
		t.Fatal("a declined candidate was re-admitted by a second operator")
	}
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != RejectDenied {
		t.Fatalf("rejections = %+v, want one RejectDenied", res.Rejected)
	}

	rows := g.Ledger()
	if len(rows) != 1 || rows[0].Key != "pruned.com" || rows[0].Reason != ReasonBelief {
		t.Fatalf("ledger = %+v, want one belief row for pruned.com", rows)
	}
}

func TestClosureRejectionDoesNotDenyTheNode(t *testing.T) {
	// A closure violation is an invariant rejection, not a truncation. The node
	// itself may still be admitted legitimately by another edge.
	g, seed := seeded(t)
	g.Apply(op("dns"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "RESOLVES_TO",
		To:   NodeRef{Type: "ip", Key: "1.2.3.4"},
	}}})
	g.Apply(op("ptr"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "ip", Key: "1.2.3.4"},
		Rel:  "RESOLVES_TO",
		To:   NodeRef{Type: "domain", Key: "other.net"},
	}}})
	g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "other.net"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "othr.net"},
	}}})

	if _, ok := g.Node(newNodeID("domain", "other.net")); !ok {
		t.Fatal("closure rejection removed a legitimately admitted node")
	}
	if len(g.Ledger()) != 0 {
		t.Fatal("a closure rejection wrote a truncation-ledger row; it must not deny the candidate")
	}
}

func TestUnknownTypeAndRelationRejected(t *testing.T) {
	g, seed := seeded(t)
	res := g.Apply(op("x"), seed, Delta{Nodes: []NodeRef{{Type: "nosuch", Key: "k"}}})
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != RejectUnknownType {
		t.Fatalf("rejections = %+v, want RejectUnknownType", res.Rejected)
	}
	res = g.Apply(op("x"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "NOSUCH",
		To:   NodeRef{Type: "domain", Key: "b.com"},
	}}})
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != RejectUnknownRel {
		t.Fatalf("rejections = %+v, want RejectUnknownRel", res.Rejected)
	}
}

func TestEdgeIsIdempotent(t *testing.T) {
	g, seed := seeded(t)
	e := EdgeRef{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exmple.com"},
	}
	g.Apply(op("a"), seed, Delta{Edges: []EdgeRef{e}})
	g.Apply(op("b"), seed, Delta{Edges: []EdgeRef{e}})
	if got := len(g.Edges()); got != 1 {
		t.Fatalf("edges = %d, want 1 — the same edge must converge", got)
	}
}

// --- status ------------------------------------------------------------------

func TestSkippedIsNotTerminal(t *testing.T) {
	// A gated pair must be able to run at a later barrier once belief rises.
	// Recording skipped terminally would make the first gate permanent.
	g, seed := seeded(t)
	g.SetStatus(seed, "whois", StatusSkipped)
	if _, ok := g.Status(seed, "whois"); ok {
		t.Fatal("skipped was recorded as terminal")
	}
	g.SetStatus(seed, "whois", StatusOK)
	if s, ok := g.Status(seed, "whois"); !ok || s != StatusOK {
		t.Fatalf("status = %v/%v, want ok", s, ok)
	}
}

func TestLiveRequiresAnOKStatus(t *testing.T) {
	g, seed := seeded(t)
	g.SetStatus(seed, "dns", StatusEmpty)
	if g.Live(seed) {
		t.Fatal("empty (confirmed absent) counted as live")
	}
	g.SetStatus(seed, "whois", StatusOK)
	if !g.Live(seed) {
		t.Fatal("ok status did not make the node live")
	}
}

// --- registry validation -----------------------------------------------------

func TestRegistryValidation(t *testing.T) {
	ok := func(string) (string, error) { return "x", nil }
	cases := []struct {
		name string
		def  NodeTypeDef
	}{
		{"no name", NodeTypeDef{Cap: Nameable, Version: 1, Canonical: ok}},
		{"no canonical", NodeTypeDef{Name: "a", Cap: Nameable, Version: 1}},
		{"no version", NodeTypeDef{Name: "a", Cap: Nameable, Canonical: ok}},
		{"bad cap", NodeTypeDef{Name: "a", Version: 1, Canonical: ok}},
		{"dup field", NodeTypeDef{Name: "a", Cap: Nameable, Version: 1, Canonical: ok,
			Fields: []FieldDef{{Name: "f", Kind: KindInt}, {Name: "f", Kind: KindInt}}}},
		{"bad kind", NodeTypeDef{Name: "a", Cap: Nameable, Version: 1, Canonical: ok,
			Fields: []FieldDef{{Name: "f"}}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := NewRegistry().AddType(c.def); err == nil {
				t.Fatal("accepted an invalid type definition")
			}
		})
	}
}

func TestDuplicateRegistrationRejected(t *testing.T) {
	r := testRegistry(t)
	ok := func(s string) (string, error) { return s, nil }
	if _, err := r.AddType(NodeTypeDef{Name: "domain", Cap: Nameable, Version: 1, Canonical: ok}); err == nil {
		t.Fatal("registered the same node type twice")
	}
	if _, err := r.AddRel(RelDef{Name: "TLD_OF", Class: Structural, Version: 1}); err == nil {
		t.Fatal("registered the same relation twice")
	}
}

func TestSeedOnlyOnce(t *testing.T) {
	g, _ := seeded(t)
	if _, err := g.Seed("domain", "other.com"); err == nil {
		t.Fatal("accepted a second seed")
	}
}

func TestReportOrderIsCanonical(t *testing.T) {
	// Nodes come back in (type, key) order regardless of admission order, which
	// is what makes two runs byte-comparable.
	g, seed := seeded(t)
	for _, k := range []string{"zeta.com", "alpha.com", "mid.com"} {
		g.Apply(op("x"), seed, Delta{Nodes: []NodeRef{{Type: "domain", Key: k}}})
	}
	g.Apply(op("x"), seed, Delta{Nodes: []NodeRef{{Type: "ip", Key: "9.9.9.9"}}})

	var got []string
	for _, n := range g.Nodes() {
		got = append(got, n.Type.Name()+":"+n.Key)
	}
	want := []string{
		"domain:alpha.com", "domain:example.com", "domain:mid.com", "domain:zeta.com",
		"ip:9.9.9.9",
	}
	if len(got) != len(want) {
		t.Fatalf("nodes = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("nodes = %v, want %v", got, want)
		}
	}
}

func TestScopeRestrictsWhichTypesRootVariants(t *testing.T) {
	// The CLI's scope positional (§12). Before this was enforced it was
	// validated, hashed into the plan and printed by --explain — and then
	// ignored by execution, so `typo username bob@example.com` produced exactly
	// the same graph as the unscoped run.
	g, seed := seeded(t)
	g.Apply(op("decompose"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "TLD_OF",
		To:   NodeRef{Type: "tld", Key: "com"},
	}}})
	g.SetScope([]string{"tld"})

	res := g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exmple.com"},
	}}})
	if len(res.Edges) != 0 {
		t.Fatal("a variant rooted on an out-of-scope type was admitted")
	}
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != RejectScope {
		t.Fatalf("rejections = %+v, want one RejectScope", res.Rejected)
	}

	// The scoped type still varies.
	res = g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "tld", Key: "com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "tld", Key: "cm"},
	}}})
	if len(res.Edges) != 1 {
		t.Fatalf("the in-scope type did not vary: %+v", res.Rejected)
	}
}

func TestEmptyScopeVariesEverything(t *testing.T) {
	// Absent, the positional must not narrow anything — including the seed.
	g, seed := seeded(t)
	g.SetScope(nil)
	res := g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exmple.com"},
	}}})
	if len(res.Edges) != 1 {
		t.Fatalf("an unscoped run refused a variant: %+v", res.Rejected)
	}
}

func TestScopeIsNotAClosureViolation(t *testing.T) {
	// Distinct reject kinds on purpose: an out-of-scope root is a legitimate
	// variant root the user chose not to expand, and reporting it as an
	// invariant violation would read as a bug in the scan.
	g, seed := seeded(t)
	g.SetScope([]string{"tld"})
	res := g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exmple.com"},
	}}})
	if res.Rejected[0].Kind == RejectClosure {
		t.Fatal("a scope exclusion was reported as a seed-closure violation")
	}
	if got := res.Rejected[0].Kind.String(); got != "outside-scope" {
		t.Fatalf("reject kind = %q", got)
	}
	// And it must not deny the candidate: scope is not truncation, so a later
	// run with a wider scope is not poisoned by this one.
	if len(g.Ledger()) != 0 {
		t.Fatal("a scope rejection wrote a truncation-ledger row")
	}
}

func TestSeedCountsAgainstTheGraphTotal(t *testing.T) {
	// Leaving the seed out of the tally makes g.total disagree with the node
	// count forever after, so `--budget 5` quietly admits six.
	g, _ := seeded(t)
	g.SetBudgets(Budgets{Global: 3})
	for _, k := range []string{"a.com", "b.com", "c.com", "d.com"} {
		g.Apply(op("omission"), g.SeedID(), Delta{Edges: []EdgeRef{{
			From: NodeRef{Type: "domain", Key: "example.com"},
			Rel:  "VARIANT_OF",
			To:   NodeRef{Type: "domain", Key: k},
		}}})
	}
	if got := len(g.Nodes()); got > 3 {
		t.Fatalf("a global budget of 3 admitted %d nodes", got)
	}
}
