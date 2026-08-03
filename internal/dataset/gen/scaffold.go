// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rangertaha/urlinsane/pkg/kb"
)

// Relation is one of the .lst files a language directory holds.
type Relation struct {
	// Name is the file's basename, which is also the dataset name it imports
	// under: "vowel" is datasets/languages/en/vowel.lst.
	Name string

	// Doc says what the file holds and how a line is shaped. It is written
	// into the file as a comment, because a scaffolded corpus is empty and the
	// filename alone does not say whether a line is one entry or a group.
	Doc string
}

// Relations is the file set of a language directory, in the order a language
// is usually curated: the character-level data the variant algorithms need
// first, then the word lists.
//
// The list is here rather than inferred from whichever files an existing
// directory happens to contain, so that a new language is scaffolded with the
// full set and a curated one can be checked against it.
var Relations = []Relation{
	{"grapheme", "Smallest functional units of the writing system, one per line."},
	{"vowel", "Syllabic characters, one per line."},
	{"homoglyph", "A character, then the characters that look like it: a à á а"},
	{"numeral", "A number, then the words denoting it: 1 one first"},
	{"homophone", "Words that sound alike, one group per line."},
	{"misspelling", "Words habitually confused for each other, one group per line."},
	{"synonym", "Words of like meaning, one group per line."},
	{"antonym", "A word and its opposite, one pair per line."},
	{"stopword", "Words carrying no distinguishing sense, one per line."},
	{"positive", "Words with a favourable connotation, one per line."},
	{"negative", "Words with an unfavourable connotation, one per line."},
	{"word", "The language's vocabulary, one word per line."},
}

// aliases are curated directories whose name is a code kb does not use, mapped
// to the code kb ships for the same language.
//
// Only what CanonicalCode cannot work out belongs here. It already folds "iw"
// onto "he", because BCP 47 retired the one in favour of the other and every
// language library will say so. It does not fold "no" onto "nb": those are a
// macrolanguage and one of its members, both current and neither a spelling of
// the other, so no amount of tag parsing turns one into the other. The pairing
// is a fact about this repo -- the curated Norwegian is written in Bokmal, and
// kb ships Bokmal -- and a fact about the repo has to be written down.
//
// The point of the table is to stop a scaffold creating an empty nb/ beside the
// curated no/, which would split one language across two directories: words
// looked up under one code, keyboard adjacency under the other.
var aliases = map[string]string{
	"no": "nb",
}

// Scaffold creates a directory of .lst files under root/languages for every
// language kb has a keyboard for and the tree does not carry yet, and returns
// the codes it created, sorted.
//
// kb is the authority on which languages exist for the same reason it seeds the
// Language table: a language the tool can reason about is one with a keyboard,
// whether or not anyone has curated a word list for it. Left to grow by hand
// the tree drifted to thirty directories against kb's hundred and ten, so a
// language would appear in --list languages with nowhere to put its data.
//
// Nothing is overwritten. An existing file is left exactly as it is, whatever
// it contains, so this is safe to run over a curated tree and safe to run
// twice; a curated directory missing one of the Relations gains that file and
// keeps the rest. The corpora themselves are left empty. They could have been
// filled from the layouts -- the characters a language's keyboards type look
// like its alphabet -- but only for some languages: kb's Japanese layout is
// Latin-only and would give Japanese the ASCII alphabet as its graphemes,
// English picks up seventy characters from the US-International variants, and
// Greek and Georgian come back carrying a full a-z. Machine-derived data that
// wrong is worse than none, because once written it is indistinguishable from
// the curated kind.
func Scaffold(root string) ([]string, error) {
	dir := filepath.Join(root, "languages")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("gen: scaffold %s: %w", dir, err)
	}

	have, err := covered(dir)
	if err != nil {
		return nil, err
	}

	var created []string
	for _, code := range kb.Languages() {
		if have[code] {
			continue
		}
		if err := os.MkdirAll(filepath.Join(dir, code), 0o750); err != nil {
			return nil, fmt.Errorf("gen: scaffold %s: %w", code, err)
		}
		created = append(created, code)
	}

	// Fill in missing files for every directory present, not only the ones
	// just made, so that a relation added to Relations reaches the curated
	// languages too.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		for _, r := range Relations {
			if err := touch(filepath.Join(dir, e.Name(), r.Name+".lst"), r.Doc); err != nil {
				return nil, err
			}
		}
	}

	sort.Strings(created)
	return created, nil
}

// Missing lists the languages Scaffold would create, sorted, without writing
// anything. A missing tree has every language missing rather than being an
// error, since that is what a scaffold of it would produce.
func Missing(root string) ([]string, error) {
	have, err := covered(filepath.Join(root, "languages"))
	if err != nil {
		return nil, err
	}

	var out []string
	for _, code := range kb.Languages() {
		if !have[code] {
			out = append(out, code)
		}
	}
	sort.Strings(out)
	return out, nil
}

// covered lists the language codes the tree already carries, including the
// kb code each directory stands in for. A directory named with a code kb does
// not use still covers the language kb knows by another name, and creating the
// other name would split the language in two.
func covered(dir string) (map[string]bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]bool{}, nil
		}
		return nil, err
	}

	have := make(map[string]bool, len(entries)*2)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		code := e.Name()
		have[code] = true
		have[CanonicalCode(code)] = true
		if alias, ok := aliases[code]; ok {
			have[alias] = true
		}
	}
	return have, nil
}

// touch writes an empty corpus with a one-line comment saying what belongs in
// it, and does nothing at all if the file is already there.
//
// O_EXCL rather than a Stat first: two builds racing on the same tree would
// both see the file missing and the second would truncate whatever the first
// had written. Here the loser is told the file exists, which is the outcome it
// wanted anyway.
func touch(path, doc string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o640)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("gen: scaffold %s: %w", path, err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "# %s\n", strings.TrimSpace(doc)); err != nil {
		return fmt.Errorf("gen: scaffold %s: %w", path, err)
	}
	return nil
}
