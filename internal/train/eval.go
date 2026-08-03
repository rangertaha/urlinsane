// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package train

import (
	"fmt"
	"sort"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Scored is one node's belief against what was observed about it.
type Scored struct {
	Key    string
	Type   string
	Belief float64
	Live   bool
}

// Quality is how well belief ordered a scan.
//
// The question the whole execution model exists to answer: if the frontier is
// bounded, are the nodes it keeps the ones that turn out to be registered? A
// model that separates cleanly but in the wrong direction scores below 0.5 and
// is worse than no model, which is why AUC is reported rather than accuracy.
type Quality struct {
	// Scored is how many nodes could be judged at all.
	Scored int
	Live   int
	Absent int
	// Skipped is nodes whose existence was unknown or untried. They are not
	// counted as negatives — that is the whole three-state discipline, and
	// scoring a model against observations nobody made would reward it for
	// agreeing with the network's failures.
	Skipped int

	// AUC is the probability that a randomly chosen live node outranks a
	// randomly chosen absent one. 0.5 is a coin toss, which is what the uniform
	// model scores by construction.
	AUC float64
	// PrecisionAt is the fraction of the top k that were live, for a few k.
	PrecisionAt map[int]float64
	// BaseRate is the fraction of judged nodes that were live — what
	// precision@k would be if the order were random. Precision above the base
	// rate is the only kind that means anything.
	BaseRate float64
}

func (q Quality) String() string {
	return fmt.Sprintf("scored %d (live %d, absent %d, skipped %d) AUC=%.3f base=%.3f p@10=%.3f p@25=%.3f p@50=%.3f",
		q.Scored, q.Live, q.Absent, q.Skipped, q.AUC, q.BaseRate,
		q.PrecisionAt[10], q.PrecisionAt[25], q.PrecisionAt[50])
}

// Rank orders a graph's nodes by the belief currently installed on it.
//
// Ties are broken by (type, key) so the order is total and two runs agree. A
// uniform model makes every belief equal, so its ranking is alphabetical — the
// honest representation of "no information", and the baseline everything else
// has to beat.
func Rank(g *graph.Graph) []Scored {
	a := g.Analyze()
	out := make([]Scored, 0, len(a.Nodes()))
	for _, n := range a.Nodes() {
		out = append(out, Scored{
			Key:    n.Key,
			Type:   n.Type.Name(),
			Belief: g.Belief(n.ID),
			Live:   a.Existence(n.ID) == graph.Live,
		})
	}
	sortScored(out)
	return out
}

// Evaluate measures how well the installed belief model ordered a scan.
//
// Only nodes with a settled observation are judged. A node whose lookups all
// failed is not evidence either way, and counting it as absent would score the
// model on the network's behaviour rather than on the name's.
func Evaluate(g *graph.Graph) Quality {
	a := g.Analyze()
	q := Quality{PrecisionAt: map[int]float64{}}

	var judged []Scored
	for _, n := range a.Nodes() {
		s := Scored{Key: n.Key, Type: n.Type.Name(), Belief: g.Belief(n.ID)}
		switch a.Existence(n.ID) {
		case graph.Live:
			s.Live = true
			q.Live++
		case graph.Absent:
			q.Absent++
		default:
			q.Skipped++
			continue
		}
		judged = append(judged, s)
	}
	sortScored(judged)

	q.Scored = len(judged)
	if q.Scored == 0 || q.Live == 0 || q.Absent == 0 {
		// One class only: AUC is undefined and any number would be a claim
		// nobody can check.
		return q
	}
	q.BaseRate = float64(q.Live) / float64(q.Scored)
	q.AUC = auc(judged)

	for _, k := range []int{10, 25, 50} {
		n := k
		if n > len(judged) {
			n = len(judged)
		}
		if n == 0 {
			continue
		}
		var hits float64
		for _, s := range judged[:n] {
			if s.Live {
				hits++
			}
		}
		q.PrecisionAt[k] = hits / float64(n)
	}
	return q
}

// auc is the probability that a randomly chosen live node outranks a randomly
// chosen absent one, ties counting a half.
//
// The half is what puts a uniform model at exactly 0.5 rather than at whatever
// its tie-break order happens to produce — so the baseline is the coin toss it
// should be, and a model below it is doing harm.
func auc(judged []Scored) float64 {
	var concordant, pairs float64
	for i := range judged {
		for j := i + 1; j < len(judged); j++ {
			if judged[i].Live == judged[j].Live {
				continue
			}
			live, dead := judged[i], judged[j]
			if !live.Live {
				live, dead = dead, live
			}
			pairs++
			switch {
			case live.Belief > dead.Belief:
				concordant++
			case live.Belief == dead.Belief:
				concordant += 0.5
			}
		}
	}
	if pairs == 0 {
		return 0
	}
	return concordant / pairs
}

func sortScored(s []Scored) {
	sort.Slice(s, func(i, j int) bool {
		if s[i].Belief != s[j].Belief {
			return s[i].Belief > s[j].Belief
		}
		if s[i].Type != s[j].Type {
			return s[i].Type < s[j].Type
		}
		return s[i].Key < s[j].Key
	})
}
