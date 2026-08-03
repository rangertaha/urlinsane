// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/rangertaha/urlinsane/internal/dataset/gen"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/xhs"
)

// crossFixture imports the shipped .lst into a scratch database and reads it
// back the way a run does.
//
// It goes through the real importer rather than parsing the file directly,
// because the round trip is where this data can break: a line is stored as a
// clique of transitions and reconstructed one group per member, so a reader
// that does not deduplicate gets each group back as many times as it has
// words. Testing the file alone would not catch that, and testing the shipped
// dataset.db would test an artifact rather than the source.
func crossFixture(t *testing.T) [][]string {
	t.Helper()
	dataset.Config(filepath.Join(t.TempDir(), "test.db"))
	if dataset.DB == nil {
		t.Skip("no scratch database")
	}
	src := filepath.Join("..", "..", "..", "..", "datasets", "phonetics", "homophone.lst")
	if err := gen.One("", "phonetics/homophone", src); err != nil {
		t.Fatalf("import %s: %v", src, err)
	}
	return dataset.Groups("phonetics/homophone")
}

// The shipped file has to survive the import as groups, not as a word list.
func TestCrossHomophoneDatasetLoads(t *testing.T) {
	groups := crossFixture(t)
	if len(groups) < 1000 {
		t.Fatalf("loaded %d cross-language groups, want thousands", len(groups))
	}
	for i, g := range groups {
		if len(g) < 2 {
			t.Fatalf("group %d has %d members; a group of one is not a group", i, len(g))
		}
	}
}

// The cases the technique is named for: one pronunciation, spelled the way
// different languages write that sound.
func TestCrossHomophoneGeneratesKnownSoundSquats(t *testing.T) {
	gen := xhs.Spec(crossFixture(t)).Gen

	for _, tc := range []struct{ name, want string }{
		{"youtube", "yutup"},     // Turkish/Indonesian spelling of the sound
		{"youtube.com", "yutup"}, // and inside a longer key
		{"boutique", "butik"},    // one word, four orthographies
		{"google", "googel"},     // a real registered squat
		{"baby", "bebi"},
	} {
		got := gen(tc.name)
		var found bool
		for _, g := range got {
			if strings.Contains(g, tc.want) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("xhs(%q) = %q, want something containing %q", tc.name, got, tc.want)
		}
	}
}

// No duplicates. A word with two pronunciations is in two groups, so the same
// substitution is reachable twice and the naive version emitted it twice.
func TestCrossHomophoneDoesNotRepeatItself(t *testing.T) {
	got := xhs.Spec(crossFixture(t)).Gen("youtube.com")
	if len(got) == 0 {
		t.Fatal("fixture produced nothing, so this proves nothing")
	}
	seen := map[string]bool{}
	for _, g := range got {
		if seen[g] {
			t.Errorf("xhs returned %q twice: %q", g, got)
		}
		seen[g] = true
	}
}

// Quiet when it has nothing, rather than inventing something.
func TestCrossHomophoneIsQuietWhenItHasNothing(t *testing.T) {
	gen := xhs.Spec(crossFixture(t)).Gen

	if got := gen("zzqxvwk"); len(got) != 0 {
		t.Errorf("xhs(zzqxvwk) = %q, want nothing", got)
	}
	if got := gen(""); len(got) != 0 {
		t.Errorf("xhs(\"\") = %q, want nothing", got)
	}
	for _, g := range gen("youtube") {
		if g == "youtube" {
			t.Error("xhs returned the input as a variant")
		}
	}
}

// Order-stable: variant order reaches admission order, which reaches the
// scan's content address.
func TestCrossHomophoneIsOrderStable(t *testing.T) {
	groups := crossFixture(t)
	first := xhs.Spec(groups).Gen("boutique")
	for i := 0; i < 10; i++ {
		got := xhs.Spec(groups).Gen("boutique")
		if len(got) != len(first) {
			t.Fatalf("run %d produced %d variants, first run produced %d", i, len(got), len(first))
		}
		for j := range got {
			if got[j] != first[j] {
				t.Fatalf("run %d differs at %d: %q vs %q", i, j, got[j], first[j])
			}
		}
	}
}
