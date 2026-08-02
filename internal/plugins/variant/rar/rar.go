// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package rar is the repetition adjacent replacement algorithm.
package rar

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/kb"
)

// repetitionAdjacentReplacement replaces a doubled character with a doubled
// neighbour: "gooogle" from "google" is insertion, "gopogle" is this.
func repetitionAdjacentReplacement(token string, adj variant.Adjacency) (out []string) {
	r := []rune(token)
	for i := 0; i+1 < len(r); i++ {
		if r[i] != r[i+1] {
			continue
		}
		for _, key := range adj(string(r[i])) {
			out = append(out, string(r[:i])+key+key+string(r[i+2:]))
		}
	}
	return
}

// Spec declares the algorithm.
func Spec(boards []*kb.Layout) variant.Spec {
	return variant.Spec{
		ID: "rar", Title: "Repetition Adjacent Replacement", Version: 1,
		Gen: variant.OverKeyboards(boards, repetitionAdjacentReplacement),
	}
}
