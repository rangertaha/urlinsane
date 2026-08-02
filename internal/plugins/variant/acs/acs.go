// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package acs is the adjacent character substitution algorithm.
package acs

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/kb"
)

// adjacentSubstitution replaces each character with a neighbouring key.
func adjacentSubstitution(token string, adj variant.Adjacency) (out []string) {
	for i, char := range token {
		for _, key := range adj(string(char)) {
			out = append(out, token[:i]+key+token[i+len(string(char)):])
		}
	}
	return
}

// Spec declares the algorithm.
func Spec(boards []*kb.Layout) variant.Spec {
	return variant.Spec{
		ID: "acs", Title: "Adjacent Character Substitution", Version: 1,
		Gen: variant.OverKeyboards(boards, adjacentSubstitution),
	}
}
