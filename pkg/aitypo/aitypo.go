// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package aitypo packages urlinsane's typo generators as learnable tasks.
//
// Every function in pkg/typo is a total, deterministic map from a name to a set
// of names: CharacterOmission("google") is exactly {gogle, googe, googl, goole,
// oogle} and nothing else, forever. That property is unusual and it is what
// makes these functions trainable. A task with an exact oracle needs no
// annotation, cannot disagree with itself, generates unlimited labelled data,
// and can be scored exactly rather than by a similarity metric someone chose.
//
// So this package is not a model and trains nothing. It is the harness around
// the functions:
//
//	Task     one generator, with its oracle and what it needs to run
//	Corpus   the names to run it over — the shipped datasets, or your own
//	Example  one (input, expected-set) pair, with a split assigned
//	Score    a model's predicted set against the oracle's, per task
//
// # It wraps pkg/typo, it does not copy it
//
// This package began as a copied directory. That is worth naming, because a
// copy of a generator set is the exact defect this repository spent a day
// retiring: acs, aci and rar had been fixed rune-safe in their own plugin
// directories while eight sibling generators kept slicing bytes, and the split
// was the bug. Four minutes after the copy was taken it was already two commits
// behind. There is one implementation, in pkg/typo, and everything here calls
// it.
//
// # What a model is actually being asked to learn
//
// Task.Needs is the interesting field, not the id. A generator with Needs ==
// NeedsNothing encodes a rule — delete one character, transpose two, insert a
// hyphen at every gap — and a model that gets it right has learned the rule.
// A generator with NeedsLanguage or NeedsKeyboard is a table lookup: homoglyphs
// for a script, the neighbours of a key on a layout. A model can only reproduce
// those by memorising the table, so success there means something different,
// and generalisation to a script the table does not carry is the real test.
//
// Splitting the two is why Needs exists. Reporting one accuracy number across
// both would average a rule a model learned with a table it memorised.
package aitypo

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
)

// Needs says what data a task's oracle reads beyond its input.
//
// It is a property of the generator, not of a run: CharacterOmission reads
// nothing whatever you pass it, and HomoglyphSwapping is a table lookup even
// when the table is empty.
type Needs uint8

const (
	// NeedsNothing is a pure string rule. The oracle is a function of the
	// input alone, so a model that reproduces it has learned the rule.
	NeedsNothing Needs = iota
	// NeedsLanguage reads a language's vocabulary: homoglyphs, homophones,
	// misspellings, vowels, graphemes, numerals.
	NeedsLanguage
	// NeedsKeyboard reads a layout's key adjacency.
	NeedsKeyboard
)

func (n Needs) String() string {
	switch n {
	case NeedsLanguage:
		return "language"
	case NeedsKeyboard:
		return "keyboard"
	}
	return "nothing"
}

// Oracle is a generator: a total, deterministic map from one name to the set of
// names it produces.
//
// Total and deterministic are load-bearing. Every scoring number in this
// package assumes calling an oracle twice with the same input gives the same
// set, so a generator that read a clock, ranged a map, or depended on
// registration order would make an experiment unrepeatable without failing
// anywhere visible.
type Oracle func(name string) []string

// Task is one generator, packaged for training.
type Task struct {
	// ID matches the algorithm id the scanner uses — co, cs, hr — so a result
	// here names the same thing a scan does.
	ID    string
	Title string
	Needs Needs
	// Oracle is the reference implementation. It is the label source, the
	// evaluation target, and the definition of correct.
	Oracle Oracle
}

