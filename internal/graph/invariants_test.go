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
	"testing"
	"time"
)

// Rejections() is declared, returned, and never appended to.
func TestRejectionsAreRecordedRunWide(t *testing.T) {
	g, seed := seeded(t)
	res := g.Apply(op("x"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "NOSUCHREL",
		To:   NodeRef{Type: "domain", Key: "b.com"},
	}}})
	if len(res.Rejected) == 0 {
		t.Fatal("setup: expected a rejection in the Result")
	}
	if len(g.Rejections()) == 0 {
		t.Errorf("Result had %d rejections but Graph.Rejections() is empty",
			len(res.Rejected))
	}
}

// the ledger denylist must compare canonical keys.
func TestLedgerDenylistComparesCanonicalKeys(t *testing.T) {
	g, seed := seeded(t)
	if err := g.Decline("domain", "EXAMPLE-DENIED.com", 0, 0, ReasonBudget, op("x")); err != nil {
		t.Fatalf("decline: %v", err)
	}
	res := g.Apply(op("y"), seed, Delta{Nodes: []NodeRef{
		{Type: "domain", Key: "example-denied.com"},
	}})
	for _, n := range g.Nodes() {
		if n.Key == "example-denied.com" {
			t.Errorf("a differently-cased spelling slipped past the denylist (%+v)", res)
		}
	}
}

// depth must be the SHORTEST observation distance, whatever order edges arrive.
func TestDepthIsShortestPathRegardlessOfArrivalOrder(t *testing.T) {
	build := func(longFirst bool) int {
		g, seed := seeded(t)
		long := []EdgeRef{
			{From: NodeRef{Type: "domain", Key: "example.com"}, Rel: "RESOLVES_TO", To: NodeRef{Type: "ip", Key: "1.1.1.1"}},
			{From: NodeRef{Type: "ip", Key: "1.1.1.1"}, Rel: "RESOLVES_TO", To: NodeRef{Type: "domain", Key: "hop.com"}},
			{From: NodeRef{Type: "domain", Key: "hop.com"}, Rel: "RESOLVES_TO", To: NodeRef{Type: "ip", Key: "9.9.9.9"}},
		}
		short := []EdgeRef{
			{From: NodeRef{Type: "domain", Key: "example.com"}, Rel: "RESOLVES_TO", To: NodeRef{Type: "ip", Key: "9.9.9.9"}},
		}
		if longFirst {
			g.Apply(op("a"), seed, Delta{Edges: long})
			g.Apply(op("b"), seed, Delta{Edges: short})
		} else {
			g.Apply(op("b"), seed, Delta{Edges: short})
			g.Apply(op("a"), seed, Delta{Edges: long})
		}
		for _, n := range g.Nodes() {
			if n.Key == "9.9.9.9" {
				return g.Depth(n.ID)
			}
		}
		t.Fatal("node missing")
		return -1
	}
	a, b := build(true), build(false)
	if a != b {
		t.Errorf("depth depends on edge arrival order: %d vs %d", a, b)
	}
	if a != 1 {
		t.Errorf("shortest observation distance = %d, want 1", a)
	}
}

// merge policy must be order-independent.
func TestMergeResultIsIndependentOfArrivalOrder(t *testing.T) {
	// Each operator always writes the SAME value, so the only thing varying
	// between the two runs is arrival order.
	val := map[string]int64{"rdap": 100, "whois": 200}
	run := func(first, second string) Value {
		g, seed := seeded(t)
		ref := NodeRef{Type: "domain", Key: "example.com"}
		g.Apply(Provenance{Operator: first, Round: 1}, seed, Delta{Props: []PropSet{
			{Node: &ref, Field: "created", Value: Time(time.Unix(val[first], 0))},
		}})
		g.Apply(Provenance{Operator: second, Round: 1}, seed, Delta{Props: []PropSet{
			{Node: &ref, Field: "created", Value: Time(time.Unix(val[second], 0))},
		}})
		for _, n := range g.Nodes() {
			if n.Key == "example.com" {
				f, _ := n.Type.Field("created")
				v, _ := n.Props.Get(f)
				return v
			}
		}
		return Value{}
	}
	a := run("rdap", "whois")
	b := run("whois", "rdap")
	if !a.equal(b) {
		t.Errorf("merge result depends on arrival order: %v vs %v", a.Num(), b.Num())
	}
}

// cache must miss when the read-set changes.
func TestCacheKeyCoversReadSetAndResourceConfig(t *testing.T) {
	c := NewCache()
	o := &fakeOp{id: "dns", ver: 1, trig: onDomain(Reads{})}
	id := newNodeID("domain", "example.com")
	k1 := c.Key(o, id, [32]byte{1})
	k2 := c.Key(o, id, [32]byte{2})
	if k1 == k2 {
		t.Error("cache key ignores the read digest")
	}
	c.SetResourceConfig("", "ns=9.9.9.9")
	if k3 := c.Key(o, id, [32]byte{1}); k3 == k1 {
		t.Error("cache key ignores resource config")
	}
}

// an out-of-scope variant must not deny the candidate for later runs.
func TestScopeRejectionWritesNoLedgerRow(t *testing.T) {
	g, seed := seeded(t)
	g.SetScope([]string{"tld"})
	g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exmple.com"},
	}}})
	if len(g.Ledger()) != 0 {
		t.Error("a scope rejection wrote a ledger row")
	}
}

