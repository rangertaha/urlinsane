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
package sp

// singularPluraliseFunc are created by making a singular domain plural and
// vice versa. For example, www.google.com becomes www.googles.com and
// www.games.co.nz becomes www.game.co.nz
// func singularPluraliseFunc(tc Result) (results []Result) {
// 	for _, pchar := range []string{"s", "ing"} {
// 		var domain string
// 		if strings.HasSuffix(tc.Original.Domain, pchar) {
// 			domain = strings.TrimSuffix(tc.Original.Domain, pchar)
// 		} else {
// 			domain = tc.Original.Domain + pchar
// 		}
// 		dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 		results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 	}
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const CODE = "sp"
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
	return "Singular Pluralize"
}

func (n *Algo) Description() string {
	return "Creates singular and plural names"
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
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
