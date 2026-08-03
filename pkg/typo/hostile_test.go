// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package typo

import (
	"strings"
	"testing"
)

// hostileTokens is what a generator may actually arrive with once a name has
// been through decomposition: empty, blank, unprintable, outside the Basic
// Multilingual Plane, combining, right-to-left, and long.
var hostileTokens = []string{
	"", " ", "\x00", "\n", "\t", ".", "-", "..", "--", ".-.",
	"a", "ab", "1", "11", "a1", "1a",
	"中", "🙂", "é", "é", "لا", "ß",
	"a.b-c_d", "example.com", "npm:lodash",
	strings.Repeat("a", 300),
	strings.Repeat("ab-", 100),
}

// hostileNumerals covers the shapes the numeral relation is handed in practice.
//
// clique is the one that matters and the one that was missing: dataset.Lang
// stores a numeral line as transitions between every pair of tokens on it, so
// each word becomes a key and "1" -> "first" and "first" -> "1" both exist.
// Against a per-frame seen set that is a two-cycle, and CardinalSwap recursed
// until the stack gave out. Every fixture in this package was keyed by digit
// alone, which cannot cycle, so nothing here could catch it.
var hostileNumerals = map[string]map[string][]string{
	"empty":     {},
	"digitOnly": {"1": {"one", "first"}, "2": {"two", "second"}},
	"clique": {
		"1": {"first", "one"}, "one": {"1", "first"}, "first": {"1", "one"},
		"2": {"second", "two"}, "two": {"2", "second"}, "second": {"2", "two"},
	},
	"selfMapping":  {"1": {"1"}, "one": {"one"}},
	"emptyStrings": {"": {""}, "1": {""}},
	"substrings":   {"1": {"one"}, "11": {"eleven"}, "111": {"onehundredeleven"}},
	"singleWord":   {"0": {"zero"}},
}

// A generator must return from every one of these. A generator that does not is
// not a slow generator: the numeral walk took the whole process down with a
// stack overflow, which no recover can catch.
func TestNoGeneratorPanicsOrHangs(t *testing.T) {
	kb := []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}
	graphemes := []string{"a", "b", "", "中"}
	vowels := []string{"a", "e", "i", "o", "u", ""}
	groups := [][]string{{"to", "two", "too"}, {"a"}, {}, {"", ""}}
	glyphs := map[string][]string{"a": {"4", "@"}, "o": {"0"}, "": {"x"}, "z": {"z"}}

	gens := map[string]func(string) []string{
		"CharacterSwapping":             CharacterSwapping,
		"CharacterOmission":             CharacterOmission,
		"CharacterRepetition":           CharacterRepetition,
		"HyphenInsertion":               HyphenInsertion,
		"HyphenOmission":                HyphenOmission,
		"DotInsertion":                  DotInsertion,
		"DotOmission":                   DotOmission,
		"DotHyphenSubstitution":         DotHyphenSubstitution,
		"TokenOrderSwap":                TokenOrderSwap,
		"SingularPluralise":             SingularPluralise,
		"AdjacentCharacterSubstitution": func(s string) []string { return AdjacentCharacterSubstitution(s, kb...) },
		"AdjacentCharacterInsertion":    func(s string) []string { return AdjacentCharacterInsertion(s, kb...) },
		"RepetitionAdjacentReplacement": func(s string) []string { return RepetitionAdjacentReplacement(s, kb...) },
		"GraphemeInsertion":             func(s string) []string { return GraphemeInsertion(s, graphemes...) },
		"GraphemeReplacement":           func(s string) []string { return GraphemeReplacement(s, graphemes...) },
		"VowelSwapping":                 func(s string) []string { return VowelSwapping(s, vowels...) },
		"BitFlipping":                   func(s string) []string { return BitFlipping(s, graphemes...) },
		"CommonMisspellings":            func(s string) []string { return CommonMisspellings(s, groups...) },
		"HomophoneSwapping":             func(s string) []string { return HomophoneSwapping(s, groups...) },
		"HomoglyphSwapping":             func(s string) []string { return HomoglyphSwapping(s, glyphs) },
		"StemSwapping":                  func(s string) []string { return StemSwapping(s, []string{"a", "", "ab"}) },
		"EmojiInsertion":                func(s string) []string { return EmojiInsertion(s, []string{"🙂", ""}) },
		"PrefixInsertion":               func(s string) []string { return PrefixInsertion(s, "py", "") },
		"SuffixInsertion":               func(s string) []string { return SuffixInsertion(s, "js", "") },
	}
	for shape, nums := range hostileNumerals {
		gens["CardinalSwap/"+shape] = func(s string) []string { return CardinalSwap(s, nums) }
		gens["OrdinalSwap/"+shape] = func(s string) []string { return OrdinalSwap(s, nums) }
	}

	for name, gen := range gens {
		for _, in := range hostileTokens {
			label := in
			if len(label) > 20 {
				label = label[:20] + "..."
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("%s(%q) panicked: %v", name, label, r)
					}
				}()
				for _, got := range gen(in) {
					if got == in {
						t.Errorf("%s(%q) returned its own input", name, label)
						break
					}
				}
			}()
		}
	}
}

// The inflection client is built once, not once per name.
//
// nlp.NewClient loads the whole rule set and compiles a regexp; building it per
// call cost 488us and 2311 allocations against 41us and 10 for the shared one,
// and SingularPluralise runs once per candidate. Comparing the pointer is the
// cheap way to say "still shared" without a timing assertion that would flake.
func TestPluralizerIsBuiltOnce(t *testing.T) {
	first := pluralizer()
	if first == nil {
		t.Fatal("pluralizer returned nil")
	}
	for i := 0; i < 100; i++ {
		if got := pluralizer(); got != first {
			t.Fatalf("call %d returned a different client; the rule set is being "+
				"rebuilt per call", i)
		}
	}
}

// SplitTokens is not a generator but every domain-shaped one is built on it.
func TestSplitTokensSurvivesHostileInput(t *testing.T) {
	for _, in := range hostileTokens {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("SplitTokens(%q) panicked: %v", in, r)
				}
			}()
			parts, seps := SplitTokens(in)
			// The contract the reordering algorithms rely on: one separator
			// between each adjacent pair, or nothing at all.
			if parts != nil && len(seps) != len(parts)-1 {
				t.Errorf("SplitTokens(%q) = %d parts and %d separators; a token "+
					"boundary needs exactly one separator between each pair",
					in, len(parts), len(seps))
			}
		}()
	}
}
