// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package cr is the character repetition algorithm.
package cr

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "cr", Title: "Character Repetition", Version: 1,
		Gen: typo.CharacterRepetition,
	}
}
