// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package gi is the grapheme insertion algorithm.
package gi

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
func Spec(langs []variant.Language) variant.Spec {
	return variant.Spec{
		ID: "gi", Title: "Grapheme Insertion", Version: 1,
		Gen: variant.OverLanguages(langs, func(l variant.Language, name string) []string {
			return typo.GraphemeInsertion(name, l.Graphemes()...)
		}),
	}
}
