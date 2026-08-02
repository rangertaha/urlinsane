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
	"fmt"
	"math"
)

// Smoothing holds the Dirichlet concentration parameters applied to every table
// row during training (§10.4). Each is the pseudo-count added to every cell of a
// row before it is normalized.
//
// Smoothing is not cosmetic here. Expansion traces are sparse: a relation may
// appear a handful of times and a prop value once. Without a prior, an
// unobserved transition gets probability zero, and a single zero anywhere on a
// path annihilates the whole forward filter — a node the model has simply never
// seen becomes impossible rather than merely unlikely, and gets pruned with
// certainty. The prior is what keeps "unseen" distinct from "ruled out".
//
// The values are also carried in the encoded model, so a retrained model is
// comparable with the one it replaces.
type Smoothing struct {
	// Init is the prior on the seed's initial distribution.
	Init float64
	// Trans is the prior on each transition row, per relation.
	Trans float64
	// Emit is the prior on each state's emission row. It is what gives
	// OOVSymbol its probability mass: the OOV symbol is never observed in
	// training, so its entire posterior weight comes from this prior.
	Emit float64
}

// DefaultSmoothing is add-one (Laplace) smoothing on every table. It is the
// weakest prior that still guarantees no zero anywhere, which is the property
// the forward filter actually needs.
func DefaultSmoothing() Smoothing {
	return Smoothing{Init: 1, Trans: 1, Emit: 1}
}

// IsZero reports an unset Smoothing, so that a decoded model can tell "no
// smoothing recorded" from "smoothing recorded as zero".
func (s Smoothing) IsZero() bool { return s == Smoothing{} }

// validate rejects priors that would defeat their own purpose. A negative
// pseudo-count can drive a probability below zero, and training requires a
// strictly positive one on every table so that no cell can be exactly zero.
func (s Smoothing) validate(training bool) error {
	for _, f := range []struct {
		name string
		v    float64
	}{{"init", s.Init}, {"trans", s.Trans}, {"emit", s.Emit}} {
		if math.IsNaN(f.v) || math.IsInf(f.v, 0) {
			return fmt.Errorf("model: %s smoothing must be finite, got %v", f.name, f.v)
		}
		if f.v < 0 {
			return fmt.Errorf("model: %s smoothing must not be negative, got %v", f.name, f.v)
		}
		if training && f.v <= 0 {
			return fmt.Errorf("model: %s smoothing must be positive for training, got %v", f.name, f.v)
		}
	}
	return nil
}

// dirichletLog turns expected counts into the log of the Dirichlet posterior
// mean, (count + alpha) / (total + n*alpha).
//
// The posterior mean rather than the mode: the mode is undefined for alpha < 1
// and puts zeros back into the table, which is precisely what the prior is here
// to prevent.
func dirichletLog(counts []float64, alpha float64) []float64 {
	n := len(counts)
	out := make([]float64, n)
	total := float64(n) * alpha
	for _, c := range counts {
		if c > 0 {
			total += c
		}
	}
	if total <= 0 {
		// No counts and no prior: nothing distinguishes the cells, so the
		// honest answer is uniform.
		u := -math.Log(float64(n))
		for i := range out {
			out[i] = u
		}
		return out
	}
	logTotal := math.Log(total)
	for i, c := range counts {
		if c < 0 {
			c = 0
		}
		v := c + alpha
		if v <= 0 {
			out[i] = LogZero
			continue
		}
		out[i] = math.Log(v) - logTotal
	}
	return out
}

// dirichletLogPrior is the penalty term a row contributes to the EM objective,
// up to a constant that does not change between iterations.
//
// The exponent is alpha, not alpha-1, because dirichletLog returns the
// posterior *mean*. p = (c+alpha)/(N+n*alpha) is the maximizer of
// sum_k (c_k + alpha) log p_k over the simplex, so the term that pairs with it
// is sum_k alpha*log p_k. Using alpha-1 — the density's own exponent — would
// pair the objective with the posterior mode instead, and the monotonicity
// claim would be off by one prior count and quietly false.
//
// Baum-Welch under a prior is MAP-EM: what it guarantees to increase
// monotonically is the log-likelihood plus this term, not the bare
// log-likelihood. Train reports both so the guarantee is testable rather than
// assumed.
func dirichletLogPrior(logRow []float64, alpha float64) float64 {
	if alpha == 0 {
		return 0 // No prior, no penalty; also avoids 0 * -Inf.
	}
	sum := 0.0
	for _, lp := range logRow {
		if math.IsInf(lp, -1) {
			// An impossible cell has zero prior density. Training never
			// produces one, because validate demands a positive alpha.
			return LogZero
		}
		sum += alpha * lp
	}
	return sum
}
