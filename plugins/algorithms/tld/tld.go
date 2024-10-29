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
package tld

// // wrongTopLevelDomain for example, www.google.co.nz becomes www.google.co.ns
// // and www.google.com becomes www.google.org. uses the 19 most common top level
// // domains.
// func wrongTopLevelDomainFunc(tc Result) (results []Result) {
// 	labels := strings.Split(tc.Original.Suffix, ".")
// 	length := len(labels)
// 	for _, suffix := range datasets.TLD {
// 		suffixLen := len(strings.Split(suffix, "."))
// 		if length == suffixLen && length == 1 {
// 			if suffix != tc.Original.Suffix {
// 				dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 			}
// 		}
// 	}
// 	return
// }
// http://www.iana.org/domains/root/db
// https://en.wikipedia.org/wiki/Country_code_top-level_domain

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "tld"
// const (
// 	CODE        = ""
// 	NAME        = ""
// 	DESCRIPTION = ""
// )


type Algo struct {
	types []string
}

func (n *Algo) Id() string {
	return CODE
}
func (n *Algo) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}
func (n *Algo) Name() string {
	return "Wrong TLD"
}

func (n *Algo) Description() string {
	return "Wrong top level domain (TLD)"
}

func (n *Algo) Fields() []string {
	return []string{}
}

func (n *Algo) Headers() []string {
	return []string{}
}

func (n *Algo) Exec(in urlinsane.Typo) (out []urlinsane.Typo) {
	out = append(out, in)
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &Algo{
			types: []string{algorithms.DOMAIN},
		}
	})
}
