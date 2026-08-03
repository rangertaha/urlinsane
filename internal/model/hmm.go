// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package model is the HMM library used by the engine and by plugins alike
// (docs/DESIGN.md §10.6). It provides the primitives — forward filtering,
// Baum-Welch, log-space arithmetic, Dirichlet smoothing — plus the dag-cbor
// artifact format that gives a model a CID.
//
// It is a library, not a model. The engine instantiates one to steer expansion;
// an operator instantiates its own to decide what it emits; neither knows about
// the other's states or alphabet.
//
// # What a model may be used for
//
// The engine's model steers execution and nothing else: frontier ordering,
// pruning, and operator gating (§10.2). It must never produce a number a user
// reads — no report score, no finding, no severity. That restriction is what
// removes the calibration requirement: a miscalibrated model that mis-ranks the
// frontier wastes network calls, whereas a miscalibrated model in a report makes
// a false accusation. Nothing in this package is named or shaped like a
// severity, and Belief deliberately exposes only what graph.BeliefModel needs.
//
// # Inference
//
// Forward filtering is the only inference offered at run time: no Viterbi, no
// backward pass, no belief propagation (§10.1). A node's belief is a pure
// function of its parent's belief and its own props, which is what makes it
// available *during* execution, where it can still save a network call. The
// backward recursion exists in train.go only, because Baum-Welch needs it
// offline; it is unexported and never runs during a scan.
//
// # Numerics
//
// All table arithmetic is in log space. An expansion path is short but an
// emission is a product over every observed prop, so linear-space
// probabilities underflow quickly and silently; log space makes the arithmetic
// stable and turns products into sums.
package model

import (
	"fmt"
	"math"
	"sort"
)

const (
	// OOVSymbol is the explicit out-of-vocabulary emission symbol. Every
	// alphabet contains it, and any symbol a featurizer produces that training
	// never saw is mapped onto it. Without an explicit symbol an unseen prop
	// would either be dropped — silently discarding evidence — or annihilate
	// the filter with a zero probability.
	OOVSymbol = "<oov>"

	// OOVRelation is the transition table used for a relation the model does
	// not know. Plans add relations over time; a model trained before one
	// existed must still be usable rather than panicking mid-scan.
	OOVRelation = "<oov>"
)

// LogZero is the log of an impossible event. Log-space arithmetic is closed
// over it: LogAdd and LogSumExp both handle it without producing NaN.
var LogZero = math.Inf(-1)

// HMM is a discrete hidden Markov model whose transition table is conditioned
// on the name of the relation taken, so that sharing an IP can transmit far
// more belief than sharing a TLD (§10).
//
// All tables are held in log space. The type is immutable once built: Forward
// allocates its own output, so one HMM is safe for concurrent use by every
// worker in a round.
type HMM struct {
	states  []string
	rels    []string
	symbols []string

	// focus are the state indices whose posterior mass makes up the scalar
	// belief the engine consumes. See Belief in belief.go for why a scalar is
	// all graph.BeliefModel can carry.
	focus    []int
	focusSet []bool

	logInit  []float64     // [state]
	logTrans [][][]float64 // [rel][from][to]
	logEmit  [][]float64   // [state][symbol]

	smooth Smoothing
	prov   Provenance

	stateIdx map[string]int
	relIdx   map[string]int
	symIdx   map[string]int
	oovRel   int
	oovSym   int
}

// Spec describes a model in ordinary probabilities, which is how a human writes
// one down. New normalizes and converts to log space; a decoded model bypasses
// this path so that a round trip is bit-exact (see codec.go).
type Spec struct {
	// States is the latent status alphabet, in a fixed order that is part of
	// the model's identity.
	States []string
	// Focus names the states whose posterior mass is reported as belief. It
	// must be non-empty. Naming every state makes belief constant, which is the
	// untrained degenerate case.
	Focus []string
	// Rels is the relation alphabet. OOVRelation is appended if absent.
	Rels []string
	// Symbols is the emission alphabet. OOVSymbol is appended if absent.
	Symbols []string

	// Init is P(state) for the seed, which has no parent. Barrier 0 assigns it.
	Init []float64
	// Trans is P(state | parent state, rel), keyed by relation name; rows are
	// indexed by parent state. A relation absent from the map falls back to the
	// OOVRelation table, and an absent OOVRelation table is uniform.
	Trans map[string][][]float64
	// Emit is P(symbol | state). A model with no emissions at all is a pure
	// transition prior, which is exactly what a candidate node gets (§10.2).
	Emit [][]float64

	// Smoothing records the Dirichlet priors the tables were fitted under. It
	// is carried for provenance and reused when the model is retrained.
	Smoothing Smoothing
	// Provenance records how the model was produced.
	Provenance Provenance
}

