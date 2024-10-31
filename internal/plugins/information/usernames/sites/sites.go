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
package sites

// // cnameLookupFunc
// func cnameLookupFunc(tr Result) (results []Result) {
// 	records, _ := net.LookupCNAME(tr.Variant.String())
// 	// tr.Variant.Meta.DNS.CName = records
// 	for _, record := range records {
// 		tr.Data["CNAME"] = strings.TrimSuffix(string(record), ".")
// 		tr.Variant.Meta.DNS.CName = append(tr.Variant.Meta.DNS.CName, string(record))
// 	}
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information/usernames"
)

const CODE = "sites"

type None struct {
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Name() string {
	return "CNAME"
}

func (n *None) Description() string {
	return "DNS CNAME record"
}

func (n *None) Headers() []string {
	return []string{"CNAME"}
}

func (n *None) Exec(in internal.Typo) (out internal.Typo) {
	in.Variant().Add("CNAME", 111111)
	in.Variant().Add("JSON", []string{"one", "two"})
	return in
}

// Register the plugin
func init() {
	usernames.Add(CODE, func() internal.Information {
		return &None{}
	})
}
