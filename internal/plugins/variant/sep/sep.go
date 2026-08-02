// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package sep is the separator substitution algorithm.
package sep

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"strings"
)

var separators = []string{"-", "_", ".", ""}

// separatorSubstitution re-joins the name's word tokens with each alternate
// separator. A name with a single token has nothing to re-separate.
func separatorSubstitution(name string) []string {
	tokens := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_' || r == '.'
	})
	if len(tokens) < 2 {
		return nil
	}
	out := make([]string, 0, len(separators))
	for _, s := range separators {
		out = append(out, strings.Join(tokens, s))
	}
	return out
}

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "sep", Title: "Separator Substitution", Version: 1,
		Types: []string{variant.TypePackage, variant.TypeRepo, variant.TypeUsername},
		Whole: true,
		Gen:   separatorSubstitution,
	}
}
