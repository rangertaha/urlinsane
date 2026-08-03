package train

import (
	"fmt"
	"testing"

	"github.com/rangertaha/urlinsane/internal/model"
)

// Fit's own output must never fail BeliefFrom's guard.
//
// The two are computed independently — Fit anchors, BeliefFrom re-derives the
// expected state and compares — so a disagreement is possible in principle, and
// it would be the worst kind: the tool refusing a model it had just built
// correctly. AnchorFocus round-trips the distributions through exp/log in specOf
// and model.New renormalises every row, so when two states emit the live symbol
// at nearly equal probability, drift could in principle move the argmax between
// the anchor and the check.
//
// Swept over 60 seeds against the real fixture: it does not. Kept as a sweep
// rather than one case because the failure would be seed-dependent by nature.
func TestFitsOwnModelAlwaysPassesTheGuard(t *testing.T) {
	g, seed := scan(t)
	bad := 0
	for s := int64(1); s <= 60; s++ {
		cfg := DefaultConfig()
		cfg.Seed = s
		res, _, err := Fit(cfg, Scan{Graph: g, Seed: seed})
		if err != nil {
			t.Logf("seed %d: fit refused: %v", s, err)
			continue
		}
		if _, err := BeliefFrom(res.Model); err != nil {
			bad++
			t.Errorf("seed %d: BeliefFrom REJECTED a model Fit produced: %v", s, err)
		}
	}
	t.Logf("%d/60 seeds produced a model Fit made and BeliefFrom refused", bad)
}

// The degenerate case the drift sweep approximates: an exact tie.
//
// With both emission rows identical there is no evidence to orient on, so the
// choice is arbitrary — but it must still be *consistent*, because the anchor
// and the guard make it separately. An anchor that picked state 0 while the
// guard expected state 1 would reject every tied model.
func TestAnchoringATieIsStableAndAccepted(t *testing.T) {
	g, seed := scan(t)
	res, _, err := Fit(DefaultConfig(), Scan{Graph: g, Seed: seed})
	if err != nil {
		t.Fatal(err)
	}
	spec := specOf(res.Model)
	// Make both states' emission rows identical -> exact tie on every symbol.
	for i := range spec.Emit[1] {
		spec.Emit[1][i] = spec.Emit[0][i]
	}
	tied, err := model.New(spec)
	if err != nil {
		t.Fatal(err)
	}
	a, focus, err := AnchorFocus(tied)
	if err != nil {
		t.Fatalf("AnchorFocus on a tie: %v", err)
	}
	t.Logf("tie -> focus %q", focus)
	if _, err := BeliefFrom(a); err != nil {
		t.Errorf("BeliefFrom rejected AnchorFocus's own output on a tie: %v", err)
	}
	// And stability: anchoring twice must agree.
	_, focus2, err := AnchorFocus(a)
	if err != nil {
		t.Fatal(err)
	}
	if focus != focus2 {
		t.Errorf("anchoring is unstable on a tie: %q then %q", focus, focus2)
	}
}

// Anchoring an already-anchored model must not move it.
//
// Fit anchors, so any model that reaches AnchorFocus again — a saved artifact
// re-anchored by a caller, say — must come back unchanged rather than flipping
// on each pass.
func TestAnchoringIsIdempotent(t *testing.T) {
	g, seed := scan(t)
	for s := int64(1); s <= 12; s++ {
		cfg := DefaultConfig()
		cfg.Seed = s
		res, _, err := Fit(cfg, Scan{Graph: g, Seed: seed})
		if err != nil {
			continue
		}
		once := fmt.Sprint(res.Model.Focus())
		again, _, err := AnchorFocus(res.Model)
		if err != nil {
			t.Errorf("seed %d: re-anchoring failed: %v", s, err)
			continue
		}
		if twice := fmt.Sprint(again.Focus()); once != twice {
			t.Errorf("seed %d: anchoring is not idempotent: %s then %s", s, once, twice)
		}
	}
}
