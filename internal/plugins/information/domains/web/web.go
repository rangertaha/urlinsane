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
package web

// // httpLookupFunc
// func httpLookupFunc(tr Result) (results []Result) {
// 	if tr := checkIP(tr); tr.Variant.Live {
// 		httpReq, gerr := http.Get("http://" + tr.Variant.String())
// 		if gerr == nil {
// 			tr.Variant.Meta.HTTP = httpLib.NewResponse(httpReq)
// 			// spew.Dump(original)

// 			str := httpReq.Request.URL.String()
// 			subdomain := domainutil.Subdomain(str)
// 			domain := domainutil.DomainPrefix(str)
// 			suffix := domainutil.DomainSuffix(str)
// 			if domain == "" {
// 				domain = str
// 			}
// 			dm := Domain{subdomain, domain, suffix, tr.Variant.Meta, true}
// 			if tr.Variant.String() != dm.String() {
// 				tr.Data["Redirect"] = dm.String()
// 				tr.Variant.Meta.Redirect = dm.String()
// 			}
// 		}
// 	}
// 	results = append(results, tr)
// 	return
// }

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
)

const (
	CODE        = "web"
	NAME        = "Download Webpage"
	DESCRIPTION = "Retrieving the web page contents"
)

type None struct {
	types []string
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Description() string {
	return DESCRIPTION
}

func (n *None) Headers() []string {
	return []string{"HTTP"}
}

func (n *None) Exec(in internal.Typo) (out internal.Typo) {
	// in.Variant().Add("HTTP", []string{"one", "two"})
	return in
}

// Register the plugin
func init() {
	domains.Add(CODE, func() internal.Information {
		return &None{}
	})
}
