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
package si

// AlgoFunc typos are created by inserting common subdomains at the begining of the domain. wwwgoogle.com and ftpgoogle.com

// func AlgoFunc(tc Result) (results []Result) {
// 	for _, str := range datasets.SUBDOMAINS {
// 		dm := Domain{tc.Original.Subdomain, str + tc.Original.Domain, tc.Original.Suffix, Meta{}, false}
// 		results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 	}
// 	return results
// }


import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)


const CODE = "si"
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
	return "Subdomain Insertion"
}

func (n *Algo) Description() string {
	return "Inserts common subdomain at the beginning of the domain"
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


// // AlgoFunc typos are created by inserting common subdomains at the begining of the domain. wwwgoogle.com and ftpgoogle.com
// func AlgoFunc(tc Result) (results []Result) {
// 	for _, str := range datasets.SUBDOMAINS {
// 		dm := Domain{tc.Original.Subdomain, str + tc.Original.Domain, tc.Original.Suffix, Meta{}, false}
// 		results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 	}
// 	return results
// }
