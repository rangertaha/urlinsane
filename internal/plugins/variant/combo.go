// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package variant

import (
	"sort"

	"github.com/rangertaha/urlinsane/datasets"
)

// comboJoins are the four ways a keyword attaches to a name, and they are the
// four shapes real combosquatting takes: hyphenated or run together, before or
// after. "paypal-login", "paypallogin", "login-paypal", "loginpaypal" are all
// registered patterns.
//
// Dotted forms — "login.paypal.com" — are deliberately absent. That is a
// subdomain, which the si operator already generates, and emitting it here
// would attribute one finding to two algorithms.
const comboJoinsPerKeyword = 4

// DefaultComboLanguage is the language whose vocabulary cb uses when the caller
// does not choose one.
//
// Unlike every other language-driven algorithm, cb must not default to the full
// registered set. Two reasons, and the second is the load-bearing one:
//
// Volume. The spelling algorithms are filtered by the name — cm fires only
// where the name contains a misspellable substring — so their cost scales with
// the name, not the dataset. cb is a pure cross-product of vocabulary and join
// forms, so it scales with the dataset alone: 263 English keywords give ~1k
// variants per name, all thirty languages give ~26k. Nothing else here emits
// more than the ~1.5k of a TLD swap.
//
// Meaning. Bolting Japanese keywords onto an English brand is not additional
// coverage, it is noise. A combosquat works because the keyword reads as
// plausible to the target's users, which makes the vocabulary a property of the
// audience rather than of the name. Callers scanning a non-English brand should
// pass Options.Combos explicitly; there is no detection here to guess it for
// them.
const DefaultComboLanguage = "en"

// DefaultComboKeywords is the vocabulary cb uses when Options.Combos is nil.
func DefaultComboKeywords() []string {
	return datasets.SynonymWords(DefaultComboLanguage)
}

// ComboKeywords is the combosquatting vocabulary for a set of languages, in
// sorted order with duplicates removed. Languages that ship no synonym set
// contribute nothing rather than erroring: a language plugin without a dataset
// is a gap in coverage, not a broken run.
//
// This is the explicit multi-language form. It is not what the default uses —
// see DefaultComboLanguage.
func ComboKeywords(langs []Language) []string {
	seen := make(map[string]bool)
	var out []string
	for _, l := range langs {
		for _, w := range datasets.SynonymWords(l.Code()) {
			if w == "" || seen[w] {
				continue
			}
			seen[w] = true
			out = append(out, w)
		}
	}
	sort.Strings(out)
	return out
}
