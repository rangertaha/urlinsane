// Copyright 2024 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
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

	var subdomains []dataset.Prefix
	if dataset.DB != nil {
		dataset.DB.Find(&subdomains)
	}

	for _, sub := range subdomains {
		domains = append(domains, Join(sub.Name, d.Name, d.Suffix))
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
