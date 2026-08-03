// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package model

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
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

// A field that should be a list but is a scalar must be an error, not a panic
// and not a silent empty.
//
// Every list reader opened with make([]T, 0, v.Length()) and then walked
// v.ListIterator(). Length() returns -1 for every non-recursive kind, so the
// make panicked with "cap out of range"; where it did not, ListIterator()
// returned nil, the loop never ran, and the field decoded as an empty list with
// no error at all. A model block whose emission table had been replaced by a
// string loaded as a model with no emissions.
func TestDecodeRejectsAScalarWhereAListBelongs(t *testing.T) {
	// Field 1 is states, field 4 symbols, 5 logInit, 6 logTrans, 7 logEmit —
	// one per reader, so every one of the five is covered.
	for _, field := range []int{1, 2, 3, 4, 5, 6, 7} {
		block := blockWithScalarAt(t, field)

		var h *HMM
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("field %d: Decode panicked instead of erroring: %v", field, r)
				}
			}()
			h, err = Decode(block)
		}()

		if err == nil {
			t.Errorf("field %d: Decode returned a model (%v) for a scalar where a list belongs", field, h != nil)
		}
	}
}

// blockWithScalarAt encodes a well-formed 10-element model block with one field
// replaced by a string.
func blockWithScalarAt(t *testing.T, field int) []byte {
	t.Helper()
	nb := basicnode.Prototype.List.NewBuilder()
	la, err := nb.BeginList(10)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		av := la.AssembleValue()
		switch {
		case i == field:
			err = av.AssignString("not-a-list")
		case i == 0:
			err = av.AssignInt(blockVersion)
		default:
			// An empty list is shape-valid for every other position; the
			// decoder is expected to fail on `field`, not before it.
			var sub datamodel.ListAssembler
			sub, err = av.BeginList(0)
			if err == nil {
				err = sub.Finish()
			}
		}
		if err != nil {
			t.Fatalf("assembling field %d: %v", i, err)
		}
	}
	if err := la.Finish(); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := dagcbor.Encode(nb.Build(), &buf); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// New is the validation boundary, so a malformed spec must come back as an
// error. The emission guard tested len(s.Emit) != 0 while the row read tested
// s.Emit != nil: an empty-but-present table passed the first and was indexed by
// the second. The transition guard beside it already used the nil form.
func TestNewRejectsAnEmptyButPresentEmissionTable(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("New panicked on a malformed spec instead of erroring: %v", r)
		}
	}()
	if _, err := New(Spec{
		States: []string{"A"},
		Focus:  []string{"A"},
		Emit:   [][]float64{},
	}); err == nil {
		t.Fatal("New accepted an emission table with zero rows")
	}
}
