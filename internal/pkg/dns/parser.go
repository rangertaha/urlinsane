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
	var tlds []dataset.Suffix
	if dataset.DB != nil {
		dataset.DB.Find(&tlds)
	}

	for _, tld := range tlds {
		parser.suffixes = append(parser.suffixes, tld.Name)
	}

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
