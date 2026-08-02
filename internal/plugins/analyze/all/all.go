// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package all composes the shipped analyzers, in run order.
//
// Separate from the analyze library so the analyzers can import that library
// without importing each other through it.
package all

import (
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze/campaign"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze/depconfusion"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze/scoring"
)

// All returns the shipped analyzers, in run order.
func All() []graph.Analyzer {
	return []graph.Analyzer{
		campaign.Campaign{},
		scoring.Scoring{},
		depconfusion.DepConfusion{},
	}
}
