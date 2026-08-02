// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package decompose

import "testing"

func TestDetect(t *testing.T) {
	for _, tc := range []struct {
		in, want, key, why string
	}{
		{in: "bob@example.com", want: TypeEmail, key: "bob@example.com",
			why: "entity.Classify called this a user because it saw an @, collapsing the address into its local part"},
		{in: "Bob@Example.COM", want: TypeEmail, key: "Bob@example.com",
			why: "the domain folds, the local part does not — RFC-significant even though providers ignore it"},

		{in: "github.com/acme/tool", want: TypeRepo, key: "github.com/acme/tool"},
		{in: "https://github.com/acme/tool", want: TypeRepo, key: "github.com/acme/tool"},
		{in: "git@github.com:acme/tool.git", want: TypeRepo, key: "github.com/acme/tool",
			why: "an scp-style remote and a browser URL name one repository and must reach one node"},
		{in: "https://github.com/acme/tool/tree/main", want: TypeRepo, key: "github.com/acme/tool",
			why: "a view of a repo is not part of its identity"},

		{in: "npm:lodash", want: TypePackage, key: "npm:lodash"},
		{in: "pypi:Foo.Bar", want: TypePackage, key: "pypi:foo-bar",
			why: "PEP 503 normalization is what makes a typosquat of Foo.Bar detectable"},

		{in: "example.com", want: TypeDomain, key: "example.com"},
		{in: "example.co.uk", want: TypeDomain, key: "example.co.uk"},
		{in: "acme.blogspot.com", want: TypeDomain, key: "acme.blogspot.com",
			why: "a private suffix is still a name a squatter can register under"},
		{in: "https://example.com", want: TypeDomain, key: "example.com",
			why: "a pasted URL is the common way to name a domain target"},
		{in: "https://example.com/about?x=1#top", want: TypeDomain, key: "example.com",
			why: "path, query and fragment carry no identity"},
		{in: "example.com:8080", want: TypeDomain, key: "example.com",
			why: "a port must not read as a package on registry example.com"},
		{in: "EXAMPLE.com.", want: TypeDomain, key: "example.com",
			why: "the root dot and case both fold"},

		{in: "lodash", want: TypeUsername, key: "lodash",
			why: "a bare word is a legal hostname, so only a real registry suffix separates a domain from a handle"},
		{in: "rangertaha", want: TypeUsername, key: "rangertaha"},
	} {
		t.Run(tc.in, func(t *testing.T) {
			got, key, err := DetectSeed(tc.in)
			if err != nil {
				t.Fatalf("DetectSeed(%q): %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("type = %q, want %q. %s", got, tc.want, tc.why)
			}
			if key != tc.key {
				t.Fatalf("key = %q, want %q. %s", key, tc.key, tc.why)
			}
		})
	}
}

func TestDetectRefusesWhatItCannotName(t *testing.T) {
	// Guessing would be worse than refusing: a mis-detected seed scans the
	// wrong namespace entirely and reports confidently about it.
	for _, bad := range []string{"", "   ", "bob@", "@example.com", "has space"} {
		if typ, err := Detect(bad); err == nil {
			t.Fatalf("Detect(%q) = %q, want an error", bad, typ)
		}
	}
}

func TestDetectIsIndependentOfScope(t *testing.T) {
	// §12: `typo username bob@example.com` and `typo bob@example.com` must read
	// the target identically. Detection takes only the target, so this is true
	// by construction — the test exists to keep it that way, because the moment
	// scope reaches detection the same string means different things on
	// different runs and the seed closure differs with it.
	typ, key, err := DetectSeed("bob@example.com")
	if err != nil {
		t.Fatalf("detect: %v", err)
	}
	if typ != TypeEmail || key != "bob@example.com" {
		t.Fatalf("detected %s:%s, want email:bob@example.com", typ, key)
	}
}

func TestDetectSeedAgreesWithTheCanonicalizer(t *testing.T) {
	// The failure this guards: detecting against one form of the string and
	// seeding with another. "https://example.com" detects as a domain, but
	// canonDomain refuses a delimiter outright — so a two-pass implementation
	// returns a type it cannot then canonicalize.
	for _, in := range []string{"https://example.com", "example.com:8080", "EXAMPLE.com."} {
		typ, key, err := DetectSeed(in)
		if err != nil {
			t.Fatalf("DetectSeed(%q): %v", in, err)
		}
		again, err := canonicalFor(typ, key)
		if err != nil {
			t.Fatalf("canonicalizing the detected key %q as %s failed: %v", key, typ, err)
		}
		if again != key {
			t.Fatalf("key %q is not canonical: recanonicalized to %q", key, again)
		}
	}
}
