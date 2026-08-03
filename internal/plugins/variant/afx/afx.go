// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package afx is the affix squatting algorithm.
package afx

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// affixes are the common ecosystem/role brackets seen in real registry
// squatting. Ported verbatim from the afx plugin.
var (
	affixPrefixes = []string{"py", "py-", "python-", "node-", "js-", "go-", "lib", "lib-", "the-", "get-"}
	affixSuffixes = []string{"2", "3", "js", "-js", "py", "-py", "-python", "-cli", "-dev",
		"-core", "-utils", "-api", "-sdk", "-lib", ".js", "-ng", "-master", "-official"}
)

// affixSquatting brackets the name with each common ecosystem affix.
func affixSquatting(name string) []string {
	out := make([]string, 0, len(affixPrefixes)+len(affixSuffixes))
	for _, p := range affixPrefixes {
		out = append(out, p+name)
	}
	for _, s := range affixSuffixes {
		out = append(out, name+s)
	}
	// An empty name would return the bare affixes, which are not variations of
	// anything.
	return variant.Clean(name, out)
}

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		// Not a whole-key algorithm, though it reads like one. The whole key
		// of a package is "npm:lodash", and prefixing that gives "py-npm:lodash"
		// — a package in a registry nobody has. The core is the name with its
		// registry or forge stripped, which is what an affix belongs in front
		// of; the split puts the qualifier back afterwards.
		ID: "afx", Title: "Affix Squatting", Version: 1,
		Types: []string{variant.TypePackage, variant.TypeRepo, variant.TypeUsername},
		Whole: false,
		Gen:   affixSquatting,
	}
}
