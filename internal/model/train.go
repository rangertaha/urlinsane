// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package model

import (
	"fmt"
	"math"
	"time"

	"github.com/ipfs/go-cid"
)

// OutcomePrefix namespaces the emission symbol built from a trace's outcome.
//
// The outcome — resolved, absent, refused — is an observation like any other
// prop, so it belongs in the emission alphabet rather than being dropped. It is
// prefixed so a featurizer can produce the identical symbol at run time
// (OutcomeSymbol below) and so it can never collide with a prop named
// "outcome".
const OutcomePrefix = "outcome="

// OutcomeSymbol renders an outcome as the emission symbol training used, so a
// run-time featurizer and the corpus agree on spelling. An empty outcome
// produces no symbol.
func OutcomeSymbol(outcome string) string {
	if outcome == "" {
		return ""
	}
	return OutcomePrefix + outcome
}

// Trace is one recorded expansion step: the (parent belief, relation, props,
// outcome) tuple that `typo --trace FILE` writes (§10.4).
//
// Recording is opt-in for a reason — it persists observation data a normal scan
// discards — so this type is what a user has consented to keep, and nothing
// more.
type Trace struct {
	// Parent is the belief the run actually acted on when this node was
	// admitted. Baum-Welch does not read it: the latent chain already carries
	// the parent's state, and conditioning on a number the previous model
	// produced would train the new model to imitate the old one. It is
	// recorded so a retrained model can be compared against the decisions the
	// scan really made, and so a future discriminative fit has it.
	Parent float64
	// Rel is the relation that admitted this node. It is what the transition
	// table is conditioned on. Ignored on the first step of a Path, which has
	// no parent.
	Rel string
	// Props are the emission symbols for this node as of the barrier —
	// typically "field=value" strings from a Featurizer. Empty is legitimate: a
	// candidate has no props yet.
	Props []string
	// Outcome is the recorded result for this node, folded into the emission
	// alphabet under OutcomePrefix. Empty means none was recorded.
	Outcome string
}

// Symbols is the full observation for this step: props plus the outcome.
func (t Trace) Symbols() []string {
	if t.Outcome == "" {
		return t.Props
	}
	return append(append([]string(nil), t.Props...), OutcomeSymbol(t.Outcome))
}

// Path is one root-to-leaf expansion path, seed first.
//
// A path and not a graph: the expansion tree gives every node exactly one
// parent, so seed → node is a sequence, which is what an HMM is defined over
// (§10.1). Steps[0] is the sequence start, drawn from the initial distribution;
// its Rel is ignored because the seed has no parent.
type Path struct {
	Steps []Trace
}

// Corpus is a set of recorded expansion paths and the blocks they came from.
type Corpus struct {
	Paths []Path
	// CIDs are the trace blocks this corpus was read from. They are copied into
	// the trained model's provenance so a model always points back at what it
	// was fitted on.
	CIDs []cid.Cid
}

// Config parameterizes Baum-Welch.
type Config struct {
	// States is the latent alphabet to fit. Unsupervised EM cannot invent
	// meaning for these, so the number is a modelling choice: see §16, "state
	// cardinality".
	States []string
	// Focus names the states whose posterior mass becomes belief.
	Focus []string
	// Iterations caps EM passes.
	Iterations int
	// Tolerance stops early once the objective improves by less than this.
	Tolerance float64
	// Seed initializes the symmetry-breaking RNG. It is recorded in the model's
	// provenance so that training reproduces exactly (§10.4).
	Seed int64
	// Smoothing is the Dirichlet prior; every component must be positive.
	Smoothing Smoothing
	// Date is stamped into provenance. Zero means "now"; tests pass a fixed
	// value so the resulting CID is stable.
	Date time.Time
}

// Result is a training run and its diagnostics.
type Result struct {
	// Model is the fitted HMM, with provenance already filled in.
	Model *HMM
	// LogLikelihood is the corpus log-likelihood before each M-step, plus the
	// value for the final model. Under a prior, EM is only guaranteed to
	// improve Objective, so this series may dip; it is reported because it is
	// the quantity a person actually reads.
	LogLikelihood []float64
	// Objective is LogLikelihood plus the Dirichlet log-prior, evaluated on the
	// same series of models. MAP-EM guarantees this is non-decreasing, and that
	// is the guarantee worth testing.
	Objective []float64
	// Iterations is the number of M-steps performed, which may be fewer than
	// Config.Iterations if the run converged.
	Iterations int
}

