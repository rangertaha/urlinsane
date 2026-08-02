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
	"math"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
)

// syntheticCorpus is two clearly separable regimes of expansion path: variants
// that resolve and keep resolving, and variants that never do. A two-state
// model should be able to tell them apart, which is all Baum-Welch is being
// asked for here.
func syntheticCorpus() Corpus {
	live := Path{Steps: []Trace{
		{Props: []string{"resolves=true"}, Outcome: "ok"},
		{Rel: "VARIANT_OF", Props: []string{"resolves=true", "mx=true"}, Outcome: "ok"},
		{Rel: "RESOLVES_TO", Props: []string{"resolves=true", "mx=true"}, Outcome: "ok"},
	}}
	dead := Path{Steps: []Trace{
		{Props: []string{"resolves=false"}, Outcome: "absent"},
		{Rel: "VARIANT_OF", Props: []string{"resolves=false"}, Outcome: "absent"},
		{Rel: "VARIANT_OF", Props: []string{"resolves=false"}, Outcome: "absent"},
	}}
	mixed := Path{Steps: []Trace{
		{Props: []string{"resolves=true"}, Outcome: "ok"},
		{Rel: "VARIANT_OF", Props: []string{"resolves=false"}, Outcome: "absent"},
		{Rel: "VARIANT_OF", Props: []string{"resolves=false"}, Outcome: "absent"},
	}}
	c := Corpus{}
	for i := 0; i < 6; i++ {
		c.Paths = append(c.Paths, live, dead)
	}
	c.Paths = append(c.Paths, mixed)
	return c
}

func trainConfig() Config {
	return Config{
		States:     []string{"dull", "interesting"},
		Focus:      []string{"interesting"},
		Iterations: 30,
		Seed:       20260801,
		Smoothing:  Smoothing{Init: 0.5, Trans: 0.5, Emit: 0.5},
		Date:       time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
	}
}

// TestBaumWelchObjectiveIsMonotone is the guarantee EM actually makes. Under a
// Dirichlet prior the fit is MAP-EM, so the non-decreasing quantity is the
// penalized objective; asserting it on the bare log-likelihood would be
// asserting something that is merely usually true.
func TestBaumWelchObjectiveIsMonotone(t *testing.T) {
	res, err := Train(syntheticCorpus(), trainConfig())
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	if len(res.Objective) < 3 {
		t.Fatalf("only %d objective samples; nothing to test", len(res.Objective))
	}
	for i := 1; i < len(res.Objective); i++ {
		// A tolerance for floating point summation only, not for a real dip.
		if res.Objective[i] < res.Objective[i-1]-1e-9 {
			t.Fatalf("objective fell at iteration %d: %v -> %v",
				i, res.Objective[i-1], res.Objective[i])
		}
	}
}

func TestBaumWelchImprovesLikelihood(t *testing.T) {
	c := syntheticCorpus()
	res, err := Train(c, trainConfig())
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	first := res.LogLikelihood[0]
	last := res.LogLikelihood[len(res.LogLikelihood)-1]
	if last <= first {
		t.Fatalf("log-likelihood did not improve: %v -> %v", first, last)
	}
	// The reported final likelihood must be the model's own.
	if got := res.Model.LogLikelihood(c); math.Abs(got-last) > 1e-9 {
		t.Fatalf("model likelihood %v disagrees with reported %v", got, last)
	}
	if res.Model.Provenance().LogLikelihood != last {
		t.Fatal("provenance records a different likelihood from the run")
	}
}

// TestBaumWelchSeparatesTheRegimes is the sanity check that the fit is doing
// something, not merely improving a number: a resolving node and an absent one
// must not come out with the same belief.
func TestBaumWelchSeparatesTheRegimes(t *testing.T) {
	res, err := Train(syntheticCorpus(), trainConfig())
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	h := res.Model
	up := h.Mass(h.Forward(h.Prior(), "VARIANT_OF", []string{"resolves=true", "mx=true"}))
	down := h.Mass(h.Forward(h.Prior(), "VARIANT_OF", []string{"resolves=false"}))
	if math.Abs(up-down) < 0.1 {
		t.Fatalf("live and absent nodes got indistinguishable beliefs %v and %v", up, down)
	}
}

// TestTrainIsReproducible: the recorded seed is what makes a model rebuildable,
// so the same corpus and Config must produce the same CID.
func TestTrainIsReproducible(t *testing.T) {
	a, err := Train(syntheticCorpus(), trainConfig())
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	b, err := Train(syntheticCorpus(), trainConfig())
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	ca, err := a.Model.CID()
	if err != nil {
		t.Fatalf("CID: %v", err)
	}
	cb, err := b.Model.CID()
	if err != nil {
		t.Fatalf("CID: %v", err)
	}
	if ca != cb {
		t.Fatalf("two runs of the same training gave %s and %s", ca, cb)
	}
}

