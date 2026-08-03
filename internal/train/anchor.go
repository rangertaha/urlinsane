// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package train

import (
	"fmt"
	"math"

	"github.com/rangertaha/urlinsane/internal/model"
)

// LiveSymbol is the emission a node carries when an observation found it.
const LiveSymbol = model.OutcomePrefix + "live"

// AnchorFocus picks the latent state that means "worth expanding", by asking
// which one actually emits live observations.
//
// Baum-Welch is unsupervised. It finds structure and has no idea which of its
// states a human would call promising — the labels in Config.States are names
// we wrote on two buckets EM filled in whatever order it converged to. That is
// label switching, and it is not a subtlety: fitted on a real scan of
// example.com, the model separated cleanly (34 distinct belief values against
// uniform's 1) and put IPv6 addresses at 1.0 with every live typosquat at 0.0.
// Ranking the frontier by that is worse than not ranking it, because it is
// confidently backwards.
//
// So Focus is chosen after fitting rather than declared before it: the state
// with the higher probability of emitting `outcome=live` is the one belief
// should report. That is the whole fix, it costs one comparison, and it turns a
// model that is arbitrarily oriented into one that is oriented on the evidence.
//
// It cannot rescue a model that learned nothing. If both states emit live
// equally the choice is arbitrary and the belief it produces is near-constant,
// which is the honest outcome for a corpus with nothing in it — and why Fit's
// summary reports the outcome counts.
func AnchorFocus(h *model.HMM) (*model.HMM, string, error) {
	want, _, err := focusOn(h)
	if err != nil {
		return h, "", err
	}

	// Refit nothing: only the reported state changes. The tables are rebuilt
	// from the model's own accessors and handed back to model.New, so the
	// result is the same distribution with a different state reported — and
	// its identity stays a function of its tables, which is what its CID is.
	spec := specOf(h)
	spec.Focus = []string{want}
	out, err := model.New(spec)
	if err != nil {
		return h, "", fmt.Errorf("train: re-focusing: %w", err)
	}
	return out, want, nil
}

// focusOn returns the state that most strongly emits a live observation, and
// its index.
//
// Shared by AnchorFocus and BeliefFrom's guard so the two cannot disagree about
// which state is the right one — a guard that computed it differently would
// either reject correct models or pass inverted ones.
func focusOn(h *model.HMM) (string, int, error) {
	states := h.States()
	if len(states) < 2 {
		return "", 0, fmt.Errorf("train: %d states; nothing to choose between", len(states))
	}
	// Membership in the alphabet, not SymbolIndex: an unknown symbol maps to
	// the out-of-vocabulary slot rather than to -1, so asking for its index
	// answers for OOV and the check would never fire. Anchoring on the OOV
	// emission would orient the model on "everything I have never seen", which
	// is worse than declining.
	var known bool
	for _, sym := range h.Symbols() {
		if sym == LiveSymbol {
			known = true
			break
		}
	}
	if !known {
		return "", 0, fmt.Errorf("train: %q is not in the model's alphabet; the corpus recorded no live observation to orient on", LiveSymbol)
	}

	best, bestLog := 0, h.LogEmission(0, []string{LiveSymbol})
	for i := 1; i < len(states); i++ {
		if l := h.LogEmission(i, []string{LiveSymbol}); l > bestLog {
			best, bestLog = i, l
		}
	}
	return states[best], best, nil
}

// specOf reconstructs a Spec from a fitted model.
//
// internal/model exposes its tables in log space and has no Spec accessor, so
// the distributions are exponentiated back. exp(log p) is exact to the last
// place for the values involved, and model.New renormalizes every row anyway —
// so this round trip changes the model's numbers by nothing that survives
// normalization.
func specOf(h *model.HMM) model.Spec {
	states, rels, syms := h.States(), h.Rels(), h.Symbols()

	init := make([]float64, len(states))
	for i, lp := range h.Prior() {
		init[i] = math.Exp(lp)
	}

	trans := make(map[string][][]float64, len(rels))
	for _, rel := range rels {
		rows := make([][]float64, len(states))
		for from := range states {
			row := make([]float64, len(states))
			for to := range states {
				row[to] = math.Exp(h.LogTransition(rel, from, to))
			}
			rows[from] = row
		}
		trans[rel] = rows
	}

	emit := make([][]float64, len(states))
	for st := range states {
		row := make([]float64, len(syms))
		for j, sym := range syms {
			row[j] = math.Exp(h.LogEmission(st, []string{sym}))
		}
		emit[st] = row
	}

	return model.Spec{
		States:     states,
		Focus:      h.Focus(),
		Rels:       rels,
		Symbols:    syms,
		Init:       init,
		Trans:      trans,
		Emit:       emit,
		Smoothing:  h.Smoothing(),
		Provenance: h.Provenance(),
	}
}
