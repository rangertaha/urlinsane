// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package dns

import (
	"strings"

	"github.com/rangertaha/urlinsane/internal/dataset"
)

// Join assembles dotted-domain parts into a single name, trimming stray dots
// and whitespace from each part.
func Join(parts ...string) (domain string) {
	// clean parts
	for i := range parts {
		parts[i] = strings.Trim(parts[i], ".")
		parts[i] = strings.TrimSpace(parts[i])
	}
	domain = strings.Join(parts, ".")
	domain = strings.Trim(domain, ".")
	return
}

// Split parses a domain into its prefix (subdomains), registrable name, and
// public suffix.
func Split(domain string) (prefix, name, suffix string) {
	parser := New()
	d := parser.Parse(domain)
	return d.Prefix, d.Name, d.Suffix
}

// PermutatePrefix returns the domain with each known subdomain prefix prepended.
func PermutatePrefix(domain string) (domains []string) {
	parser := New()
	d := parser.Parse(domain)

	for _, sub := range dataset.Tokens("domains/prefix") {
		domains = append(domains, Join(sub, d.Name, d.Suffix))
	}

	return
}

// PermutateName returns the domain with its registrable name replaced by each
// of the given names (keeping prefix and suffix).
func PermutateName(domain string, names []string) (domains []string) {
	prefix, _, suffix := Split(domain)
	for _, n := range names {
		domains = append(domains, Join(prefix, n, suffix))
	}
	return
}

// PermutateSuffix returns the domain with its public suffix replaced by each
// known TLD/suffix.
func PermutateSuffix(domain string) (domains []string) {
	parser := New()
	d := parser.Parse(domain)

	for _, suf := range parser.suffixes {
		domains = append(domains, Join(d.Prefix, d.Name, suf))
	}
	return
}