// budgets must not be circumventable via bare Delta.Nodes.
func TestBareNodesRespectBudgets(t *testing.T) {
	g, seed := seeded(t)
	g.SetBudgets(Budgets{Global: 3})
	var refs []NodeRef
	for _, k := range []string{"a.com", "b.com", "c.com", "d.com", "e.com"} {
		refs = append(refs, NodeRef{Type: "domain", Key: k})
	}
	g.Apply(op("x"), seed, Delta{Nodes: refs})
	if got := len(g.Nodes()); got > 3 {
		t.Errorf("bare Delta.Nodes bypassed the budget: %d nodes admitted", got)
	}
}

// stateful is a model whose latent state cannot be reconstructed from its
// scalar: two different states can yield the same belief.
type stateful struct{ seen [][]int }

func (m *stateful) Initial() (float64, State) { return 1, []int{1, 0, 0} }

func (m *stateful) Step(parent State, _ string, _ View) (float64, State) {
	p, _ := parent.([]int)
	m.seen = append(m.seen, p)
	next := []int{0, 0, 0}
	for i := range p {
		next[(i+1)%3] = p[i]
	}
	// Belief is deliberately constant, so a scalar-only interface would carry
	// no information about the rotation at all.
	return 1, next
}

func TestBeliefModelStateReachesChildrenIntact(t *testing.T) {
	// §10.1 specifies an HMM, and forward filtering propagates a distribution
	// over latent states. The interface used to pass only the parent's scalar,
	// so a model had to guess a distribution consistent with one number —
	// exact for two states, wrong for three or more, and wrong quietly.
	g, seed := seeded(t)
	m := &stateful{}
	g.SetBeliefModel(m)

	g.Apply(op("decompose"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "TLD_OF",
		To:   NodeRef{Type: "tld", Key: "com"},
	}}})
	g.recomputeBelief()

	if len(m.seen) == 0 {
		t.Fatal("Step never ran")
	}
	got := m.seen[0]
	if len(got) != 3 {
		t.Fatalf("Step received %v; the model's own state must arrive intact", got)
	}
	if got[0] != 1 || got[1] != 0 || got[2] != 0 {
		t.Fatalf("Step received state %v, want the parent's [1 0 0]", got)
	}
}

func TestBeliefStateSurvivesAcrossGenerations(t *testing.T) {
	// A grandchild must see the child's state, not the seed's — otherwise the
	// chain is not a chain and forward filtering means nothing.
	g, seed := seeded(t)
	m := &stateful{}
	g.SetBeliefModel(m)

	g.Apply(op("dns"), seed, Delta{Edges: []EdgeRef{
		{From: NodeRef{Type: "domain", Key: "example.com"}, Rel: "RESOLVES_TO",
			To: NodeRef{Type: "ip", Key: "1.2.3.4"}},
		{From: NodeRef{Type: "ip", Key: "1.2.3.4"}, Rel: "RESOLVES_TO",
			To: NodeRef{Type: "domain", Key: "hop.com"}},
	}})
	g.recomputeBelief()

	if len(m.seen) < 2 {
		t.Fatalf("Step ran %d times, want at least 2", len(m.seen))
	}
	// One rotation apart: the second call must see what the first returned.
	if m.seen[1][1] != 1 {
		t.Fatalf("grandchild stepped from %v; it must inherit the child's state, not the seed's",
			m.seen[1])
	}
}

func TestANodeIsNotAVariantOfItself(t *testing.T) {
	// Bit-flipping "google.com" flips a case bit to "Google.com", which folds
	// back onto the origin at canonicalization. Operators emit raw keys and
	// cannot see canonical form, so their own dedupe cannot catch this — only
	// the applier can. Left in, the seed becomes a variant of itself and gets
	// scored as a live typosquat of itself.
	g, seed := seeded(t)
	res := g.Apply(op("bf"), seed, Delta{Edges: []EdgeRef{{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "EXAMPLE.com"}, // canonicalizes onto the origin
	}}})
	if len(res.Edges) != 0 {
		t.Fatal("admitted a VARIANT_OF self-edge")
	}
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != RejectSelfVariant {
		t.Fatalf("rejections = %+v, want one RejectSelfVariant", res.Rejected)
	}
	for _, e := range g.Edges() {
		if e.Rel.Name() == VariantRel && e.From == e.To {
			t.Fatal("a self-variant edge survived in the graph")
		}
	}
}

func TestExistenceCountsOnlyOperatorsThatLookedSomethingUp(t *testing.T) {
	// A decomposer returns ok when it successfully *parses* a name, which says
	// nothing about whether that name exists. Counting it made every
	// syntactically valid variant read as live: "-google.com" was reported as a
	// live typosquat on the strength of having been parsed, with every DNS and
	// whois lookup against it empty.
	g, seed := seeded(t)
	g.SetObservers([]string{"dns", "whois"})

	g.SetStatus(seed, "decompose.domain", StatusOK)
	g.SetStatus(seed, "dns", StatusEmpty)
	g.SetStatus(seed, "whois", StatusEmpty)

	if got := g.Analyze().Existence(seed); got != Absent {
		t.Fatalf("existence = %v, want absent — only the decomposer said ok", got)
	}

	g.SetStatus(seed, "dns", StatusOK)
	if got := g.Analyze().Existence(seed); got != Live {
		t.Fatalf("existence = %v, want live once a lookup succeeded", got)
	}
}

func TestExistenceCountsEverythingWhenNoObserversDeclared(t *testing.T) {
	// The honest fallback: a caller that never said which operators observe has
	// given us nothing to discriminate on. Answering "no" would make every node
	// read as untried.
	g, seed := seeded(t)
	g.SetStatus(seed, "anything", StatusOK)
	if got := g.Analyze().Existence(seed); got != Live {
		t.Fatalf("existence = %v, want live", got)
	}
}
