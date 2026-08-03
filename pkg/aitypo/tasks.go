// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package aitypo

import "github.com/rangertaha/urlinsane/pkg/typo"

// Data is the reference data the table-driven tasks read.
//
// Every field is optional and a nil one silences its tasks rather than
// erroring: a corpus for the pure rules should not require a language database
// to exist. Tasks whose data is missing generate empty expectations, which
// Emit drops, so an absent table shows up as a task with no examples rather
// than as a task that appears to have been learned perfectly.
//
// The shapes match pkg/typo's parameters exactly so nothing has to be adapted
// here; fill them from internal/dataset, from pkg/kb, or by hand in a test.
type Data struct {
	// Graphemes is the alphabet for insertion and replacement.
	Graphemes []string
	// Vowels feeds vowel swapping.
	Vowels []string
	// Homoglyphs maps a character to the characters that look like it.
	Homoglyphs map[string][]string
	// Homophones are groups of words that sound alike. Cross-language groups
	// go here too — the generator does not care which it is given, and the
	// difference is worth recording in the corpus Source rather than in a
	// second task.
	Homophones [][]string
	// Misspellings are groups whose members are habitually confused.
	Misspellings [][]string
	// Numerals maps a digit to the words denoting it.
	Numerals map[string][]string
	// Keyboard is a layout as rows of characters, the shape
	// pkg/typo's adjacency helpers expect.
	Keyboard []string
}