func TestADifferentSeedGivesADifferentModel(t *testing.T) {
	cfg := trainConfig()
	cfg.Iterations = 2 // stop well short of convergence, where seeds still differ
	a, err := Train(syntheticCorpus(), cfg)
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	cfg.Seed = 99
	b, err := Train(syntheticCorpus(), cfg)
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	ca, _ := a.Model.CID()
	cb, _ := b.Model.CID()
	if ca == cb {
		t.Fatal("two seeds produced the same model; the seed is not being used")
	}
}

func TestTrainRecordsProvenance(t *testing.T) {
	c := syntheticCorpus()
	c.CIDs = []cid.Cid{testCID(t, "corpus")}
	cfg := trainConfig()
	res, err := Train(c, cfg)
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	p := res.Model.Provenance()
	if p.Algorithm != AlgorithmBaumWelch {
		t.Fatalf("algorithm = %q", p.Algorithm)
	}
	if p.Seed != cfg.Seed {
		t.Fatalf("seed = %d, want %d", p.Seed, cfg.Seed)
	}
	if !p.Date.Equal(cfg.Date) {
		t.Fatalf("date = %v, want %v", p.Date, cfg.Date)
	}
	if len(p.Corpus) != 1 || p.Corpus[0] != c.CIDs[0] {
		t.Fatalf("corpus = %v, want %v", p.Corpus, c.CIDs)
	}
	if p.Iterations != res.Iterations || p.Iterations == 0 {
		t.Fatalf("iterations = %d, result reports %d", p.Iterations, res.Iterations)
	}
}

func TestTrainConvergesEarly(t *testing.T) {
	cfg := trainConfig()
	cfg.Iterations = 500
	cfg.Tolerance = 1e-6
	res, err := Train(syntheticCorpus(), cfg)
	if err != nil {
		t.Fatalf("Train: %v", err)
	}
	if res.Iterations >= cfg.Iterations {
		t.Fatalf("ran %d iterations; tolerance never bound", res.Iterations)
	}
}

func TestTrainRejectsBadInput(t *testing.T) {
	good := trainConfig()
	cases := map[string]struct {
		c   Corpus
		cfg Config
	}{
		"no states":      {syntheticCorpus(), Config{Iterations: 1}},
		"no iterations":  {syntheticCorpus(), Config{States: []string{"A"}}},
		"empty corpus":   {Corpus{}, good},
		"zero smoothing": {syntheticCorpus(), Config{States: []string{"A"}, Iterations: 1, Smoothing: Smoothing{Init: 0, Trans: 1, Emit: 1}}},
	}
	for name, c := range cases {
		if _, err := Train(c.c, c.cfg); err == nil {
			t.Fatalf("%s: Train accepted bad input", name)
		}
	}
}

// TestAlphabetsAreDerivedInSortedOrder guards the determinism of training
// itself: index position is part of the model's identity, so deriving it from
// map iteration would give two runs different CIDs.
func TestAlphabetsAreDerivedInSortedOrder(t *testing.T) {
	paths := []Path{{Steps: []Trace{
		{Props: []string{"z", "a"}},
		{Rel: "ZED", Props: []string{"m"}},
		{Rel: "ALPHA", Props: []string{"b"}, Outcome: "ok"},
	}}}
	for i := 0; i < 20; i++ {
		rels, syms := alphabets(paths)
		if len(rels) != 2 || rels[0] != "ALPHA" || rels[1] != "ZED" {
			t.Fatalf("rels = %v, want sorted", rels)
		}
		want := []string{"a", "b", "m", OutcomePrefix + "ok", "z"}
		if len(syms) != len(want) {
			t.Fatalf("symbols = %v, want %v", syms, want)
		}
		for j := range want {
			if syms[j] != want[j] {
				t.Fatalf("symbols = %v, want %v", syms, want)
			}
		}
	}
}

func TestRNGIsSelfContainedAndSeeded(t *testing.T) {
	a, b := newRNG(7), newRNG(7)
	for i := 0; i < 100; i++ {
		x, y := a.float(), b.float()
		if x != y {
			t.Fatalf("same seed diverged at draw %d: %v != %v", i, x, y)
		}
		if x < 0 || x >= 1 {
			t.Fatalf("draw %d out of range: %v", i, x)
		}
	}
	if newRNG(7).float() == newRNG(8).float() {
		t.Fatal("different seeds gave the same first draw")
	}
}

func TestTraceSymbolsIncludeOutcome(t *testing.T) {
	tr := Trace{Props: []string{"a=1"}, Outcome: "absent"}
	got := tr.Symbols()
	if len(got) != 2 || got[1] != OutcomePrefix+"absent" {
		t.Fatalf("Symbols() = %v", got)
	}
	// The props slice must not be aliased into the result.
	if &got[0] == &tr.Props[0] {
		t.Fatal("Symbols aliased the caller's slice")
	}
	if len(Trace{Props: []string{"a=1"}}.Symbols()) != 1 {
		t.Fatal("an empty outcome should add no symbol")
	}
}
