// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package sld is the wrong second-level domain algorithm.
//
// Under a country-code TLD that delegates a second level, the registrable name
// sits one label deeper than it does under .com: bbc.co.uk, not bbc.uk. Those
// second-level labels are a small, fixed, semantically loaded set — co, org,
// ac, gov, net — and swapping one for another produces a name that reads as the
// same organisation in a different category, which is a more plausible lure
// than any character-level edit of the same length.
//
// `tld` cannot produce these. It replaces the whole public suffix, so
// example.co.uk becomes example.com or example.de — a different country. This
// keeps the country and changes the category, which is the variation a reader
// checking "is this the UK site?" will not catch.
//
// Both URLCrazy (wrong_sld) and ail-typo-squatting (WrongSld) implement it; it
// was the one domain-shaped generator in either that urlinsane had no answer
// for.
package sld

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// siblingSuffixes indexes the suffix list by its final label, keeping only the
// multi-label entries.
//
// Built once at construction rather than scanned per call: the list is the
// whole public suffix list, and an operator that walked it for every candidate
// would turn a string edit into a linear scan of ten thousand entries.
func siblingSuffixes(suffixes []string) map[string][]string {
	byTLD := make(map[string][]string)
	for _, s := range suffixes {
		s = strings.Trim(strings.TrimSpace(s), ".")
		// Single-label suffixes have no second level to be wrong about, and a
		// wildcard entry is not a name anyone registers under.
		if s == "" || strings.Contains(s, "*") || strings.Count(s, ".") != 1 {
			continue
		}
		tld := s[strings.LastIndex(s, ".")+1:]
		byTLD[tld] = append(byTLD[tld], s)
	}
	return byTLD
}

// wrongSLD swaps the second-level label for every sibling under the same TLD.
func wrongSLD(suffixes []string) variant.Generate {
	byTLD := siblingSuffixes(suffixes)

	return func(name string) []string {
		prefix, core, suffix := variant.SplitDomain(name)
		// Only a name whose own suffix is two labels has a second level to
		// change. example.com has none, and inventing one would be `tli`'s job.
		if core == "" || strings.Count(suffix, ".") != 1 {
			return nil
		}
		tld := suffix[strings.LastIndex(suffix, ".")+1:]

		siblings := byTLD[tld]
		out := make([]string, 0, len(siblings))
		for _, s := range siblings {
			if s == suffix {
				continue
			}
			out = append(out, variant.JoinDomain(prefix, core, s))
		}
		return out
	}
}

// Spec declares the algorithm.
func Spec(suffixes []string) variant.Spec {
	return variant.Spec{
		ID: "sld", Title: "Wrong Second-Level Domain", Version: 1,
		Types: []string{variant.TypeDomain},
		Whole: true,
		Gen:   wrongSLD(suffixes),
	}
}
