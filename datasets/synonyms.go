// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package datasets

import (
	"embed"
	"sort"
	"strings"
	"sync"
)

// synonymFS carries the per-language synonym sets. They are embedded rather
// than generated into Go source because the .lst files are the source of truth
// — `urlinsane-datasets sync languages` writes them — and a generated copy
// would be a second place for the same data to be wrong.
//
//go:embed languages/*/synonym.lst
var synonymFS embed.FS

// A synonym set is one line: whitespace-separated words that mean the same
// thing in that language. Despite the name these are not general-purpose
// thesaurus entries. They are the brand-adjacent vocabulary an attacker bolts
// onto a target name — login, verify, pay, invoice, delivery — grouped so that
// picking one theme gives every word that expresses it. See defaultSynonyms in
// cmd/datasets/sync_languages.go for how they are produced.
var (
	synonymOnce sync.Once
	synonymSets map[string][][]string
)

func loadSynonyms() {
	synonymSets = make(map[string][][]string)
	entries, err := synonymFS.ReadDir("languages")
	if err != nil {
		// The embed directive guarantees the tree exists; a failure here is a
		// build problem, not a runtime condition worth propagating to every
		// caller of Synonyms.
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		lang := e.Name()
		b, err := synonymFS.ReadFile("languages/" + lang + "/synonym.lst")
		if err != nil {
			continue
		}
		var sets [][]string
		for _, line := range strings.Split(string(b), "\n") {
			words := strings.Fields(line)
			if len(words) == 0 {
				continue
			}
			sets = append(sets, words)
		}
		if len(sets) > 0 {
			synonymSets[lang] = sets
		}
	}
}

// Synonyms returns the synonym sets for a language id ("en", "fr", ...), or nil
// if that language ships none. The returned slices are the package's own: treat
// them as read-only.
//
// Order is the file's order, which is the order defaultSynonyms declares them
// in — stable across runs, which is what the operator cache assumes.
func Synonyms(lang string) [][]string {
	synonymOnce.Do(loadSynonyms)
	return synonymSets[lang]
}

// SynonymLanguages returns every language id that ships a synonym set, sorted.
func SynonymLanguages() []string {
	synonymOnce.Do(loadSynonyms)
	out := make([]string, 0, len(synonymSets))
	for lang := range synonymSets {
		out = append(out, lang)
	}
	sort.Strings(out)
	return out
}

// SynonymWords returns every distinct word across a language's synonym sets,
// sorted. This is the flat vocabulary form the combosquatting algorithm wants:
// which theme a word belongs to does not change how it attaches to a name.
func SynonymWords(lang string) []string {
	sets := Synonyms(lang)
	if len(sets) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(sets)*4)
	out := make([]string, 0, len(sets)*4)
	for _, set := range sets {
		for _, w := range set {
			if w == "" || seen[w] {
				continue
			}
			seen[w] = true
			out = append(out, w)
		}
	}
	sort.Strings(out)
	return out
}
