// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package tld is the wrong tld algorithm.
package tld

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// tldSwap replaces the public suffix with every other known suffix, keeping any
// subdomain prefix. This is the one algorithm allowed to change the registry a
// name lives in; every other algorithm preserves the suffix precisely so this
// one owns that axis.
func tldSwap(suffixes []string) variant.Generate {
	return func(name string) []string {
		prefix, core, _ := variant.SplitDomain(name)
		if core == "" {
			return nil
		}
		out := make([]string, 0, len(suffixes))
		for _, s := range suffixes {
			out = append(out, variant.JoinDomain(prefix, core, s))
		}
		return out
	}
}

// Spec declares the algorithm.
func Spec(suffixes []string) variant.Spec {
	return variant.Spec{
		ID: "tld", Title: "Wrong TLD", Version: 1,
		Types: []string{variant.TypeDomain},
		Whole: true,
		Gen:   tldSwap(suffixes),
	}
}
