// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/datasets"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/xhs"
)

// The dataset has to be there and be groups, not a word list.
func TestCrossHomophoneDatasetLoads(t *testing.T) {
	groups := datasets.CrossHomophones()
	if len(groups) < 1000 {
		t.Fatalf("loaded %d cross-language groups, want thousands — the embed or the parse is wrong",
			len(groups))
	}
	for i, g := range groups {
		if len(g) < 2 {
			t.Fatalf("group %d has %d members; a group of one is not a group", i, len(g))
		}
		for _, w := range g {
			if strings.ContainsAny(w, " \t#") {
				t.Fatalf("group %d member %q was not tokenised", i, w)
			}
		}
	}
}

// The cases the technique is named for. Each is one pronunciation spelled the
// way a different language writes that sound.
func TestCrossHomophoneGeneratesKnownSoundSquats(t *testing.T) {
	gen := xhs.Spec().Gen

	for _, tc := range []struct{ name, want string }{
		{"youtube", "yutup"},     // Turkish/Indonesian spelling of the sound
		{"youtube.com", "yutup"}, // and inside a longer key
		{"boutique", "butik"},    // one word, four orthographies
		{"sushi", "susyi"},
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

// It must not return the input, and must leave a name it knows nothing about
// alone rather than inventing something.
func TestCrossHomophoneIsQuietWhenItHasNothing(t *testing.T) {
	gen := xhs.Spec().Gen

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

// Deterministic: the groups are snapshotted at construction, so two runs of the
// same build produce the same list in the same order. Variant order reaches
// admission order, which reaches the scan's content address.
func TestCrossHomophoneIsOrderStable(t *testing.T) {
	first := xhs.Spec().Gen("boutique")
	for i := 0; i < 10; i++ {
		got := xhs.Spec().Gen("boutique")
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