// Expect returns the oracle's answer as a sorted, deduplicated set.
//
// Sorted because a set has no order and a corpus written in map order would
// diff against itself; deduplicated because a duplicate is not extra
// information and would weight one variant twice in training.
func (t Task) Expect(name string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, v := range t.Oracle(name) {
		if v == "" || v == name || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// Example is one training record: an input and the complete set of outputs the
// oracle produces for it.
//
// The whole set, not one pair. These generators are set-valued — omission of
// "google" is five names — and a corpus of (input, one-output) rows would teach
// a model to produce one variant and give it no way to learn when to stop.
type Example struct {
	Task   string   `json:"task"`
	Input  string   `json:"input"`
	Expect []string `json:"expect"`
	Split  string   `json:"split,omitempty"`
	// Source names where the input came from — a dataset path, a scan, a
	// hand-written case. A corpus that cannot say what it is made of cannot be
	// reproduced.
	Source string `json:"source,omitempty"`
}

// Group is the unit a split must not cut across: the input name.
//
// One name yields an example per task, and those examples share almost all
// their structure. Putting "google" in train for omission and in test for
// transposition measures memorisation of the name, not of the rule.
func (e Example) Group() string { return e.Input }

// hash64 is a stable hash. Deterministic across processes and machines, unlike
// Go's map hash, which is seeded per process — a split built on that would put
// a name in train today and test tomorrow, and two experiments would differ by
// their corpora as well as their models.
func hash64(s string) uint64 {
	sum := sha256.Sum256([]byte(s))
	return binary.BigEndian.Uint64(sum[:8])
}

// Split names the three partitions.
const (
	SplitTrain = "train"
	SplitVal   = "val"
	SplitTest  = "test"
)

// Ratio is the fraction of *groups* — not examples — in each partition.
type Ratio struct{ Train, Val, Test float64 }

// DefaultRatio is 80/10/10.
var DefaultRatio = Ratio{Train: 0.8, Val: 0.1, Test: 0.1}

// Assign labels every example with a split, grouped by input name.
//
// Deterministic and stateless: the partition depends only on the name and the
// salt, so two machines building the same corpus agree, and adding a name never
// moves an existing one. A shuffle-and-slice would move every name each time
// the corpus grew, quietly invalidating every earlier measurement.
//
// salt lets one corpus be re-split independently of another — pass the run or
// experiment name. An empty salt is fine and gives the canonical split.
func Assign(examples []Example, r Ratio, salt string) {
	train, val := r.Train, r.Train+r.Val
	if total := r.Train + r.Val + r.Test; total > 0 {
		train, val = r.Train/total, (r.Train+r.Val)/total
	}
	for i := range examples {
		// 53 bits is the mantissa of a float64; taking the whole 64 and
		// dividing would lose the low bits to rounding and bias the buckets.
		f := float64(hash64(salt+"\x00"+examples[i].Group())>>11) / float64(1<<53)
		switch {
		case f < train:
			examples[i].Split = SplitTrain
		case f < val:
			examples[i].Split = SplitVal
		default:
			examples[i].Split = SplitTest
		}
	}
}

// Registry is a set of tasks, addressable by id.
type Registry map[string]Task

// IDs returns the registered ids, sorted.
func (r Registry) IDs() []string {
	out := make([]string, 0, len(r))
	for id := range r {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// Select returns the named tasks, or every task when no ids are given.
//
// An unknown id is an error rather than a silent omission: a training run that
// quietly dropped a task would report a mean over fewer tasks than the operator
// asked for, and nothing would say so.
func (r Registry) Select(ids ...string) ([]Task, error) {
	if len(ids) == 0 {
		out := make([]Task, 0, len(r))
		for _, id := range r.IDs() {
			out = append(out, r[id])
		}
		return out, nil
	}
	var out []Task
	var unknown []string
	// Deduplicated: a repeated id would emit the same example twice and give
	// that task two votes in the macro average, so `--task co --task co,cs`
	// would quietly weight co double.
	picked := map[string]bool{}
	var named int
	for _, raw := range ids {
		for _, id := range strings.Split(raw, ",") {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			named++
			t, ok := r[id]
			if !ok {
				unknown = append(unknown, id)
				continue
			}
			if picked[id] {
				continue
			}
			picked[id] = true
			out = append(out, t)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return nil, fmt.Errorf("aitypo: unknown task(s): %s; have: %s",
			strings.Join(unknown, ", "), strings.Join(r.IDs(), ", "))
	}
	// A caller that passed something naming nothing gets an error, not silence.
	// Returning no tasks would build an empty corpus and report success, which
	// is the quiet drop this function's contract exists to prevent — and it is
	// what `--task ""` produced.
	if named == 0 {
		return nil, fmt.Errorf("aitypo: %q names no task; pass none to select all of: %s",
			ids, strings.Join(r.IDs(), ", "))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