// New builds an HMM from a Spec, normalizing every distribution. Rows that sum
// to zero become uniform rather than an error: a relation nobody has observed is
// a normal state of affairs, and the honest prior for it is "no information".
func New(s Spec) (*HMM, error) {
	if len(s.States) == 0 {
		return nil, fmt.Errorf("model: a model needs at least one state")
	}
	h := &HMM{
		states:  append([]string(nil), s.States...),
		rels:    withReserved(s.Rels, OOVRelation),
		symbols: withReserved(s.Symbols, OOVSymbol),
		smooth:  s.Smoothing,
		prov:    s.Provenance,
	}
	var err error
	if h.stateIdx, err = index("state", h.states); err != nil {
		return nil, err
	}
	if h.relIdx, err = index("relation", h.rels); err != nil {
		return nil, err
	}
	if h.symIdx, err = index("symbol", h.symbols); err != nil {
		return nil, err
	}
	h.oovRel, h.oovSym = h.relIdx[OOVRelation], h.symIdx[OOVSymbol]

	if err := h.setFocus(s.Focus); err != nil {
		return nil, err
	}

	ns, nr, nk := len(h.states), len(h.rels), len(h.symbols)

	if len(s.Init) != 0 && len(s.Init) != ns {
		return nil, fmt.Errorf("model: initial distribution has %d entries, want %d", len(s.Init), ns)
	}
	h.logInit = logNormalized(s.Init, ns)

	h.logTrans = make([][][]float64, nr)
	for ri, rel := range h.rels {
		rows := s.Trans[rel]
		if rows != nil && len(rows) != ns {
			return nil, fmt.Errorf("model: transition table for %q has %d rows, want %d", rel, len(rows), ns)
		}
		h.logTrans[ri] = make([][]float64, ns)
		for i := 0; i < ns; i++ {
			var row []float64
			if rows != nil {
				row = rows[i]
				if len(row) != 0 && len(row) != ns {
					return nil, fmt.Errorf("model: transition row %d of %q has %d entries, want %d", i, rel, len(row), ns)
				}
			}
			h.logTrans[ri][i] = logNormalized(row, ns)
		}
	}

	// `s.Emit != nil`, matching the transition guard above, not `len(s.Emit) != 0`.
	// The two say different things about an empty-but-present table: a caller
	// passing `Emit: [][]float64{}` has a table with zero rows, which is a
	// malformed spec, but the length form read it as "no table supplied" and
	// waved it through — and the row read below tests `s.Emit != nil`, which is
	// true for an empty slice, so it indexed straight past the end and panicked
	// out of a constructor whose whole job is to reject malformed specs.
	if s.Emit != nil && len(s.Emit) != ns {
		return nil, fmt.Errorf("model: emission table has %d rows, want %d", len(s.Emit), ns)
	}
	h.logEmit = make([][]float64, ns)
	for i := 0; i < ns; i++ {
		var row []float64
		if s.Emit != nil {
			row = s.Emit[i]
			if len(row) != 0 && len(row) != nk {
				return nil, fmt.Errorf("model: emission row %d has %d entries, want %d", i, len(row), nk)
			}
		}
		h.logEmit[i] = logNormalized(row, nk)
	}
	return h, nil
}

// Uniform is the untrained model. It has one state, so its posterior is a point
// mass and its belief is exactly 1 for every node, every relation and every
// observation — bit-identical to graph's default model.
//
// This is the property §10.5 rests on: a uniform model reduces expansion to
// breadth-first and unranked, so the engine ships and runs correctly before any
// model exists, and a model that turns out poor can be dropped without
// invalidating a single result. A single state also makes the reduction exact
// rather than approximate: with N uniform states the belief would be the
// constant |focus|/N, equal for every node and therefore still unranked, but no
// longer the literal 1 that graph's own default returns.
func Uniform() *HMM {
	h, err := New(Spec{
		States:     []string{"any"},
		Focus:      []string{"any"},
		Provenance: Provenance{Algorithm: AlgorithmUniform},
	})
	if err != nil {
		// Unreachable: the spec is a constant and New only rejects malformed
		// ones. Panicking beats returning an error nobody can act on.
		panic("model: uniform model is invalid: " + err.Error())
	}
	return h
}

func (h *HMM) setFocus(names []string) error {
	h.focusSet = make([]bool, len(h.states))
	if len(names) == 0 {
		return fmt.Errorf("model: focus states must be named; belief is their posterior mass")
	}
	seen := map[int]bool{}
	for _, n := range names {
		i, ok := h.stateIdx[n]
		if !ok {
			return fmt.Errorf("model: focus state %q is not a state", n)
		}
		if seen[i] {
			continue
		}
		seen[i] = true
		h.focusSet[i] = true
	}
	for i := range h.states {
		if h.focusSet[i] {
			h.focus = append(h.focus, i)
		}
	}
	return nil
}

