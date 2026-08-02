// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package co is the character omission algorithm.
package co

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "co", Title: "Character Omission", Version: 1,
		Gen: typo.CharacterOmission,
	}
}
