// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package dataset

import (
	"fmt"
	"path/filepath"
	"testing"
)

// A relation bigger than the driver's host-parameter limit must still read
// back whole.
//
// transitions() binds one parameter per vocabulary id. SQLite caps those at
// 32766, and older builds at 999; past the cap the query errors, transitions
// reports !ok, and edges() and groups() return nil. Nothing anywhere reports a
// failure -- the relation is simply empty, so the algorithm reading it goes
// quiet. The shipped es/misspelling relation crossed the line at 36,650 tokens
// and Spanish common misspellings generated nothing at all.
//
// 40,000 ids is comfortably past 32766, so this fails outright on an
// unchunked read.
func TestGroupsSurviveMoreIDsThanTheDriverBinds(t *testing.T) {
	Config(filepath.Join(t.TempDir(), "big.db"))
	if DB == nil {
		t.Fatal("no database")
	}

	ds := Dataset{Name: "probe/big"}
	if err := DB.Create(&ds).Error; err != nil {
		t.Fatal(err)
	}

	const pairs = 20000 // 40,000 vocabulary ids
	vocab := make([]Vocabulary, 0, pairs*2)
	for i := 0; i < pairs; i++ {
		vocab = append(vocab,
			Vocabulary{Dataset: ds.ID, Token: fmt.Sprintf("w%da", i)},
			Vocabulary{Dataset: ds.ID, Token: fmt.Sprintf("w%db", i)})
	}
	if err := DB.CreateInBatches(&vocab, 2000).Error; err != nil {
		t.Fatal(err)
	}

	edges := make([]Transition, 0, pairs*2)
	for i := 0; i < pairs; i++ {
		a, b := vocab[i*2].ID, vocab[i*2+1].ID
		edges = append(edges,
			Transition{Src: a, Dest: b, Probability: 1},
			Transition{Src: b, Dest: a, Probability: 1})
	}
	if err := DB.CreateInBatches(&edges, 2000).Error; err != nil {
		t.Fatal(err)
	}

	// Both directions of each pair reduce to one group, which is what Groups
	// deduplicates for, so the count is pairs rather than ids.
	got := Groups("probe/big")
	if len(got) != pairs {
		t.Fatalf("Groups returned %d groups, want %d; the read was truncated or "+
			"silently failed", len(got), pairs)
	}
	for i, g := range got {
		if len(g) != 2 {
			t.Fatalf("group %d has %d members, want 2", i, len(g))
		}
	}
}

// The order groups come back in has to be stable, because it reaches admission
// order in the engine and so the content address of a scan. Chunking the read
// moved the sort out of SQL and into Go; this is what says it still happens.
func TestGroupsAreOrderedAcrossChunkBoundaries(t *testing.T) {
	Config(filepath.Join(t.TempDir(), "ord.db"))
	if DB == nil {
		t.Fatal("no database")
	}
	ds := Dataset{Name: "probe/ord"}
	if err := DB.Create(&ds).Error; err != nil {
		t.Fatal(err)
	}

	// Several times the 900-row chunk, so ordering spans many boundaries.
	const pairs = 5000
	vocab := make([]Vocabulary, 0, pairs*2)
	for i := 0; i < pairs; i++ {
		vocab = append(vocab,
			Vocabulary{Dataset: ds.ID, Token: fmt.Sprintf("k%05da", i)},
			Vocabulary{Dataset: ds.ID, Token: fmt.Sprintf("k%05db", i)})
	}
	if err := DB.CreateInBatches(&vocab, 2000).Error; err != nil {
		t.Fatal(err)
	}
	edges := make([]Transition, 0, pairs)
	for i := 0; i < pairs; i++ {
		edges = append(edges, Transition{Src: vocab[i*2].ID, Dest: vocab[i*2+1].ID, Probability: 1})
	}
	if err := DB.CreateInBatches(&edges, 2000).Error; err != nil {
		t.Fatal(err)
	}

	first := Groups("probe/ord")
	for i := 0; i < 3; i++ {
		again := Groups("probe/ord")
		if len(again) != len(first) {
			t.Fatalf("read %d returned %d groups, first returned %d", i, len(again), len(first))
		}
		for j := range first {
			if first[j][0] != again[j][0] {
				t.Fatalf("group %d is %q on one read and %q on the next; order is unstable",
					j, first[j][0], again[j][0])
			}
		}
	}
}
