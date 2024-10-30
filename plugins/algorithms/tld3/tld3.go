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
package tld3

// // wrongThirdLevelDomainFunc uses an alternate, valid third level domain.
// func wrongThirdLevelDomainFunc(tc Result) (results []Result) {
// 	labels := strings.Split(tc.Original.Suffix, ".")
// 	length := len(labels)
// 	for _, suffix := range datasets.TLD {
// 		suffixLbl := strings.Split(suffix, ".")
// 		suffixLen := len(suffixLbl)
// 		if length == suffixLen && length == 3 {
// 			if suffixLbl[1] == labels[1] && suffixLbl[2] == labels[2] {
// 				if suffix != tc.Original.Suffix {
// 					dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const (
	CODE        = "tld3"
	NAME        = "Wrong TLD3"
	DESCRIPTION = "Wrong third level domain (TLD3)"
)


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
	return NAME
}

func (n *Algo) Description() string {
	return DESCRIPTION
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