// States returns the state alphabet in index order.
func (h *HMM) States() []string { return append([]string(nil), h.states...) }

// Rels returns the relation alphabet in index order, including OOVRelation.
func (h *HMM) Rels() []string { return append([]string(nil), h.rels...) }

// Symbols returns the emission alphabet in index order, including OOVSymbol.
func (h *HMM) Symbols() []string { return append([]string(nil), h.symbols...) }

// Focus returns the names of the states whose mass is reported as belief.
func (h *HMM) Focus() []string {
	out := make([]string, 0, len(h.focus))
	for _, i := range h.focus {
		out = append(out, h.states[i])
	}
	return out
}

// Smoothing returns the Dirichlet priors this model was fitted under.
func (h *HMM) Smoothing() Smoothing { return h.smooth }

// Provenance returns the training provenance recorded in the model block.
func (h *HMM) Provenance() Provenance { return h.prov }

// Prior returns the seed's log distribution — the one node with no parent.
func (h *HMM) Prior() []float64 { return append([]float64(nil), h.logInit...) }

// RelIndex resolves a relation name, falling back to OOVRelation.
func (h *HMM) RelIndex(rel string) int {
	if i, ok := h.relIdx[rel]; ok {
		return i
	}
	return h.oovRel
}

// SymbolIndex resolves an emission symbol, falling back to OOVSymbol.
func (h *HMM) SymbolIndex(sym string) int {
	if i, ok := h.symIdx[sym]; ok {
		return i
	}
	return h.oovSym
}

// LogTransition returns log P(to | from, rel).
func (h *HMM) LogTransition(rel string, from, to int) float64 {
	return h.logTrans[h.RelIndex(rel)][from][to]
}

// LogEmission returns log P(obs | state) for a whole observation.
//
// An observation is a set of symbols, one per prop the node has as of the
// current barrier, and they are combined as a product — a naive-Bayes
// factorization. Modelling the joint distribution over prop combinations
// instead would need a table exponential in the number of props and a corpus to
// match, for a quantity that only has to rank a frontier.
//
// An empty observation contributes nothing, which is the correct answer for a
// candidate: it has no props yet, so its belief is the transition prior alone
// (§10.2).
func (h *HMM) LogEmission(state int, obs []string) float64 {
	sum := 0.0
	for _, s := range obs {
		sum += h.logEmit[state][h.SymbolIndex(s)]
	}
	return sum
}

// Forward is the forward filtering step, and the only inference this package
// performs at run time (§10.1).
//
// prev is the parent's normalized log distribution; the result is this node's,
// also normalized. Normalizing every step is what keeps a long path from
// underflowing and makes the result independent of any constant factor in the
// emission tables.
//
// If the observation is impossible under every state the filter would collapse
// to NaN, so the result falls back to uniform. That can only happen to a
// hand-written model with hard zeros: Dirichlet smoothing plus an explicit OOV
// symbol keeps every trained probability strictly positive.
func (h *HMM) Forward(prev []float64, rel string, obs []string) []float64 {
	ns := len(h.states)
	out := make([]float64, ns)
	if len(prev) != ns {
		prev = h.logInit
	}
	trans := h.logTrans[h.RelIndex(rel)]
	scratch := make([]float64, ns)
	for j := 0; j < ns; j++ {
		for i := 0; i < ns; i++ {
			scratch[i] = prev[i] + trans[i][j]
		}
		out[j] = LogSumExp(scratch) + h.LogEmission(j, obs)
	}
	normalizeLog(out)
	return out
}

// Mass returns the focus states' share of a log distribution — the scalar the
// engine calls belief. It is clamped to [0,1] so that accumulated floating
// point error can never hand the scheduler a value outside the range its
// thresholds assume.
func (h *HMM) Mass(logDist []float64) float64 {
	if len(logDist) != len(h.states) {
		return 0
	}
	terms := make([]float64, 0, len(h.focus))
	for _, i := range h.focus {
		terms = append(terms, logDist[i])
	}
	m := math.Exp(LogSumExp(terms))
	switch {
	case math.IsNaN(m):
		return 0
	case m < 0:
		return 0
	case m > 1:
		return 1
	}
	return m
}

