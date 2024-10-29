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
package ld

// var levenshteinDistance = Module{
// 	Code:        "LD",
// 	Name:        "Levenshtein Distance",
// 	Description: "The Levenshtein distance between strings",
// 	Exe:         levenshteinDistanceFunc,
// 	Fields:      []string{"LD"},
// }

// // levenshteinDistanceFunc
// func levenshteinDistanceFunc(tr Result) (results []Result) {
// 	domain := tr.Original.String()
// 	variant := tr.Variant.String()
// 	tr.Data["LD"] = strconv.Itoa(nlpLib.Levenshtein(domain, variant))
// 	tr.Variant.Meta.Levenshtein = nlpLib.Levenshtein(domain, variant)
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	"github.com/rangertaha/urlinsane/plugins/information"
)

const CODE = "ld"

type None struct {
	types []string
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Name() string {
	return "None"
}

func (n *None) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *None) Description() string {
	return "Nothing"
}

func (n *None) Fields() []string {
	return []string{}
}

func (n *None) Headers() []string {
	return []string{"LD"}
}

func (n *None) Exec(in urlinsane.Typo) (out urlinsane.Typo) {
	in.Variant().Add("LD", 111111)
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() urlinsane.Information {
		return &None{
			types: []string{urlinsane.ENTITY, urlinsane.DOMAIN},
		}
	})
}
