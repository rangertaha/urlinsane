// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package gen loads a datasets tree from disk into the database.
//
// One pass over every .lst, into Vocabulary and Transition. There used to be an
// importer per relation — Synonyms, Antonyms, Homoglyphs and fifteen more —
// each filling its own table through a many2many association. Those tables
// recorded only that two words were related; the two that replaced them record
// how strongly, and one generic walk fills them for every dataset, including
// the ones no per-relation importer covered.
//
// It is a package rather than a //go:build ignore generator because it has
// enough logic to be worth testing, and an ignored file cannot be imported or
// tested. cmd/datasets is the thin shell over it.
//
// It is under internal/ because it is build-time only: it produces the
// database that ships, and nothing depending on this module should be able to
// import it. The datasets tree it reads is found by path at run time, so
// living outside datasets/ costs nothing.
//
// Build is the entry point — see build.go. This file is the walk itself.
package gen

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/rangertaha/urlinsane/internal/dataset"
)

// Progress, if set, is called after each file is imported. Nil is silent.
//
// A hook rather than a log line, because a library that printed to stdout would
// make its caller's output format its business — and one of those callers
// writes machine-readable output.
var Progress func(name, file string, words, transitions int)

// Extract reads a .lst into one string slice per line, collapsing runs of
// whitespace so that a line is its words.
// The error is returned rather than printed: an unreadable file used to be
// reported to stdout and then treated as empty, so a corrupt corpus imported as
// zero words and the import still claimed success.
func Extract(file string) (lines [][]string, err error) {
	readFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	space := regexp.MustCompile(`\s+`)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	// Word lists can be long and homoglyph lines carry many multi-byte runes;
	// the default 64K token limit is not obviously enough.
	fileScanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	for fileScanner.Scan() {
		line := strings.TrimSpace(space.ReplaceAllString(fileScanner.Text(), " "))
		if line == "" {
			continue
		}
		lines = append(lines, strings.Split(line, " "))
	}
	return lines, fileScanner.Err()
}

// languageID resolves a dataset directory name ("en") to a Language row id.
// Empty gives zero, which is what the list-shaped corpora carry: packages and
// surnames belong to no language.
func languageID(code string) uint {
	if code == "" {
		return 0
	}
	lng := dataset.Language{Code: code}
	dataset.DB.FirstOrCreate(&lng, dataset.Language{Code: code})
	return lng.ID
}

// datasetID resolves a dataset name ("synonym", "packages/npm") to a row id.
func datasetID(name string) uint {
	ds := dataset.Dataset{Name: name}
	dataset.DB.FirstOrCreate(&ds, dataset.Dataset{Name: name})
	return ds.ID
}

