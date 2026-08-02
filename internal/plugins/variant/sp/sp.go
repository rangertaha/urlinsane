// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package sp is the singular pluralise algorithm.
package sp

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "sp", Title: "Singular Pluralise", Version: 1,
		Gen: typo.SingularPluralise,
	}
}
