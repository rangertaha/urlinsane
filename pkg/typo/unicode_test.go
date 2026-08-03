// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package typo

import (
	"testing"
	"unicode/utf8"
)

// The class test for this package.
//
// Every generator here takes a name and returns names. A name is a sequence of
// characters, and a character is a rune — so a generator that indexes bytes
// cuts multi-byte characters in half and returns strings that are not names at
// all. That is not a hypothetical: eight of these functions did exactly that,
// and `CharacterOmission("яндекс")` returned six variants, all six invalid
// UTF-8, which the graph then admitted as node keys.
//
// This walks every generator with non-Latin input and asserts the two
// properties that make an output a name: it is valid UTF-8, and it is not
// simply the input. It exists so the next generator written byte-wise fails the
// build here rather than shipping variants nobody can read.
//
// BitFlipping is deliberately absent. It models a bit flipping in a resolver's
// memory rather than a human mistyping, so byte-level output is its whole
// point and invalid UTF-8 is a legitimate result.
func TestGeneratorsAreRuneSafe(t *testing.T) {
	// One Cyrillic, one with a combining mark, one CJK, one Latin control.
	inputs := []string{"яндекс", "münchen", "日本語", "example"}

	// Keyboard-driven generators need a layout; this is the shape
	// adjacentCharacters expects.
	layout := []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}

	graphemes := []string{"а", "b", "ü"}
	homoglyphs := map[string][]string{"я": {"r"}, "e": {"е"}, "n": {"п"}}

	gens := map[string]func(string) []string{
		"CharacterSwapping":             CharacterSwapping,
		"CharacterOmission":             CharacterOmission,
		"CharacterRepetition":           CharacterRepetition,
		"HyphenInsertion":               HyphenInsertion,
		"HyphenOmission":                HyphenOmission,
		"DotInsertion":                  DotInsertion,
		"DotOmission":                   DotOmission,
		"DotHyphenSubstitution":         DotHyphenSubstitution,
		"AdjacentCharacterSubstitution": func(s string) []string { return AdjacentCharacterSubstitution(s, layout...) },
		"AdjacentCharacterInsertion":    func(s string) []string { return AdjacentCharacterInsertion(s, layout...) },
		"RepetitionAdjacentReplacement": func(s string) []string { return RepetitionAdjacentReplacement(s, layout...) },
		"GraphemeInsertion":             func(s string) []string { return GraphemeInsertion(s, graphemes...) },
		"GraphemeReplacement":           func(s string) []string { return GraphemeReplacement(s, graphemes...) },
		"HomoglyphSwapping":             func(s string) []string { return HomoglyphSwapping(s, homoglyphs) },
	}

	for name, gen := range gens {
		for _, in := range inputs {
			t.Run(name+"/"+in, func(t *testing.T) {
				for _, got := range gen(in) {
					if !utf8.ValidString(got) {
						t.Errorf("%s(%q) produced invalid UTF-8: %q — indexing bytes where it should index runes",
							name, in, got)
					}
					if got == in {
						t.Errorf("%s(%q) returned the input as a variant", name, in)
					}
					if got == "" {
						t.Errorf("%s(%q) returned an empty variant", name, in)
					}
				}
			})
		}
	}
}

// Variants come back in a stable order.
//
// Several generators deduplicated through a map and ranged it, so the order
// changed between runs of the same input. That order reaches admission order in
// the engine, and admission order decides which candidates survive a frontier
// or a budget — so an unstable order here reaches the content address of a
// scan, and two identical scans stop being identical.
func TestGeneratorOrderIsStable(t *testing.T) {
	gens := map[string]func(string) []string{
		"DotInsertion":          DotInsertion,
		"DotOmission":           DotOmission,
		"HyphenOmission":        HyphenOmission,
		"DotHyphenSubstitution": DotHyphenSubstitution,
		"CharacterOmission":     CharacterOmission,
	}
	const in = "one.two-three.four"

	for name, gen := range gens {
		first := gen(in)
		for i := 0; i < 20; i++ {
			got := gen(in)
			if len(got) != len(first) {
				t.Fatalf("%s returned %d variants, then %d", name, len(first), len(got))
			}
			for j := range got {
				if got[j] != first[j] {
					t.Fatalf("%s is not order-stable: position %d was %q, now %q",
						name, j, first[j], got[j])
				}
			}
		}
	}
}

// The gap between the last two characters is a real insertion point.
//
// Every insertion generator special-cased its final iteration to move the
// separator to the end instead of using the last interior gap, so "abc" never
// yielded "ab-c" — and for a two-character name that was the only interior gap
// there was.
func TestInsertionCoversEveryGap(t *testing.T) {
	for _, tc := range []struct {
		name string
		got  []string
		want string
	}{
		{"HyphenInsertion", HyphenInsertion("abc"), "ab-c"},
		{"DotInsertion", DotInsertion("abc"), "ab.c"},
		{"GraphemeInsertion", GraphemeInsertion("abc", "z"), "abzc"},
	} {
		var found bool
		for _, v := range tc.got {
			if v == tc.want {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: %q is missing from %q", tc.name, tc.want, tc.got)
		}
	}

	// A two-character name has exactly one interior gap, and dropping it left
	// DotInsertion with nothing to return but the input.
	if got := DotInsertion("ab"); len(got) != 1 || got[0] != "a.b" {
		t.Errorf("DotInsertion(%q) = %q, want [a.b]", "ab", got)
	}
}

// Both directions of the separator swap.
func TestDotHyphenSubstitutionSwapsBothWays(t *testing.T) {
	want := map[string]bool{"one-two.three": true, "one.two-three": true, "one-two-three": true}
	got := DotHyphenSubstitution("one.two.three")
	if len(got) != len(want) {
		t.Errorf("DotHyphenSubstitution(one.two.three) = %q, want the %d single- and all-swaps", got, len(want))
	}
	for _, v := range got {
		if !want[v] {
			t.Errorf("unexpected variant %q — a substitution should not omit a separator", v)
		}
	}

	// A name with hyphens and no dots produced nothing at all, which is every
	// package and repo named like "my-lib".
	if got := DotHyphenSubstitution("my-lib"); len(got) != 1 || got[0] != "my.lib" {
		t.Errorf("DotHyphenSubstitution(my-lib) = %q, want [my.lib]", got)
	}
}
