// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package train

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/model"
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
	// And it must not *print* as 0.000, which reads as the worst possible score
	// rather than as an undefined one.
	if q.Comparable != 0 {
		t.Errorf("Comparable = %d, want 0: there are no (live, absent) pairs", q.Comparable)
	}
	if !strings.Contains(q.String(), "AUC=n/a") {
		t.Errorf("String() = %q, want AUC reported as n/a", q.String())
	}
}

// Fit must hand back an oriented model, and BeliefFrom must refuse one that is
// not.
//
// This is the defect that made the whole learned-belief path worse than no
// belief at all. Baum-Welch is unsupervised, so which latent state means "live"
// is decided by initialisation jitter rather than by the names in Config.States;
// with the seed fixed for reproducibility, an arbitrary orientation becomes a
// reproducible one. AnchorFocus was written for exactly this and its own comment
// says an unanchored model is "confidently backwards" — and nothing called it.
//
// Measured over three saved scans, leave-one-out: held-out AUC 0.150 and 0.210
// unanchored against a 0.500 uniform prior, 0.867 and 0.802 anchored. The same
// fit, read from opposite axes.
func TestFitAnchorsAndBeliefFromRefusesAnInvertedModel(t *testing.T) {
	g, seed := scan(t)

	// Seed 3 deliberately. Over seeds 1..12 on this fixture, five (3, 6, 8, 9,
	// 12) leave Config.Focus naming the wrong state and seven do not — a coin
	// flip, which is the whole point. Pinning to a seed that inverts is what
	// makes this a regression test rather than one that passes because the
	// arbitrary choice happened to land right.
	cfg := DefaultConfig()
	cfg.Seed = 3
	res, _, err := Fit(cfg, Scan{Graph: g, Seed: seed})
	if err != nil {
		t.Fatalf("fit: %v", err)
	}

	// What Fit returns is oriented on the evidence, whatever Config.Focus said.
	want, _, err := focusOn(res.Model)
	if err != nil {
		t.Fatalf("focusOn: %v", err)
	}
	if got := res.Model.Focus(); len(got) != 1 || got[0] != want {
		t.Errorf("Fit returned focus %v, want [%s]: the model is not oriented on "+
			"the state that emits %q", got, want, LiveSymbol)
	}
	if _, err := BeliefFrom(res.Model); err != nil {
		t.Errorf("BeliefFrom rejected a model Fit produced: %v", err)
	}

	// And an inverted one cannot be served. Invert by pointing Focus at the
	// other state — the exact thing an unanchored fit does by chance.
	spec := specOf(res.Model)
	for _, s := range res.Model.States() {
		if s != want {
			spec.Focus = []string{s}
		}
	}
	bad, err := model.New(spec)
	if err != nil {
		t.Fatal(err)
	}
	bm, err := BeliefFrom(bad)
	if err == nil {
		t.Fatal("BeliefFrom served a model whose belief is inverted; the scan " +
			"would spend its budget on the names least likely to exist")
	}
	if bm != nil {
		t.Error("a belief model was returned alongside the error")
	}
	if !strings.Contains(err.Error(), "inverted") {
		t.Errorf("err = %v, want it to say belief would be inverted", err)
	}
}

// A model must arrive at a scan still oriented.
//
// Fit anchors and BeliefFrom checks, but between them lies the artifact: a
// model is encoded to dag-cbor, addressed by CID, and decoded somewhere else.
// If Focus did not survive that, every model loaded from disk would be rejected
// by the guard — or worse, accepted while inverted. The path is fit -> Addressed
// -> Decode -> BeliefFrom, and it is the one a --model flag will use, so it is
// worth asserting rather than assuming.
//
// Seed 3 because it is one of the seeds whose unanchored orientation is wrong
// (see TestFitAnchorsAndBeliefFromRefusesAnInvertedModel), so this also proves
// the anchoring is what survives, not a coincidence.
func TestAModelSurvivesTheArtifactRoundTripOriented(t *testing.T) {
	g, seed := scan(t)
	cfg := DefaultConfig()
	cfg.Seed = 3
	res, _, err := Fit(cfg, Scan{Graph: g, Seed: seed})
	if err != nil {
		t.Fatalf("fit: %v", err)
	}

	block, want, err := res.Model.Addressed()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	back, err := model.Decode(block)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got, orig := back.Focus(), res.Model.Focus(); len(got) != len(orig) ||
		(len(got) == 1 && got[0] != orig[0]) {
		t.Errorf("Focus = %v after the round trip, was %v", got, orig)
	}
	if _, err := BeliefFrom(back); err != nil {
		t.Errorf("BeliefFrom rejected a model that came back off its own artifact: %v", err)
	}

	// The CID is the model's identity, so it must not depend on having been
	// decoded rather than fitted.
	got, err := back.CID()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("CID = %s after the round trip, was %s", got, want)
	}
}
