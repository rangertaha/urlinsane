// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package fsd is the delegated-subdomain algorithm.
//
// Some hosts hand out subdomains to anyone: duckdns.org, github.io,
// herokuapp.com, vercel.app, ddns.net. Putting a target's name under one of
// them produces paypal.duckdns.org or paypal.com.duckdns.org — a working
// lookalike that cost nothing, took a minute, and left no registration record.
//
// The distinction from `tld` is risk, not shape. `tld` swaps the suffix for any
// of the ~7,000 ICANN entries, where taking a name costs money and leaves a
// WHOIS trail; this uses only the ~3,000 private entries, where it does not.
// A scan that wants the cheap attacks first asks for `fsd` and gets the subset
// that needs no registrar at all, rather than sifting them out of tld's output
// afterwards.
//
// ail-typo-squatting calls the narrower version of this AddDynamicDns; the
// public suffix list's own private section is the general form, and it stays
// current without anyone curating a provider list.
package fsd

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// delegated puts the name under each provider, two ways.
//
// Both forms are real and they catch different readers. The core alone —
// paypal.duckdns.org — is what someone skimming the leftmost label sees. The
// whole name — paypal.com.duckdns.org — is what survives the check "does the
// address contain paypal.com", which is the check most people actually make.
func delegated(providers []string) variant.Generate {
	clean := make([]string, 0, len(providers))
	seen := make(map[string]bool, len(providers))
	for _, p := range providers {
		p = strings.Trim(strings.TrimSpace(p), ".")
		if p == "" || strings.HasPrefix(p, "#") || strings.ContainsAny(p, "*! ") || seen[p] {
			continue
		}
		seen[p] = true
		clean = append(clean, p)
	}

	return func(name string) []string {
		name = strings.Trim(strings.TrimSpace(name), ".")
		if name == "" {
			return nil
		}
		_, core, suffix := variant.SplitDomain(name)
		if core == "" || suffix == "" {
			return nil
		}
		out := make([]string, 0, len(clean)*2)
		for _, p := range clean {
			// A name already under this provider is not a variant of itself.
			if p == suffix {
				continue
			}
			out = append(out, core+"."+p, name+"."+p)
		}
		return out
	}
}

// Spec declares the algorithm.
func Spec(providers []string) variant.Spec {
	return variant.Spec{
		ID: "fsd", Title: "Delegated Subdomain", Version: 1,
		Types: []string{variant.TypeDomain},
		Whole: true,
		Gen:   delegated(providers),
	}
}