// One imports one .lst into the Vocabulary and Transition tables.
//
// Every word on a line joins the vocabulary whichever shape the line has. Two
// or more words are also associated with each other; a single word is not,
// because there is nothing to associate it with. That is the difference between
// the language datasets, whose lines are groups, and the list-shaped ones under
// domains/, entities/ and packages/, which are plain vocabularies.
//
// Idempotent: the rows for this (language, dataset) pair are cleared before
// reinserting, so a re-import replaces rather than duplicates.
func One(language, name, file string) error {
	lang, ds := languageID(language), datasetID(name)

	// Transitions name words by id, so clear the ones leaving this table's
	// vocabulary before the vocabulary itself is replaced -- otherwise the old
	// edges would dangle against ids that no longer exist.
	dataset.DB.Where("src IN (?)",
		dataset.DB.Model(&dataset.Vocabulary{}).Select("id").
			Where("language = ? AND dataset = ?", lang, ds),
	).Delete(&dataset.Transition{})
	dataset.DB.Where("language = ? AND dataset = ?", lang, ds).Delete(&dataset.Vocabulary{})

	words := map[string]bool{}
	// counts[src][dest] is how often the pair was seen. Counting first and
	// dividing at the end weights a word by how often it was associated, rather
	// than by whichever group happened to come last.
	counts := map[string]map[string]float64{}

	parsed, err := Extract(file)
	if err != nil {
		return err
	}
	for _, line := range parsed {
		uniq := make([]string, 0, len(line))
		seen := map[string]bool{}
		for _, w := range line {
			// A leading # marks a comment in the configuration-shaped datasets
			// and nothing else uses it.
			if w == "" || seen[w] || strings.HasPrefix(w, "#") {
				continue
			}
			seen[w] = true
			uniq = append(uniq, w)
		}
		if len(uniq) == 0 {
			continue
		}
		for _, w := range uniq {
			words[w] = true
		}
		if len(uniq) < 2 {
			continue
		}
		for _, src := range uniq {
			for _, dst := range uniq {
				if src == dst {
					continue
				}
				if counts[src] == nil {
					counts[src] = map[string]float64{}
				}
				counts[src][dst]++
			}
		}
	}

	// Vocabulary first, sorted so ids are assigned in a stable order and a
	// re-import of unchanged data produces the same rows.
	vocab := make([]string, 0, len(words))
	for w := range words {
		vocab = append(vocab, w)
	}
	sort.Strings(vocab)

	rows := make([]*dataset.Vocabulary, 0, len(vocab))
	for _, w := range vocab {
		rows = append(rows, &dataset.Vocabulary{Language: lang, Dataset: ds, Token: w})
	}
	if len(rows) > 0 {
		if err := dataset.DB.Create(&rows).Error; err != nil {
			return err
		}
	}

	// Transitions name words by their Vocabulary id, so the ids just assigned
	// are what the edges point at.
	id := make(map[string]uint, len(rows))
	for _, r := range rows {
		id[r.Token] = r.ID
	}

	srcs := make([]string, 0, len(counts))
	for s := range counts {
		srcs = append(srcs, s)
	}
	sort.Strings(srcs)

	var edges []*dataset.Transition
	for _, src := range srcs {
		var total float64
		for _, c := range counts[src] {
			total += c
		}
		dests := make([]string, 0, len(counts[src]))
		for d := range counts[src] {
			dests = append(dests, d)
		}
		sort.Strings(dests)
		for _, dst := range dests {
			p := 0.0
			if total > 0 {
				p = counts[src][dst] / total
			}
			edges = append(edges, &dataset.Transition{
				Src:         id[src],
				Dest:        id[dst],
				Probability: p,
			})
		}
	}
	if len(edges) > 0 {
		if err := dataset.DB.Create(&edges).Error; err != nil {
			return err
		}
	}

	if Progress != nil {
		Progress(name, file, len(rows), len(edges))
	}
	return nil
}

// All walks a datasets tree and imports every .lst into the two tables.
//
// languages/<code>/<relation>.lst is keyed by the code and the relation;
// anything else, <group>/<name>.lst, has no language and is named by its path,
// so packages/npm.lst becomes language 0, dataset "packages/npm".
func All(root string) error {
	groups, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if !g.IsDir() {
			continue
		}
		children, err := os.ReadDir(filepath.Join(root, g.Name()))
		if err != nil {
			return err
		}
		for _, c := range children {
			if c.IsDir() {
				files, err := os.ReadDir(filepath.Join(root, g.Name(), c.Name()))
				if err != nil {
					return err
				}
				for _, f := range files {
					if f.IsDir() || !strings.HasSuffix(f.Name(), ".lst") {
						continue
					}
					rel := strings.TrimSuffix(f.Name(), ".lst")
					p := filepath.Join(root, g.Name(), c.Name(), f.Name())
					if err := One(c.Name(), rel, p); err != nil {
						return err
					}
				}
				continue
			}
			if !strings.HasSuffix(c.Name(), ".lst") {
				continue
			}
			name := g.Name() + "/" + strings.TrimSuffix(c.Name(), ".lst")
			if err := One("", name, filepath.Join(root, g.Name(), c.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
