// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package vs is the vowel swapping algorithm.
package vs

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
func Spec(langs []variant.Language) variant.Spec {
	return variant.Spec{
		ID: "vs", Title: "Vowel Swapping", Version: 1,
		Gen: variant.OverLanguages(langs, func(l variant.Language, name string) []string {
			return typo.VowelSwapping(name, l.Vowels()...)
		}),
	}
}
