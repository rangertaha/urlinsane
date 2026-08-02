// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package variant

import (
	"github.com/rangertaha/urlinsane/pkg/kb"
)

// Adjacency answers "what is next to this character" for one layout. Taking a
// function rather than a *kb.Layout keeps the three generators testable with a
// literal, and independent of how Adjacency is computed.
type Adjacency func(char string) []string

// overKeyboards folds a generator over every selected layout.
func OverKeyboards(boards []*kb.Layout, fn func(string, Adjacency) []string) Generate {
	// Snapshot at construction. Resolving layouts per call would make an
	// operator's output depend on registration order at call time, which the
	// scheduler's cache assumes cannot happen.
	adjs := make([]Adjacency, 0, len(boards))
	for _, b := range boards {
		if b != nil {
			adjs = append(adjs, b.Adjacent)
		}
	}
	return func(name string) []string {
		var out []string
		for _, adj := range adjs {
			out = append(out, fn(name, adj)...)
		}
		return out
	}
}

// overLanguages folds one per-language generator over every selected language.
func OverLanguages(langs []Language, fn func(Language, string) []string) Generate {
	// Snapshot for the same reason overKeyboards does.
	selected := append([]Language(nil), langs...)
	return func(name string) []string {
		var out []string
		for _, l := range selected {
			out = append(out, fn(l, name)...)
		}
		return out
	}
}