// Train fits an HMM to recorded expansion traces by Baum-Welch (§10.4).
//
// Training is where the backward recursion lives. §10.1 rules it out of run
// time — a model that only runs after expansion cannot save a network call —
// but offline, with the whole path in hand, forward-backward is what EM needs
// and costs nothing at scan time.
//
// Everything that could vary between two runs over the same corpus is pinned:
// the alphabets are derived in sorted order rather than from map iteration, and
// parameter initialization comes from an explicit seeded PRNG rather than the
// global one. Training the same corpus with the same Config twice gives the
// same tables and therefore the same CID.
func Train(c Corpus, cfg Config) (*Result, error) {
	if len(cfg.States) == 0 {
		return nil, fmt.Errorf("model: training needs at least one state")
	}
	if cfg.Iterations <= 0 {
		return nil, fmt.Errorf("model: iterations must be positive")
	}
	sm := cfg.Smoothing
	if sm.IsZero() {
		sm = DefaultSmoothing()
	}
	if err := sm.validate(true); err != nil {
		return nil, err
	}
	paths := usablePaths(c)
	if len(paths) == 0 {
		return nil, fmt.Errorf("model: corpus has no paths")
	}

	rels, symbols := alphabets(paths)
	h, err := initialModel(cfg, rels, symbols, sm)
	if err != nil {
		return nil, err
	}

	// Symbols resolved once: a path is re-scanned every iteration, and name
	// lookups would otherwise dominate the inner loop.
	obs := encodePaths(h, paths)

	res := &Result{}
	for iter := 0; ; iter++ {
		ll, cnt := expectation(h, obs)
		obj := ll + logPrior(h, sm)
		res.LogLikelihood = append(res.LogLikelihood, ll)
		res.Objective = append(res.Objective, obj)

		if iter >= cfg.Iterations {
			break
		}
		if iter > 0 && obj-res.Objective[iter-1] < cfg.Tolerance {
			break
		}
		maximization(h, cnt, sm)
		res.Iterations++
	}

	when := cfg.Date
	if when.IsZero() {
		when = time.Now()
	}
	h.prov = Provenance{
		Algorithm:     AlgorithmBaumWelch,
		Seed:          cfg.Seed,
		Date:          when.UTC(),
		Corpus:        append([]cid.Cid(nil), c.CIDs...),
		Iterations:    res.Iterations,
		LogLikelihood: res.LogLikelihood[len(res.LogLikelihood)-1],
	}
	res.Model = h
	return res, nil
}

// usablePaths drops empty paths, which carry no evidence and would divide by
// zero in the E-step.
func usablePaths(c Corpus) []Path {
	out := make([]Path, 0, len(c.Paths))
	for _, p := range c.Paths {
		if len(p.Steps) > 0 {
			out = append(out, p)
		}
	}
	return out
}

// alphabets derives the relation and symbol vocabularies from the corpus in
// sorted order. Sorted, not observed-order: index position is part of the
// model's encoded identity, so deriving it from map iteration would give two
// runs over the same corpus different CIDs.
//
// Relations are collected from steps 1 onward only. Step 0 is the sequence
// start and has no transition, so a relation seen only there would get a table
// no evidence ever touches.
func alphabets(paths []Path) (rels, symbols []string) {
	rs := map[string]bool{}
	ss := map[string]bool{}
	for _, p := range paths {
		for i, st := range p.Steps {
			if i > 0 && st.Rel != "" {
				rs[st.Rel] = true
			}
			for _, s := range st.Symbols() {
				if s != "" {
					ss[s] = true
				}
			}
		}
	}
	return sortedKeys(rs), sortedKeys(ss)
}

