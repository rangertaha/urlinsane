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
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

func testCID(t *testing.T, s string) cid.Cid {
	t.Helper()
	h, err := multihash.Sum([]byte(s), multihash.SHA2_256, -1)
	if err != nil {
		t.Fatalf("multihash: %v", err)
	}
	return cid.NewCidV1(cid.DagCBOR, h)
}

func provenanced(t *testing.T) *HMM {
	t.Helper()
	h := twoState(t)
	h.prov = Provenance{
		Algorithm:     AlgorithmBaumWelch,
		Seed:          0x5eed,
		Date:          time.Date(2026, 8, 1, 12, 0, 0, 123456789, time.UTC),
		Corpus:        []cid.Cid{testCID(t, "trace-a"), testCID(t, "trace-b")},
		Iterations:    17,
		LogLikelihood: -42.5,
	}
	return h
}

// TestCIDIsStableAcrossEncodes is the property the plan hash rests on: a model
// pinned by CID must encode to the same bytes every time, or a pinned plan
// stops reproducing (§10.4).
func TestCIDIsStableAcrossEncodes(t *testing.T) {
	h := provenanced(t)
	block, first, err := h.Addressed()
	if err != nil {
		t.Fatalf("Addressed: %v", err)
	}
	for i := 0; i < 50; i++ {
		b, c, err := h.Addressed()
		if err != nil {
			t.Fatalf("Addressed %d: %v", i, err)
		}
		if !bytes.Equal(b, block) {
			t.Fatalf("encode %d produced different bytes", i)
		}
		if c != first {
			t.Fatalf("encode %d produced CID %s, want %s", i, c, first)
		}
	}
}

// TestTwoIdenticalModelsShareACID: identity is content, not object identity.
func TestTwoIdenticalModelsShareACID(t *testing.T) {
	a, b := provenanced(t), provenanced(t)
	ca, err := a.CID()
	if err != nil {
		t.Fatalf("CID: %v", err)
	}
	cb, err := b.CID()
	if err != nil {
		t.Fatalf("CID: %v", err)
	}
	if ca != cb {
		t.Fatalf("identical models have CIDs %s and %s", ca, cb)
	}
}

func TestRoundTripPreservesCIDAndTables(t *testing.T) {
	h := provenanced(t)
	block, want, err := h.Addressed()
	if err != nil {
		t.Fatalf("Addressed: %v", err)
	}
	got, err := Decode(block)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	reblock, c, err := got.Addressed()
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if !bytes.Equal(block, reblock) {
		t.Fatal("round trip changed the block bytes")
	}
	if c != want {
		t.Fatalf("round trip changed the CID: %s -> %s", want, c)
	}

	// Behaviour, not just bytes: the decoded model must filter identically.
	a := h.Mass(h.Forward(h.Prior(), "R", []string{"x"}))
	b := got.Mass(got.Forward(got.Prior(), "R", []string{"x"}))
	if a != b {
		t.Fatalf("decoded model gives belief %v, want %v", b, a)
	}
}

func TestRoundTripPreservesProvenance(t *testing.T) {
	h := provenanced(t)
	block, _, err := h.Addressed()
	if err != nil {
		t.Fatalf("Addressed: %v", err)
	}
	got, err := Decode(block)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	want := h.Provenance()
	p := got.Provenance()
	if p.Algorithm != want.Algorithm || p.Seed != want.Seed || p.Iterations != want.Iterations {
		t.Fatalf("provenance = %+v, want %+v", p, want)
	}
	if !p.Date.Equal(want.Date) {
		t.Fatalf("date = %v, want %v", p.Date, want.Date)
	}
	if p.LogLikelihood != want.LogLikelihood {
		t.Fatalf("log-likelihood = %v, want %v", p.LogLikelihood, want.LogLikelihood)
	}
	if len(p.Corpus) != len(want.Corpus) {
		t.Fatalf("corpus has %d CIDs, want %d", len(p.Corpus), len(want.Corpus))
	}
	for i := range p.Corpus {
		if p.Corpus[i] != want.Corpus[i] {
			t.Fatalf("corpus[%d] = %s, want %s", i, p.Corpus[i], want.Corpus[i])
		}
	}
	if got.Smoothing() != h.Smoothing() {
		t.Fatalf("smoothing = %+v, want %+v", got.Smoothing(), h.Smoothing())
	}
	if len(got.Focus()) != 1 || got.Focus()[0] != "A" {
		t.Fatalf("focus = %v, want [A]", got.Focus())
	}
}

// TestRoundTripSurvivesImpossibleCells covers the sentinel: dag-cbor's decoder
// rejects infinities outright, so a hard zero in a hand-written table has to
// travel as a finite stand-in and come back as -Inf.
func TestRoundTripSurvivesImpossibleCells(t *testing.T) {
	h, err := New(Spec{
		States:  []string{"A", "B"},
		Focus:   []string{"B"},
		Rels:    []string{"R"},
		Symbols: []string{"x"},
		Init:    []float64{1, 0}, // B is impossible at the seed
		Trans:   map[string][][]float64{"R": {{0, 1}, {1, 0}}},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !math.IsInf(h.Prior()[1], -1) {
		t.Fatalf("prior[B] = %v, want -Inf", h.Prior()[1])
	}
	block, want, err := h.Addressed()
	if err != nil {
		t.Fatalf("Addressed: %v", err)
	}
	got, err := Decode(block)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if !math.IsInf(got.Prior()[1], -1) {
		t.Fatalf("decoded prior[B] = %v, want -Inf", got.Prior()[1])
	}
	_, c, err := got.Addressed()
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if c != want {
		t.Fatalf("round trip changed the CID: %s -> %s", want, c)
	}
}

func TestUniformModelRoundTrips(t *testing.T) {
	h := Uniform()
	block, want, err := h.Addressed()
	if err != nil {
		t.Fatalf("Addressed: %v", err)
	}
	got, err := Decode(block)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	_, c, err := got.Addressed()
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if c != want {
		t.Fatalf("uniform model CID changed: %s -> %s", want, c)
	}
	if b, _ := NewBelief(got, nil).Initial(); b != 1 {
		t.Fatalf("decoded uniform model belief = %v, want exactly 1", b)
	}
}

func TestDecodeRejectsGarbage(t *testing.T) {
	if _, err := Decode([]byte{0x00, 0x01}); err == nil {
		t.Fatal("Decode accepted garbage")
	}
	if _, err := Decode(nil); err == nil {
		t.Fatal("Decode accepted an empty block")
	}
}

func TestDecodeRejectsAMismatchedShape(t *testing.T) {
	h := twoState(t)
	// Corrupt the emission table so its width no longer matches the alphabet.
	h.logEmit[1] = h.logEmit[1][:1]
	block, _, err := h.Addressed()
	if err != nil {
		t.Fatalf("Addressed: %v", err)
	}
	if _, err := Decode(block); err == nil {
		t.Fatal("Decode accepted a table whose shape disagrees with the alphabet")
	}
}
