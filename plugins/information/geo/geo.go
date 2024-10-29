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
package geo

// // geoIPLookupFunc
// func geoIPLookupFunc(tr Result) (results []Result) {
// 	tr = checkIP(tr)
// 	if tr.Variant.Live {
// 		_, ok := tr.Data["IPv4"]
// 		if ok {
// 			for _, ip4 := range tr.Variant.Meta.DNS.IPv4 {
// 				if ip4 != "" {
// 					record, _ := geoLib.GeoIP(ip4)
// 					// if err != nil {
// 					// 	fmt.Print(err)
// 					// }
// 					tr.Data["GEO"] = fmt.Sprint(record.Country.Names["en"])
// 					tr.Variant.Meta.Geo = *record
// 				}
// 			}
// 		}
// 	}

// 	// If you are using strings that may be invalid, check that ip is not nil
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	"github.com/rangertaha/urlinsane/plugins/information"
)

const CODE = "geo"

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
	return []string{"GEO"}
}

func (n *None) Exec(in urlinsane.Typo) (out urlinsane.Typo) {
	in.Variant().Add("GEO", "1000, 22333")
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
