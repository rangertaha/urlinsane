// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package hi is the hyphen insertion algorithm.
package hi

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "hi", Title: "Hyphen Insertion", Version: 1,
		Gen: typo.HyphenInsertion,
	}
}
