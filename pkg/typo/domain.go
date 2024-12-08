// Copyright (C) 2024 Rangertaha
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
package typo

import (
	"strings"
)

// TopLevelDomain typos occur when the top-level domain (TLD) of a domain name is replaced with a common or similar TLD.
// This type of typo exploits the familiarity of well-known TLDs, creating domains that look almost identical to the original,
// but with a small variation in the TLD. This can make the altered domain appear legitimate at first glance.
// For example, "www.example.com" could be mistyped as "www.example.net", "www.example.org", or "www.example.co",
// where the TLD is swapped for another commonly used extension, potentially leading to confusion or typosquatting.
func TopLevelDomain(tld string, tlds ...string) (records []string) {
	for _, suffix := range tlds {
		if len(strings.Split(suffix, ".")) == 1 {
			records = append(records, suffix)
		}
	}
	return
}

// SecondLevelDomain
// SecondLevelDomain typos occur when the second-level domain (SLD) of a domain name is replaced with another common or similar SLD,
// while keeping the original TLD intact. These types of typos are often subtle, involving the substitution of one part of the domain
// for another, resulting in a domain that appears similar but leads to a different destination. For example, "www.example.com"
// could be mistyped as "www.examples.com" or "www.exampel.com", where the second-level domain is altered but the TLD remains unchanged,
// creating a deceptive resemblance to the original domain.
func SecondLevelDomain(tld string, tlds ...string) (records []string) {
	for _, suffix := range tlds {
		if len(strings.Split(suffix, ".")) == 2 {
			records = append(records, suffix)
		}
	}
	return

}

// ThirdLevelDomain typos occur when the third-level domain (TLD3) of a domain name is replaced with another common or similar TLD3,
// while keeping the original TLD3 intact. These types of typos are often subtle, involving the substitution of one part of the domain
// for another, resulting in a domain that appears similar but leads to a different destination. For example, "www.example.com"
// could be mistyped as "www.examples.com" or "www.exampel.com", where the second-level domain is altered but the TLD remains unchanged,
// creating a deceptive resemblance to the original domain.
func ThirdLevelDomain(tld string, tlds ...string) (records []string) {
	for _, suffix := range tlds {
		if len(strings.Split(suffix, ".")) == 3 {
			records = append(records, suffix)
		}
	}
	return
}
