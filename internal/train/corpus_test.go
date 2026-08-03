// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package train

import (
	"fmt"
	"testing"
	"time"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Finalize must not lose the observer set.
//
// It runs a scheduler with no operators, and NewScheduler *rebuilds* the
// observer set from the operators it is given — so with none, the set comes out
// empty, and an empty set means "everything observes". Every operator's status
// then attests that the name exists, including a variant generator's "I produced
// this", and a node whose only lookup returned an authoritative absence is
// relabelled live.
//
// That is not a cosmetic wrong label. Finalize is the documented entry point for
// a graph rehydrated from the store, so it ran over the training corpus and the
// held-out corpus alike: both were labelled all-live, and the AUC computed from
// them measured nothing.
func TestFinalizeKeepsTheObserverSet(t *testing.T) {
	g, _ := scan(t)
	id := node(t, g, "exampl.com")

	before := Outcome(g.Analyze(), id)
	if before != "absent" {
		t.Fatalf("fixture changed: exampl.com is %q before Finalize, want %q",
			before, "absent")
	}

	if err := Finalize(g); err != nil {
		t.Fatalf("Finalize: %v", err)
	}

	if got := g.Observers(); len(got) != 1 || got[0] != "dns" {
		t.Errorf("Observers() = %v after Finalize, want [dns]", got)
	}
	if got := Outcome(g.Analyze(), id); got != before {
		t.Errorf("Outcome = %q after Finalize, was %q before: a node whose only "+
			"lookup returned NXDOMAIN is now labelled as existing", got, before)
	}
}

// Paths must not be quadratic in the size of the scan.
//
// Two ways it was: nodeView.EdgeProps answered from Analysis.Edges(), which
// allocates and sorts the whole edge list on every call and Features asks six
// times per node; and an interior node was re-featurized once per path through
// it rather than once. Over this fixture the difference is 24.1s against 141ms.
//
// The bound is deliberately loose. It is not a performance target: 2s is 14×
// the fixed cost and 6% of the quadratic one, so it fails on a return to
// quadratic rather than on a slow machine.
func TestPathsScalesLinearly(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a 2000-node graph")
	}
	g, seed := wideScan(t, 2000)

	start := time.Now()
	paths := Paths(g, seed)
	elapsed := time.Since(start)

	if len(paths) == 0 {
		t.Fatal("no paths from a 2000-node graph")
	}
	t.Logf("Paths over 2000 nodes: %v", elapsed)
	if elapsed > 2*time.Second {
		t.Errorf("Paths over 2000 nodes took %v, want well under 2s "+
			"(it was 24.1s when EdgeProps re-sorted every edge per lookup)", elapsed)
	}
}

// wideScan builds a seed with n resolved variants, the shape a real scan has.
func wideScan(t *testing.T, n int) (*graph.Graph, graph.NodeID) {
	t.Helper()
	g := graph.New(registry(t))
	seed, err := g.Seed("domain", "example.com")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	origin := graph.NodeRef{Type: "domain", Key: "example.com"}
	by := graph.Provenance{Operator: "co", Round: 1}

	for i := 0; i < n; i++ {
		key := fmt.Sprintf("v%04d.com", i)
		ref := graph.NodeRef{Type: "domain", Key: key}
		e := graph.EdgeRef{From: origin, Rel: graph.VariantRel, To: ref}
		g.Apply(by, seed, graph.Delta{
			Edges: []graph.EdgeRef{e},
			Props: []graph.PropSet{
				{Edge: &e, Field: "algorithm", Value: graph.String("co")},
				{Edge: &e, Field: "distance", Value: graph.Int(1)},
			},
		})
		id := node(t, g, key)
		g.SetStatus(id, "dns", graph.StatusOK)
		g.Apply(graph.Provenance{Operator: "dns", Round: 2}, id, graph.Delta{
			Edges: []graph.EdgeRef{{From: ref, Rel: "RESOLVES_TO",
				To: graph.NodeRef{Type: "ip", Key: fmt.Sprintf("10.0.%d.%d", i/256, i%256)}}},
		})
	}
	g.SetStatus(seed, "dns", graph.StatusOK)
	barrier(t, g)
	g.SetObservers([]string{"dns"})
	return g, seed
}

// A scan where every judged name resolved must still report a base rate.
//
// AUC is genuinely undefined with one class — it is a statement about pairs of
// nodes with opposite labels, and there are none — but the base rate and
// precision@k are not. Returning early skipped all three, so such a scan printed
// "AUC=0.000 base=0.000 p@10=0.000", which reads as a model that ranked every
// live name below every absent one. The truth is the opposite: everything
// resolved, and p@10 is 1.000.
func TestEvaluateReportsWhatIsDefinedWithOneClass(t *testing.T) {
	g, _ := wideScan(t, 12) // every variant resolves; nothing absent

	q := Evaluate(g)

	if q.Absent != 0 {
		t.Fatalf("fixture changed: %d absent nodes, want a one-class scan", q.Absent)
	}
	if q.BaseRate != 1 {
		t.Errorf("BaseRate = %.3f, want 1.000: every judged node is live", q.BaseRate)
	}
	if got := q.PrecisionAt[10]; got != 1 {
		t.Errorf("PrecisionAt[10] = %.3f, want 1.000: the top ten are all live", got)
	}
	if q.AUC != 0 {
		t.Errorf("AUC = %.3f, want 0: it is undefined with one class and must not "+
			"be invented", q.AUC)
	}
}