// Lift reconstructs a log distribution from a scalar belief.
//
// graph.BeliefModel carries a float64 between parent and child, not a
// distribution, so the parent's full posterior is not available to Step and has
// to be reconstructed. Lift puts mass b on the focus states and 1-b on the
// rest, shaped within each group by the initial distribution — the
// maximum-entropy choice consistent with both the scalar and the model's own
// prior.
//
// The reconstruction is exact when the focus set is a single state and there is
// one non-focus state, which is the canonical binary model, and exact when
// every state is in focus, which is the uniform model. It is an approximation
// only for a model with three or more states, where the scalar genuinely cannot
// carry the shape. See the note in belief.go.
func (h *HMM) Lift(b float64) []float64 {
	ns := len(h.states)
	if len(h.focus) == ns {
		return append([]float64(nil), h.logInit...)
	}
	if b < 0 {
		b = 0
	} else if b > 1 {
		b = 1
	}
	in := make([]float64, 0, ns)
	out := make([]float64, 0, ns)
	for i := 0; i < ns; i++ {
		if h.focusSet[i] {
			in = append(in, h.logInit[i])
		} else {
			out = append(out, h.logInit[i])
		}
	}
	logIn, logOut := LogSumExp(in), LogSumExp(out)
	dist := make([]float64, ns)
	for i := 0; i < ns; i++ {
		if h.focusSet[i] {
			dist[i] = share(h.logInit[i], logIn, len(in)) + math.Log(b)
		} else {
			dist[i] = share(h.logInit[i], logOut, len(out)) + math.Log(1-b)
		}
	}
	normalizeLog(dist)
	return dist
}

// share is a state's log share within its group, falling back to uniform when
// the group carries no prior mass at all.
func share(logP, logGroup float64, n int) float64 {
	if math.IsInf(logGroup, -1) {
		return -math.Log(float64(n))
	}
	return logP - logGroup
}

// LogAdd returns log(exp(a) + exp(b)) without leaving log space. Subtracting
// the larger term first keeps exp in range, so LogZero inputs give LogZero
// rather than NaN.
func LogAdd(a, b float64) float64 {
	if a < b {
		a, b = b, a
	}
	if math.IsInf(a, -1) {
		return LogZero
	}
	return a + math.Log1p(math.Exp(b-a))
}

// LogSumExp returns log(sum(exp(x))) using the max-shift trick. An empty or
// all-LogZero input gives LogZero.
func LogSumExp(xs []float64) float64 {
	max := LogZero
	for _, x := range xs {
		if x > max {
			max = x
		}
	}
	if math.IsInf(max, -1) {
		return LogZero
	}
	sum := 0.0
	for _, x := range xs {
		sum += math.Exp(x - max)
	}
	return max + math.Log(sum)
}

// normalizeLog scales a log distribution to sum to one in place and returns the
// log evidence it divided out. A distribution with no mass anywhere becomes
// uniform, because the alternative is propagating NaN into the scheduler.
func normalizeLog(dst []float64) float64 {
	z := LogSumExp(dst)
	if math.IsInf(z, -1) || math.IsNaN(z) {
		u := -math.Log(float64(len(dst)))
		for i := range dst {
			dst[i] = u
		}
		return LogZero
	}
	for i := range dst {
		dst[i] -= z
	}
	return z
}

// logNormalized converts probabilities to a normalized log distribution of
// length n. A nil, empty or zero-sum input yields the uniform distribution.
func logNormalized(p []float64, n int) []float64 {
	out := make([]float64, n)
	sum := 0.0
	for _, v := range p {
		if v > 0 {
			sum += v
		}
	}
	if len(p) != n || sum <= 0 {
		u := -math.Log(float64(n))
		for i := range out {
			out[i] = u
		}
		return out
	}
	for i := 0; i < n; i++ {
		if p[i] > 0 {
			out[i] = math.Log(p[i] / sum)
		} else {
			out[i] = LogZero
		}
	}
	return out
}

// withReserved appends a reserved name to an alphabet if it is missing. The
// alphabet's order is otherwise preserved, because index order is part of the
// model's encoded identity.
func withReserved(names []string, reserved string) []string {
	out := append([]string(nil), names...)
	for _, n := range out {
		if n == reserved {
			return out
		}
	}
	return append(out, reserved)
}

// index builds a name-to-position map and rejects duplicates, which would
// otherwise make a table row unreachable and impossible to debug.
func index(what string, names []string) (map[string]int, error) {
	m := make(map[string]int, len(names))
	for i, n := range names {
		if n == "" {
			return nil, fmt.Errorf("model: %s %d has no name", what, i)
		}
		if _, dup := m[n]; dup {
			return nil, fmt.Errorf("model: %s %q declared twice", what, n)
		}
		m[n] = i
	}
	return m, nil
}

// sortedKeys returns a map's keys in a fixed order. Training derives alphabets
// from a corpus, and map iteration order would otherwise make two runs over the
// same corpus produce different state indices and therefore different CIDs.
func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
