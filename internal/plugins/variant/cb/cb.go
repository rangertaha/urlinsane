// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package cb is the combo squatting algorithm.
package cb

import (
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"sort"
)

// JoinsPerKeyword is how many join forms each combo keyword yields.
const JoinsPerKeyword = 4

// comboSquatting brackets the name with each keyword, in each join form.
//
// The vocabulary is snapshotted and deduped at construction rather than read
// per call: resolving it at Exec time would make the operator's output depend
// on plugin registration order, which the scheduler's cache assumes cannot
// happen.
func ComboSquatting(keywords []string) variant.Generate {
	kws := make([]string, 0, len(keywords))
	seen := make(map[string]bool, len(keywords))
	for _, k := range keywords {
		if k == "" || seen[k] {
			continue
		}
		seen[k] = true
		kws = append(kws, k)
	}
	sort.Strings(kws)

	return func(name string) []string {
		if name == "" {
			return nil
		}
		out := make([]string, 0, len(kws)*JoinsPerKeyword)
		for _, k := range kws {
			if k == name {
				// "paypal-paypal" is a doubling, not a combination, and the
				// repetition algorithms own that shape.
				continue
			}
			out = append(out,
				name+"-"+k,
				name+k,
				k+"-"+name,
				k+name,
			)
		}
		return out
	}
}

// Spec declares the algorithm.
func Spec(keywords []string) variant.Spec {
	return variant.Spec{
		ID: "cb", Title: "Combo Squatting", Version: 1,
		Gen: ComboSquatting(keywords),
	}
}
