// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/fsd"
	"github.com/rangertaha/urlinsane/internal/plugins/variant/tos"
	"github.com/rangertaha/urlinsane/pkg/typo"
)

// Word order, which no character-level generator reaches: `cs` transposes two
// adjacent characters and turns shop-online into shpo-online, never into
// online-shop.
func TestTokenOrderSwapsWords(t *testing.T) {
	gen := tos.Spec().Gen

	for _, tc := range []struct{ in, want string }{
		{"shop-online", "online-shop"},
		{"node-fetch", "fetch-node"},
		{"2024example", "example2024"},   // the letter/digit boundary is a boundary
		{"münchen-shop", "shop-münchen"}, // runes, not bytes
	} {
		var found bool
		for _, g := range gen(tc.in) {
			if g == tc.want {
				found = true
			}
		}
		if !found {
			t.Errorf("tos(%q) = %q, missing %q", tc.in, gen(tc.in), tc.want)
		}
	}

	// A single token has no order to change.
	if got := gen("example"); len(got) != 0 {
		t.Errorf("tos(example) = %q, want nothing", got)
	}
	// And it never returns the input.
	for _, g := range gen("my_cool_lib") {
		if g == "my_cool_lib" {
			t.Error("tos returned the input as a variant")
		}
	}
}

// Separators stay where they are; only the words move. Moving the separators
// with the words would be a different edit, and one dhs already owns.
func TestTokenOrderKeepsSeparatorPositions(t *testing.T) {
	parts, seps := typo.SplitTokens("shop-online24")
	if len(parts) != 3 || len(seps) != 2 {
		t.Fatalf("SplitTokens = %q / %q, want 3 parts and 2 separators", parts, seps)
	}
	if seps[0] != "-" || seps[1] != "" {
		t.Errorf("separators = %q, want [- \"\"] — the digit boundary has no character", seps)
	}
	// A leading or trailing separator is not a token boundary.
	for _, s := range []string{"-lead", "trail-", "-", ""} {
		if p, _ := typo.SplitTokens(s); len(p) != 0 {
			t.Errorf("SplitTokens(%q) = %q, want nothing to reorder", s, p)
		}
	}
}

// fsd puts the name under hosts that give subdomains away, in both the forms
// that fool a reader: the bare core, and the whole name.
func TestDelegatedSubdomainUsesBothForms(t *testing.T) {
	gen := fsd.Spec([]string{"duckdns.org", "github.io", "com"}).Gen
	got := gen("paypal.com")

	for _, want := range []string{
		"paypal.duckdns.org",     // what someone skimming the first label sees
		"paypal.com.duckdns.org", // what survives "does it contain paypal.com"
		"paypal.github.io",
	} {
		var found bool
		for _, g := range got {
			if g == want {
				found = true
			}
		}
		if !found {
			t.Errorf("fsd(paypal.com) = %q, missing %q", got, want)
		}
	}

	// A provider equal to the name's own suffix would produce paypal.com.com.
	for _, g := range got {
		if strings.HasSuffix(g, ".com.com") || g == "paypal.com" {
			t.Errorf("fsd produced %q", g)
		}
	}
}

// Both new domain algorithms take the whole key and bind to domains only.
func TestNewAlgorithmsScoping(t *testing.T) {
	f := fsd.Spec([]string{"duckdns.org"})
	if len(f.Types) != 1 || f.Types[0] != variant.TypeDomain || !f.Whole {
		t.Errorf("fsd scoping wrong: types=%v whole=%v", f.Types, f.Whole)
	}
	if got := f.Gen("npm:lodash"); len(got) != 0 {
		t.Errorf("fsd(npm:lodash) = %q, want nothing", got)
	}
	// tos is deliberately unrestricted: a hyphenated domain, a scoped package
	// and a handle all carry word order.
	if len(tos.Spec().Types) != 0 {
		t.Errorf("tos restricts types to %v; word order is not a domain property", tos.Spec().Types)
	}
}