// initialModel builds the starting point for EM: uniform tables plus a small
// seeded jitter.
//
// The jitter is not decoration. Exactly uniform tables are a stationary point
// of EM — every state has identical parameters, so every posterior is identical
// and the M-step reproduces the same uniform tables forever. Breaking the
// symmetry is what lets the states differentiate at all, and doing it from a
// recorded seed is what keeps the result reproducible.
func initialModel(cfg Config, rels, symbols []string, sm Smoothing) (*HMM, error) {
	r := newRNG(cfg.Seed)
	ns := len(cfg.States)
	rels = withReserved(rels, OOVRelation)
	symbols = withReserved(symbols, OOVSymbol)

	spec := Spec{
		States:    append([]string(nil), cfg.States...),
		Focus:     append([]string(nil), cfg.Focus...),
		Rels:      rels,
		Symbols:   symbols,
		Smoothing: sm,
		Trans:     map[string][][]float64{},
	}
	if len(spec.Focus) == 0 {
		// A model must have somewhere to put belief; without an explicit
		// choice, take the last state, so a two-state fit reads as
		// "uninteresting, interesting".
		spec.Focus = []string{cfg.States[ns-1]}
	}
	spec.Init = jitterRow(r, ns)
	for _, rel := range rels {
		rows := make([][]float64, ns)
		for i := range rows {
			rows[i] = jitterRow(r, ns)
		}
		spec.Trans[rel] = rows
	}
	spec.Emit = make([][]float64, ns)
	for i := range spec.Emit {
		spec.Emit[i] = jitterRow(r, len(symbols))
	}
	return New(spec)
}

// jitterRow is uniform mass perturbed by up to 10%.
func jitterRow(r *rng, n int) []float64 {
	row := make([]float64, n)
	for i := range row {
		row[i] = 1 + 0.1*r.float()
	}
	return row
}

// encoded is a path with its names already resolved to table indices.
type encoded struct {
	rel []int   // rel[t] is the transition into step t; rel[0] is unused
	obs [][]int // obs[t] are the emission symbol indices for step t
}

func encodePaths(h *HMM, paths []Path) []encoded {
	out := make([]encoded, 0, len(paths))
	for _, p := range paths {
		e := encoded{
			rel: make([]int, len(p.Steps)),
			obs: make([][]int, len(p.Steps)),
		}
		for t, st := range p.Steps {
			e.rel[t] = h.RelIndex(st.Rel)
			for _, s := range st.Symbols() {
				if s == "" {
					continue
				}
				e.obs[t] = append(e.obs[t], h.SymbolIndex(s))
			}
		}
		out = append(out, e)
	}
	return out
}

// counts are the expected sufficient statistics accumulated by the E-step.
type counts struct {
	init  []float64
	trans [][][]float64
	emit  [][]float64
}

func newCounts(ns, nr, nk int) *counts {
	c := &counts{
		init:  make([]float64, ns),
		trans: make([][][]float64, nr),
		emit:  make([][]float64, ns),
	}
	for r := range c.trans {
		c.trans[r] = make([][]float64, ns)
		for i := range c.trans[r] {
			c.trans[r][i] = make([]float64, ns)
		}
	}
	for i := range c.emit {
		c.emit[i] = make([]float64, nk)
	}
	return c
}

// expectation runs forward-backward over every path and returns the corpus
// log-likelihood together with the expected counts.
func expectation(h *HMM, paths []encoded) (float64, *counts) {
	ns, nr, nk := len(h.states), len(h.rels), len(h.symbols)
	c := newCounts(ns, nr, nk)
	total := 0.0

	for _, p := range paths {
		T := len(p.obs)
		emit := make([][]float64, T)
		for t := 0; t < T; t++ {
			emit[t] = make([]float64, ns)
			for i := 0; i < ns; i++ {
				emit[t][i] = h.logEmitIdx(i, p.obs[t])
			}
		}

		alpha := forwardPath(h, p, emit)
		beta := backwardPath(h, p, emit)
		logZ := LogSumExp(alpha[T-1])
		if math.IsInf(logZ, -1) || math.IsNaN(logZ) {
			// A path impossible under the current tables contributes nothing.
			// Smoothing makes this unreachable in practice; the guard is here
			// so a hand-built starting point cannot poison the counts with NaN.
			continue
		}
		total += logZ

		scratch := make([]float64, ns)
		for t := 0; t < T; t++ {
			for i := 0; i < ns; i++ {
				scratch[i] = alpha[t][i] + beta[t][i] - logZ
			}
			for i := 0; i < ns; i++ {
				g := math.Exp(scratch[i])
				if g <= 0 {
					continue
				}
				if t == 0 {
					c.init[i] += g
				}
				for _, k := range p.obs[t] {
					c.emit[i][k] += g
				}
			}
			if t+1 >= T {
				continue
			}
			r := p.rel[t+1]
			for i := 0; i < ns; i++ {
				for j := 0; j < ns; j++ {
					xi := alpha[t][i] + h.logTrans[r][i][j] + emit[t+1][j] + beta[t+1][j] - logZ
					if v := math.Exp(xi); v > 0 {
						c.trans[r][i][j] += v
					}
				}
			}
		}
	}
	return total, c
}

