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

package variant

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// KeyboardSpecs are the algorithms driven by physical keyboard adjacency. The
// layouts come from the language plugins, which stay exactly as they are: their
// datasets are algorithm input, not text to be corrected.
//
// Each generator folds every selected keyboard's output together. The adapter
// dedupes and sorts, so two layouts that agree on a variation produce one edge
// and the result does not depend on which layout was consulted first.
func KeyboardSpecs(boards []internal.Keyboard) []Spec {
	return []Spec{
		{
			ID: "aci", Title: "Adjacent Character Insertion", Version: 1,
			Gen: overKeyboards(boards, typo.AdjacentCharacterInsertion),
		},
		{
			ID: "acs", Title: "Adjacent Character Substitution", Version: 1,
			Gen: overKeyboards(boards, typo.AdjacentCharacterSubstitution),
		},
		{
			ID: "rar", Title: "Repetition Adjacent Replacement", Version: 1,
			Gen: overKeyboards(boards, typo.RepetitionAdjacentReplacement),
		},
	}
}

// overKeyboards adapts a pkg/typo function that takes a keyboard layout into a
// Generate over every selected keyboard.
func overKeyboards(boards []internal.Keyboard, fn func(string, ...string) []string) Generate {
	// Snapshot the layouts at construction. Resolving them per call would make
	// an operator's output depend on plugin registration order at call time,
	// which the scheduler's cache assumes cannot happen.
	layouts := make([][]string, 0, len(boards))
	for _, b := range boards {
		if l := b.Layouts(); len(l) > 0 {
			layouts = append(layouts, l)
		}
	}
	return func(name string) []string {
		var out []string
		for _, l := range layouts {
			out = append(out, fn(name, l...)...)
		}
		return out
	}
}

// LanguageSpecs are the algorithms driven by a language's own data: its vowels,
// graphemes, homoglyphs, homophones, numerals and habitual misspellings.
//
// Every one of them varies the registrable name rather than the whole key. The
// old plugins were inconsistent here — common misspellings and cardinal swaps
// ran over the whole domain — which let a substring match reach into the public
// suffix and silently produce a different-registry variant attributed to a
// spelling algorithm. Suffix changes belong to the tld operator alone.
func LanguageSpecs(langs []internal.Language) []Spec {
	return []Spec{
		{
			ID: "vs", Title: "Vowel Swapping", Version: 1,
			Gen: overLanguages(langs, func(l internal.Language, name string) []string {
				return typo.VowelSwapping(name, l.Vowels()...)
			}),
		},
		{
			ID: "hr", Title: "Homoglyph Replacement", Version: 1,
			Gen: overLanguages(langs, func(l internal.Language, name string) []string {
				return typo.HomoglyphSwapping(name, l.Homoglyphs())
			}),
		},
		{
			ID: "hs", Title: "Homophone Substitution", Version: 1,
			Gen: overLanguages(langs, func(l internal.Language, name string) []string {
				return typo.HomophoneSwapping(name, l.Homophones()...)
			}),
		},
		{
			ID: "cm", Title: "Common Misspellings", Version: 1,
			Gen: overLanguages(langs, func(l internal.Language, name string) []string {
				return typo.CommonMisspellings(name, l.Misspellings()...)
			}),
		},
		{
			ID: "gi", Title: "Grapheme Insertion", Version: 1,
			Gen: overLanguages(langs, func(l internal.Language, name string) []string {
				return typo.GraphemeInsertion(name, l.Graphemes()...)
			}),
		},
		{
			ID: "gr", Title: "Grapheme Replacement", Version: 1,
			Gen: overLanguages(langs, func(l internal.Language, name string) []string {
				return typo.GraphemeReplacement(name, l.Graphemes()...)
			}),
		},
		{
			ID: "cns", Title: "Cardinal Substitution", Version: 1,
			Gen: overLanguages(langs, func(l internal.Language, name string) []string {
				return typo.CardinalSwap(name, l.Numerals())
			}),
		},
		{
			ID: "ons", Title: "Ordinal Substitution", Version: 1,
			Gen: overLanguages(langs, func(l internal.Language, name string) []string {
				return typo.OrdinalSwap(name, l.Numerals())
			}),
		},
	}
}

// overLanguages folds one per-language generator over every selected language.
func overLanguages(langs []internal.Language, fn func(internal.Language, string) []string) Generate {
	// Snapshot for the same reason overKeyboards does.
	selected := append([]internal.Language(nil), langs...)
	return func(name string) []string {
		var out []string
		for _, l := range selected {
			out = append(out, fn(l, name)...)
		}
		return out
	}
}
