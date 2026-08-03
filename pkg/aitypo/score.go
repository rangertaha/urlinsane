// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package aitypo

import (
	"fmt"
	"sort"
	"strings"
)

// Result is one prediction scored against its oracle.
//
// Both counts and the sets are kept. The counts are what a training loop reads;
// the sets are what a human reads when the number is disappointing, and losing
// them means re-running the model to find out what it actually said.
type Result struct {
	Task  string
	Input string

	Expect  []string
	Predict []string

	Hit    []string // predicted and expected
	Missed []string // expected, not predicted
	Spuri  []string // predicted, not expected
}

// Exact reports whether the prediction is the expected set exactly.
//
// The headline metric for these tasks, and a strict one on purpose. A
// generator's contract is the whole set — omission of "google" is five names,
// not four of them and something plausible — so a model that gets four right
// and invents a fifth has not learned the rule. Precision and recall are there
// to say how it failed.
func (r Result) Exact() bool { return len(r.Missed) == 0 && len(r.Spuri) == 0 }

// Precision is the fraction of predictions that were expected. An empty
// prediction is defined as 1: a model that said nothing said nothing wrong.
func (r Result) Precision() float64 {
	if len(r.Predict) == 0 {
		return 1
	}
	return float64(len(r.Hit)) / float64(len(r.Predict))
}

// Recall is the fraction of expectations that were predicted. An empty
// expectation is defined as 1: there was nothing to find.
func (r Result) Recall() float64 {
	if len(r.Expect) == 0 {
		return 1
	}
	return float64(len(r.Hit)) / float64(len(r.Expect))
}

// F1 is the harmonic mean, 0 when either side is 0.
func (r Result) F1() float64 {
	p, rc := r.Precision(), r.Recall()
	if p+rc == 0 {
		return 0
	}
	return 2 * p * rc / (p + rc)
}

// Score compares one prediction against a task's oracle.
//
// The prediction is deduplicated and the input itself is dropped from it before
// comparison, exactly as Expect does to the oracle's output. Otherwise a model
// could inflate recall by emitting the input and every variant twice, and the
// two sides would be measured under different rules.
func Score(t Task, input string, predict []string) Result {
	expect := t.Expect(input)

	pred := make([]string, 0, len(predict))
	seen := make(map[string]bool, len(predict))
	for _, p := range predict {
		if p == "" || p == input || seen[p] {
			continue
		}
		seen[p] = true
		pred = append(pred, p)
	}
	sort.Strings(pred)

	inExpect := make(map[string]bool, len(expect))
	for _, e := range expect {
		inExpect[e] = true
	}
	r := Result{Task: t.ID, Input: input, Expect: expect, Predict: pred}
	for _, p := range pred {
		if inExpect[p] {
			r.Hit = append(r.Hit, p)
		} else {
			r.Spuri = append(r.Spuri, p)
		}
	}
	for _, e := range expect {
		if !seen[e] {
			r.Missed = append(r.Missed, e)
		}
	}
	return r
}

// Summary aggregates results for one task.
type Summary struct {
	Task  string
	Needs Needs
	N     int

	ExactMatch float64 // fraction of inputs answered exactly
	Precision  float64 // macro-averaged over inputs
	Recall     float64
	F1         float64
}

func (s Summary) String() string {
	return fmt.Sprintf("%-5s %-9s n=%-6d exact=%.3f  P=%.3f R=%.3f F1=%.3f",
		s.Task, s.Needs, s.N, s.ExactMatch, s.Precision, s.Recall, s.F1)
}

// Summarize aggregates per task, macro-averaged over inputs.
//
// It takes the Registry rather than a needs map, and errors on a result naming
// a task the registry does not have. Both halves are needed and only the first
// was here: Registry is itself a map, so reg[id] on a miss returned the zero
// Task, whose Needs is NeedsNothing — exactly the silent mislabelling the type
// was supposed to prevent.
//
// It is reachable in normal use. Building a corpus with language and keyboard
// data and then scoring it against a registry built without them — which is
// what Tasks(Data{}) gives, and it deregisters hr, hs, acs, aci, vs and gi —
// reported a memorised homoglyph table and a memorised keyboard layout as
// learned rules, at a perfect 1.000. That is the one distinction this package
// exists to keep, so it is an error rather than a footnote.
//
// Macro, not micro: micro-averaging would weight an input by how many variants
// it happens to produce, so a ten-character name would count ten times as much
// as a three-character one and the score would mostly measure the corpus's
// length distribution. Every input is one vote.
//
// Results are grouped by task and returned in task order, so two runs of the
// same evaluation print identically.
func Summarize(results []Result, reg Registry) ([]Summary, error) {
	type acc struct {
		n              int
		exact, p, r, f float64
	}
	by := map[string]*acc{}
	for _, res := range results {
		a := by[res.Task]
		if a == nil {
			a = &acc{}
			by[res.Task] = a
		}
		a.n++
		if res.Exact() {
			a.exact++
		}
		a.p += res.Precision()
		a.r += res.Recall()
		a.f += res.F1()
	}

	ids := make([]string, 0, len(by))
	for id := range by {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var unknown []string
	for _, id := range ids {
		if _, ok := reg[id]; !ok {
			unknown = append(unknown, id)
		}
	}
	if len(unknown) > 0 {
		return nil, fmt.Errorf(
			"aitypo: results name task(s) the registry does not have: %s; "+
				"their Needs would default to %v and a memorised table would be "+
				"reported as a learned rule — score against the registry the "+
				"corpus was built from (have: %s)",
			strings.Join(unknown, ", "), NeedsNothing, strings.Join(reg.IDs(), ", "))
	}

	out := make([]Summary, 0, len(ids))
	for _, id := range ids {
		a := by[id]
		n := float64(a.n)
		out = append(out, Summary{
			Task:       id,
			Needs:      reg[id].Needs,
			N:          a.n,
			ExactMatch: a.exact / n,
			Precision:  a.p / n,
			Recall:     a.r / n,
			F1:         a.f / n,
		})
	}
	return out, nil
}
