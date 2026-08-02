// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package domain splits a domain into its registrable name and public suffix.
package domain

import (
	"strings"

	"golang.org/x/net/publicsuffix"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
)

// Domain pulls a domain's registry suffix out as its own node. The suffix is
// worth converging on: a hundred variants under one unusual TLD is a campaign,
// and that is invisible if the TLD only ever exists as a substring of a key.
func Domain() graph.Operator {
	return decompose.New("decompose.domain", 1, decompose.TypeDomain,
		graph.Effects{
			Nodes: []string{decompose.TypeTLD},
			Rels:  []string{decompose.RelTLDOf},
		},
		splitDomain)
}

// splitDomain finds the suffix under which the name was registered. The public
// suffix list is used rather than the last label because "example.co.uk" is
// registered under co.uk, not under uk, and keying on uk would group it with
// names that share nothing with it.
//
// The key is already lowercase punycode, which is what the list expects.
func splitDomain(key string) graph.Delta {
	suffix, icann := publicsuffix.PublicSuffix(key)
	if !icann {
		// A private suffix (blogspot.com) or an unrecognised one. Neither is a
		// registry, and admitting them as tld nodes would put a squatter's own
		// hosting domain on equal footing with a real TLD; fall back to the
		// last label, which always is one.
		if i := strings.LastIndex(suffix, "."); i >= 0 {
			suffix = suffix[i+1:]
		}
	}
	// A bare TLD, or a single-label host: there is no suffix to separate out.
	if suffix == "" || suffix == key {
		return graph.Delta{}
	}
	self := graph.NodeRef{Type: decompose.TypeDomain, Key: key}
	tld := graph.NodeRef{Type: decompose.TypeTLD, Key: suffix}
	return graph.Delta{
		Nodes: []graph.NodeRef{tld},
		Edges: []graph.EdgeRef{{From: self, Rel: decompose.RelTLDOf, To: tld}},
	}
}
