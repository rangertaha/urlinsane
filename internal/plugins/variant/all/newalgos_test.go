// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/sld"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/tli"
)

var testSuffixes = []string{
	"com", "net", "org", "de", "br", "it",
	"co.uk", "org.uk", "ac.uk", "gov.uk",
	"com.au", "net.au",
}

func has(got []string, want string) bool {
	for _, g := range got {
		if g == want {
			return true
		}
	}
	return false
}

// sld keeps the country and changes the category, which is the variation `tld`
// cannot produce: it replaces the whole suffix and so changes the country.
func TestWrongSLDSwapsTheSecondLevelOnly(t *testing.T) {
	gen := sld.Spec(testSuffixes).Gen
	got := gen("bbc.co.uk")

	for _, want := range []string{"bbc.org.uk", "bbc.ac.uk", "bbc.gov.uk"} {
		if !has(got, want) {
			t.Errorf("sld(bbc.co.uk) = %q, missing %q", got, want)
		}
	}
	if has(got, "bbc.co.uk") {
		t.Error("sld returned the input as a variant")
	}
	// Never crosses to another country's suffixes.
	for _, g := range got {
		if !strings.HasSuffix(g, ".uk") {
			t.Errorf("sld(bbc.co.uk) produced %q, which is not under .uk — that is tld's axis", g)
		}
	}
}

// A name whose suffix is a single label has no second level to be wrong about.
func TestWrongSLDIgnoresSingleLabelSuffixes(t *testing.T) {
	gen := sld.Spec(testSuffixes).Gen
	if got := gen("example.com"); len(got) != 0 {
		t.Errorf("sld(example.com) = %q, want nothing — .com delegates no second level", got)
	}
}

// tli keeps the whole target name intact and makes it somebody else's
// subdomain, so there is no misspelling for a reader to catch.
func TestTLDInsertionKeepsTheWholeNameIntact(t *testing.T) {
	gen := tli.Spec(testSuffixes).Gen
	got := gen("example.com")

	for _, want := range []string{"example.com.br", "example.com.it", "example.com.de"} {
		if !has(got, want) {
			t.Errorf("tli(example.com) = %q, missing %q", got, want)
		}
	}
	for _, g := range got {
		if !strings.HasPrefix(g, "example.com.") {
			t.Errorf("tli produced %q, which does not contain the target name in full", g)
		}
	}
	// example.com.com is a repetition, not a level.
	if has(got, "example.com.com") {
		t.Error("tli appended the suffix the name already ends in")
	}
}

// Both bind to domains only, and both take the whole key rather than a part —
// they operate on the suffix, which is not in the varyable portion of a name.
func TestNewDomainAlgorithmsAreScopedToDomains(t *testing.T) {
	for _, s := range []variant.Spec{sld.Spec(testSuffixes), tli.Spec(testSuffixes)} {
		if len(s.Types) != 1 || s.Types[0] != variant.TypeDomain {
			t.Errorf("%s applies to %v, want [domain]", s.ID, s.Types)
		}
		if !s.Whole {
			t.Errorf("%s must take the whole key: it edits the suffix", s.ID)
		}
		if got := s.Gen("npm:lodash"); len(got) != 0 {
			t.Errorf("%s(npm:lodash) = %q, want nothing", s.ID, got)
		}
	}
}
