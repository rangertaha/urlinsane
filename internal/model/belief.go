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

package model

import (
	"encoding/hex"
	"sort"
	"strconv"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Featurizer turns a node's view into emission symbols — the observation the
// forward step conditions on.
//
// It is the whole of what a model knows about the graph. Keeping it a function
// supplied by the caller is what makes this package a library rather than the
// engine's model (§10.6): the engine featurizes the props its plan produces, a
// variant operator featurizes whatever its own model needs, and neither
// package has to know the other's alphabet.
//
// A featurizer must be a pure function of the view. Anything it reads from
// outside — a clock, a counter, a map iteration — makes belief depend on
// something other than (parent belief, props), and the reproducibility the
// whole barrier design exists to protect is gone.
type Featurizer func(graph.View) []string

// Belief adapts an HMM to graph.BeliefModel.
//
// It exists only to steer execution: frontier ordering, pruning, and
// BeliefAbove gates (§10.2). Nothing it returns is a report score, a finding or
// a severity, and analyzers never see it. That restriction is what makes a
// miscalibrated model a waste of network calls rather than a false accusation.
type Belief struct {
	h        *HMM
	features Featurizer
}

// NewBelief wraps a model as the engine's (or a plugin's) belief model. A nil
// featurizer means no observations at all, which leaves belief as the pure
// transition prior — exactly the state a candidate is in.
func NewBelief(h *HMM, f Featurizer) *Belief {
	if h == nil {
		h = Uniform()
	}
	return &Belief{h: h, features: f}
}

// UniformBelief is the model the engine runs with before one has been trained.
// It returns 1 for every node, so the frontier sorts on (depth, type, key)
// alone and no candidate is ever pruned: breadth-first, unranked, which is
// exactly the pre-model behaviour (§10.5).
func UniformBelief() graph.BeliefModel { return NewBelief(Uniform(), nil) }

// Model returns the underlying HMM, so a caller can pin its CID into the plan
// hash.
func (b *Belief) Model() *HMM { return b.h }

// Initial is the seed's prior. The seed has no parent, so barrier 0 takes it
// straight from the initial distribution — and carries that distribution
// forward as the state its children step from.
func (b *Belief) Initial() (float64, graph.State) {
	return b.h.Mass(b.h.logInit), b.dist(b.h.logInit)
}

// Step is the forward filtering step: the parent's belief pushed through the
// relation that admitted this node, then conditioned on the node's props as of
// this barrier.
//
// It is recomputed from scratch at every barrier rather than accumulated, which
// is what makes it independent of which operator returned first — only which
// round an edge belongs to matters, and that is deterministic (§10.3).
//
// The parent's full posterior travels in graph.State, so nothing is
// reconstructed and nothing is lost. An earlier version received only the
// parent's scalar and had to Lift it back into a distribution — exact for two
// states, and for three or more the maximum-entropy distribution consistent
// with that one number, which is not the posterior the parent actually had.
func (b *Belief) Step(parent graph.State, rel string, v graph.View) (float64, graph.State) {
	var obs []string
	if b.features != nil && v != nil {
		obs = b.features(v)
	}
	next := b.h.Forward(b.prior(parent), rel, obs)
	return b.h.Mass(next), b.dist(next)
}

// prior recovers the parent's log-space distribution.
//
// A nil state means the parent has not been stepped yet, which the interface
// defines as "use the initial distribution". A state of the wrong width means
// the model was swapped mid-run — impossible today, but falling back beats
// indexing past the end of a slice.
func (b *Belief) prior(s graph.State) []float64 {
	d, ok := s.([]float64)
	if !ok || len(d) != len(b.h.states) {
		return b.h.logInit
	}
	return d
}

// dist hands out a copy. The graph holds this value for the node's lifetime and
// Forward allocates a fresh slice per call today, but a caller that later
// reused a buffer would silently rewrite every descendant's recorded state.
func (b *Belief) dist(d []float64) graph.State {
	return append([]float64(nil), d...)
}

// PropFeatures builds a Featurizer that renders the named fields as
// "field=value" emission symbols, skipping fields that are unset.
//
// Skipping rather than emitting "field=" is deliberate: emissions are a product
// over the symbols present, so an absent prop contributes nothing and a node
// observed in round g is not penalized for props that only arrive in round g+1.
//
// Field names are sorted so the symbol order is fixed. Order does not change
// the product, but a stable order keeps traces recorded by one run comparable
// with another's byte for byte.
//
// Continuous fields — a registration timestamp, a response size — should be
// bucketed by the caller before they reach here. Rendered exactly, every
// distinct value becomes its own symbol, the alphabet explodes and every symbol
// falls to OOV.
func PropFeatures(fields ...string) Featurizer {
	names := append([]string(nil), fields...)
	sort.Strings(names)
	return func(v graph.View) []string {
		if v == nil {
			return nil
		}
		out := make([]string, 0, len(names))
		for _, f := range names {
			val, ok := v.Prop(f)
			if !ok {
				continue
			}
			out = append(out, f+"="+renderValue(val))
		}
		return out
	}
}

// renderValue formats a prop value as a symbol. Every kind gets an exact,
// locale-free rendering so that the same value always produces the same symbol.
func renderValue(v graph.Value) string {
	switch v.Kind() {
	case graph.KindString:
		return v.Str()
	case graph.KindInt:
		return strconv.FormatInt(v.Num(), 10)
	case graph.KindFloat:
		return strconv.FormatFloat(v.Real(), 'g', -1, 64)
	case graph.KindBool:
		return strconv.FormatBool(v.Flag())
	case graph.KindBytes:
		return hex.EncodeToString(v.Raw())
	case graph.KindTime:
		return strconv.FormatInt(v.Num(), 10)
	}
	return ""
}
