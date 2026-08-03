// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package tli is the TLD insertion algorithm.
//
// The whole target name, suffix included, becomes the leading labels of a name
// under some other suffix: example.com -> example.com.br. Nothing is misspelled
// and nothing is substituted, so there is no character-level edit for a reader
// to catch — the name they are checking for is present, in full, spelled
// correctly, and it is a subdomain of somebody else's registration.
//
// It defeats the check most people actually perform, which is "does the address
// contain the name I expect". It is also what a truncating address bar shows
// first, which is why the technique is called levelsquatting in the literature.
//
// ail-typo-squatting calls this AddTld. It differs from `si` (subdomain
// insertion), which puts a new label in front of a name the target still owns;
// here the target owns nothing.
package tli

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// tldInsertion appends each public suffix to the whole name.
func tldInsertion(suffixes []string) variant.Generate {
	clean := make([]string, 0, len(suffixes))
	for _, s := range suffixes {
		s = strings.Trim(strings.TrimSpace(s), ".")
		if s == "" || strings.Contains(s, "*") {
			continue
		}
		clean = append(clean, s)
	}

	return func(name string) []string {
		name = strings.Trim(strings.TrimSpace(name), ".")
		if name == "" {
			return nil
		}
		_, core, suffix := variant.SplitDomain(name)
		if core == "" || suffix == "" {
			return nil
		}
		out := make([]string, 0, len(clean))
		for _, s := range clean {
			// Appending the suffix the name already ends in yields
			// example.com.com, which is a repetition rather than a level.
			if s == suffix {
				continue
			}
			out = append(out, name+"."+s)
		}
		return out
	}
}

// Spec declares the algorithm.
func Spec(suffixes []string) variant.Spec {
	return variant.Spec{
		ID: "tli", Title: "TLD Insertion", Version: 1,
		Types: []string{variant.TypeDomain},
		Whole: true,
		Gen:   tldInsertion(suffixes),
	}
}
