// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package xhs is the cross-language homophone algorithm.
//
// `hs` swaps a word for one that sounds like it *in the same language*: base
// and bass, to and too. This swaps for one that sounds like it to a speaker of
// a different language — youtube and yutup, boutique and boetiek and butik.
//
// The distinction matters because the two have different victims. A
// within-language homophone catches someone who heard the name; a
// cross-language one catches someone who heard it in a language whose
// orthography spells that sound differently, and who then writes down what they
// heard. That is most of the world for a name coined in English, and it is the
// case a curated English homophone list cannot reach by construction.
//
// The technique is measured in X-squatter (Valentim et al., ACM TOPS 2024,
// docs/papers/3663569.pdf): about 15% of cross-language sound-squatting
// candidates were found to have TLS certificates, against 7% for other
// squatting types — so these names are not merely registrable, they are
// registered and provisioned more often than the rest.
package xhs

import (
	"github.com/rangertaha/urlinsane/datasets"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// crossHomophones substitutes each group member for every other one it finds.
//
// The groups are snapshotted at construction rather than read per call. Every
// variant operator is cached against a digest of what it read, and a generator
// that reached for package state at Exec time would make its output depend on
// initialisation order that the cache assumes cannot change.
func crossHomophones(groups [][]string) variant.Generate {
	return func(name string) []string {
		if name == "" || len(groups) == 0 {
			return nil
		}
		return typo.HomophoneSwapping(name, groups...)
	}
}

// Spec declares the algorithm.
//
// It takes no language argument, which is the point: the groups are not a
// language's data. Passing the run's languages would suggest the output could
// be narrowed to them, and it cannot — a group is only in the file because it
// crosses a language boundary.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "xhs", Title: "Cross-language Homophone", Version: 1,
		Gen: crossHomophones(datasets.CrossHomophones()),
	}
}
