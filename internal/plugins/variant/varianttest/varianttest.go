// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package varianttest provides fixtures for testing variant algorithms.
//
// The algorithms are data-driven: eight of them read a Language and three read
// a keyboard Layout. Testing them against the shipped datasets would assert on
// curated linguistic data that is expected to grow, so a test written that way
// fails whenever somebody adds a homoglyph. These fixtures are small, fixed and
// obviously wrong as real language data, which is the point -- an expectation
// written against Lang stays true forever.
package varianttest

import (
	"sort"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/kb"
)

// Lang is a Language with just enough data to exercise every relation an
// algorithm reads, and no more.
type Lang struct{}

func (Lang) Code() string        { return "xx" }
func (Lang) Name() string        { return "Test" }
func (Lang) Vowels() []string    { return []string{"a", "e", "i", "o", "u"} }
func (Lang) Graphemes() []string { return []string{"a", "b", "c"} }

// Numerals maps a digit to the words for it, which is what the cardinal and
// ordinal algorithms substitute in both directions.
func (Lang) Numerals() map[string][]string {
	return map[string][]string{
		"1": {"one", "first"},
		"2": {"two", "second"},
	}
}

// Homoglyphs maps a character to what it can be mistaken for. Digit-for-letter
// is the clearest case to assert on.
func (Lang) Homoglyphs() map[string][]string {
	return map[string][]string{"o": {"0"}, "l": {"1"}}
}

func (Lang) Homophones() [][]string   { return [][]string{{"to", "two", "too"}} }
func (Lang) Misspellings() [][]string { return [][]string{{"seperate", "separate"}} }

// Langs is the one-element slice the language-driven Specs take.
func Langs() []variant.Language { return []variant.Language{Lang{}} }

// Boards is a one-element slice holding the US layout. The keyboard algorithms
// are tested against a real layout rather than a literal adjacency map because
// the interesting cases -- that "a" neighbours "s" and "q" -- are facts about
// the geometry, and a hand-written map would only assert that the test author
// and the code agree.
func Boards() []*kb.Layout { return []*kb.Layout{kb.MustGet("kbdus")} }

// Sorted copies and sorts, so an expectation does not depend on the order a
// generator happens to emit in.
func Sorted(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

// Has reports whether want is in got.
func Has(got []string, want string) bool {
	for _, g := range got {
		if g == want {
			return true
		}
	}
	return false
}

// Missing returns the elements of want that are not in got.
func Missing(got, want []string) []string {
	var out []string
	for _, w := range want {
		if !Has(got, w) {
			out = append(out, w)
		}
	}
	return out
}
