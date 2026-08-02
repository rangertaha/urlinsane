// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package eval

import (
	"sort"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// Candidate is one generated name and the algorithm that produced it. The same
// name reached by two algorithms is two candidates, because attribution is the
// point: knowing which algorithm found a real squat is what tells you whether
// it earns its cost.
type Candidate struct {
	Key       string
	Core      string
	Algorithm string
}

// Generate applies every spec once to the seed, exactly as the operator adapter
// would: whole-key algorithms see the whole key, the rest see the registrable
// core and have their output rejoined.
//
// "Once" is the limitation that matters. In a real scan a variant is itself a
// nameable node inside the seed closure, so the engine composes algorithms
// across rounds — a combosquat of a TLD swap is reachable, and this is not.
// Recall measured here is therefore a lower bound on what a scan finds, which
// is the right direction for a metric to be wrong in.
func Generate(seedType, seed string, specs []variant.Spec) []Candidate {
	parts := variant.DefaultSplit(seedType, seed)
	out := make([]Candidate, 0, 4096)

	for _, s := range specs {
		if !appliesTo(s, seedType) {
			continue
		}
		source := parts.Core
		if s.Whole {
			source = seed
		}
		if source == "" || s.Gen == nil {
			continue
		}
		seen := map[string]bool{seed: true}
		for _, raw := range s.Gen(source) {
			if raw == "" || raw == source {
				continue
			}
			key := raw
			if !s.Whole {
				key = parts.Join(raw)
			}
			if key == "" || seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, Candidate{
				Key:       key,
				Core:      coreOf(seedType, key),
				Algorithm: s.ID,
			})
		}
	}
	return out
}

func appliesTo(s variant.Spec, seedType string) bool {
	if len(s.Types) == 0 {
		return true
	}
	for _, t := range s.Types {
		if t == seedType {
			return true
		}
	}
	return false
}

// AlgorithmStat is one algorithm's contribution for a brand.
type AlgorithmStat struct {
	ID string
	// Candidates is how many distinct names this algorithm produced.
	Candidates int
	// Hits is how many truth cores it produced.
	Hits int
	// Unique is how many truth cores *only* it produced. An algorithm with
	// many candidates and no unique hits is paying for coverage another
	// algorithm already provides.
	Unique int
}

// Report is the recall measurement for one brand.
type Report struct {
	Brand     string
	Type      string
	Truth     int
	Generated int

	// CoreHits counts truth records whose registrable name was generated. This
	// is the primary metric: it asks whether the algorithms can reach the
	// squatted *name*, without also requiring that the TLD swap needed to
	// reach the exact domain was composed in the same pass.
	CoreHits int
	// ExactHits counts truth records matched as a full key. Always <= CoreHits.
	ExactHits int

	ByAlgorithm []AlgorithmStat
	// Missed is the truth names no algorithm reached, sorted. This is the list
	// worth reading: it is where the next algorithm comes from.
	Missed []string
}

// CoreRecall is the fraction of truth names whose registrable name was reached.
func (r Report) CoreRecall() float64 {
	if r.Truth == 0 {
		return 0
	}
	return float64(r.CoreHits) / float64(r.Truth)
}

// ExactRecall is the fraction matched as full keys.
func (r Report) ExactRecall() float64 {
	if r.Truth == 0 {
		return 0
	}
	return float64(r.ExactHits) / float64(r.Truth)
}

// CandidatesPerHit is the cost side of the measurement. Every candidate becomes
// a node the engine may spend a DNS lookup on, so an algorithm set that doubles
// recall by decupling candidates is not obviously a win.
func (r Report) CandidatesPerHit() float64 {
	if r.CoreHits == 0 {
		return 0
	}
	return float64(r.Generated) / float64(r.CoreHits)
}

// Score measures one brand's truth records against generated candidates.
func Score(brand, seedType string, truth []Record, cands []Candidate) Report {
	rep := Report{Brand: brand, Type: seedType, Truth: len(truth)}

	// core -> set of algorithms that produced it, and distinct keys per algorithm.
	byCore := map[string]map[string]bool{}
	exact := map[string]bool{}
	perAlgo := map[string]map[string]bool{}

	for _, c := range cands {
		if byCore[c.Core] == nil {
			byCore[c.Core] = map[string]bool{}
		}
		byCore[c.Core][c.Algorithm] = true
		exact[c.Key] = true
		if perAlgo[c.Algorithm] == nil {
			perAlgo[c.Algorithm] = map[string]bool{}
		}
		perAlgo[c.Algorithm][c.Key] = true
	}
	distinct := map[string]bool{}
	for _, c := range cands {
		distinct[c.Key] = true
	}
	rep.Generated = len(distinct)

	hits := map[string]int{}   // algorithm -> truth cores reached
	unique := map[string]int{} // algorithm -> truth cores only it reached

	for _, t := range truth {
		core := t.Core()
		algos := byCore[core]
		if len(algos) == 0 {
			rep.Missed = append(rep.Missed, t.Name)
			continue
		}
		rep.CoreHits++
		if exact[t.Name] {
			rep.ExactHits++
		}
		for a := range algos {
			hits[a]++
		}
		if len(algos) == 1 {
			for a := range algos {
				unique[a]++
			}
		}
	}
	sort.Strings(rep.Missed)

	for id, keys := range perAlgo {
		rep.ByAlgorithm = append(rep.ByAlgorithm, AlgorithmStat{
			ID:         id,
			Candidates: len(keys),
			Hits:       hits[id],
			Unique:     unique[id],
		})
	}
	// Most hits first, then most unique, then id — so the useful algorithms
	// sort to the top and ties stay stable across runs.
	sort.Slice(rep.ByAlgorithm, func(i, j int) bool {
		a, b := rep.ByAlgorithm[i], rep.ByAlgorithm[j]
		if a.Hits != b.Hits {
			return a.Hits > b.Hits
		}
		if a.Unique != b.Unique {
			return a.Unique > b.Unique
		}
		return a.ID < b.ID
	})
	return rep
}

// ScoreAll scores every brand in a truth set against the given specs.
func ScoreAll(records []Record, specs []variant.Spec) []Report {
	grouped := ByBrand(records)
	var out []Report
	for _, brand := range Brands(records) {
		recs := grouped[brand]
		seedType := recs[0].Type
		out = append(out, Score(brand, seedType, recs, Generate(seedType, brand, specs)))
	}
	return out
}
