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
package mx

// // mxLookupFunc
// func mxLookupFunc(tr Result) (results []Result) {
// 	records, _ := net.LookupMX(tr.Variant.String())
// 	tr.Variant.Meta.DNS.MX = dnsLib.NewMX(records...)
// 	for _, record := range records {
// 		record := strings.TrimSuffix(record.Host, ".")
// 		if !strings.Contains(tr.Data["MX"], record) {
// 			tr.Data["MX"] = strings.TrimSpace(tr.Data["MX"] + "\n" + record)
// 		}
// 	}
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	"github.com/rangertaha/urlinsane/plugins/information"
)

const CODE = "mx"

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
	return []string{"MX"}
}

func (n *None) Exec(in urlinsane.Typo) (out urlinsane.Typo) {
	in.Variant().Add("MX", "MXMXMX")
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
