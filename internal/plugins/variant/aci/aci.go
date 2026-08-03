// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package aci is the adjacent character insertion algorithm.
package aci

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/kb"
)

// adjacentInsertion inserts a neighbouring key beside each character, on both
// sides: "google" gives "gooogle" and "goopgle" from the neighbours of "g".
func adjacentInsertion(token string, adj variant.Adjacency) (out []string) {
	for i, char := range token {
		for _, key := range adj(string(char)) {
			out = append(out,
				token[:i]+key+string(char)+token[i+len(string(char)):],
				token[:i]+string(char)+key+token[i+len(string(char)):])
		}
	}
	// Inserting a neighbour before a character and after the one behind it are
	// the same edit, so every interior position is reachable twice.
	return variant.Clean(token, out)
}

// Spec declares the algorithm.
func Spec(boards []*kb.Layout) variant.Spec {
	return variant.Spec{
		ID: "aci", Title: "Adjacent Character Insertion", Version: 1,
		Gen: variant.OverKeyboards(boards, adjacentInsertion),
	}
}
