// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package bf is the bit flipping algorithm.
package bf

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		// Bitsquatting: a single flipped bit in a transmitted name. It is
		// the one algorithm here that models hardware error rather than
		// human error, which is why it ignores keyboards and languages.
		ID: "bf", Title: "Bit Flipping", Version: 1,
		Gen: func(name string) []string { return typo.BitFlipping(name) },
	}
}
