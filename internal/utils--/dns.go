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
package pkg

// import (
// 	"index/suffixarray"
// 	"strings"

// 	"github.com/rangertaha/urlinsane/internal/dataset"
// )

// type Parser struct {
// 	sa *suffixarray.Index
// }

// type xDomain struct {
// 	Prefix string
// 	Name   string
// 	Suffix string
// }

// type Domain struct {
// 	Subdomain string
// 	Domain    string
// 	TLD       string
// }

// func NewDomainParser() Parser {
// 	var tlds dataset.Suffix
// 	dataset.DB.Find(&tlds)

// 	// data, err := ioutil.ReadFile("/tmp/.tlds")
// 	// if err != nil {
// 	// 	data, _ = download()
// 	// 	ioutil.WriteFile("/tmp/.tlds", data, 0644)
// 	// }

// 	// tlds := strings.Split(string(data), "\n")

// 	sa := CreateTLDIndex(tlds)
// 	return Parser{
// 		sa: sa,
// 	}
// }

// func CreateTLDIndex(tlds []string) *suffixarray.Index {
// 	data := []byte("\x00" + strings.Join(tlds, "\x00") + "\x00")
// 	sa := suffixarray.New(data)
// 	return sa
// }

// func (p *Parser) FindTldOffset(domain_parts []string) int {
// 	counter := 2
// 	for counter > 0 {
// 		start_point := len(domain_parts) - counter
// 		if start_point < 0 {
// 			return 0
// 		}
// 		tld_parts := strings.Join(domain_parts[len(domain_parts)-counter:], ".")

// 		indicies := p.sa.Lookup([]byte(tld_parts), -1)
// 		if len(indicies) > 0 {
// 			offset := (len(domain_parts) - (counter + 1))
// 			if offset >= 0 {
// 				return offset
// 			}
// 		}
// 		counter--
// 	}

// 	return 0

// }

// func (p *Parser) ParseDomain(domain string) Domain {
// 	domain_parts := strings.Split(domain, ".")
// 	offset := p.FindTldOffset(domain_parts)
// 	return Domain{
// 		Subdomain: strings.Join(domain_parts[:offset], "."),
// 		Domain:    domain_parts[offset],
// 		TLD:       strings.Join(domain_parts[offset+1:], "."),
// 	}
// }

// func (p *Parser) GetDomain(domain string) string {
// 	domain_parts := strings.Split(domain, ".")
// 	offset := p.FindTldOffset(domain_parts)
// 	return domain_parts[offset]

// }

// func (p *Parser) GetSubdomain(domain string) string {
// 	domain_parts := strings.Split(domain, ".")
// 	offset := p.FindTldOffset(domain_parts)
// 	return strings.Join(domain_parts[:offset], ".")
// }

// func (p *Parser) GetFQDN(domain string) string {
// 	domain_parts := strings.Split(domain, ".")
// 	offset := p.FindTldOffset(domain_parts)
// 	return strings.Join(domain_parts[offset:], ".")
// }

// func (p *Parser) GetTld(domain string) string {
// 	domain_parts := strings.Split(domain, ".")
// 	offset := p.FindTldOffset(domain_parts)
// 	return strings.Join(domain_parts[offset+1:], ".")
// }
