// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package dns

import (
	"index/suffixarray"
	"strings"

	"github.com/rangertaha/urlinsane/internal/dataset"
)

// Parser splits domains into prefix/name/suffix using the public suffix list
// loaded from the dataset.
type Parser struct {
	sa       *suffixarray.Index
	suffixes []string
}

// Domain is a parsed domain: its subdomain Prefix, registrable Name, and public
// Suffix.
type Domain struct {
	Prefix string
	Name   string
	Suffix string
}

// New builds a Parser, loading the public-suffix list from the dataset.
func New() (parser Parser) {
	parser.suffixes = dataset.Tokens("domains/suffix")

	data := []byte("\x00" + strings.Join(parser.suffixes, "\x00") + "\x00")
	parser.sa = suffixarray.New(data)
	return
}

// Offset returns the index, within the dot-separated parts, of the registrable
// name (the label immediately left of the public suffix).
func (p *Parser) Offset(parts []string) int {
	counter := 2
	for counter > 0 {
		start_point := len(parts) - counter
		if start_point < 0 {
			return 0
		}
		tld_parts := strings.Join(parts[len(parts)-counter:], ".")

		indicies := p.sa.Lookup([]byte(tld_parts), -1)
		if len(indicies) > 0 {
			offset := (len(parts) - (counter + 1))
			if offset >= 0 {
				return offset
			}
		}
		counter--
	}

	return 0
}

// Parse splits a domain into its prefix, registrable name, and suffix.
func (p *Parser) Parse(domain string) Domain {
	parts := strings.Split(domain, ".")
	offset := p.Offset(parts)
	return Domain{
		Prefix: strings.Join(parts[:offset], "."),
		Name:   parts[offset],
		Suffix: strings.Join(parts[offset+1:], "."),
	}
}

// GetDomain returns the registrable name label of the domain.
func (p *Parser) GetDomain(domain string) string {
	parts := strings.Split(domain, ".")
	offset := p.Offset(parts)
	return parts[offset]
}

// GetPrefix returns the subdomain prefix of the domain (may be empty).
func (p *Parser) GetPrefix(domain string) string {
	parts := strings.Split(domain, ".")
	offset := p.Offset(parts)
	return strings.Join(parts[:offset], ".")
}

// GetFQDN returns the registrable domain (name plus suffix, without subdomains).
func (p *Parser) GetFQDN(domain string) string {
	parts := strings.Split(domain, ".")
	offset := p.Offset(parts)
	return strings.Join(parts[offset:], ".")
}

// GetSuffix returns the public suffix (TLD) of the domain.
func (p *Parser) GetSuffix(domain string) string {
	parts := strings.Split(domain, ".")
	offset := p.Offset(parts)
	return strings.Join(parts[offset+1:], ".")
}
