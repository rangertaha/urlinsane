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
package txt

// // txtLookupFunc
// func txtLookupFunc(tr Result) (results []Result) {
// 	records, _ := net.LookupTXT(tr.Variant.String())
// 	tr.Variant.Meta.DNS.TXT = records
// 	for _, record := range records {
// 		tr.Data["TXT"] = strings.TrimSpace(tr.Data["TXT"] + "\n" + record)
// 		tr.Variant.Meta.DNS.TXT = append(tr.Variant.Meta.DNS.TXT, record)
// 	}
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const CODE = "txt"

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
	return []string{"TXT"}
}

func (n *None) Exec(in internal.Typo) (out internal.Typo) {
	in.Variant().Add("TXT", "Text records")
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &None{
			types: []string{internal.ENTITY, internal.DOMAIN},
		}
	})
}
