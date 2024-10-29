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
package dcr

// AlgoFunc are created by replacing identical, consecutive
// letters of the domain name with adjacent letters on the keyboard.
// For example, www.gppgle.com and www.giigle.com
// func AlgoFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for i, char := range tc.Original.Domain {
// 			if i < len(tc.Original.Domain)-1 {
// 				if tc.Original.Domain[i] == tc.Original.Domain[i+1] {
// 					for _, key := range keyboard.Adjacent(string(char)) {
// 						domain := tc.Original.Domain[:i] + string(key) + string(key) + tc.Original.Domain[i+2:]
// 						dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 					}
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

const CODE = "dcr"
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
	return "Double Character Replacement"
}

func (n *Algo) Description() string {
	return "Created by replacing identical, consecutive letters in the name."
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
func (n *Algo) Domain(in urlinsane.Typo) (out []urlinsane.Typo) {
	return n.Exec(in)
}

func (n *Algo) Username(in urlinsane.Typo) (out []urlinsane.Typo) {
	return n.Exec(in)
}

func (n *Algo) Emmail(in urlinsane.Typo) (out []urlinsane.Typo) {
	return n.Exec(in)
}

func (n *Algo) Entity(in urlinsane.Typo) (out []urlinsane.Typo) {
	return n.Exec(in)
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &Algo{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