// Tasks builds the registry from d.
//
// Every entry wraps a pkg/typo function directly. Nothing is reimplemented
// here, and nothing should be: the oracle has to be the same code the scanner
// runs, or a model trained to match this package would be trained against a
// second definition of the truth.
func Tasks(d Data) Registry {
	r := Registry{}
	add := func(id, title string, needs Needs, o Oracle) {
		r[id] = Task{ID: id, Title: title, Needs: needs, Oracle: o}
	}

	// --- rules: a model that reproduces these has learned a transformation ---

	add("co", "Character Omission", NeedsNothing, typo.CharacterOmission)
	add("cs", "Character Swapping", NeedsNothing, typo.CharacterSwapping)
	add("cr", "Character Repetition", NeedsNothing, typo.CharacterRepetition)
	add("hi", "Hyphen Insertion", NeedsNothing, typo.HyphenInsertion)
	add("ho", "Hyphen Omission", NeedsNothing, typo.HyphenOmission)
	add("di", "Dot Insertion", NeedsNothing, typo.DotInsertion)
	add("do", "Dot Omission", NeedsNothing, typo.DotOmission)
	add("dhs", "Dot Hyphen Substitution", NeedsNothing, typo.DotHyphenSubstitution)
	add("tos", "Token Order Swap", NeedsNothing, typo.TokenOrderSwap)

	// Bit flipping is a rule, but not a rule about characters: it models a bit
	// flipping in a resolver's memory, so its output is byte-level and may not
	// be valid UTF-8. It is included because it is learnable and excluded from
	// any claim about plausible human typing.
	add("bf", "Bit Flipping", NeedsNothing, func(s string) []string {
		return typo.BitFlipping(s)
	})

	// Singular/plural inflection is a table, not a rule, and it sat in the block
	// above until this comment was written. typo.SingularPluralise goes through
	// nlp.NewClient, which loads every irregular, plural, singular and
	// uncountable entry it knows — so mouse/mice, goose/geese, child/children,
	// index/indices, ox/oxen, person/people, datum/data, cactus/cacti and the
	// uncountables (news -> nothing) are dictionary lookups, derivable from the
	// input by no rule at all. Reported as NeedsNothing, a model scoring 1.000
	// on sp was credited with a rule it had in fact memorised as an English
	// dictionary, which is the one distinction Needs exists to keep.
	//
	// It is not gated on a Data field like the tables below, because its table
	// is compiled into the library rather than supplied by the caller. That has
	// a consequence worth stating: sp is the one task that is *not* silenced on
	// a non-English corpus, and it will happily emit English inflection as
	// ground truth for other languages — sp("maison") is ["maisons"] and
	// sp("hund") is ["hunds"]. Treat sp results on a non-English corpus as
	// measuring English morphology applied to foreign strings.
	add("sp", "Singular Pluralise", NeedsLanguage, typo.SingularPluralise)

	// --- tables: reproducing these means memorising a lookup ---

	if len(d.Graphemes) > 0 {
		add("gi", "Grapheme Insertion", NeedsLanguage, func(s string) []string {
			return typo.GraphemeInsertion(s, d.Graphemes...)
		})
		add("gr", "Grapheme Replacement", NeedsLanguage, func(s string) []string {
			return typo.GraphemeReplacement(s, d.Graphemes...)
		})
	}
	if len(d.Vowels) > 0 {
		add("vs", "Vowel Swapping", NeedsLanguage, func(s string) []string {
			return typo.VowelSwapping(s, d.Vowels...)
		})
	}
	if len(d.Homoglyphs) > 0 {
		add("hr", "Homoglyph Replacement", NeedsLanguage, func(s string) []string {
			return typo.HomoglyphSwapping(s, d.Homoglyphs)
		})
	}
	if len(d.Homophones) > 0 {
		add("hs", "Homophone Substitution", NeedsLanguage, func(s string) []string {
			return typo.HomophoneSwapping(s, d.Homophones...)
		})
	}
	if len(d.Misspellings) > 0 {
		add("cm", "Common Misspellings", NeedsLanguage, func(s string) []string {
			return typo.CommonMisspellings(s, d.Misspellings...)
		})
	}
	if len(d.Numerals) > 0 {
		add("cns", "Cardinal Substitution", NeedsLanguage, func(s string) []string {
			return typo.CardinalSwap(s, d.Numerals)
		})
		add("ons", "Ordinal Substitution", NeedsLanguage, func(s string) []string {
			return typo.OrdinalSwap(s, d.Numerals)
		})
	}
	if len(d.Keyboard) > 0 {
		add("acs", "Adjacent Character Substitution", NeedsKeyboard, func(s string) []string {
			return typo.AdjacentCharacterSubstitution(s, d.Keyboard...)
		})
		add("aci", "Adjacent Character Insertion", NeedsKeyboard, func(s string) []string {
			return typo.AdjacentCharacterInsertion(s, d.Keyboard...)
		})
		add("rar", "Repetition Adjacent Replacement", NeedsKeyboard, func(s string) []string {
			return typo.RepetitionAdjacentReplacement(s, d.Keyboard...)
		})
	}
	return r
}

// Emit runs each task over each input and returns one example per pair that
// produced anything.
//
// Inputs that produce nothing are dropped rather than written with an empty
// expectation. An empty set is a real answer — "example" has no doubled
// characters, so rar produces nothing — but a corpus of mostly-empty rows
// trains a model to answer nothing, and the honest place to teach "sometimes
// the answer is empty" is a deliberate sample rather than an accident of which
// names the corpus happened to hold. EmitWithEmpty is that deliberate sample.
func Emit(tasks []Task, corpus []string, source string) []Example {
	return emit(tasks, corpus, source, false)
}

// EmitWithEmpty keeps the pairs whose expectation is empty.
func EmitWithEmpty(tasks []Task, corpus []string, source string) []Example {
	return emit(tasks, corpus, source, true)
}

func emit(tasks []Task, corpus []string, source string, keepEmpty bool) []Example {
	out := make([]Example, 0, len(tasks)*len(corpus))
	for _, in := range corpus {
		if in == "" {
			continue
		}
		for _, t := range tasks {
			exp := t.Expect(in)
			if len(exp) == 0 && !keepEmpty {
				continue
			}
			out = append(out, Example{
				Task:   t.ID,
				Input:  in,
				Expect: exp,
				Source: source,
			})
		}
	}
	return out
}
