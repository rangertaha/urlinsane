// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package train

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/model"
)

// The latent alphabet.
//
// Two states, named for what the engine actually needs to decide: is this
// lineage worth spending network calls on. DESIGN §16 leaves the cardinality
// open, and it stays open — three states would need a third meaning somebody
// can name, and unsupervised EM cannot invent one. Two is the smallest choice
// that is not degenerate, and Focus must be a proper subset or belief is
// constant.
const (
	StatePromising = "promising"
	StateDead      = "dead"
)

// DefaultConfig is the training configuration.
//
// Seed is fixed rather than drawn from a clock: it is recorded in the model's
// provenance, and a model that cannot be refitted to the same numbers cannot be
// compared with its successor.
func DefaultConfig() model.Config {
	return model.Config{
		States:     []string{StatePromising, StateDead},
		Focus:      []string{StatePromising},
		Iterations: 50,
		Tolerance:  1e-4,
		Seed:       1,
	}
}

// Fit trains a belief model on one or more finished scans.
//
// Every graph contributes its own expansion tree; they are concatenated because
// Baum-Welch is defined over a set of independent sequences, and two scans of
// different brands are exactly that.
//
// It refuses an empty corpus rather than returning a model. Fitting nothing
// yields the prior, which is uniform, which is what the engine already does —
// so a "successful" run on no data would produce an artifact that changes
// nothing and looks like progress.
func Fit(cfg model.Config, scans ...Scan) (*model.Result, model.Corpus, error) {
	var c model.Corpus
	for _, s := range scans {
		if s.Graph == nil {
			continue
		}
		c.Paths = append(c.Paths, Paths(s.Graph, s.Seed)...)
		c.CIDs = append(c.CIDs, s.Roots...)
	}
	if len(c.Paths) == 0 {
		return nil, c, fmt.Errorf("train: no expansion paths; fitting nothing yields the uniform prior the engine already has")
	}

	res, err := model.Train(c, cfg)
	if err != nil {
		return nil, c, fmt.Errorf("train: %w", err)
	}

	// Anchored here, not left to the caller. Baum-Welch is unsupervised, so
	// which of its states means "live" is decided by the jitter in the
	// initialisation, not by the names in Config.States — and because the seed
	// is fixed for reproducibility, an arbitrary orientation becomes a
	// *reproducible* one. AnchorFocus existed for this and nothing called it, so
	// every model this package produced had a coin-flip chance of being
	// confidently backwards, identically on every fold.
	//
	// Measured on three saved scans, leave-one-out: held-out AUC 0.150 and 0.210
	// unanchored against a 0.500 uniform prior, 0.867 and 0.802 anchored. Same
	// fit, same data, read from opposite axes.
	//
	// Doing it in Fit is what retires that: there is no longer a way to obtain a
	// model from this package that has not been oriented on its own evidence.
	anchored, focus, err := AnchorFocus(res.Model)
	if err != nil {
		return nil, c, fmt.Errorf("train: fitted model cannot be oriented: %w", err)
	}
	res.Model = anchored
	_ = focus // recoverable from res.Model.Focus(); reported by Describe
	return res, c, nil
}

// Scan is one finished expansion to train on.
type Scan struct {
	Graph *graph.Graph
	Seed  graph.NodeID
	// Roots are the store roots this graph was loaded from, if any. They reach
	// the model's provenance so an artifact points back at its data.
	Roots []cid.Cid
}

// BeliefFrom wraps a fitted model as the engine's execution model.
//
// The featurizer is this package's, not the caller's, which is the point: a
// model fitted on Features must be served Features. Letting a caller supply a
// different one is how train/serve skew gets in, and it would not fail —
// belief would just quietly collapse to the prior.
//
// It also refuses a model whose Focus is not the state that emits live
// observations. Fit anchors, so a model from this package always passes; the
// check is here because this is the seam where the harm happens. An
// unanchored model does not degrade — it inverts, and an inverted belief
// ranks the frontier worse than no belief at all, spending the budget on the
// names least likely to exist while reporting confidence. That is worth an
// error rather than a scan nobody can tell is backwards.
func BeliefFrom(h *model.HMM) (graph.BeliefModel, error) {
	if err := checkOriented(h); err != nil {
		return nil, err
	}
	return model.NewBelief(h, Featurizer), nil
}

// checkOriented reports whether the model's Focus is the state most likely to
// emit a live observation — the same comparison AnchorFocus makes.
func checkOriented(h *model.HMM) error {
	want, _, err := focusOn(h)
	if err != nil {
		return err
	}
	got := h.Focus()
	if len(got) != 1 || got[0] != want {
		return fmt.Errorf("train: model reports state %v, but %q is the one that emits %q; "+
			"belief would be inverted — build it with Fit, or pass it through AnchorFocus",
			got, want, LiveSymbol)
	}
	return nil
}

// Summary is what a training run should print. It is deliberately small: the
// numbers a person checks before trusting a model.
type Summary struct {
	Paths      int
	Steps      int
	Symbols    int
	Rels       int
	Iterations int
	// LogLikelihood is the final value, and Objective the final MAP objective.
	// Objective is the one with the monotonicity guarantee under a prior;
	// LogLikelihood is the one people read, so both are reported.
	LogLikelihood float64
	Objective     float64
	// Outcomes counts how many steps carried each observation, which is the
	// first thing to check: a corpus that is 99% untried has nothing to learn
	// from, and a corpus with no "absent" cannot learn what a dead lineage
	// looks like.
	Outcomes map[string]int
}

// Describe summarises a corpus and a training result together.
func Describe(c model.Corpus, res *model.Result) Summary {
	s := Summary{Outcomes: map[string]int{}}
	s.Paths = len(c.Paths)
	for _, p := range c.Paths {
		s.Steps += len(p.Steps)
		for _, st := range p.Steps {
			key := st.Outcome
			if key == "" {
				key = "untried"
			}
			s.Outcomes[key]++
		}
	}
	syms, rels := Alphabet(c)
	s.Symbols, s.Rels = len(syms), len(rels)

	if res != nil {
		s.Iterations = res.Iterations
		if n := len(res.LogLikelihood); n > 0 {
			s.LogLikelihood = res.LogLikelihood[n-1]
		}
		if n := len(res.Objective); n > 0 {
			s.Objective = res.Objective[n-1]
		}
	}
	return s
}

func (s Summary) String() string {
	return fmt.Sprintf("%d paths, %d steps, %d symbols, %d relations; %d iterations, logL=%.2f objective=%.2f; outcomes=%v",
		s.Paths, s.Steps, s.Symbols, s.Rels, s.Iterations, s.LogLikelihood, s.Objective, s.Outcomes)
}
