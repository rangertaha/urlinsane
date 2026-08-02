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
)

const eps = 1e-12

func closeTo(t *testing.T, got, want float64, what string) {
	t.Helper()
	if math.Abs(got-want) > eps {
		t.Fatalf("%s = %v, want %v", what, got, want)
	}
}

// twoState is the hand-worked example used across these tests.
//
//	init         A=0.6  B=0.4
//	P(.|.,"R")   A->A=0.7 A->B=0.3   B->A=0.2 B->B=0.8
//	emit A       x=0.9 y=0.05 oov=0.05
//	emit B       x=0.1 y=0.80 oov=0.10
func twoState(t *testing.T) *HMM {
	t.Helper()
	h, err := New(Spec{
		States:  []string{"A", "B"},
		Focus:   []string{"A"},
		Rels:    []string{"R", "S"},
		Symbols: []string{"x", "y"},
		Init:    []float64{0.6, 0.4},
		Trans: map[string][][]float64{
			"R": {{0.7, 0.3}, {0.2, 0.8}},
			"S": {{0.5, 0.5}, {0.5, 0.5}},
		},
		Emit: [][]float64{
			{0.9, 0.05, 0.05},
			{0.1, 0.80, 0.10},
		},
		Smoothing:  DefaultSmoothing(),
		Provenance: Provenance{Algorithm: AlgorithmManual},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return h
}

// --- log-space arithmetic ---------------------------------------------------

func TestLogAddMatchesLinearArithmetic(t *testing.T) {
	for _, c := range []struct{ a, b float64 }{
		{0.5, 0.25}, {1e-9, 1e-9}, {0.999, 1e-300}, {1, 1},
	} {
		got := LogAdd(math.Log(c.a), math.Log(c.b))
		want := math.Log(c.a + c.b)
		if math.Abs(got-want) > 1e-9 {
			t.Fatalf("LogAdd(log %v, log %v) = %v, want %v", c.a, c.b, got, want)
		}
	}
}

func TestLogAddIsClosedOverLogZero(t *testing.T) {
	if got := LogAdd(LogZero, LogZero); !math.IsInf(got, -1) {
		t.Fatalf("LogAdd(0,0) = %v, want -Inf", got)
	}
	if got := LogAdd(LogZero, math.Log(0.25)); math.Abs(got-math.Log(0.25)) > eps {
		t.Fatalf("LogAdd(0, 0.25) = %v, want log 0.25", got)
	}
	if math.IsNaN(LogAdd(LogZero, LogZero)) {
		t.Fatal("LogAdd produced NaN")
	}
}

func TestLogSumExpSurvivesUnderflow(t *testing.T) {
	// In linear space every term is 0 and the sum is 0; in log space the answer
	// is exact. This is the whole reason the tables are log-space.
	xs := []float64{-800, -801, -802}
	got := LogSumExp(xs)
	want := -800 + math.Log(1+math.Exp(-1)+math.Exp(-2))
	closeTo(t, got, want, "LogSumExp")
	if math.Exp(xs[0]) != 0 {
		t.Skip("platform has extended range; underflow premise does not hold")
	}
}

func TestLogSumExpEmptyAndZero(t *testing.T) {
	if !math.IsInf(LogSumExp(nil), -1) {
		t.Fatal("LogSumExp(nil) should be LogZero")
	}
	if !math.IsInf(LogSumExp([]float64{LogZero, LogZero}), -1) {
		t.Fatal("LogSumExp(all zero) should be LogZero")
	}
}

func TestNormalizeLogSumsToOne(t *testing.T) {
	d := []float64{math.Log(2), math.Log(3), math.Log(5)}
	z := normalizeLog(d)
	closeTo(t, z, math.Log(10), "log evidence")
	sum := 0.0
	for _, v := range d {
		sum += math.Exp(v)
	}
	closeTo(t, sum, 1, "normalized mass")
}

func TestNormalizeLogFallsBackToUniform(t *testing.T) {
	d := []float64{LogZero, LogZero}
	normalizeLog(d)
	for i, v := range d {
		closeTo(t, math.Exp(v), 0.5, "uniform fallback")
		if math.IsNaN(v) {
			t.Fatalf("entry %d is NaN", i)
		}
	}
}

// --- forward step -----------------------------------------------------------

// TestForwardHandWorkedTwoState checks the forward step against arithmetic done
// by hand:
//
//	predict A = 0.6*0.7 + 0.4*0.2 = 0.50
//	predict B = 0.6*0.3 + 0.4*0.8 = 0.50
//	observe x: A = 0.50*0.9 = 0.45, B = 0.50*0.1 = 0.05
//	normalize: A = 0.9, B = 0.1
func TestForwardHandWorkedTwoState(t *testing.T) {
	h := twoState(t)
	got := h.Forward(h.Prior(), "R", []string{"x"})
	closeTo(t, math.Exp(got[0]), 0.9, "P(A|x)")
	closeTo(t, math.Exp(got[1]), 0.1, "P(B|x)")
	closeTo(t, h.Mass(got), 0.9, "belief")
}

func TestForwardSecondStepHandWorked(t *testing.T) {
	h := twoState(t)
	// From (0.9, 0.1): predict A = 0.9*0.7 + 0.1*0.2 = 0.65, B = 0.35.
	// observe y: A = 0.65*0.05 = 0.0325, B = 0.35*0.80 = 0.28.
	first := h.Forward(h.Prior(), "R", []string{"x"})
	got := h.Forward(first, "R", []string{"y"})
	total := 0.0325 + 0.28
	closeTo(t, math.Exp(got[0]), 0.0325/total, "P(A|x,y)")
	closeTo(t, math.Exp(got[1]), 0.28/total, "P(B|x,y)")
}

// TestForwardWithoutObservationIsTransitionPrior is the candidate case from
// §10.2: a candidate has no props, so its belief must be the transition prior
// alone.
func TestForwardWithoutObservationIsTransitionPrior(t *testing.T) {
	h := twoState(t)
	got := h.Forward(h.Prior(), "R", nil)
	closeTo(t, math.Exp(got[0]), 0.5, "predicted A")
	closeTo(t, math.Exp(got[1]), 0.5, "predicted B")
}

// TestForwardIsConditionedOnRelation is the point of a relation-conditioned
// transition table: the same parent and the same props give different beliefs
// down different relations.
func TestForwardIsConditionedOnRelation(t *testing.T) {
	h := twoState(t)
	parent := h.Lift(0.9)
	r := h.Mass(h.Forward(parent, "R", nil))
	s := h.Mass(h.Forward(parent, "S", nil))
	closeTo(t, r, 0.9*0.7+0.1*0.2, "belief down R")
	closeTo(t, s, 0.5, "belief down the uninformative S")
	if math.Abs(r-s) < 1e-9 {
		t.Fatalf("relations R and S gave the same belief %v", r)
	}
}

func TestForwardIsDeterministic(t *testing.T) {
	h := twoState(t)
	first := h.Mass(h.Forward(h.Prior(), "R", []string{"x", "y"}))
	for i := 0; i < 100; i++ {
		if got := h.Mass(h.Forward(h.Prior(), "R", []string{"x", "y"})); got != first {
			t.Fatalf("run %d gave %v, want %v", i, got, first)
		}
	}
}

func TestForwardDoesNotMutateInput(t *testing.T) {
	h := twoState(t)
	prev := h.Prior()
	before := append([]float64(nil), prev...)
	h.Forward(prev, "R", []string{"x"})
	for i := range prev {
		if prev[i] != before[i] {
			t.Fatalf("Forward mutated its input at %d", i)
		}
	}
}

// --- OOV --------------------------------------------------------------------

func TestUnseenSymbolFallsBackToOOV(t *testing.T) {
	h := twoState(t)
	if got := h.SymbolIndex("never-trained"); got != h.SymbolIndex(OOVSymbol) {
		t.Fatalf("unseen symbol mapped to %d, want the OOV index", got)
	}
	// The point of an explicit OOV symbol: an unseen prop must not annihilate
	// the filter.
	d := h.Forward(h.Prior(), "R", []string{"never-trained"})
	for i, v := range d {
		if math.IsInf(v, -1) || math.IsNaN(v) {
			t.Fatalf("state %d collapsed to %v on an unseen symbol", i, v)
		}
	}
}

func TestUnseenRelationFallsBackToOOV(t *testing.T) {
	h := twoState(t)
	if got := h.RelIndex("NOT_IN_MODEL"); got != h.RelIndex(OOVRelation) {
		t.Fatalf("unseen relation mapped to %d, want the OOV index", got)
	}
	// The OOV relation table was never specified, so it is uniform: the belief
	// is the emission alone.
	d := h.Forward(h.Prior(), "NOT_IN_MODEL", nil)
	closeTo(t, math.Exp(d[0]), 0.5, "OOV relation is uninformative")
}

func TestEmissionIsAProductOverSymbols(t *testing.T) {
	h := twoState(t)
	both := h.LogEmission(0, []string{"x", "y"})
	closeTo(t, both, math.Log(0.9)+math.Log(0.05), "product emission")
	closeTo(t, h.LogEmission(0, nil), 0, "empty observation")
}

// --- lift / mass ------------------------------------------------------------

// TestLiftIsExactForABinaryModel: for the shape §10.2 actually needs, the
// scalar round trip through graph.BeliefModel loses nothing.
func TestLiftIsExactForABinaryModel(t *testing.T) {
	h := twoState(t)
	for _, b := range []float64{0, 0.1, 0.25, 0.5, 0.75, 0.9, 1} {
		if got := h.Mass(h.Lift(b)); math.Abs(got-b) > 1e-12 {
			t.Fatalf("Mass(Lift(%v)) = %v", b, got)
		}
	}
}

func TestMassIsClamped(t *testing.T) {
	h := twoState(t)
	if got := h.Mass([]float64{0, 0}); got != 1 {
		t.Fatalf("Mass of an over-full distribution = %v, want 1", got)
	}
	if got := h.Mass(nil); got != 0 {
		t.Fatalf("Mass of a mismatched distribution = %v, want 0", got)
	}
}

// --- construction -----------------------------------------------------------

func TestNewRejectsMalformedSpecs(t *testing.T) {
	cases := map[string]Spec{
		"no states":     {},
		"no focus":      {States: []string{"A"}},
		"unknown focus": {States: []string{"A"}, Focus: []string{"Z"}},
		"duplicate state": {
			States: []string{"A", "A"}, Focus: []string{"A"},
		},
		"short init": {
			States: []string{"A", "B"}, Focus: []string{"A"}, Init: []float64{1},
		},
		"wrong transition shape": {
			States: []string{"A", "B"}, Focus: []string{"A"}, Rels: []string{"R"},
			Trans: map[string][][]float64{"R": {{1, 0}}},
		},
	}
	for name, s := range cases {
		if _, err := New(s); err == nil {
			t.Fatalf("%s: New accepted a malformed spec", name)
		}
	}
}

func TestNewAppendsReservedAlphabetEntries(t *testing.T) {
	h := twoState(t)
	if got := h.Symbols(); got[len(got)-1] != OOVSymbol {
		t.Fatalf("symbols = %v, want %q appended", got, OOVSymbol)
	}
	if got := h.Rels(); got[len(got)-1] != OOVRelation {
		t.Fatalf("rels = %v, want %q appended", got, OOVRelation)
	}
}

func TestUnspecifiedTablesAreUniform(t *testing.T) {
	h, err := New(Spec{States: []string{"A", "B"}, Focus: []string{"A"}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	closeTo(t, math.Exp(h.Prior()[0]), 0.5, "unspecified init")
	closeTo(t, h.LogTransition("anything", 0, 1), math.Log(0.5), "unspecified transition")
}

func TestSmoothingValidation(t *testing.T) {
	if err := (Smoothing{Init: -1, Trans: 1, Emit: 1}).validate(false); err == nil {
		t.Fatal("negative smoothing accepted")
	}
	if err := (Smoothing{Init: 0, Trans: 1, Emit: 1}).validate(true); err == nil {
		t.Fatal("zero smoothing accepted for training")
	}
	if err := DefaultSmoothing().validate(true); err != nil {
		t.Fatalf("default smoothing rejected: %v", err)
	}
}

func TestDirichletLogNeverProducesZero(t *testing.T) {
	// The property the forward filter depends on: an unobserved cell is
	// unlikely, not impossible.
	row := dirichletLog([]float64{10, 0, 0}, 1)
	for i, lp := range row {
		if math.IsInf(lp, -1) {
			t.Fatalf("cell %d is impossible under a positive prior", i)
		}
	}
	closeTo(t, math.Exp(row[0]), 11.0/13.0, "smoothed cell")
	closeTo(t, math.Exp(row[1]), 1.0/13.0, "unobserved cell")
}
