// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Language data: what the variant algorithms run over.
//
// A language is not a plugin. It was one — thirty Go files, each restating its
// vowels, graphemes, homoglyphs and misspellings as literals — and every one of
// them had to be edited, compiled and shipped to correct a word. The same data
// now lives in datasets/languages/<code>/*.lst, is imported into the dataset
// database, and is read from there: one implementation for every language,
// and adding a language is adding a directory.
//
// Nothing here holds behaviour. The algorithms in internal/plugins/variant
// decide what to do with vowels; this decides only what the vowels are.
package dataset

import (
	"fmt"
	"sort"
	"strings"
)

// The dataset names this package reads, matching the .lst basenames.
const (
	dsVowel       = "vowel"
	dsGrapheme    = "grapheme"
	dsNumeral     = "numeral"
	dsHomoglyph   = "homoglyph"
	dsHomophone   = "homophone"
	dsMisspelling = "misspelling"
)

// Lang is one language's data, read on demand.
//
// A value type holding only its identity: the data is fetched per call rather
// than loaded up front, because a scan restricted to two languages should not
// pay to read thirty. Callers that read repeatedly should hold the result, not
// the Lang.
type Lang struct {
	code string
	name string
}

// Code is the dataset directory name: "en", "fr", "ru".
func (l Lang) Code() string { return l.code }

// Name is the language in English, or the code when the dataset names no other.
func (l Lang) Name() string {
	if l.name == "" {
		return l.code
	}
	return l.name
}

// Vowels are the syllabic characters, for vowel swapping.
func (l Lang) Vowels() []string { return tokens(l.code, dsVowel) }

// Graphemes are the smallest functional units of the writing system.
func (l Lang) Graphemes() []string { return tokens(l.code, dsGrapheme) }

// Numerals maps a number to the words denoting it: "1" -> {"one", "first"}.
func (l Lang) Numerals() map[string][]string { return edges(l.code, dsNumeral) }

// Homoglyphs maps a character to the characters that look like it.
func (l Lang) Homoglyphs() map[string][]string { return edges(l.code, dsHomoglyph) }

// Homophones are groups of words that sound alike.
func (l Lang) Homophones() [][]string { return groups(l.code, dsHomophone) }

// Misspellings are groups whose members are habitually confused.
func (l Lang) Misspellings() [][]string { return groups(l.code, dsMisspelling) }

// LanguageCodes lists every language the dataset carries, sorted.
//
// Sorted rather than in insertion order because an operator built from these is
// cached against a digest of what it read: an unstable order would invalidate
// that cache on every run.
func LanguageCodes() []string {
	if DB == nil {
		return nil
	}
	var rows []Language
	if err := DB.Order("code").Find(&rows).Error; err != nil {
		return nil
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		if r.Code != "" {
			out = append(out, r.Code)
		}
	}
	return out
}

// Languages returns every language in the dataset, in code order.
func Languages() []Lang {
	codes := LanguageCodes()
	out := make([]Lang, 0, len(codes))
	for _, c := range codes {
		out = append(out, Lang{code: c, name: nameOf(c)})
	}
	return out
}

// LanguageOf returns one language, or false if the dataset does not carry it.
func LanguageOf(code string) (Lang, bool) {
	for _, c := range LanguageCodes() {
		if c == code {
			return Lang{code: c, name: nameOf(c)}, true
		}
	}
	return Lang{}, false
}

