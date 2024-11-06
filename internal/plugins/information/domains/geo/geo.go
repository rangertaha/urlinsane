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

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
)

const (
	ORDER       = 5
	CODE        = "geo"
	NAME        = "GEOIP Lookup"
	DESCRIPTION = "Retrieves Location of IP addresses"
)

type None struct {
	types []string
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Order() int {
	return ORDER
}

func (n *None) Name() string {
	return NAME
}

func (n *None) Description() string {
	return DESCRIPTION
}

func (n *None) Headers() []string {
	return []string{"GEO"}
}

func (n *None) Exec(in internal.Typo) (out internal.Typo) {
	// db, err := maxminddb.Open("../GeoLite2-City.mmdb")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// addr := netip.MustParseAddr("81.2.69.142")

	// var record struct {
	// 	Country struct {
	// 		ISOCode string `maxminddb:"iso_code"`
	// 	} `maxminddb:"country"`
	// } // Or any appropriate struct

	// err = db.Lookup(addr).Decode(&record)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// fmt.Print(record.Country.ISOCode)
	// // Output:
	// // GB
	return in
}

// Register the plugin
func init() {
	domains.Add(CODE, func() internal.Information {
		return &None{}
	})
}
