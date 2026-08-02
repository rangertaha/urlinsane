// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package dhs is the dot hyphen substitution algorithm.
package dhs

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "dhs", Title: "Dot Hyphen Substitution", Version: 1,
		Gen: typo.DotHyphenSubstitution,
	}
}