// forwardPath is the training-time forward recursion. It is unnormalized, so
// the last column sums to the path likelihood; the run-time Forward normalizes
// every step instead, because it needs a distribution rather than an evidence
// term.
func forwardPath(h *HMM, p encoded, emit [][]float64) [][]float64 {
	ns, T := len(h.states), len(p.obs)
	alpha := make([][]float64, T)
	alpha[0] = make([]float64, ns)
	for i := 0; i < ns; i++ {
		alpha[0][i] = h.logInit[i] + emit[0][i]
	}
	scratch := make([]float64, ns)
	for t := 1; t < T; t++ {
		alpha[t] = make([]float64, ns)
		trans := h.logTrans[p.rel[t]]
		for j := 0; j < ns; j++ {
			for i := 0; i < ns; i++ {
				scratch[i] = alpha[t-1][i] + trans[i][j]
			}
			alpha[t][j] = LogSumExp(scratch) + emit[t][j]
		}
	}
	return alpha
}

// backwardPath is the training-only backward recursion (see Train's doc).
func backwardPath(h *HMM, p encoded, emit [][]float64) [][]float64 {
	ns, T := len(h.states), len(p.obs)
	beta := make([][]float64, T)
	beta[T-1] = make([]float64, ns) // log 1
	scratch := make([]float64, ns)
	for t := T - 2; t >= 0; t-- {
		beta[t] = make([]float64, ns)
		trans := h.logTrans[p.rel[t+1]]
		for i := 0; i < ns; i++ {
			for j := 0; j < ns; j++ {
				scratch[j] = trans[i][j] + emit[t+1][j] + beta[t+1][j]
			}
			beta[t][i] = LogSumExp(scratch)
		}
	}
	return beta
}

// maximization replaces the tables with the Dirichlet posterior means of the
// expected counts.
func maximization(h *HMM, c *counts, sm Smoothing) {
	h.logInit = dirichletLog(c.init, sm.Init)
	for r := range h.logTrans {
		for i := range h.logTrans[r] {
			h.logTrans[r][i] = dirichletLog(c.trans[r][i], sm.Trans)
		}
	}
	for i := range h.logEmit {
		h.logEmit[i] = dirichletLog(c.emit[i], sm.Emit)
	}
}

// logPrior is the Dirichlet log-density of every table, the penalty term that
// makes the EM objective monotone.
func logPrior(h *HMM, sm Smoothing) float64 {
	sum := dirichletLogPrior(h.logInit, sm.Init)
	for r := range h.logTrans {
		for i := range h.logTrans[r] {
			sum += dirichletLogPrior(h.logTrans[r][i], sm.Trans)
		}
	}
	for i := range h.logEmit {
		sum += dirichletLogPrior(h.logEmit[i], sm.Emit)
	}
	if math.IsNaN(sum) {
		return LogZero
	}
	return sum
}

// logEmitIdx is LogEmission over already-resolved symbol indices.
func (h *HMM) logEmitIdx(state int, obs []int) float64 {
	sum := 0.0
	for _, k := range obs {
		sum += h.logEmit[state][k]
	}
	return sum
}

// LogLikelihood returns the corpus log-likelihood under this model. It is a
// training diagnostic — comparing two models on the same corpus — and never a
// number a user sees.
func (h *HMM) LogLikelihood(c Corpus) float64 {
	ll, _ := expectation(h, encodePaths(h, usablePaths(c)))
	return ll
}

// rng is splitmix64, carried here rather than taken from math/rand so that a
// recorded seed reproduces the same model on any Go version and any platform.
// Reproducibility is the whole point of recording the seed; borrowing a
// generator whose stream is only promised for the current major version would
// undermine it.
type rng struct{ state uint64 }

func newRNG(seed int64) *rng { return &rng{state: uint64(seed)} }

func (r *rng) next() uint64 {
	r.state += 0x9e3779b97f4a7c15
	z := r.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// float returns a value in [0,1) from the top 53 bits, which is where the
// generator's quality is highest.
func (r *rng) float() float64 {
	return float64(r.next()>>11) / float64(uint64(1)<<53)
}
