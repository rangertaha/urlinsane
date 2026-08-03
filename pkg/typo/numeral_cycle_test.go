// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package typo

import (
	"sort"
	"strings"
	"testing"
)

// cliqueNumerals is the shape dataset.Lang.Numerals actually returns.
//
// A numeral line is imported as a clique of transitions, so every token on it
// becomes a key and the words map back to the digit as readily as the digit
// maps to the words. The fixture the other numeral tests use is keyed by digit
// alone, which cannot cycle -- which is why those tests passed while
// `urlinsane typo -a cns` crashed on the shipped English numerals.
var cliqueNumerals = map[string][]string{
	"1":      {"first", "one"},
	"one":    {"1", "first"},
	"first":  {"1", "one"},
	"2":      {"second", "two"},
	"two":    {"2", "second"},
	"second": {"2", "two"},
}

// A word that maps back to its own digit is a two-cycle: "1" becomes "first",
// "first" becomes "1", and a seen set rebuilt in each frame never notices. The
// walk ran until the goroutine stack limit and took the process down -- a
// fatal error, not a recoverable panic, so this test kills the whole binary if
// it regresses rather than merely failing.
func TestNumeralSwapTerminatesOnCliqueData(t *testing.T) {
	for _, name := range []string{"shop1", "one", "first", "1and2"} {
		got := CardinalSwap(name, cliqueNumerals)
		if len(got) == 0 {
			t.Errorf("CardinalSwap(%q) returned nothing", name)
		}
		for _, g := range got {
			if g == name {
				t.Errorf("CardinalSwap(%q) returned its own input", name)
			}
		}
		if got := OrdinalSwap(name, cliqueNumerals); len(got) == 0 {
			t.Errorf("OrdinalSwap(%q) returned nothing", name)
		}
	}
}

// Ranging a map is randomised, and the order variants come back in reaches
// admission order in the engine and so the content address of a scan.
func TestNumeralSwapIsDeterministic(t *testing.T) {
	first := strings.Join(CardinalSwap("shop1", cliqueNumerals), ",")
	for i := 0; i < 20; i++ {
		if again := strings.Join(CardinalSwap("shop1", cliqueNumerals), ","); again != first {
			t.Fatalf("run %d gave %q, first gave %q; order is not stable", i, again, first)
		}
	}
}

// Whatever the shape of the data, a variant is produced once.
func TestNumeralSwapDoesNotRepeatItself(t *testing.T) {
	got := CardinalSwap("1and2", cliqueNumerals)
	sorted := append([]string(nil), got...)
	sort.Strings(sorted)
	for i := 1; i < len(sorted); i++ {
		if sorted[i] == sorted[i-1] {
			t.Errorf("%q appears more than once in %q", sorted[i], got)
		}
	}
}
