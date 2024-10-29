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
package ssd

// ssdeep is a program for computing context triggered piecewise hashes (CTPH). Also called fuzzy hashes, CTPH can match inputs that have homologies. Such inputs have sequences of identical bytes in the same order, although bytes in between these sequences may be different in both content and length.

// func ssdeepFunc(tr Result) (results []Result) {
// 	tr = checkIP(tr)
// 	if tr.Original.Live {
// 		var h1, h2 string
// 		h1, _ = ssdeep.FuzzyBytes([]byte(tr.Original.Meta.HTTP.Body))
// 		tr.Original.Meta.SSDeep = h1

// 		if tr.Variant.Live {
// 			h2, _ = ssdeep.FuzzyBytes([]byte(tr.Variant.Meta.HTTP.Body))
// 			tr.Variant.Meta.SSDeep = h2
// 		}

// 		if compare, err := ssdeep.Distance(h1, h2); err == nil {
// 			tr.Data["SIM"] = fmt.Sprintf("%d%s", compare, "%")
// 			tr.Variant.Meta.Similarity = compare
// 		}
// 	}
// 	results = append(results, tr)
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	"github.com/rangertaha/urlinsane/plugins/information"
)

const CODE = "ssd"

type None struct {
	types []string
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Name() string {
	return "SSDeep"
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
	return []string{"SSDEEP"}
}

func (n *None) Exec(in urlinsane.Typo) (out urlinsane.Typo) {
	in.Variant().Add("SSDEEP", []string{"one", "two"})
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