// SelectLanguages returns the named languages in code order, or every language when no
// codes are given.
//
// An unknown code is an error rather than a silent omission. A run that quietly
// dropped a language would generate fewer candidates and still report the
// result as complete — and the failure this data exists to prevent is a
// squatted name nobody generated.
func SelectLanguages(codes ...string) ([]Lang, error) {
	known := LanguageCodes()
	if len(codes) == 0 {
		return Languages(), nil
	}

	want := map[string]bool{}
	for _, raw := range codes {
		for _, c := range strings.Split(raw, ",") {
			if c = strings.TrimSpace(c); c != "" {
				want[c] = true
			}
		}
	}

	have := map[string]bool{}
	for _, c := range known {
		have[c] = true
	}
	var unknown []string
	for c := range want {
		if !have[c] {
			unknown = append(unknown, c)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return nil, fmt.Errorf("unknown language(s): %s; known: %s",
			strings.Join(unknown, ", "), strings.Join(known, ", "))
	}

	out := make([]Lang, 0, len(want))
	for _, c := range known {
		if want[c] {
			out = append(out, Lang{code: c, name: nameOf(c)})
		}
	}
	return out, nil
}

// --- dataset queries --------------------------------------------------------

func nameOf(code string) string {
	if DB == nil {
		return ""
	}
	var row Language
	if err := DB.Where("code = ?", code).First(&row).Error; err != nil {
		return ""
	}
	return row.Name
}

// ids resolves a (language, dataset) pair to the row ids the tables are keyed
// by. Zero ids mean the pair is not in the dataset, which is a normal state:
// not every language ships every relation.
func ids(code, name string) (lang, ds uint, ok bool) {
	if DB == nil {
		return 0, 0, false
	}
	var l Language
	if err := DB.Where("code = ?", code).First(&l).Error; err != nil {
		return 0, 0, false
	}
	var d Dataset
	if err := DB.Where("name = ?", name).First(&d).Error; err != nil {
		return 0, 0, false
	}
	return l.ID, d.ID, true
}

// tokens returns one table's vocabulary, in id order — which is the order the
// importer inserted it, which is sorted.
func tokens(code, name string) []string {
	lang, ds, ok := ids(code, name)
	if !ok {
		return nil
	}
	var rows []Vocabulary
	if err := DB.Where("language = ? AND dataset = ?", lang, ds).
		Order("id").Find(&rows).Error; err != nil {
		return nil
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.Token)
	}
	return out
}

// edges returns the transitions of one table as source -> targets.
//
// This is the natural shape for the relations that *are* a mapping: a homoglyph
// table says what "o" can be confused with, and a numeral table says what "1"
// can be written as. Both were map[string][]string when they were Go literals,
// and a transition row is exactly one entry of one.
func edges(code, name string) map[string][]string {
	rows, vocab, ok := transitions(code, name)
	if !ok {
		return nil
	}
	out := map[string][]string{}
	for _, t := range rows {
		src, sok := vocab[t.Src]
		dst, dok := vocab[t.Dest]
		if sok && dok {
			out[src] = append(out[src], dst)
		}
	}
	for k := range out {
		sort.Strings(out[k])
	}
	return out
}

// groups returns the transitions of one table as {source, targets...} groups.
//
// The relations that are *sets* rather than mappings — homophones, habitual
// misspellings — were [][]string when they were literals, one slice per line of
// the .lst. A line was a clique when it was imported, so every member points at
// every other; taking each source with its targets recovers the line, once per
// member. The algorithms swap within a group, so seeing the same set from each
// of its members costs a duplicate the caller already dedupes and avoids
// reconstructing cliques, which would over-merge two lines that share a word.
func groups(code, name string) [][]string {
	rows, vocab, ok := transitions(code, name)
	if !ok {
		return nil
	}
	by := map[string][]string{}
	for _, t := range rows {
		src, sok := vocab[t.Src]
		dst, dok := vocab[t.Dest]
		if sok && dok {
			by[src] = append(by[src], dst)
		}
	}
	keys := make([]string, 0, len(by))
	for k := range by {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([][]string, 0, len(keys))
	for _, k := range keys {
		members := by[k]
		sort.Strings(members)
		out = append(out, append([]string{k}, members...))
	}
	return out
}

// transitions reads one table's edges and the vocabulary they index.
func transitions(code, name string) ([]Transition, map[uint]string, bool) {
	lang, ds, ok := ids(code, name)
	if !ok {
		return nil, nil, false
	}
	var vocab []Vocabulary
	if err := DB.Where("language = ? AND dataset = ?", lang, ds).
		Find(&vocab).Error; err != nil {
		return nil, nil, false
	}
	if len(vocab) == 0 {
		return nil, nil, false
	}
	byID := make(map[uint]string, len(vocab))
	ids := make([]uint, 0, len(vocab))
	for _, v := range vocab {
		byID[v.ID] = v.Token
		ids = append(ids, v.ID)
	}
	var rows []Transition
	if err := DB.Where("src IN ?", ids).Order("src, dest").
		Find(&rows).Error; err != nil {
		return nil, nil, false
	}
	return rows, byID, true
}
